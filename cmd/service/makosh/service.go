package makosh

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"

	"github.com/godverv/makosh/internal/config"
	"github.com/godverv/makosh/internal/store/in_memory"
	"github.com/godverv/makosh/internal/transport/grpc"
	"github.com/godverv/makosh/internal/transport/grpc/makosh"
	"github.com/godverv/makosh/internal/utils/closer"
)

type App struct {
	server *grpc.Server
}

func New() (App, error) {
	logrus.Println("starting app")

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return App{}, errors.Wrap(err, "error reading config")
	}

	if cfg.GetAppInfo().StartupDuration == 0 {
		return App{}, errors.New("no startup duration in config")
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.GetAppInfo().StartupDuration)
	closer.Add(func() error {
		cancel()
		return nil
	})

	grpcServer, err := cfg.GetServers().GRPC(config.ServerGrpc)
	if err != nil {
		return App{}, errors.Wrap(err, "error getting grpc server from config")
	}

	inMemoryStore := in_memory.New()

	a := App{}

	a.server, err = grpc.NewServer(grpcServer, makosh.New(cfg, inMemoryStore))
	if err != nil {
		return App{}, errors.Wrap(err, "error initializing server")
	}

	return a, nil
}

func (a *App) Start() error {
	ctx := context.Background()

	err := a.server.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "error starting server")
	}

	waitingForTheEnd()

	logrus.Println("shutting down the app")

	err = closer.Close()
	if err != nil {
		logrus.Fatalf("errors while shutting down application %s", err.Error())
	}

	return err
}

// rscli comment: an obligatory function for tool to work properly.
// must be called in the main function above
// also this is the LP's song name reference, so no linting rules can be applied to the function name
func waitingForTheEnd() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
