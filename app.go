package kratos_c

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"kratos_c/log"
	"kratos_c/registry"
	"kratos_c/transport"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

// app信息
type AppInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

type App struct {
	opts options

	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

func New(opts ...Option) *App {
	o := options{
		ctx:              context.Background(),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: time.Second * 10,
	}
	// 唯一id
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.logger != nil {
		log.SetLogger(o.logger)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

// ID 返回应用实例ID
func (a *App) ID() string { return a.opts.id }

// Name 返回服务名称
func (a *App) Name() string { return a.opts.name }

// Version 返回应用版本
func (a *App) Version() string { return a.opts.version }

// Metadata 返回服务元数据
func (a *App) Metadata() map[string]string { return a.opts.metadata }

// Endpoint 返回端点列表
func (a *App) Endpoint() []string {
	// 如果实例不为nil，返回实例的端点列表
	if a.instance != nil {
		return a.instance.Endpoints
	}
	// 否则返回nil
	return nil
}

func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0, len(a.opts.endpoints))
	for _, o := range a.opts.endpoints {
		endpoints = append(endpoints, o.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range a.opts.servers {
			if r, ok := srv.(transport.Endpointer); ok {
				e, err := r.Endpoint()
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, e.String())
			}
		}
	}
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
}

func (a *App) Run() error {
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()

	// 此上下文主要用于钩子函数
	sCtx := NewContext(a.ctx, a)
	eg, ctx := errgroup.WithContext(sCtx)
	wg := sync.WaitGroup{}
	// 启动前
	for _, fn := range a.opts.beforeStart {
		if err = fn(sCtx); err != nil {
			return err
		}
	}
	// 创建操作上下文, 所以这里使用的是传进来的
	oCtx := NewContext(a.opts.ctx, a)
	for _, srv := range a.opts.servers {
		server := srv
		eg.Go(func() error {
			// 接收到停止信号
			<-ctx.Done()
			stopCtx := oCtx
			if a.opts.stopTimeout > 0 {
				var cancel context.CancelFunc
				stopCtx, cancel = context.WithTimeout(oCtx, a.opts.stopTimeout)
				defer cancel()
			}
			return server.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done() // 异步start的关键
			return server.Start(oCtx)
		})
	}
	wg.Wait()
	// 服务注册
	if a.opts.registrar != nil {
		rCtx, rCancel := context.WithTimeout(ctx, a.opts.registrarTimeout)
		defer rCancel()
		if err = a.opts.registrar.Register(rCtx, instance); err != nil {
			return err
		}
	}
	// 执行启动后钩子
	for _, fn := range a.opts.afterStart {
		if err = fn(sCtx); err != nil {
			return err
		}
	}

	// 监听停止信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})
	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	err = nil
	for _, fn := range a.opts.afterStop {
		err = fn(sCtx)
	}
	return err
}

func (a *App) Stop() (err error) {
	sCtx := NewContext(a.ctx, a)
	for _, fn := range a.opts.beforeStop {
		if err = fn(sCtx); err != nil {
			return err
		}
	}
	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	if a.opts.registrar != nil && instance != nil {
		rCtx, rCancel := context.WithTimeout(a.ctx, a.opts.registrarTimeout)
		defer rCancel()
		if err = a.opts.registrar.Deregister(rCtx, instance); err != nil {
			return err
		}
	}

	if a.cancel != nil {
		a.cancel()
	}
	return err

}
