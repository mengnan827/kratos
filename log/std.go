package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// var _ Logger = ()
type stdLogger struct {
	w         io.Writer
	isDiscard bool
	mu        sync.Mutex
	pool      *sync.Pool
}

func NewStdLogger(w io.Writer) Logger {
	return &stdLogger{
		w:         w,
		isDiscard: w == io.Discard,
		pool: &sync.Pool{
			New: func() any {
				return new(bytes.Buffer)
			},
		},
	}
}

func (l *stdLogger) Log(level Level, keyVals ...any) error {
	if l.isDiscard || len(keyVals) == 0 {
		return nil
	}
	if len(keyVals)&1 == 1 {
		keyVals = append(keyVals, "KEYVALS UNPAIRED")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	defer l.pool.Put(buf)
	// 先输入级别
	buf.WriteString(level.String())
	for i := 0; i < len(keyVals); i += 2 {
		fmt.Fprintf(buf, " %s=%v", keyVals[i], keyVals[i+1])
	}
	buf.WriteByte('\n')
	l.mu.Lock()
	defer l.mu.Unlock()
	_, err := l.w.Write(buf.Bytes())
	return err
}
