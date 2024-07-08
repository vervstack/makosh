package makosh_resolver

import (
	"github.com/sirupsen/logrus"
)

func WithMakoshURL(url string) opt {
	return func(b *Builder) {
		b.makoshUrl = url
	}
}

func WithMakoshSecret(secret string) opt {
	return func(b *Builder) {
		b.secret = secret
	}
}

func WithLogger(logger logrus.StdLogger) opt {
	return func(b *Builder) {
		b.logger = logger
	}
}
