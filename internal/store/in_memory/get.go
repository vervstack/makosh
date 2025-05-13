package in_memory

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"

	"go.vervstack.ru/makosh/internal/domain"
)

func (d *InMemoryDb) Get(_ context.Context, serviceName string) (domain.Endpoint, error) {
	d.m.RLock()
	endpoint := d.data[serviceName]
	d.m.RUnlock()

	if endpoint == nil {
		return domain.Endpoint{}, errors.Wrap(domain.ErrNotFound, "service "+serviceName+" not found")
	}

	endpointCopy := domain.Endpoint{
		ServiceName: serviceName,
		Addrs:       make([]string, 0, len(endpoint.Addrs)),
	}

	for _, a := range endpoint.Addrs {
		endpointCopy.Addrs = append(endpointCopy.Addrs, a)
	}

	return endpointCopy, nil
}
