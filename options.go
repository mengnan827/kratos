package kratos_c

import (
	"net/url"
	"os"
	"time"
	"kratos_c/log"
	"kratos_c/registry"
	"kratos_c/transport"

	"golang.org/x/net/context"
)

type Option func(o *options)

type options struct {
	id       string
	name     string
	version  string
	metadata map[string]string

	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	logger           log.Logger
	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	servers          []transport.Server

	beforeStart []func(ctx context.Context) error
	beforeStop  []func(ctx context.Context) error
	afterStart  []func(ctx context.Context) error
	afterStop   []func(ctx context.Context) error
}

// ID 设置服务ID
func ID(id string) Option {
	return func(o *options) { o.id = id }
}

// Name 设置服务名称
func Name(name string) Option {
	return func(o *options) { o.name = name }
}

// Version 设置服务版本
func Version(version string) Option {
	return func(o *options) { o.version = version }
}

// Metadata 设置服务元数据
func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

// Endpoint 设置服务端点
func Endpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

// Context 设置服务上下文
func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// Logger 设置服务日志记录器
func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// Server 设置传输服务器
func Server(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

// Signal 设置退出信号
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Registrar 设置服务注册器
func Registrar(r registry.Registrar) Option {
	return func(o *options) { o.registrar = r }
}

// RegistrarTimeout 设置注册超时时间
func RegistrarTimeout(t time.Duration) Option {
	return func(o *options) { o.registrarTimeout = t }
}

// StopTimeout 设置应用停止超时时间
func StopTimeout(t time.Duration) Option {
	return func(o *options) { o.stopTimeout = t }
}

// BeforeStart 在应用启动前执行函数
func BeforeStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

// BeforeStop 在应用停止前执行函数
func BeforeStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

// AfterStart 在应用启动后执行函数
func AfterStart(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

// AfterStop 在应用停止后执行函数
func AfterStop(fn func(context.Context) error) Option {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}
