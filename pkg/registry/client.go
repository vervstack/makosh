package registry

import (
	"context"

	"google.golang.org/grpc"

	"go.vervstack.ru/makosh/pkg/makosh_be"
)

type client struct {
	reg     *Registry
	version string
}

// NewClient adapts a Registry to the generated client interface so an embedding process talks to the registry in-process, with the same errors as the gRPC server.
func NewClient(r *Registry, version string) makosh_be.MakoshBeAPIClient {
	return &client{
		reg:     r,
		version: version,
	}
}

func (c *client) Version(_ context.Context, _ *makosh_be.Version_Request, _ ...grpc.CallOption,
) (*makosh_be.Version_Response, error) {
	return &makosh_be.Version_Response{Version: c.version}, nil
}

func (c *client) ListEndpoints(ctx context.Context, req *makosh_be.ListEndpoints_Request, _ ...grpc.CallOption,
) (*makosh_be.ListEndpoints_Response, error) {
	response, err := c.reg.List(ctx, req.GetServiceName(), req.GetScope())
	if err != nil {
		return nil, ToStatusError(err)
	}

	return response, nil
}

func (c *client) UpsertEndpoints(ctx context.Context, req *makosh_be.UpsertEndpoints_Request, _ ...grpc.CallOption,
) (*makosh_be.UpsertEndpoints_Response, error) {
	err := c.reg.Upsert(ctx, req.GetEndpoints())
	if err != nil {
		return nil, ToStatusError(err)
	}

	return &makosh_be.UpsertEndpoints_Response{}, nil
}

func (c *client) DeleteEndpoints(ctx context.Context, req *makosh_be.DeleteEndpoints_Request, _ ...grpc.CallOption,
) (*makosh_be.DeleteEndpoints_Response, error) {
	err := c.reg.Delete(ctx, req.GetServiceName(), req.GetScope())
	if err != nil {
		return nil, ToStatusError(err)
	}

	return &makosh_be.DeleteEndpoints_Response{}, nil
}
