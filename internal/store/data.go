package store

import (
	"context"

	"github.com/godverv/makosh-be/internal/domain"
)

type Data interface {
	Get(ctx context.Context, serviceName string) (domain.Endpoint, error)
	Save(ctx context.Context, endpoints ...domain.Endpoint) error
}
