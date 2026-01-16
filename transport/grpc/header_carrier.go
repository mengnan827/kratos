package grpc

import "google.golang.org/grpc/metadata"

type headerCarrier metadata.MD

func (mc headerCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (mc headerCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

func (mc headerCarrier) Add(key string, value string) {
	metadata.MD(mc).Append(key, value)
}

func (mc headerCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}

func (mc headerCarrier) Values(key string) []string {
	return metadata.MD(mc).Get(key)
}
