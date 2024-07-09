package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	errors "github.com/Red-Sock/trace-errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/godverv/makosh/internal/domain"
	"github.com/godverv/makosh/pkg/makosh_be"
	"github.com/godverv/makosh/pkg/makosh_resolver"
)

const (
	makoshEndpoint = "0.0.0.0:53892"
	makoshSecret   = "makosh_secret"
)

const (
	testService1 = "test_service_1"
	testService2 = "test_service_2"

	firstServerResponse  = "123"
	secondServerResponse = "321"
)

var log = logrus.New()

func TestMain(m *testing.M) {
	initEnv()

	code := m.Run()
	for _, execute := range shutDownOnExit {
		execute()
	}
	os.Exit(code)
}

var resolverBuilder *makosh_resolver.Builder

var makoshClient makosh_be.MakoshBeAPIClient

var (
	examples = []domain.Endpoint{
		{
			ServiceName: testService1,
			Addrs:       []string{createFakeBackend(firstServerResponse), createFakeBackend(firstServerResponse)},
		},
		{
			ServiceName: testService2,
			Addrs:       []string{createFakeBackend(secondServerResponse)},
		},
	}
)

func initEnv() {
	ctx := context.Background()
	var err error
	makoshClient, err = prepareMakoshClient(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	upsertReq := &makosh_be.UpsertEndpoints_Request{}
	for _, endpoint := range examples {
		upsertReq.Endpoints = append(upsertReq.Endpoints,
			&makosh_be.Endpoint{
				ServiceName: endpoint.ServiceName,
				Addrs:       endpoint.Addrs,
			})
	}
	_, err = makoshClient.UpsertEndpoints(ctx, upsertReq)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	resolverBuilder, err = makosh_resolver.New(
		makosh_resolver.WithMakoshURL("http://"+makoshEndpoint),
		makosh_resolver.WithMakoshSecret(makoshSecret),
	)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func prepareMakoshClient(ctx context.Context) (makosh_be.MakoshBeAPIClient, error) {
	dial, err := grpc.NewClient(makoshEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrap(err, "error connecting to makosh")
	}

	client := makosh_be.NewMakoshBeAPIClient(dial)

	_, err = client.Version(ctx, &makosh_be.Version_Request{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting makosh version")
	}

	return client, nil
}

var shutDownOnExit []func()

func createFakeBackend(response string) (addr string) {
	bytesResp := []byte(response)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(bytesResp)
		if err != nil {
			log.Fatal("error writing response from fake server", err)
		}
		return
	}))

	shutDownOnExit = append(shutDownOnExit, server.Close)

	return server.URL
}
