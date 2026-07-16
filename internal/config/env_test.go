//go:build unit

package config_test

import (
	"crypto/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/authorization/internal/config"
)

func TestLoadEnv(t *testing.T) {
	cleanup := func() {
		_ = os.Unsetenv("VALKEY_URL")
		_ = os.Unsetenv("PORT")
		_ = os.Unsetenv("ENVIRONMENT")
		_ = os.Unsetenv("TELEMETRY_URL")
		_ = os.Unsetenv("SERVICE_NAME")
	}

	t.Run("Success", func(t *testing.T) {
		cleanup()

		t.Setenv("VALKEY_URL", "localhost:6379")
		t.Setenv("PORT", "8080")
		t.Setenv("SERVICE_NAME", "Test Service")
		t.Setenv("ENVIRONMENT", "production")
		t.Setenv("TELEMETRY_URL", "127.0.0.1:4317")

		env, envErr := config.LoadEnv()
		require.NoError(t, envErr)
		assert.NotNil(t, env)
		assert.Equal(t, "8080", env.Port)
	})

	t.Run("Default", func(t *testing.T) {
		cleanup()
		t.Setenv("VALKEY_URL", "localhost:6379")
		t.Setenv("TELEMETRY_URL", "127.0.0.1:4317")

		env, envErr := config.LoadEnv()
		require.NoError(t, envErr)
		assert.NotNil(t, env)
		assert.Equal(t, "50051", env.Port)
	})

	t.Run("Invalid Port", func(t *testing.T) {
		cleanup()
		t.Setenv("VALKEY_URL", "localhost:6379")
		t.Setenv("TELEMETRY_URL", "127.0.0.1:4317")
		t.Setenv("PORT", "79")

		env, err := config.LoadEnv()
		require.Error(t, err)
		assert.Nil(t, env)
	})

	t.Run("Invalid environment", func(t *testing.T) {
		t.Setenv("ENVIRONMENT", "staging")
		env, err := config.LoadEnv()
		require.Error(t, err)
		assert.Nil(t, env)
	})

	t.Run("Missing Required Environment", func(t *testing.T) {
		requiredVariables := []string{
			"VALKEY_URL",
			"TELEMETRY_URL",
		}

		for _, v := range requiredVariables {
			t.Run("Missing "+v, func(t *testing.T) {
				for _, all := range requiredVariables {
					t.Setenv(all, rand.Text())
				}

				_ = os.Unsetenv(v)

				env, err := config.LoadEnv()
				require.Error(t, err)
				assert.Nil(t, env)
				assert.Contains(t, err.Error(), v)
			})
		}
	})
}
