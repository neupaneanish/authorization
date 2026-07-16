//go:build integration

package service_test

import (
	"testing"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoyService "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/protobuf/types/known/structpb"
	"neupaneanish.com.np/authorization/internal/redis"
	"neupaneanish.com.np/authorization/internal/service"
)

func TestCheck(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		userID := uuid.NewString()
		key := uuid.NewString()

		seedSession(t, key, userID)

		claims := map[string]interface{}{
			"sub":  userID,
			"role": "test",
			"jti":  key,
		}

		req := buildCheckRequest("/test", claims, "")

		res, err := authClient.Check(t.Context(), req)
		require.NoError(t, err)
		assert.NotNil(t, res.GetOkResponse())

		ok := res.GetOkResponse()

		headers := map[string]string{}
		for _, h := range ok.Headers {
			headers[h.Header.Key] = h.Header.Value
		}

		assert.Equal(t, userID, headers["x-user-id"])
		assert.Equal(t, "test", headers["x-role"])
		assert.Equal(t, key, headers["x-jti"])
	})

	t.Run("Empty claims", func(t *testing.T) {
		t.Parallel()

		claims := map[string]interface{}{
			"sub":  "",
			"role": "",
			"jti":  "",
		}

		req := buildCheckRequest("/test", claims, "")

		res, err := authClient.Check(t.Context(), req)
		require.NoError(t, err)
		assert.Equal(t, int32(code.Code_UNAUTHENTICATED), res.Status.Code)
	})

	t.Run("Invalid sub", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewString()
		key := uuid.NewString()

		seedSession(t, key, userID)

		claims := map[string]interface{}{
			"sub":  uuid.NewString(),
			"role": "test",
			"jti":  key,
		}
		req := buildCheckRequest("/test", claims, "")
		res, err := authClient.Check(t.Context(), req)
		require.NoError(t, err)
		assert.Equal(t, int32(code.Code_UNAUTHENTICATED), res.Status.Code)
	})

	t.Run("Session Expired", func(t *testing.T) {
		t.Parallel()

		userID := uuid.NewString()
		key := uuid.NewString()

		claims := map[string]interface{}{
			"sub":  userID,
			"role": "test",
			"jti":  key,
		}
		req := buildCheckRequest("/test", claims, "")
		res, err := authClient.Check(t.Context(), req)
		require.NoError(t, err)
		assert.Equal(t, int32(code.Code_UNAUTHENTICATED), res.Status.Code)
	})

	t.Run("Invalid metadata", func(t *testing.T) {
		t.Parallel()

		req := &envoyService.CheckRequest{
			Attributes: &envoyService.AttributeContext{
				Request: &envoyService.AttributeContext_Request{
					Http: &envoyService.AttributeContext_HttpRequest{
						Path: "/test",
						Body: "",
					},
				},
			},
		}
		res, err := authClient.Check(t.Context(), req)
		require.NoError(t, err)
		assert.Equal(t, int32(code.Code_UNAUTHENTICATED), res.Status.Code)
	})
}

func seedSession(t *testing.T, key, userID string) {
	t.Helper()

	data := &service.LoginAccessSession{
		Key:    key,
		ExAt:   time.Now().Add(15 * time.Minute),
		UserID: userID,
	}
	err := redis.HSet[service.LoginAccessSession](t.Context(), service.LoginAccessSessionPrefix, data, cfg.Client)
	require.NoError(t, err)
}

func buildCheckRequest(path string, jwtClaims map[string]interface{}, body string) *envoyService.CheckRequest {
	var fields map[string]*structpb.Value
	if jwtClaims != nil {
		s, _ := structpb.NewStruct(jwtClaims)
		fields = s.GetFields()
	}

	req := &envoyService.CheckRequest{
		Attributes: &envoyService.AttributeContext{
			Request: &envoyService.AttributeContext_Request{
				Http: &envoyService.AttributeContext_HttpRequest{
					Path: path,
					Body: body,
				},
			},
		},
	}

	if jwtClaims != nil {
		req.Attributes.MetadataContext = &corev3.Metadata{
			FilterMetadata: map[string]*structpb.Struct{
				"envoy.filters.http.jwt_authn": {
					Fields: fields,
				},
			},
		}
	}

	return req
}
