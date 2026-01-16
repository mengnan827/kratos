package grpc

import (
	"kratos_c/selector"
	"kratos_c/transport"
)

type Transport struct {
	endpoint    string
	operation   string
	reqHeader   headerCarrier
	replyHeader headerCarrier
	nodeFilters []selector.NodeFilter
}

func (tr *Transport) Kind() transport.Kind {
	return transport.KindGRPC
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) Operation() string {
	return tr.operation
}

func (tr *Transport) RequestHeader() transport.Header {
	return tr.reqHeader
}

func (tr *Transport) ReplyHeader() transport.Header {
	return tr.replyHeader
}

func (tr *Transport) NodeFilters() []selector.NodeFilter {
	return tr.nodeFilters
}
