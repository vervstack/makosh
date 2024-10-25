package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	"github.com/godverv/makosh/internal/interceptors"
	"github.com/godverv/makosh/pkg/resolver"
	"github.com/godverv/makosh/pkg/resolver/makosh_resolver"
)

func Test_Resolving(t *testing.T) {
	t.Parallel()

	type testCase struct {
		makoshEndPoint string
		makoshSecret   string

		targetName string

		getUpdaterAndResult func() (updateFunc func([]string) error, resultSlice *[]string)
		expectedResult      []string
		// Error for service discovery error responses
		expectedServiceErr map[string]any
	}

	testCases := map[string]testCase{
		"OK": {
			makoshEndPoint: httpMakoshEndpoint,
			makoshSecret:   makoshSecret,
			targetName:     testService1,
			getUpdaterAndResult: func() (updateFunc func([]string) error, resultSlice *[]string) {
				resultSlice = &[]string{}
				return func(addrs []string) error {
					*resultSlice = addrs
					return nil
				}, resultSlice
			},
			expectedResult: examples[0].Addrs,
		},
		"INVALID_AUTH_ERROR": {
			makoshEndPoint: httpMakoshEndpoint,
			makoshSecret:   "fake_secret",
			targetName:     testService1,
			getUpdaterAndResult: func() (updateFunc func([]string) error, resultSlice *[]string) {
				resultSlice = &[]string{}
				return func(addrs []string) error {
					*resultSlice = addrs
					return nil
				}, resultSlice
			},
			expectedResult: []string{},
			expectedServiceErr: map[string]any{
				"code":    float64(codes.PermissionDenied),
				"message": interceptors.InvalidAuthErrMessage,
			},
		},
	}

	for name, tc := range testCases {
		testName, test := name, tc
		t.Run(testName, func(t *testing.T) {

			makoshBuilder, err := makosh_resolver.NewBuilder(
				makosh_resolver.WithURL(test.makoshEndPoint),
				makosh_resolver.WithSecret(test.makoshSecret),
			)
			require.NoError(t, err)

			resolverBuilder, err := resolver.NewLocalServiceDiscovery(
				resolver.WithResolverBuilder(makoshBuilder),
			)

			updaterCallback, result := test.getUpdaterAndResult()
			resolverPtr, err := resolverBuilder.GetResolver(test.targetName)
			require.NoError(t, err)
			resolver := *resolverPtr.Load()
			resolver.AddSubscribers(updaterCallback)

			err = resolver.Resolve()
			if test.expectedServiceErr == nil {
				require.NoError(t, err)
			} else {
				m := make(map[string]any)
				marshErr := json.Unmarshal([]byte(err.Error()), &m)
				require.NoError(t, marshErr)

				require.Equal(t, test.expectedServiceErr, m)
			}

			require.Equal(t, *result, test.expectedResult)
		})
	}
}
