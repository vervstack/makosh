package resolver

import (
	"sync"
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh/pkg/resolver/grpc"
	"github.com/godverv/makosh/pkg/resolver/makosh_resolver"
)

const DefaultSchema = "verv"

type ServiceDiscovery struct {
	schema string

	m         sync.Mutex
	resolvers map[string]*atomic.Pointer[makosh_resolver.EndpointsResolver]

	resolverBuilder makosh_resolver.ResolverBuilder
}

type opt func(b *ServiceDiscovery)

// NewLocalServiceDiscovery creates builder for local service discovery
func NewLocalServiceDiscovery(opts ...opt) (*ServiceDiscovery, error) {
	b := &ServiceDiscovery{
		schema:    DefaultSchema,
		resolvers: map[string]*atomic.Pointer[makosh_resolver.EndpointsResolver]{},
	}

	for _, o := range opts {
		o(b)
	}

	if b.resolverBuilder == nil {
		var err error
		b.resolverBuilder, err = makosh_resolver.NewBuilder()
		if err != nil {
			return nil, errors.Wrap(err, "error creating default resolver (verv makosh")
		}
	}

	return b, nil
}

func (b *ServiceDiscovery) SetCustomResolver(
	newResolver makosh_resolver.Resolver, serviceNames ...string) error {
	b.m.Lock()
	defer b.m.Unlock()

	for _, serviceName := range serviceNames {
		rPtr := b.resolvers[serviceName]
		if rPtr == nil {
			rPtr = &atomic.Pointer[makosh_resolver.EndpointsResolver]{}
			b.resolvers[serviceName] = rPtr
		} else {
			oldResolver := *rPtr.Load()
			newResolver.AddSubscribers(oldResolver.GetSubscribers()...)
		}

		rPtr.Store(&makosh_resolver.EndpointsResolver{Resolver: newResolver})
	}

	return nil
}

func (b *ServiceDiscovery) GrpcBuilder() *grpc.Builder {
	return grpc.NewBuilder(b, b.schema)
}

func (b *ServiceDiscovery) GetResolver(target string) (*atomic.Pointer[makosh_resolver.EndpointsResolver], error) {
	b.m.Lock()
	defer b.m.Unlock()

	rPtr := b.resolvers[target]
	if rPtr != nil {
		return rPtr, nil
	}

	r, err := b.resolverBuilder.NewResolver(target)
	if err != nil {
		return nil, errors.Wrap(err, "error building resolver")
	}

	rPtr = &atomic.Pointer[makosh_resolver.EndpointsResolver]{}
	rPtr.Store(&r)

	b.resolvers[target] = rPtr

	return rPtr, nil
}
