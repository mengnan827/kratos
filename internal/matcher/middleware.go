package matcher

import (
	"sort"
	"strings"

	"kratos_c/middleware"
)

type Matcher interface {
	Use(ms ...middleware.Middleware)
	Add(selector string, ms ...middleware.Middleware)
	Match(operation string) []middleware.Middleware
}

type matcher struct {
	prefix   []string
	defaults []middleware.Middleware
	matches  map[string][]middleware.Middleware
}

func New() Matcher {
	return &matcher{
		matches: make(map[string][]middleware.Middleware),
	}
}

func (m *matcher) Use(ms ...middleware.Middleware) {
	m.defaults = ms
}

func (m *matcher) Add(selector string, ms ...middleware.Middleware) {
	if strings.HasSuffix(selector, "*") {
		selector = strings.TrimSuffix(selector, "*")
		m.prefix = append(m.prefix, selector)
		sort.Slice(m.prefix, func(i, j int) bool { return m.prefix[i] > m.prefix[j] })
	}
	m.matches[selector] = ms
}

func (m *matcher) Match(operation string) []middleware.Middleware {
	// 添加默认的
	ms := make([]middleware.Middleware, 0, len(m.defaults))
	if len(m.defaults) > 0 {
		ms = append(ms, m.defaults...)
	}
	if next, ok := m.matches[operation]; ok {
		ms = append(ms, next...)
	} else {
		for _, prefix := range m.prefix {
			if strings.HasPrefix(operation, prefix) {
				ms = append(ms, m.matches[prefix]...)
			}
		}
	}
	return nil
}
