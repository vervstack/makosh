package makosh_be_impl

import (
	"context"

	"go.vervstack.ru/makosh/pkg/makosh_be"
	"go.vervstack.ru/makosh/pkg/registry"
)

func (impl *Impl) DeleteEndpoints(ctx context.Context, req *makosh_be.DeleteEndpoints_Request,
) (*makosh_be.DeleteEndpoints_Response, error) {
	err := impl.reg.Delete(ctx, req.GetServiceName(), req.GetScope())
	if err != nil {
		return nil, registry.ToStatusError(err)
	}

	return &makosh_be.DeleteEndpoints_Response{}, nil
}
