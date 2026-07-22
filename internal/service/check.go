package service

import (
	"context"
	"strings"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyService "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/valkey-io/valkey-go/om"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"neupaneanish.com.np/authorization/internal/redis"
)

type LoginAccessSession struct {
	Key    string    `json:"key"     valkey:",key"`
	Ver    int64     `json:"ver"     valkey:",ver"`
	ExAt   time.Time `json:"exat"    valkey:",exat"`
	UserID string    `json:"user_id"`
	Role   string    `json:"role"`
}

const (
	LoginAccessSessionPrefix = "login:access:session"
	RoleRoot                 = "root"
)

func (s *AuthorizationService) Check(
	ctx context.Context,
	req *envoyService.CheckRequest,
) (*envoyService.CheckResponse, error) {
	metadata := req.GetAttributes().GetMetadataContext().GetFilterMetadata()
	method := req.GetAttributes().GetRequest().GetHttp().GetPath()

	jwtMetadata, ok := metadata["envoy.filters.http.jwt_authn"]
	if !ok {
		s.cfg.Logger.WarnContext(ctx, "JWT metadata missing from Envoy request")
		return authError(code.Code_UNAUTHENTICATED, "Session expired")
	}

	fields := jwtMetadata.GetFields()

	sub := fields["sub"].GetStringValue()
	jti := fields["jti"].GetStringValue()

	if sub == "" || jti == "" {
		s.cfg.Logger.ErrorContext(ctx, "Envoy pass empty jwt")
		return authError(code.Code_UNAUTHENTICATED, "Session expired")
	}

	session, err := redis.HGet[LoginAccessSession](ctx, LoginAccessSessionPrefix, jti, s.cfg.Client)
	if err != nil {
		if om.IsRecordNotFound(err) {
			s.cfg.Logger.WarnContext(ctx, "Session not found", "userID", sub)
			return authError(code.Code_UNAUTHENTICATED, "Session expired")
		}
		s.cfg.Logger.ErrorContext(ctx, "Valkey get", "error", err)
		return authError(code.Code_INTERNAL, "Internal server error")
	}

	if session.UserID != sub {
		s.cfg.Logger.WarnContext(ctx, "UserID not match", "userID", sub)
		return authError(code.Code_UNAUTHENTICATED, "Session expired")
	}

	if strings.HasPrefix(method, "/root.") && session.Role != RoleRoot {
		return authError(code.Code_PERMISSION_DENIED, "Permission denied")
	}

	return &envoyService.CheckResponse{
		Status: &status.Status{
			Code: int32(code.Code_OK),
		},
		HttpResponse: &envoyService.CheckResponse_OkResponse{
			OkResponse: &envoyService.OkHttpResponse{
				Headers: []*corev3.HeaderValueOption{
					{
						Header: &corev3.HeaderValue{
							Key:   "x-user-id",
							Value: sub,
						},
						AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
					},
					{
						Header: &corev3.HeaderValue{
							Key:   "x-role",
							Value: session.Role,
						},
						AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
					},
					{
						Header: &corev3.HeaderValue{
							Key:   "x-jti",
							Value: jti,
						},
						AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD,
					},
				},
			},
		},
	}, nil
}

func authError(c code.Code, message string) (*envoyService.CheckResponse, error) {
	return &envoyService.CheckResponse{
		Status: &status.Status{
			Code:    int32(c),
			Message: message,
		},
	}, nil
}
