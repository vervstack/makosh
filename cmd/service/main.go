package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	//_transport_imports

	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"

	"github.com/godverv/makosh-be/internal/config"
	"github.com/godverv/makosh-be/internal/store/in_memory"
	"github.com/godverv/makosh-be/internal/transport/grpc"
	"github.com/godverv/makosh-be/internal/transport/grpc/makosh"
	"github.com/godverv/makosh-be/internal/utils/closer"
)

func main() {
	err := start()
	if err != nil {
		logrus.Fatal(err)
	}
}

func start() error {
	logrus.Println("starting app")

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "error reading config")
	}

	if cfg.GetAppInfo().StartupDuration == 0 {
		return errors.New("no startup duration in config")
	}

	ctx, cancel := context.WithTimeout(ctx, cfg.GetAppInfo().StartupDuration)
	closer.Add(func() error {
		cancel()
		return nil
	})

	grpcServer, err := cfg.GetServers().GRPC(config.ServerGrpc)
	if err != nil {
		return errors.Wrap(err, "error getting grpc server from config")
	}

	inMemoryStore := in_memory.New()

	server, err := grpc.NewServer(grpcServer, makosh.New(cfg, inMemoryStore))
	if err != nil {
		return errors.Wrap(err, "error initializing server")
	}

	err = server.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "error starting server")
	}

	waitingForTheEnd()

	logrus.Println("shutting down the app")

	if err = closer.Close(); err != nil {
		logrus.Fatalf("errors while shutting down application %s", err.Error())
	}

	return nil
}

// rscli comment: an obligatory function for tool to work properly.
// must be called in the main function above
// also this is the LP's song name reference, so no linting rules can be applied to the function name
func waitingForTheEnd() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
}
