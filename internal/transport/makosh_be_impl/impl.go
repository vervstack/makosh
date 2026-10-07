package makosh_be_impl

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"go.vervstack.ru/makosh/internal/config"
	"go.vervstack.ru/makosh/pkg/makosh_be"
	"go.vervstack.ru/makosh/pkg/registry"
)

type Impl struct {
	makosh_be.UnimplementedMakoshBeAPIServer

	version string
	reg     *registry.Registry
}

func New(cfg config.Config, reg *registry.Registry) *Impl {
	return &Impl{
		version: cfg.AppInfo.Version,
		reg:     reg,
	}
}

func (impl *Impl) Register(server grpc.ServiceRegistrar) {
	makosh_be.RegisterMakoshBeAPIServer(server, impl)
}

func (impl *Impl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (route string, handler http.Handler) {
	gwHttpMux := runtime.NewServeMux()

	err := makosh_be.RegisterMakoshBeAPIHandlerFromEndpoint(
		ctx,
		gwHttpMux,
		endpoint,
		opts,
	)
	if err != nil {
		logrus.Errorf("error registering grpc2http handler: %s", err)
	}

	return "/api/", gwHttpMux
}
