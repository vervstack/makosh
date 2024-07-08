package resolver

import (
	"os"
	"sync"

	errors "github.com/Red-Sock/trace-errors"
	"google.golang.org/grpc/resolver"
)

type Builder struct {
	srcUrl string
}

func New() (*Builder, error) {
	srcUrl := os.Getenv(MakoshSourceURL)
	if srcUrl == "" {
		return nil, errors.Wrap(ErrNoMakoshSourceURL, "Need env variable to resolve")
	}

	return &Builder{
		srcUrl: srcUrl,
	}, nil
}

func (b *Builder) Build(target resolver.Target, clientConn resolver.ClientConn, _ resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		target:     target.Endpoint(),
		wg:         sync.WaitGroup{},
		clientConn: clientConn,
	}
	r.wg.Add(1)

	go r.watch()
	return r, nil
}
