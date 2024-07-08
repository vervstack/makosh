package makosh

import (
	"context"

	"github.com/godverv/makosh-be/pkg/makosh_be"
)

func (impl *Implementation) Version(_ context.Context, _ *makosh_be.Version_Request) (*makosh_be.Version_Response, error) {
	return &makosh_be.Version_Response{
		Version: impl.version,
	}, nil
}
