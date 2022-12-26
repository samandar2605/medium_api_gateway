package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	pbu "github.com/samandar2605/medium_api_gateway/genproto/user_service"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

type Payload struct {
	ID        string `json:"id"`
	UserID    int64  `json:"user_id"`
	Email     string `json:"email"`
	UserType  string `json:"type"`
	IssuedAt  string `json:"issued_at"`
	ExpiredAt string `json:"expired_at"`
}

func (h *handlerV1) AuthMiddleware(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader(authorizationHeaderKey)
		fmt.Println(c.Request.URL.Path)
		if len(accessToken) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := h.grpcClient.AuthService().VerifyToken(context.Background(), &pbu.VerifyTokenRequest{
			AccessToken: accessToken,
			Resource:    resource,
			Action:      action,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if !payload.HasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, errorResponse(ErrNotAllowed))
			return
		}

		c.Set(authorizationPayloadKey, Payload{
			ID:        payload.Id,
			UserID:    payload.UserId,
			Email:     payload.Email,
			UserType:  payload.UserType,
			IssuedAt:  payload.IssuedAt,
			ExpiredAt: payload.ExpiredAt,
		})
		c.Next()
	}
}

func (m *handlerV1) GetAuthPayload(ctx *gin.Context) (*Payload, error) {
	i, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("")
	}

	payload, ok := i.(Payload)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return &payload, nil
}
