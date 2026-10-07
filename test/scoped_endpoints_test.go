package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"go.vervstack.ru/makosh/internal/interceptors"
	"go.vervstack.ru/makosh/pkg/makosh_be"
)

func Test_ScopedEndpoints_UpsertListDelete(t *testing.T) {
	t.Parallel()

	const serviceName = "test_scoped_service"

	md := metadata.New(map[string]string{interceptors.AuthHeader: makoshSecret})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	upsertReq := newScopedUpsertRequest(serviceName)
	_, err := makoshClient.UpsertEndpoints(ctx, upsertReq)
	require.NoError(t, err)

	dockerReq := &makosh_be.ListEndpoints_Request{ServiceName: serviceName, Scope: makosh_be.AddressScope_DOCKER}
	dockerResp, err := makoshClient.ListEndpoints(ctx, dockerReq)
	require.NoError(t, err)
	require.Equal(t, []string{"container:8080"}, dockerResp.GetUrls())
	require.Len(t, dockerResp.GetAddresses(), 1)
	require.Equal(t, makosh_be.AddressScope_DOCKER, dockerResp.GetAddresses()[0].GetScope())

	deleteVcnReq := &makosh_be.DeleteEndpoints_Request{ServiceName: serviceName, Scope: makosh_be.AddressScope_VCN}
	_, err = makoshClient.DeleteEndpoints(ctx, deleteVcnReq)
	require.NoError(t, err)

	allReq := &makosh_be.ListEndpoints_Request{ServiceName: serviceName}
	allResp, err := makoshClient.ListEndpoints(ctx, allReq)
	require.NoError(t, err)
	require.Equal(t, []string{"container:8080"}, allResp.GetUrls())

	deleteAllReq := &makosh_be.DeleteEndpoints_Request{ServiceName: serviceName}
	_, err = makoshClient.DeleteEndpoints(ctx, deleteAllReq)
	require.NoError(t, err)

	_, err = makoshClient.ListEndpoints(ctx, allReq)
	require.Equal(t, codes.NotFound, status.Code(err))
}

func newScopedUpsertRequest(serviceName string) *makosh_be.UpsertEndpoints_Request {
	return &makosh_be.UpsertEndpoints_Request{
		Endpoints: []*makosh_be.Endpoint{
			{
				ServiceName: serviceName,
				Addresses: []*makosh_be.Address{
					{Addr: "container:8080", Scope: makosh_be.AddressScope_DOCKER, Name: "http"},
					{Addr: "100.64.0.1:8080", Scope: makosh_be.AddressScope_VCN, Name: "http"},
				},
			},
		},
	}
}
