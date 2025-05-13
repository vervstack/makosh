package resolver

import (
	"sync"

	errors "github.com/Red-Sock/trace-errors"
	"google.golang.org/grpc/resolver"

	"go.vervstack.ru/makosh/pkg/resolver/makosh_resolver"
)

var (
	defaultServiceDiscovery         *ServiceDiscovery
	initDefaultServiceDiscoveryOnce sync.Once
)

func Init() (_ *ServiceDiscovery, err error) {
	initDefaultServiceDiscoveryOnce.Do(func() {
		var remoteResolver *makosh_resolver.RemoteResolverBuilder
		remoteResolver, err = makosh_resolver.NewBuilder()
		if err != nil {
			err = errors.Wrap(err, "error creating remote resolver")
			return
		}

		defaultServiceDiscovery, err = NewLocalServiceDiscovery(
			WithResolverBuilder(remoteResolver),
		)
		if err != nil {
			err = errors.Wrap(err, "error creating makosh resolver")
			return
		}

		resolver.Register(defaultServiceDiscovery.GrpcBuilder())
	})

	return defaultServiceDiscovery, err
}
