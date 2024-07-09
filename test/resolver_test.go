package test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	"github.com/godverv/makosh/internal/interceptors"
	"github.com/godverv/makosh/pkg/makosh_resolver"
)

func Test_Resolver(t *testing.T) {
	t.Parallel()

	type testCase struct {
		makoshEndPoint string
		makoshSecret   string

		targetName string

		getUpdaterAndResult func() (updateFunc func([]string) error, resultSlice *[]string)
		expectedResult      []string
		expectedErr         map[string]any
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
			expectedErr: map[string]any{
				"code":    float64(codes.PermissionDenied),
				"message": interceptors.InvalidAuthErrMessage,
			},
		},
	}

	for name, tc := range testCases {
		testName, test := name, tc
		t.Run(testName, func(t *testing.T) {
			resolverBuilder, err := makosh_resolver.New(
				makosh_resolver.WithMakoshURL(test.makoshEndPoint),
				makosh_resolver.WithMakoshSecret(test.makoshSecret),
			)
			require.NoError(t, err)

			updater, result := test.getUpdaterAndResult()
			resolver, err := resolverBuilder.BuildHTTPResolver(test.targetName, updater)
			require.NoError(t, err)

			err = resolver.Resolve()
			if test.expectedErr != nil {
				m := make(map[string]any)
				marshErr := json.Unmarshal([]byte(err.Error()), &m)
				require.NoError(t, marshErr)

				require.Equal(t, test.expectedErr, m)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, *result, test.expectedResult)
		})
	}
}
