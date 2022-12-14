package v1

import (
	"context"
	"errors"
	"net/http"

	pbu "github.com/samandar2605/medium_api_gateway/genproto/user_service"
	"github.com/gin-gonic/gin"
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

func (h *handlerV1) AuthMiddleware(c *gin.Context) {
	accessToken := c.GetHeader(authorizationHeaderKey)

	if len(accessToken) == 0 {
		err := errors.New("authorization header is not provided")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	payload, err := h.grpcClient.AuthService().VerifyToken(context.Background(), &pbu.VerifyTokenRequest{
		AccessToken: accessToken,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
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

func (m *handlerV1) GetAuthPayload(ctx *gin.Context) (*Payload, error) {
	i, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("")
	}

	payload, ok := i.(*Payload)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return payload, nil
}