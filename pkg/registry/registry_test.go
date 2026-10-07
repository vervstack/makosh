package registry_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.vervstack.ru/makosh/pkg/makosh_be"
	"go.vervstack.ru/makosh/pkg/registry"
)

const svc = "svc"

func docker(addr string) *makosh_be.Address {
	return &makosh_be.Address{Addr: addr, Scope: makosh_be.AddressScope_DOCKER}
}

func vcn(addr string) *makosh_be.Address {
	return &makosh_be.Address{Addr: addr, Scope: makosh_be.AddressScope_VCN}
}

func upsert(t *testing.T, r *registry.Registry, addrs []string, addresses ...*makosh_be.Address) {
	t.Helper()

	endpoint := &makosh_be.Endpoint{ServiceName: svc, Addrs: addrs, Addresses: addresses}
	err := r.Upsert(context.Background(), []*makosh_be.Endpoint{endpoint})
	require.NoError(t, err)
}

func Test_Registry_LegacyOnlyUpsertThenList(t *testing.T) {
	r := registry.New()
	upsert(t, r, []string{"a:1", "b:2"})

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Equal(t, []string{"a:1", "b:2"}, resp.GetUrls())
	require.Empty(t, resp.GetAddresses())
}

func Test_Registry_LegacyAddrsReplacedEvenWhenEmpty(t *testing.T) {
	r := registry.New()
	upsert(t, r, []string{"a:1"})
	upsert(t, r, nil)

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Empty(t, resp.GetUrls())
}

func Test_Registry_ScopedUpsertReplacesOnlySameScope(t *testing.T) {
	r := registry.New()
	upsert(t, r, nil, docker("d1"), vcn("v1"))
	upsert(t, r, nil, docker("d2"))

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Equal(t, []string{"v1", "d2"}, resp.GetUrls())
	require.Len(t, resp.GetAddresses(), 2)
}

func Test_Registry_ListFiltersByScope(t *testing.T) {
	cases := []struct {
		name  string
		scope makosh_be.AddressScope
		want  []string
	}{
		{"all scopes", makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED, []string{"legacy", "d1", "v1"}},
		{"docker", makosh_be.AddressScope_DOCKER, []string{"d1"}},
		{"vcn", makosh_be.AddressScope_VCN, []string{"v1"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := registry.New()
			upsert(t, r, []string{"legacy", "d1"}, docker("d1"), vcn("v1"))

			resp, err := r.List(context.Background(), svc, tc.scope)
			require.NoError(t, err)
			require.Equal(t, tc.want, resp.GetUrls())
		})
	}
}

func Test_Registry_DeleteOneScopeKeepsOthers(t *testing.T) {
	r := registry.New()
	upsert(t, r, []string{"legacy"}, docker("d1"), vcn("v1"))

	err := r.Delete(context.Background(), svc, makosh_be.AddressScope_VCN)
	require.NoError(t, err)

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Equal(t, []string{"legacy", "d1"}, resp.GetUrls())
}

func Test_Registry_DeleteLastScopedAddressKeepsRecord(t *testing.T) {
	r := registry.New()
	upsert(t, r, nil, vcn("v1"))

	err := r.Delete(context.Background(), svc, makosh_be.AddressScope_VCN)
	require.NoError(t, err)

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Empty(t, resp.GetUrls())
}

func Test_Registry_DeleteUnspecifiedRemovesService(t *testing.T) {
	r := registry.New()
	upsert(t, r, []string{"legacy"}, docker("d1"))

	err := r.Delete(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)

	_, err = r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.ErrorIs(t, err, registry.ErrNotFound)
}

func Test_Registry_DeleteMissingServiceIsNotAnError(t *testing.T) {
	r := registry.New()

	for _, scope := range []makosh_be.AddressScope{
		makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED,
		makosh_be.AddressScope_DOCKER,
	} {
		err := r.Delete(context.Background(), "missing", scope)
		require.NoError(t, err)
	}
}

func Test_Registry_ListMissingReturnsErrNotFound(t *testing.T) {
	r := registry.New()

	_, err := r.List(context.Background(), "missing", makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.ErrorIs(t, err, registry.ErrNotFound)
}

func Test_Registry_ClientListMissingIsNotFoundStatus(t *testing.T) {
	client := registry.NewClient(registry.New(), "v")

	_, err := client.ListEndpoints(context.Background(), &makosh_be.ListEndpoints_Request{ServiceName: "missing"})
	require.Equal(t, codes.NotFound, status.Code(err))
}

func Test_Registry_ClientRoundTrip(t *testing.T) {
	ctx := context.Background()
	client := registry.NewClient(registry.New(), "v1.2.3")

	version, err := client.Version(ctx, &makosh_be.Version_Request{})
	require.NoError(t, err)
	require.Equal(t, "v1.2.3", version.GetVersion())

	upsertResp, err := client.UpsertEndpoints(ctx, &makosh_be.UpsertEndpoints_Request{
		Endpoints: []*makosh_be.Endpoint{{ServiceName: svc, Addresses: []*makosh_be.Address{docker("d1")}}},
	})
	require.NoError(t, err)
	require.NotNil(t, upsertResp)

	listResp, err := client.ListEndpoints(ctx, &makosh_be.ListEndpoints_Request{ServiceName: svc})
	require.NoError(t, err)
	require.Equal(t, []string{"d1"}, listResp.GetUrls())

	deleteResp, err := client.DeleteEndpoints(ctx, &makosh_be.DeleteEndpoints_Request{ServiceName: svc})
	require.NoError(t, err)
	require.NotNil(t, deleteResp)
}

func Test_Registry_MutatingReturnedResponseDoesNotChangeStore(t *testing.T) {
	r := registry.New()
	upsert(t, r, []string{"legacy"}, docker("d1"))

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	resp.Addresses[0].Addr = "mutated"
	resp.Urls[0] = "mutated"

	again, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Equal(t, "d1", again.GetAddresses()[0].GetAddr())
	require.Equal(t, []string{"legacy", "d1"}, again.GetUrls())
}

func Test_Registry_MutatingUpsertedInputDoesNotChangeStore(t *testing.T) {
	r := registry.New()
	address := docker("d1")
	upsert(t, r, nil, address)
	address.Addr = "mutated"

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Equal(t, []string{"d1"}, resp.GetUrls())
}

func Test_Registry_ScopedOnlyUpsertKeepsLegacyAddrs(t *testing.T) {
	r := registry.New()
	upsert(t, r, []string{"a:1"})
	upsert(t, r, nil, docker("d1"))

	resp, err := r.List(context.Background(), svc, makosh_be.AddressScope_ADDRESS_SCOPE_UNSPECIFIED)
	require.NoError(t, err)
	require.Equal(t, []string{"a:1", "d1"}, resp.GetUrls())
}
