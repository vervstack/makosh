package store

import (
	"context"

	"github.com/godverv/makosh/internal/domain"
)

type EndpointsStorage interface {
	Get(ctx context.Context, serviceName string) (domain.Endpoint, error)
	Save(ctx context.Context, endpoints ...domain.Endpoint) error
}
