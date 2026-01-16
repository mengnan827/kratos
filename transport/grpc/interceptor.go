package grpc

import (
	"context"
	ic "kratos_c/internal/context"
	"kratos_c/middleware"
	"kratos_c/transport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		ctx, cancel := ic.Merge(ctx, s.baseCtx)
		defer cancel()
		md, _ := metadata.FromIncomingContext(ctx)
		replyHeader := metadata.MD{}
		tr := &Transport{
			operation:   info.FullMethod,
			reqHeader:   headerCarrier(md),
			replyHeader: headerCarrier(replyHeader),
		}
		if s.endpoint != nil {
			tr.endpoint = s.endpoint.String()
		}
		ctx = transport.NewServerContext(ctx, tr)
		if s.timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, s.timeout)
			defer cancel()
		}
		h := func(ctx context.Context, req any) (any, error) {
			return handler(ctx, req)
		}
		if next := s.middleware.Match(tr.Operation()); len(next) > 0 {
			h = middleware.Chain(next...)(h)
		}
		reply, err := h(ctx, req)
		if len(replyHeader) > 0 {
			_ = grpc.SetHeader(ctx, replyHeader)
		}
		return reply, err
	}
}
