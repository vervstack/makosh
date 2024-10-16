package makosh

import (
	"github.com/godverv/makosh/pkg/makosh/makosh_resolver"
)

func WithSchema(schemaName string) opt {
	return func(b *LocalServiceDiscovery) {
		b.schema = schemaName
	}
}

func WithResolverBuilder(builder makosh_resolver.ResolverBuilder) opt {
	return func(b *LocalServiceDiscovery) {
		b.resolverBuilder = builder
	}
}
