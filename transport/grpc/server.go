package grpc

import (
	"context"
	"crypto/tls"
	"net"
	"net/url"
	"time"

	"zone/test_demo/kratos_c/internal/matcher"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
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
	metadata     *map[string]string
}
