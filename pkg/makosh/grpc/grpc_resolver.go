package grpc

import (
	"sync/atomic"

	"google.golang.org/grpc/resolver"

	"github.com/godverv/makosh/pkg/makosh/makosh_resolver"
)

type Resolver struct {
	resolverPtr *atomic.Pointer[makosh_resolver.EndpointsResolver]
}

func NewGrpcResolver(
	resolverPtr *atomic.Pointer[makosh_resolver.EndpointsResolver],
) *Resolver {
	return &Resolver{
		resolverPtr: resolverPtr,
	}
}

func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {
	_ = (*r.resolverPtr.Load()).Resolve()
}

func (r *Resolver) Close() {}

func updateGrpcCallback(cc resolver.ClientConn) makosh_resolver.UpdateAddresses {
	return func(addrs []string) error {
		var state resolver.State

		state.Addresses = make([]resolver.Address, 0, len(addrs))

		for _, a := range addrs {
			state.Addresses = append(state.Addresses, resolver.Address{Addr: a})
		}

		return cc.UpdateState(state)
	}
}
