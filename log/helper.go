package log

var DefaultMessageKey = "msg"

type Option func()

type Helper struct {
	logger  Logger
	msgKey  string
	sprint  func(...any) string
	sprintf func(format string, a ...any) string
}

func NewHelper(logger Logger, opts ...Option) *Helper {
	return nil
}
