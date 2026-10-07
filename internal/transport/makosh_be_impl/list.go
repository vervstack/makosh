package makosh_be_impl

import (
	"context"

	"go.vervstack.ru/makosh/pkg/makosh_be"
	"go.vervstack.ru/makosh/pkg/registry"
)

func (impl *Impl) ListEndpoints(ctx context.Context, req *makosh_be.ListEndpoints_Request) (*makosh_be.ListEndpoints_Response, error) {
	response, err := impl.reg.List(ctx, req.GetServiceName(), req.GetScope())
	if err != nil {
		return nil, registry.ToStatusError(err)
	}

	return response, nil
}
