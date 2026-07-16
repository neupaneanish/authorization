package transport

import (
	"log/slog"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	interceptorTimeout = 30 * time.Second
	maxTimeout         = 5 * time.Minute
)

func NewOptions(logger *slog.Logger) ([]grpc.ServerOption, error) {
	oTelHandler := otelgrpc.NewServerHandler(
		otelgrpc.WithFilter(filters.Not(filters.HealthCheck())),
	)

	recoveryOpt := recovery.WithRecoveryHandler(func(p any) error {
		logger.Error("panic recovered in gRPC handler", "panic", p)
		return status.Error(codes.Internal, "Internal Server Error")
	})

	opts := []grpc.ServerOption{
		grpc.StatsHandler(oTelHandler),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recoveryOpt),
			UnaryTimeoutInterceptor(interceptorTimeout),
			logging.UnaryServerInterceptor(
				LoggerInterceptor(logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(recoveryOpt),
			StreamTimeoutInterceptor(maxTimeout),
			logging.StreamServerInterceptor(
				LoggerInterceptor(logger),
				logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
			),
		),
	}

	return opts, nil
}
