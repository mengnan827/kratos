package log

import "sync"

var global = &loggerAppliance{}

type loggerAppliance struct {
	lock sync.RWMutex
	Logger
}

func (a *loggerAppliance) SetLogger(in Logger) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.Logger = in
}

func SetLogger(logger Logger) {
	global.SetLogger(logger)
}
