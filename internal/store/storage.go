package store

import (
	"context"

	"go.vervstack.ru/makosh/internal/domain"
)

type EndpointsStorage interface {
	Get(ctx context.Context, serviceName string) (domain.Endpoint, error)
	Save(ctx context.Context, endpoints ...domain.Endpoint) error
}
