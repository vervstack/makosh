package in_memory

import (
	"sync"

	"github.com/godverv/makosh/internal/domain"
)

type InMemoryDb struct {
	data map[string]*domain.Endpoint
	m    sync.RWMutex
}

func New() *InMemoryDb {
	return &InMemoryDb{
		data: make(map[string]*domain.Endpoint),
	}
}
