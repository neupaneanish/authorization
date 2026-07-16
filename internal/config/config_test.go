//go:build integration

package config_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"neupaneanish.com.np/authorization/internal/config"
	"neupaneanish.com.np/authorization/tests"
)

var (
	valkeyURL string
)

func TestMain(m *testing.M) {
	vkURL, vkCleanup, vkErr := tests.Valkey()
	if vkErr != nil {
		panic(vkErr)
	}

	valkeyURL = vkURL
	code := m.Run()

	vkCleanup()

	os.Exit(code)
}

func TestNewConfig(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		env := &config.Env{
			ValkeyURL:    valkeyURL,
			ServiceName:  "Test",
			TelemetryURL: "127.0.0.1:4317",
			Environment:  "production",
			Port:         "50051",
		}

		cfg, cfgErr := config.NewConfig(t.Context(), env, logger)
		require.NoError(t, cfgErr)
		assert.NotNil(t, cfg)
		cfg.Close()
	})

	t.Run("Invalid Client", func(t *testing.T) {
		t.Parallel()

		env := &config.Env{
			ValkeyURL: "invalid",
		}

		cfg, cfgErr := config.NewConfig(t.Context(), env, logger)
		require.Error(t, cfgErr)
		assert.Nil(t, cfg)
	})
}
