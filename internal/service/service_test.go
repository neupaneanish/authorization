//go:build integration

package service_test

import (
	"context"
	"log/slog"
	"net"
	"os"
	"testing"

	envoyService "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"neupaneanish.com.np/authorization/internal/config"
	"neupaneanish.com.np/authorization/internal/service"
	"neupaneanish.com.np/authorization/internal/telemetry"
	"neupaneanish.com.np/authorization/internal/transport"
	"neupaneanish.com.np/authorization/tests"
)

var (
	cfg        *config.Config
	authClient envoyService.AuthorizationClient
)

type container struct {
	vkURL            string
	vkCleanup        func()
	telemetryURL     string
	telemetryCleanup func()
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	testContainer := setupContainer()

	testEnv := setupEnv(testContainer.vkURL)
	logger, loggerCleanup, loggerErr := telemetry.NewTelemetry(
		ctx,
		testContainer.telemetryURL,
		testEnv.ServiceName,
		testEnv.Environment,
	)

	if loggerErr != nil {
		slog.Error("Failed to start telemetry", "error", loggerErr)
		os.Exit(1)
	}

	testCfg, testCfgErr := config.NewConfig(ctx, testEnv, logger)
	if testCfgErr != nil {
		slog.Error("Failed to setup config", "error", testCfgErr)
		os.Exit(1)
	}

	cfg = testCfg

	client, server, testClientServerErr := testClientServer(testCfg)
	if testClientServerErr != nil {
		slog.Error("Failed to start client / server", "error", testCfgErr)
		os.Exit(1)
	}

	authClient = envoyService.NewAuthorizationClient(client)

	code := m.Run()
	loggerCleanupErr := loggerCleanup(ctx)
	if loggerCleanupErr != nil {
		slog.Error("Failed to cleanup logger", "error", loggerCleanupErr)
		os.Exit(1)
	}

	err := client.Close()
	if err != nil {
		slog.Error("Failed to close client", "error", err)
		os.Exit(1)
	}

	server.GracefulStop()
	testContainer.vkCleanup()
	testContainer.telemetryCleanup()

	os.Exit(code)
}

func setupContainer() *container {
	vkURL, vkCleanup, vkErr := tests.Valkey()
	if vkErr != nil {
		slog.Error("Failed to start valkey container", "error", vkErr)
		os.Exit(1)
	}

	telemetryURL, telemetryCleanup, telemetryErr := tests.OpenTelemetry()
	if telemetryErr != nil {
		slog.Error("Failed to start telemetry container", "error", telemetryErr)
		os.Exit(1)
	}

	return &container{
		vkURL:            vkURL,
		vkCleanup:        vkCleanup,
		telemetryURL:     telemetryURL,
		telemetryCleanup: telemetryCleanup,
	}
}

func setupEnv(vk string) *config.Env {
	return &config.Env{
		ValkeyURL:   vk,
		ServiceName: "Test",
		Environment: "test",
	}
}

func testClientServer(cfg *config.Config) (*grpc.ClientConn, *grpc.Server, error) {
	listen := bufconn.Listen(1024 * 1024)

	opts, optsErr := transport.NewOptions(cfg.Logger)

	if optsErr != nil {
		return nil, nil, optsErr
	}

	server := grpc.NewServer(opts...)

	authzService := service.NewAuthorizationService(cfg)
	envoyService.RegisterAuthorizationServer(server, authzService)

	go func() {
		if err := server.Serve(listen); err != nil {
			slog.Error("Failed to serve server", "error", err)
			os.Exit(1)
		}
	}()

	client, clientErr := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listen.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if clientErr != nil {
		return nil, nil, clientErr
	}

	return client, server, nil
}
