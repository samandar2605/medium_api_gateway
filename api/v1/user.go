package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/samandar2605/medium_api_gateway/api/models"
	pbu "github.com/samandar2605/medium_api_gateway/genproto/user_service"
)

// @Security ApiKeyAuth
// @Router /users [post]
// @Summary Create a user
// @Description Create a user
// @Tags user
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User"
// @Success 201 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		req models.CreateUserRequest
	)

	payload, err := h.GetAuthPayload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := h.grpcClient.UserService().Create(context.Background(), &pbu.CreateUser{
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Email:           req.Email,
		Gender:          req.Gender,
		Password:        req.Password,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
		Type:            req.Type,
		PayloadType:     payload.UserType,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.User{
		ID:              user.Id,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Email:           user.Email,
		Gender:          user.Gender,
		Username:        user.Username,
		Password:        user.Password,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	})
}

// @Router /users/{id} [get]
// @Summary Get user by id
// @Description Get user by id
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	user, err := h.grpcClient.UserService().Get(context.Background(), &pbu.IdRequest{Id: int64(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.User{
		ID:              user.Id,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Password:        user.Password,
		Email:           user.Email,
		Gender:          user.Gender,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	})
}

// @Router /users [get]
// @Summary Get all users
// @Description Get all users
// @Tags user
// @Accept json
// @Produce json
// @Param filter query models.GetAllUserParams false "Filter"
// @Success 200 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	req, err := validateGetAllParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := h.grpcClient.UserService().GetAll(context.Background(), &pbu.GetAllUsersRequest{
		Page:   req.Page,
		Limit:  req.Limit,
		Search: req.Search,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, getUsersResponse(result))
}

func getUsersResponse(data *pbu.GetAllUsersResponse) *models.GetAllUsersResponse {
	response := models.GetAllUsersResponse{
		Users: make([]models.User, 0),
		Count: data.Count,
	}

	for _, user := range data.Users {
		u := parseUserModel(user)
		response.Users = append(response.Users, u)
	}

	return &response
}

func parseUserModel(user *pbu.User) models.User {
	return models.User{
		ID:              user.Id,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Password:        user.Password,
		Email:           user.Email,
		Gender:          user.Gender,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	}
}

// @Security ApiKeyAuth
// @Summary Update a user
// @Description Update a users
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.UpdateUserRequest true "user"
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [put]
func (h *handlerV1) UpdateUser(ctx *gin.Context) {
	var (
		req models.User
	)

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	user, err := h.grpcClient.UserService().Update(context.Background(), &pbu.UpdateUser{
		Id:              int64(id),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		PhoneNumber:     req.PhoneNumber,
		Gender:          req.Gender,
		Username:        req.Username,
		ProfileImageUrl: req.ProfileImageUrl,
		PayloadType:     payload.UserType,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.User{
		ID:              user.Id,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		PhoneNumber:     user.PhoneNumber,
		Password:        user.Password,
		Email:           user.Email,
		Gender:          user.Gender,
		Username:        user.Username,
		ProfileImageUrl: user.ProfileImageUrl,
		Type:            user.Type,
		CreatedAt:       user.CreatedAt,
	})
}

// @Security ApiKeyAuth
// @Summary Delete a User
// @Description Delete a user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Failure 500 {object} models.ErrorResponse
// @Router /users/{id} [delete]
func (h *handlerV1) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}

	payload, err := h.GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	_, err = h.grpcClient.UserService().Delete(context.Background(), &pbu.DeleteUserRequest{
		Id:          int64(id),
		PayloadType: payload.UserType,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful delete method",
	})
}
