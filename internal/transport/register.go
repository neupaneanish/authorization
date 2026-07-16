package transport

import (
	envoyService "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"google.golang.org/grpc"
	"neupaneanish.com.np/authorization/internal/config"
	"neupaneanish.com.np/authorization/internal/service"
)

func register(server *grpc.Server, cfg *config.Config) {
	authorizationService := service.NewAuthorizationService(cfg)
	envoyService.RegisterAuthorizationServer(server, authorizationService)
}
