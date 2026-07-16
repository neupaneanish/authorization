package service

import (
	envoyService "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"neupaneanish.com.np/authorization/internal/config"
)

type AuthorizationService struct {
	envoyService.UnimplementedAuthorizationServer

	cfg *config.Config
}

func NewAuthorizationService(cfg *config.Config) *AuthorizationService {
	return &AuthorizationService{
		cfg: cfg,
	}
}
