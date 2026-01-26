package group

import "sync"

type Factory[T any] func() T

type Group[T any] struct {
	factory func() T
	vals    map[string]T
	sync.RWMutex
}

func NewGroup[T any](factory Factory[T]) *Group[T] {
	if factory == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	return &Group[T]{
		factory: factory,
		vals:    make(map[string]T),
	}
}

func (g *Group[T]) Get(key string) T {
	g.RWMutex.RLock()
	v, ok := g.vals[key]
	if ok {
		g.RUnlock()
		return v
	}
	g.RUnlock()

	g.Lock()
	defer g.Unlock()
	v, ok = g.vals[key]
	if ok {
		return v
	}
	v = g.factory()
	g.vals[key] = v
	return v
}

func (g *Group[T]) Reset(factory Factory[T]) {
	if factory == nil {
		panic("container.group: can't assign a nil to the new function")
	}
	g.Lock()
	g.factory = factory
	g.Unlock()
	g.Clear()
}

func (g *Group[T]) Clear() {
	g.Lock()
	g.vals = make(map[string]T)
	g.Unlock()
}
