package makosh_resolver

type opt func(b *RemoteResolverBuilder)

func WithURL(url string) opt {
	return func(b *RemoteResolverBuilder) {
		b.remoteServiceDiscoveryURL = url
	}
}

func WithSecret(secret string) opt {
	return func(b *RemoteResolverBuilder) {
		b.secret = secret
	}
}

func WithPublicServiceDiscovery() opt {
	return func(b *RemoteResolverBuilder) {
		b.isPublic = true
	}
}
