package resolver

import (
	"github.com/godverv/makosh/pkg/resolver/makosh_resolver"
)

func WithSchema(schemaName string) opt {
	return func(b *ServiceDiscovery) {
		b.schema = schemaName
	}
}

func WithResolverBuilder(builder makosh_resolver.ResolverBuilder) opt {
	return func(b *ServiceDiscovery) {
		b.resolverBuilder = builder
	}
}
