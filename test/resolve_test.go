package test

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OK(t *testing.T) {
	var addrs []string
	updater := func(a []string) error {
		addrs = a
		return nil
	}

	resolver, err := resolverBuilder.BuildHTTPResolver(testService1, updater)
	require.NoError(t, err)

	err = resolver.Resolve()
	require.NoError(t, err)

	require.Equal(t, examples[0].Addrs, addrs)

	for _, addr := range addrs {
		var httpResp *http.Response
		httpResp, err = http.Get(addr)
		require.NoError(t, err)

		var resp []byte
		resp, err = io.ReadAll(httpResp.Body)
		require.NoError(t, err)

		require.Equal(t, httpResp.StatusCode, http.StatusOK)
		require.Equal(t, firstServerResponse, string(resp))
	}
}
