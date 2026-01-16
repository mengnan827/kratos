package grpc

import (
	"context"
	"crypto/tls"
	"kratos_c/log"
	"kratos_c/middleware"
	"net"
	"net/url"
	"time"

	"kratos_c/internal/matcher"

	"google.golang.org/grpc"
	"google.golang.org/grpc/admin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	*grpc.Server

	baseCtx  context.Context
	tlsConf  *tls.Config
	lis      net.Listener
	err      error
	network  string
	address  string
	endpoint *url.URL
	timeout  time.Duration

	middleware       matcher.Matcher
	streamMiddleware matcher.Matcher

	unaryInts  []grpc.UnaryServerInterceptor
	streamInts []grpc.StreamServerInterceptor

	grpcOpts     []grpc.ServerOption
	health       *health.Server
	customHealth bool
	// metadata          *map[string]string
	adminClean        func()
	disableReflection bool
}

type ServerOption func(o *Server)

// Network with server network.
func Network(network string) ServerOption {
	return func(s *Server) {
		s.network = network
	}
}

// Address with server address.
func Address(addr string) ServerOption {
	return func(s *Server) {
		s.address = addr
	}
}

// Endpoint with server address.
func Endpoint(endpoint *url.URL) ServerOption {
	return func(s *Server) {
		s.endpoint = endpoint
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.timeout = timeout
	}
}

// Logger with server logger.
// Deprecated: use global logger instead.
func Logger(logger log.Logger) ServerOption {
	return func(*Server) {}
}

// Middleware with server middleware.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.middleware.Use(m...)
	}
}

func StreamMiddleware(m ...middleware.Middleware) ServerOption {
	return func(s *Server) {
		s.streamMiddleware.Use(m...)
	}
}

// CustomHealth Checks server.
func CustomHealth() ServerOption {
	return func(s *Server) {
		s.customHealth = true
	}
}

// TLSConfig with TLS config.
func TLSConfig(c *tls.Config) ServerOption {
	return func(s *Server) {
		s.tlsConf = c
	}
}

// Listener with server lis
func Listener(lis net.Listener) ServerOption {
	return func(s *Server) {
		s.lis = lis
	}
}

// UnaryInterceptor returns a ServerOption that sets the UnaryServerInterceptor for the server.
func UnaryInterceptor(in ...grpc.UnaryServerInterceptor) ServerOption {
	return func(s *Server) {
		s.unaryInts = in
	}
}

// StreamInterceptor returns a ServerOption that sets the StreamServerInterceptor for the server.
func StreamInterceptor(in ...grpc.StreamServerInterceptor) ServerOption {
	return func(s *Server) {
		s.streamInts = in
	}
}

// DisableReflection disable grpc reflection.
func DisableReflection() ServerOption {
	return func(s *Server) {
		s.disableReflection = true
	}
}

// Options with grpc options.
func Options(opts ...grpc.ServerOption) ServerOption {
	return func(s *Server) {
		s.grpcOpts = opts
	}
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		baseCtx:          context.Background(),
		network:          "tcp",
		address:          ":0",
		timeout:          time.Second,
		health:           health.NewServer(),
		middleware:       matcher.New(),
		streamMiddleware: matcher.New(),
	}
	for _, o := range opts {
		o(srv)
	}

	unaryInts := []grpc.UnaryServerInterceptor{
		srv.unaryServerInterceptor(),
	}
	if len(srv.unaryInts) > 0 {
		unaryInts = append(unaryInts, srv.unaryInts...)
	}
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(unaryInts...),
	}
	if srv.tlsConf != nil {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(srv.tlsConf)))
	}
	if len(srv.grpcOpts) > 0 {
		grpcOpts = append(grpcOpts, srv.grpcOpts...)
	}
	srv.Server = grpc.NewServer(grpcOpts...)
	if !srv.disableReflection {
		reflection.Register(srv.Server)
	}
	srv.adminClean, _ = admin.Register(srv.Server)
	return srv
}
