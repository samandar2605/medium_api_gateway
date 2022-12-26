package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/samandar2605/medium_api_gateway/api/models"
	pbp "github.com/samandar2605/medium_api_gateway/genproto/post_service"
)

// @Router /comments/{id} [get]
// @Summary Get comment by id
// @Description Get comment by id
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Success 200 {object} models.Comment
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	resp, err := h.grpcClient.CommentService().Get(context.Background(), &pbp.IdWithRequest{
		Id: int64(id),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.Comment{
		Id:          int(resp.Id),
		PostId:      int(resp.PostId),
		UserId:      int(resp.UserId),
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	})
}

// @Security ApiKeyAuth
// @Router /comments [post]
// @Summary Create a comment
// @Description Create a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.CreateComment true "comment"
// @Success 201 {object} models.Comment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) CreateComment(c *gin.Context) {
	var (
		req models.CreateComment
	)

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
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

	resp, err := h.grpcClient.CommentService().Create(context.Background(), &pbp.CreateCommentRequest{
		PostId:      int64(req.PostId),
		UserId:      int64(payload.UserID),
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.Comment{
		Id:          int(resp.Id),
		PostId:      int(resp.PostId),
		UserId:      int(resp.UserId),
		Description: resp.Description,
		CreatedAt:   resp.CreatedAt,
	})
}

// @Router /comments [get]
// @Summary Get all comments
// @Description Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Param filter query models.GetAllCommentsParams false "Filter"
// @Success 200 {object} models.GetAllCommentsParams
// @Failure 500 {object} models.ErrorResponse
func (h *handlerV1) GetAllComment(c *gin.Context) {
	req, err := commentsParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := h.grpcClient.CommentService().GetAll(context.Background(), &pbp.GetCommentQuery{
		Page:       int64(req.Page),
		Limit:      int64(req.Limit),
		PostId:     int64(req.PostID),
		SortByDate: req.SortByDate,
		UserId:     int64(req.UserID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, commentsResponse(h, result))
}

func commentsParams(c *gin.Context) (*models.GetAllCommentsParams, error) {
	var (
		limit          int = 10
		page           int = 1
		err            error
		sortByDate     string
		PostId, UserId int
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

	if c.Query("sort_by_date") != "" &&
		(c.Query("sort_by_date") == "desc" || c.Query("sort_by_date") == "asc" || c.Query("sort_by_date") == "none") {
		sortByDate = c.Query("sort_by_date")
	}

	if c.Query("post_id") != "" {
		PostId, err = strconv.Atoi(c.Query("post_id"))
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

	return &models.GetAllCommentsParams{
		Limit:      limit,
		Page:       page,
		SortByDate: sortByDate,
		PostID:     PostId,
		UserID:     UserId,
	}, nil
}

func commentsResponse(h *handlerV1, data *pbp.GetAllCommentsResult) *models.GetAllCommentsResponse {
	response := models.GetAllCommentsResponse{
		Comments: make([]*models.Comment, 0),
		Count:    int(data.Count),
	}

	for _, comment := range data.Comments {
		p := parseCommentModel(comment)
		response.Comments = append(response.Comments, &p)
	}

	return &response
}

func parseCommentModel(Comment *pbp.Comment) models.Comment {
	return models.Comment{
		Id:          int(Comment.Id),
		PostId:      int(Comment.PostId),
		UserId:      int(Comment.UserId),
		Description: Comment.Description,
		CreatedAt:   Comment.CreatedAt,
		UpdatedAt:   Comment.UpdatedAt,
	}
}

// @Security ApiKeyAuth
// @Summary Update a comment
// @Description Update a comments
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param comment body models.UpdateComment true "comment"
// @Success 200 {object} models.Comment
// @Failure 500 {object} models.ErrorResponse
// @Router /comments/{id} [put]
func (h *handlerV1) UpdateComment(ctx *gin.Context) {
	var b models.UpdateComment
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
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	comment, err := h.grpcClient.CommentService().Update(context.Background(), &pbp.Comment{
		Id:          int64(id),
		Description: b.Description,
		UserId:      payload.UserID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Comment{
		Id:          int(comment.Id),
		PostId:      int(comment.PostId),
		UserId:      int(payload.UserID),
		Description: comment.Description,
		CreatedAt:   comment.CreatedAt,
		UpdatedAt:   comment.UpdatedAt,
	})
}

// @Security ApiKeyAuth
// @Summary Delete a comment
// @Description Delete a comment
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Failure 500 {object} models.ErrorResponse
// @Router /comments/{id} [delete]
func (h *handlerV1) DeleteComment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to convert",
		})
		return
	}

	_, err = h.grpcClient.CommentService().Delete(context.Background(), &pbp.IdWithRequest{Id: int64(id)})
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
