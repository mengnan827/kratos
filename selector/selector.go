package selector

type Node interface {
	Scheme() string
	Address() string
	ServiceName() string
	InitialWeight() *int64
	Version() string
	Metadata() map[string]string
}
