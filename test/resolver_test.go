package test

import (
	"testing"

	"github.com/stretchr/testify/require"

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
			require.NoError(t, err)

			require.Equal(t, *result, test.expectedResult)
		})
	}
}
