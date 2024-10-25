package makosh_resolver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/godverv/makosh/pkg/makosh_be"
)

func Test_Resolver(t *testing.T) {
	t.Run("OK", func(t *testing.T) {

		resp := &makosh_be.ListEndpoints_Response{
			Urls: []string{"https://matreshka.verv.ru", "https://matreshka.verv.online"},
		}
		serverCalledTimes := 0
		srv := httptest.NewServer(
			http.HandlerFunc(
				func(writer http.ResponseWriter, request *http.Request) {
					serverCalledTimes++
					respB, err := json.Marshal(resp)
					require.NoError(t, err)

					_, err = writer.Write(respB)
					require.NoError(t, err)
				}))

		defer srv.Close()

		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		require.NoError(t, err)

		resolver := NewMakoshResolver(req)

		callbackCalledTimes := 0

		resolver.AddSubscribers(func(_ []string) error {
			callbackCalledTimes++
			return nil
		})

		err = resolver.Resolve()
		require.NoError(t, err)

		require.Equal(t, 1, callbackCalledTimes)
		require.Equal(t, 1, serverCalledTimes)
	})

}
