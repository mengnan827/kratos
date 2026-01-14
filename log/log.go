package log

type Logger interface {
	Log(level Level, keyVals ...any) error
}
