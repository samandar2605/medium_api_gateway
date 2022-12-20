package api

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/samandar2605/medium_api_gateway/api/v1"
	"github.com/samandar2605/medium_api_gateway/config"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "github.com/samandar2605/medium_api_gateway/api/docs" // for swagger

	grpcPkg "github.com/samandar2605/medium_api_gateway/pkg/grpc_client"
)

type RouterOptions struct {
	Cfg        *config.Config
	GrpcClient grpcPkg.GrpcClientI
}

// @title           Swagger for blog api
// @version         1.0
// @description     This is a blog service api.
// @host      localhost:8000
// @BasePath  /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HandlerV1Options{
		Cfg:        opt.Cfg,
		GrpcClient: &opt.GrpcClient,
	})

	apiV1 := router.Group("/v1")

	// Category
	apiV1.GET("/categories", handlerV1.GetCategoryAll)
	apiV1.GET("/categories/:id", handlerV1.GetCategory)
	apiV1.POST("/categories", handlerV1.CreateCategory)
	apiV1.PUT("/categories/:id", handlerV1.UpdateCategory)
	apiV1.DELETE("/categories/:id", handlerV1.DeleteCategory)

	// Like
	apiV1.POST("/likes", handlerV1.CreateOrUpdateLike)
	apiV1.GET("/likes/user-post", handlerV1.GetLike)

	// User
	apiV1.POST("/users", handlerV1.CreateUser)
	apiV1.GET("/users", handlerV1.GetAllUsers)
	apiV1.GET("/users/:id", handlerV1.GetUser)
	apiV1.PUT("/users/:id", handlerV1.UpdateUser)
	apiV1.DELETE("/users/:id", handlerV1.DeleteUser)

	// Comment
	apiV1.GET("/comments", handlerV1.GetAllComment)
	apiV1.GET("/comments/:id", handlerV1.GetComment)
	apiV1.POST("/comments", handlerV1.CreateComment)
	apiV1.PUT("/comments/:id", handlerV1.UpdateComment)
	apiV1.DELETE("/comments/:id", handlerV1.DeleteComment)

	// Post
	apiV1.GET("/posts", handlerV1.GetAllPost)
	apiV1.GET("/posts/:id", handlerV1.GetPost)
	apiV1.POST("/posts", handlerV1.CreatePost)
	apiV1.PUT("/posts/:id", handlerV1.UpdatePost)
	apiV1.DELETE("/posts/:id", handlerV1.DeletePost)

	// Register
	apiV1.POST("/auth/register", handlerV1.Register)
	apiV1.POST("/auth/verify", handlerV1.Verify)
	apiV1.POST("/auth/login", handlerV1.Login)
	apiV1.POST("/auth/forgot-password", handlerV1.ForgotPassword)
	apiV1.POST("/auth/verify-forgot-password", handlerV1.VerifyForgotPassword)
	apiV1.POST("/auth/update-password", handlerV1.UpdatePassword)


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
