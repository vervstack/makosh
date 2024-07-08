package in_memory

import (
	"context"

	"github.com/godverv/makosh-be/internal/domain"
)

func (d *InMemoryDb) Save(_ context.Context, endpoints ...domain.Endpoint) error {
	d.m.Lock()
	for i := range endpoints {
		d.data[endpoints[i].ServiceName] = &endpoints[i]
	}
	d.m.Unlock()
	return nil
}
