package grpc

import (
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"
	"google.golang.org/grpc/resolver"

	"go.vervstack.ru/makosh/pkg/resolver/makosh_resolver"
)

type serviceDiscovery interface {
	GetResolver(target string) (*atomic.Pointer[makosh_resolver.EndpointsResolver], error)
}

type Builder struct {
	localSD serviceDiscovery
	schema  string
}

func NewBuilder(src serviceDiscovery, schema string) *Builder {
	return &Builder{
		localSD: src,
		schema:  schema,
	}
}

// Build - implements grpc resolver.Builder
// Builds resolver for each host name
func (b *Builder) Build(t resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	resolverPtr, err := b.localSD.GetResolver(t.URL.Host)
	if err != nil {
		return nil, errors.Wrap(err, "error getting resolver from local service discovery")
	}

	reslvr := resolverPtr.Load()

	(*reslvr).AddSubscribers(updateGrpcCallback(cc))
	err = reslvr.Resolve()
	if err != nil {
		return nil, errors.Wrap(err, "error resolving endpoint")
	}
	return NewGrpcResolver(resolverPtr), nil
}

func (b *Builder) Scheme() string {
	return b.schema
}
