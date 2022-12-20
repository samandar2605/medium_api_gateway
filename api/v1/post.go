package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samandar2605/medium_api_gateway/api/models"
	pb "github.com/samandar2605/medium_api_gateway/genproto/post_service"
)

// @Router /posts/{id} [get]
// @Summary Get post by id
// @Description Get post by id
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Post
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.grpcClient.PostService().Get(context.Background(), &pb.GetPostRequest{Id: int64(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	post := parsePostModel(resp)
	c.JSON(http.StatusOK, post)
}

// @Router /posts [post]
// @Summary Create a post
// @Description Create a post
// @Tags post
// @Accept json
// @Produce json
// @Param post body models.CreatePostRequest true "post"
// @Success 201 {object} models.Post
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		req models.CreatePostRequest
	)
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fmt.Println("<<<Api ichida 1>>>")
	fmt.Println(req)
	fmt.Println("<<<Api ichida 2 >>>")

	resp, err := h.grpcClient.PostService().Create(context.Background(), &pb.CreatePost{
		Title:       req.Title,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		UserId:      1,
		CategoryId:  req.CategoryID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	post := parsePostModel(resp)
	c.JSON(http.StatusCreated, post)
}

func parsePostModel(post *pb.Post) models.Post {
	return models.Post{
		ID:          post.Id,
		Title:       post.Title,
		Description: post.Description,
		ImageUrl:    post.ImageUrl,
		UserID:      post.UserId,
		CategoryID:  post.CategoryId,
		ViewsCount:  int32(post.ViewsCount),
	}
}

// @Router /posts [get]
// @Summary Get all posts
// @Description Get all posts
// @Tags post
// @Accept json
// @Produce json
// @Param filter query models.GetAllPostsParams false "Filter"
// @Success 200 {object} models.GetAllPostsResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllPost(c *gin.Context) {
	req, err := postsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := h.grpcClient.PostService().GetAll(context.Background(), &pb.GetAllPostsRequest{
		Page:       req.Page,
		Limit:      req.Limit,
		CategoryId: int32(req.CategoryID),
		UserId:     req.UserID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	res, _ := postsResponse(h, result)
	c.JSON(http.StatusOK, res)
}

func postsParams(c *gin.Context) (*models.GetAllPostsParams, error) {
	var (
		limit              int = 10
		page               int = 1
		err                error
		SortByDate         string
		CategoryId, UserId int
	)

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("category_id") != "" {
		CategoryId, err = strconv.Atoi(c.Query("category_id"))
		if err != nil {
			return nil, err
		}
	}

	if c.Query("user_id") != "" {
		UserId, err = strconv.Atoi(c.Query("user_id"))
		if err != nil {
			return nil, err
		}
	}
	if c.Query("sort_by_date") != "" &&
		(c.Query("sort_by_date") == "desc" || c.Query("sort_by_date") == "asc" || c.Query("sort_by_date") == "none") {
		SortByDate = c.Query("sort_by_date")
	}

	return &models.GetAllPostsParams{
		Limit:      int32(limit),
		Page:       int32(page),
		CategoryID: int64(CategoryId),
		UserID:     int64(UserId),
		SortByData: SortByDate,
	}, nil
}

func postsResponse(h *handlerV1, data *pb.GetAllPostsResponse) (*models.GetAllPostsResponse, error) {
	response := models.GetAllPostsResponse{
		Posts: make([]*models.Post, 0),
		Count: int32(data.Count),
	}

	for _, post := range data.Posts {
		if _, err := h.grpcClient.PostService().ViewInc(context.Background(), &pb.GetPostRequest{
			Id: post.Id,
		}); err != nil {
			return nil, err
		}

		post.ViewsCount++

		response.Posts = append(response.Posts, &models.Post{
			ID:          post.Id,
			Title:       post.Title,
			Description: post.Description,
			ImageUrl:    post.Title,
			UserID:      post.UserId,
			CategoryID:  post.CategoryId,
			UpdatedAt:   post.UpdatedAt,
			ViewsCount:  post.ViewsCount,
			CreatedAt:   post.CreatedAt,
		})
	}

	fmt.Println("api gateway ichida ")
	return &response, nil
}

// @Security ApiKeyAuth
// @Summary Update a post
// @Description Update a post
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param user body models.Post true "post"
// @Success 200 {object} models.Post
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [put]
func (h *handlerV1) UpdatePost(ctx *gin.Context) {
	var b models.Post

	err := ctx.ShouldBindJSON(&b)
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
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	post, err := h.grpcClient.PostService().Update(context.Background(), &pb.Post{
		Id:          int64(id),
		Title:       b.Title,
		Description: b.Description,
		ImageUrl:    b.ImageUrl,
		UserId:      payload.UserID,
		CategoryId:  b.CategoryID,
		ViewsCount:  b.ViewsCount,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, post)
}

// @Security ApiKeyAuth
// @Summary Delete a posts
// @Description Delete a posts
// @Tags post
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Failure 500 {object} models.ErrorResponse
// @Router /posts/{id} [delete]
func (h *handlerV1) DeletePost(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}

	_, err = h.grpcClient.PostService().Delete(context.Background(), &pb.GetPostRequest{
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
