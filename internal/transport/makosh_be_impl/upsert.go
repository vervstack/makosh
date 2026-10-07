package makosh_be_impl

import (
	"context"

	"go.vervstack.ru/makosh/pkg/makosh_be"
	"go.vervstack.ru/makosh/pkg/registry"
)

func (impl *Impl) UpsertEndpoints(ctx context.Context, req *makosh_be.UpsertEndpoints_Request,
) (*makosh_be.UpsertEndpoints_Response, error) {
	err := impl.reg.Upsert(ctx, req.GetEndpoints())
	if err != nil {
		return nil, registry.ToStatusError(err)
	}

	return &makosh_be.UpsertEndpoints_Response{}, nil
}
