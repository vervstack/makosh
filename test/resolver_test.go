package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.vervstack.ru/makosh/internal/interceptors"
	"go.vervstack.ru/makosh/pkg/resolver"
	"go.vervstack.ru/makosh/pkg/resolver/makosh_resolver"
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
		expectedServiceErr error
	}

	testCases := map[string]testCase{
		"OK": {
			makoshEndPoint: makoshEndpoint,
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
			makoshEndPoint: makoshEndpoint,
			makoshSecret:   "fake_secret",
			targetName:     testService1,
			getUpdaterAndResult: func() (updateFunc func([]string) error, resultSlice *[]string) {
				resultSlice = &[]string{}
				return func(addrs []string) error {
					*resultSlice = addrs
					return nil
				}, resultSlice
			},
			expectedResult:     []string{},
			expectedServiceErr: status.Error(codes.PermissionDenied, interceptors.InvalidAuthErrMessage),
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
			rslvr := *resolverPtr.Load()
			rslvr.AddSubscribers(updaterCallback)

			err = rslvr.Resolve()
			if test.expectedServiceErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)

				require.ErrorIs(t, err, test.expectedServiceErr)
			}

			require.Equal(t, *result, test.expectedResult)
		})
	}
}
