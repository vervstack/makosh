package grpc

import (
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"
	"google.golang.org/grpc/resolver"

	"github.com/godverv/makosh/pkg/makosh/makosh_resolver"
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
	// If already have resolver for this target - return it
	resolverPtr, err := b.localSD.GetResolver(t.URL.Host)
	if err != nil {
		return nil, errors.Wrap(err, "error getting resolver from local service discovery")
	}

	reslvr := resolverPtr.Load()

	(*reslvr).AddUpdateCallbacks(updateGrpcCallback(cc))

	return NewGrpcResolver(resolverPtr), nil
}

func (b *Builder) Scheme() string {
	return b.schema
}
