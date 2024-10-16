package makosh

import (
	"sync"
	"sync/atomic"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh/pkg/makosh/grpc"
	"github.com/godverv/makosh/pkg/makosh/makosh_resolver"
)

const DefaultSchema = "verv"

type LocalServiceDiscovery struct {
	schema    string
	overrides map[string][]string

	m         sync.Mutex
	resolvers map[string]*atomic.Pointer[makosh_resolver.EndpointsResolver]

	resolverBuilder makosh_resolver.ResolverBuilder
}

type opt func(b *LocalServiceDiscovery)

// NewLocalServiceDiscovery creates builder for local service discovery
func NewLocalServiceDiscovery(opts ...opt) (*LocalServiceDiscovery, error) {
	b := &LocalServiceDiscovery{
		overrides: make(map[string][]string),

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

func (b *LocalServiceDiscovery) SetCustomResolver(
	newResolver makosh_resolver.EndpointsResolver, serviceNames ...string) error {
	b.m.Lock()
	defer b.m.Unlock()

	for _, serviceName := range serviceNames {
		rPtr := b.resolvers[serviceName]
		if rPtr == nil {
			rPtr = &atomic.Pointer[makosh_resolver.EndpointsResolver]{}
			b.resolvers[serviceName] = rPtr
		} else {
			oldResolver := *rPtr.Load()
			newResolver.AddUpdateCallbacks(oldResolver.GetUpdaters()...)
		}

		rPtr.Store(&newResolver)
	}

	return nil
}

func (b *LocalServiceDiscovery) GrpcBuilder() *grpc.Builder {
	return grpc.NewBuilder(b, b.schema)
}

func (b *LocalServiceDiscovery) GetResolver(target string) (*atomic.Pointer[makosh_resolver.EndpointsResolver], error) {
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
