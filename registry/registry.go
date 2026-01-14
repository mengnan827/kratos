package registry

import "golang.org/x/net/context"

type ServiceInstance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Metadata  map[string]string `json:"metadata"`
	Endpoints []string          `json:"endpoints"`
}

type Registrar interface {
	Register(ctx context.Context, service *ServiceInstance) error
	Deregister(ctx context.Context, service *ServiceInstance) error
}
