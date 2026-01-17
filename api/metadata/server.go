package metadata

import (
	"sync"

	"google.golang.org/grpc"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

type Server struct {
	UnimplementedMetadataServer

	lock     sync.Mutex
	srv      *grpc.Server
	services map[string]*dpb.FileDescriptorSet
	methods  map[string][]string
}

// NewServer 创建一个元数据服务
func NewServer(srv *grpc.Server) *Server {
	return &Server{
		srv:      srv,
		services: make(map[string]*dpb.FileDescriptorSet),
		methods:  make(map[string][]string),
	}
}

func (s *Server) load() error {
}
