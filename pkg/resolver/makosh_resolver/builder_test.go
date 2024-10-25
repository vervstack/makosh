package makosh_resolver

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ResolverBuilder_New(t *testing.T) {
	t.Parallel()

	t.Run("OK_NO_TRAILING_SLASH", func(t *testing.T) {
		t.Parallel()
		b, err := NewBuilder(
			WithURL("https://makosh.redsock.online/v1/endpoints"),
			WithPublicServiceDiscovery(),
		)
		require.NoError(t, err)
		require.NotNil(t, b)
	})

	t.Run("OK_WITH_TRAILING_SLASH", func(t *testing.T) {
		t.Parallel()
		b, err := NewBuilder(
			WithURL("https://makosh.redsock.online/v1/endpoints/"),
			WithPublicServiceDiscovery(),
		)
		require.NoError(t, err)
		require.NotNil(t, b)
	})

	t.Run("ERR_NO_URL", func(t *testing.T) {
		t.Parallel()
		b, err := NewBuilder()
		require.ErrorIs(t, err, ErrNoServiceDiscoveryURL)
		require.Nil(t, b)
	})

	t.Run("ERR_NO_SECRET_FOR_PRIVATE_SD", func(t *testing.T) {
		t.Parallel()
		b, err := NewBuilder(
			WithURL("guffy_url"),
		)
		require.ErrorIs(t, err, ErrNoMakoshSecret)
		require.Nil(t, b)
	})

	t.Run("ERR_INVALID_URL", func(t *testing.T) {
		t.Parallel()
		b, err := NewBuilder(
			WithURL("guffy_url"),
			WithPublicServiceDiscovery(),
		)
		require.Contains(t, err.Error(), "invalid service discovery url")
		require.Nil(t, b)
	})
}

func Test_ResolverBuilder_NewResolver(t *testing.T) {
	t.Parallel()

	t.Run("OK_PUBLIC", func(t *testing.T) {
		b, err := NewBuilder(
			WithURL("https://makosh.redsock.online/v1/endpoints/"),
			WithPublicServiceDiscovery(),
		)
		require.NoError(t, err)
		require.NotNil(t, b)

		resolver, err := b.NewResolver("matreshka")
		require.NoError(t, err)
		require.NotNil(t, resolver)
	})

	t.Run("OK_PRIVATE", func(t *testing.T) {
		const secret = "secret"
		b, err := NewBuilder(
			WithURL("https://makosh.redsock.online/v1/endpoints/"),
			WithSecret(secret),
		)
		require.NoError(t, err)
		require.NotNil(t, b)

		resolver, err := b.NewResolver("matreshka")
		require.NoError(t, err)
		require.NotNil(t, resolver)
	})

	t.Run("ERR_INVALID_TARGET", func(t *testing.T) {
		const secret = "secret"
		b, err := NewBuilder(
			WithURL("https://makosh.redsock.online/v1/endpoints/"),
			WithSecret(secret),
		)
		require.NoError(t, err)
		require.NotNil(t, b)

		resolver, err := b.NewResolver("%")
		require.Error(t, err)
		require.Nil(t, resolver)
	})
}
