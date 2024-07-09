package makosh

import (
	"context"

	errors "github.com/Red-Sock/trace-errors"

	"github.com/godverv/makosh/pkg/makosh_be"
)

func (impl *Implementation) ListEndpoints(ctx context.Context, req *makosh_be.ListEndpoints_Request) (*makosh_be.ListEndpoints_Response, error) {
	endpoints, err := impl.data.Get(ctx, req.ServiceName)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &makosh_be.ListEndpoints_Response{
		Urls: endpoints.Addrs,
	}, nil
}
