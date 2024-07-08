package makosh

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/godverv/makosh-be/internal/config"
	"github.com/godverv/makosh-be/internal/store"
	"github.com/godverv/makosh-be/pkg/makosh_be"
)

type Implementation struct {
	makosh_be.UnimplementedMakoshBeAPIServer

	data    store.Data
	version string
}

func New(cfg config.Config, data store.Data) *Implementation {
	return &Implementation{
		data:    data,
		version: cfg.GetAppInfo().Version,
	}
}

func (impl *Implementation) Register(server grpc.ServiceRegistrar) {
	makosh_be.RegisterMakoshBeAPIServer(server, impl)
}

func (impl *Implementation) RegisterGw(ctx context.Context, mux *runtime.ServeMux, addr string) error {
	return makosh_be.RegisterMakoshBeAPIHandlerFromEndpoint(
		ctx,
		mux,
		addr,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		})
}
