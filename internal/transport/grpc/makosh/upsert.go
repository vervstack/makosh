package makosh

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh-be/internal/domain"
	"github.com/godverv/makosh-be/pkg/makosh_be"
)

func (impl *Implementation) UpsertEndpoints(ctx context.Context, req *makosh_be.UpsertEndpoints_Request,
) (*makosh_be.UpsertEndpoints_Response, error) {

	endpoints := make([]domain.Endpoint, 0, len(req.Endpoints))

	for _, e := range req.GetEndpoints() {
		endpoints = append(endpoints,
			domain.Endpoint{
				ServiceName: e.ServiceName,
				Addrs:       e.Addrs,
			})
	}

	err := impl.data.Save(ctx, endpoints...)
	if err != nil {
		return nil, errors.Wrap(err, "error saving endpoints")
	}

	return nil, nil
}
