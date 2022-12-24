package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/samandar2605/medium_api_gateway/api/models"

	"github.com/gin-gonic/gin"
	pbp "github.com/samandar2605/medium_api_gateway/genproto/post_service"
)

// @Router /categories [post]
// @Summary Create a category
// @Description Create a category
// @Tags category
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category"
// @Success 201 {object} models.Category
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateCategory(c *gin.Context) {
	var (
		req models.CreateCategoryRequest
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	resp, err := h.grpcClient.CategoryService().Create(context.Background(), &pbp.Category{
		Title: req.Title,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.Category{
		Id:        resp.Id,
		Title:     resp.Title,
		CreatedAt: resp.CreatedAt,
	})
}

// @Router /categories/{id} [get]
// @Summary Get category by id
// @Description Get category by id
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Category
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.grpcClient.CategoryService().Get(context.Background(), &pbp.IdByRequest{
		Id: int64(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Category{
		Id:        resp.Id,
		Title:     resp.Title,
		CreatedAt: resp.CreatedAt,
	})
}

// @Summary Get Category
// @Description Get Category
// @Tags category
// @Accept json
// @Produce json
// @Param filter query models.GetAllCategoriesRequest false "Filter"
// @Success 200 {object} models.GetAllCategoriesResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /categories [get]
func (h *handlerV1) GetCategoryAll(ctx *gin.Context) {
	queryParams, err := validateGetCategoryQuery(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	resp, err := h.grpcClient.CategoryService().GetAll(context.Background(), &pbp.GetCategoryRequest{
		Page:   queryParams.Page,
		Limit:  queryParams.Limit,
		Search: queryParams.Search,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	result := models.GetAllCategoriesResponse{
		Count:      resp.Count,
	}
	for _, i := range resp.Categories {
		result.Categories = append(result.Categories, &models.Category{
			Id:        i.Id,
			Title:     i.Title,
			CreatedAt: i.CreatedAt,
		})
	}
	if result.Categories==nil{
		result.Categories=[]*models.Category{}
	}
	ctx.JSON(http.StatusOK, result)
}

func validateGetCategoryQuery(ctx *gin.Context) (*models.GetAllCategoriesRequest, error) {
	var (
		limit int = 10
		page  int = 1
		err   error
	)
	if ctx.Query("limit") != "" {
		limit, err = strconv.Atoi(ctx.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if ctx.Query("page") != "" {
		page, err = strconv.Atoi(ctx.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	return &models.GetAllCategoriesRequest{
		Limit:  int32(limit),
		Page:   int32(page),
		Search: ctx.Query("search")}, nil
}

// @Summary Update a Category
// @Description Update a Category
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.CreateCategoryRequest true "Category"
// @Success 200 {object} models.Category
// @Failure 500 {object} models.ErrorResponse
// @Router /categories/{id} [put]
func (h *handlerV1) UpdateCategory(ctx *gin.Context) {
	var b models.Category

	err := ctx.ShouldBindJSON(&b)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	category, err := h.grpcClient.CategoryService().Update(context.Background(), &pbp.Category{
		Id:    int64(id),
		Title: b.Title,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to create category",
		})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

// @Summary Delete a categories
// @Description Delete a categories
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Failure 500 {object} models.ErrorResponse
// @Router /categories/{id} [delete]
func (h *handlerV1) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}

	_, err = h.grpcClient.CategoryService().Delete(context.Background(), &pbp.IdByRequest{
		Id: int64(id),
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to Delete method",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "successful delete method",
	})
}
