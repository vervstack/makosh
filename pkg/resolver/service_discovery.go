package resolver

import (
	"sync"
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"
	"go.verv.tech/matreshka/service_discovery"

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

func (sd *ServiceDiscovery) SetCustomResolver(newResolver makosh_resolver.Resolver, serviceName string) {
	sd.setCustomResolver(newResolver, serviceName)
}

func (sd *ServiceDiscovery) GrpcBuilder() *grpc.Builder {
	return grpc.NewBuilder(sd, sd.schema)
}

func (sd *ServiceDiscovery) GetResolver(target string) (*atomic.Pointer[makosh_resolver.EndpointsResolver], error) {
	sd.m.Lock()
	defer sd.m.Unlock()

	rPtr := sd.resolvers[target]
	if rPtr != nil {
		return rPtr, nil
	}

	r, err := sd.resolverBuilder.NewResolver(target)
	if err != nil {
		return nil, errors.Wrap(err, "error building resolver")
	}

	rPtr = &atomic.Pointer[makosh_resolver.EndpointsResolver]{}
	rPtr.Store(&r)

	sd.resolvers[target] = rPtr

	return rPtr, nil
}

func (sd *ServiceDiscovery) SetOverrides(overrides service_discovery.Overrides) {
	for _, o := range overrides {
		resolver := makosh_resolver.NewStaticResolver(o.Urls...)
		sd.setCustomResolver(resolver, o.ServiceName)
	}
}

func (sd *ServiceDiscovery) setCustomResolver(newResolver makosh_resolver.Resolver, serviceName string) {
	sd.m.Lock()
	defer sd.m.Unlock()

	rPtr := sd.resolvers[serviceName]
	if rPtr == nil {
		rPtr = &atomic.Pointer[makosh_resolver.EndpointsResolver]{}
		sd.resolvers[serviceName] = rPtr
	} else {
		oldResolver := *rPtr.Load()
		newResolver.AddSubscribers(oldResolver.GetSubscribers()...)
	}

	rPtr.Store(&makosh_resolver.EndpointsResolver{Resolver: newResolver})
}
