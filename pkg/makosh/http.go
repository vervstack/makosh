package makosh

import (
	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh/pkg/makosh/makosh_resolver"
)

func (b *LocalServiceDiscovery) BuildHTTPResolver(targetName string, addressUpdater makosh_resolver.UpdateAddresses) (
	makosh_resolver.EndpointsResolver, error) {
	b.m.Lock()

	rPtr, err := b.resolverBuilder.NewResolver(targetName)
	if err != nil {
		return nil, errors.Wrap(err, "error building resolver")
	}
	rPtr.AddUpdateCallbacks(addressUpdater)

	return rPtr, nil
}
