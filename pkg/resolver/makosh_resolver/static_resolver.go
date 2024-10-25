package makosh_resolver

import (
	errors "github.com/Red-Sock/trace-errors"
)

type StaticResolver struct {
	EndpointsContainer
}

func NewStaticResolver(addrs ...string) *StaticResolver {
	return &StaticResolver{
		EndpointsContainer{
			addrsCache: addrs,
		},
	}
}

func (s *StaticResolver) Resolve() error {
	err := s.notifySubscribers()
	if err != nil {
		return errors.Wrap(err, "error notifying subscribers")
	}

	return nil
}
