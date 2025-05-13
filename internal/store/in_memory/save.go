package in_memory

import (
	"context"

	"go.vervstack.ru/makosh/internal/domain"
)

func (d *InMemoryDb) Save(_ context.Context, endpoints ...domain.Endpoint) error {
	d.m.Lock()
	for i := range endpoints {
		d.data[endpoints[i].ServiceName] = &endpoints[i]
	}
	d.m.Unlock()
	return nil
}
