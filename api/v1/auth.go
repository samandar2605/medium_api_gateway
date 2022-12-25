package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samandar2605/medium_api_gateway/api/models"
	pbu "github.com/samandar2605/medium_api_gateway/genproto/user_service"
	"github.com/samandar2605/medium_api_gateway/pkg/utils"
	"google.golang.org/grpc/status"
)

// @Router /auth/register [post]
// @Summary Register a user
// @Description Register a user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.RegisterRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Register(c *gin.Context) {
	var (
		req models.RegisterRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, _ := h.grpcClient.UserService().GetByEmail(context.Background(), &pbu.GetByEmailRequest{
		Email: req.Email,
	})

	if user != nil {
		c.JSON(http.StatusBadRequest, errorResponse(ErrEmailExists))
		return
	}

	_, err = h.grpcClient.AuthService().Register(context.Background(), &pbu.RegisterRequest{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
		Type:      req.Type,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, models.ResponseOK{
		Message: "success",
	})
}

// @Router /auth/verify [post]
// @Summary Verify user
// @Description Verify user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Verify(c *gin.Context) {
	var (
		req models.VerifyRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.grpcClient.AuthService().Verify(context.Background(), &pbu.VerifyRequest{
		Email: req.Email,
		Code:  req.Code,
	})

	if err != nil {
		s, _ := status.FromError(err)
		if s.Message() == "incorrect_code" {
			fmt.Println(req.Code)
			c.JSON(http.StatusBadRequest, errorResponse(ErrIncorrectCode))
			return
		} else if s.Message() == "code_expired" {
			c.JSON(http.StatusBadRequest, errorResponse(ErrCodeExpired))
			return
		} else {
			c.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		ID:          result.Id,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Password:    result.Password,
		Username:    result.Username,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt,
		AccessToken: result.AccessToken,
	})
}

// @Router /auth/login [post]
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) Login(c *gin.Context) {
	var (
		req models.LoginRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := h.grpcClient.AuthService().Login(context.Background(), &pbu.VerifyRequest{
		Email: req.Email,
		Code:  req.Password,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
			return
		}

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	err = utils.CheckPassword(req.Password, result.Password)
	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(ErrWrongEmailOrPass))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserId:   result.Id,
		UserType: result.Type,
		Username: result.Username,
		Email:    result.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		ID:          result.Id,
		FirstName:   result.FirstName,
		LastName:    result.LastName,
		Email:       result.Email,
		Username:    result.Username,
		Password:    result.Password,
		Type:        result.Type,
		CreatedAt:   result.CreatedAt,
		AccessToken: token,
	})
}

// @Router /auth/forgot-password [post]
// @Summary Forgot password
// @Description Forgot password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.ForgotPasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) ForgotPassword(c *gin.Context) {
	var (
		req models.ForgotPasswordRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	_, err = h.grpcClient.AuthService().ForgotPassword(context.Background(), &pbu.UserEmail{
		Email: req.Email,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Verification code has been sent!",
	})
}

// @Router /auth/verify-forgot-password [post]
// @Summary Verify forgot password
// @Description Verify forgot password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.VerifyRequest true "Data"
// @Success 200 {object} models.AuthResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) VerifyForgotPassword(c *gin.Context) {
	var (
		req models.VerifyRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	res, err := h.grpcClient.AuthService().VerifyForgotPassword(context.Background(), &pbu.VerifyRequest{
		Code:  req.Code,
		Email: req.Email,
	})

	if err != nil {
		c.JSON(http.StatusForbidden, errorResponse(ErrCodeExpired))
		return
	}

	token, _, err := utils.CreateToken(h.cfg, &utils.TokenParams{
		UserId:   res.Id,
		UserType: res.Type,
		Username: res.Username,
		Email:    res.Email,
		Duration: time.Hour * 24,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		ID:          res.Id,
		FirstName:   res.FirstName,
		LastName:    res.LastName,
		Email:       res.Email,
		Username:    res.Username,
		Type:        res.Type,
		CreatedAt:   res.CreatedAt,
		AccessToken: token,
	})
}

// @Security ApiKeyAuth
// @Router /auth/update-password [post]
// @Summary Update password
// @Description Update password
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.UpdatePasswordRequest true "Data"
// @Success 200 {object} models.ResponseOK
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) UpdatePassword(c *gin.Context) {
	var (
		req models.UpdatePasswordRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	_, err = h.grpcClient.AuthService().UpdatePassword(context.Background(), &pbu.NewPassword{
		UserId:      payload.UserID,
		Password:    hashedPassword,
		PayloadType: payload.UserType,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.ResponseOK{
		Message: "Password has been updated!",
	})
}
