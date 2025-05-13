package makosh_be_impl

import (
	"context"

	"go.vervstack.ru/makosh/pkg/makosh_be"
)

func (impl *Impl) Version(_ context.Context, _ *makosh_be.Version_Request) (*makosh_be.Version_Response, error) {
	return &makosh_be.Version_Response{
		Version: impl.version,
	}, nil
}
