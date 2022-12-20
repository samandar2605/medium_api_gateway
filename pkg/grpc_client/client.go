package grpc_client

import (
	"fmt"

	"github.com/samandar2605/medium_api_gateway/config"
	pbp "github.com/samandar2605/medium_api_gateway/genproto/post_service"
	pbu "github.com/samandar2605/medium_api_gateway/genproto/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClientI interface {
	UserService() pbu.UserServiceClient
	AuthService() pbu.AuthServiceClient
	PostService() pbp.PostServiceClient
	CategoryService() pbp.CategoryServiceClient
	LikeService() pbp.LikeServiceClient
	CommentService() pbp.CommentServiceClient
}

type GrpcClient struct {
	cfg         config.Config
	connections map[string]interface{}
}

func New(cfg config.Config) (GrpcClientI, error) {
	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.UserServiceHost, cfg.UserServiceGrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("user service dial host: %s port:%s err: %v",
			cfg.UserServiceHost, cfg.UserServiceGrpcPort, err)
	}

	connPostService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.PostServiceHost, cfg.PostServiceGrpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("post service dial host: %s port:%s err: %v",
			cfg.PostServiceHost, cfg.PostServiceGrpcPort, err)
	}

	return &GrpcClient{
		cfg: cfg,
		connections: map[string]interface{}{
			"user_service":     pbu.NewUserServiceClient(connUserService),
			"auth_service":     pbu.NewAuthServiceClient(connUserService),
			"post_service":     pbp.NewPostServiceClient(connPostService),
			"category_service": pbp.NewCategoryServiceClient(connPostService),
			"like_service":     pbp.NewLikeServiceClient(connPostService),
			"comment_service":  pbp.NewCommentServiceClient(connPostService),
		},
	}, nil
}

func (g *GrpcClient) UserService() pbu.UserServiceClient {
	return g.connections["user_service"].(pbu.UserServiceClient)
}

func (g *GrpcClient) AuthService() pbu.AuthServiceClient {
	return g.connections["auth_service"].(pbu.AuthServiceClient)
}

func (g *GrpcClient) PostService() pbp.PostServiceClient {
	return g.connections["post_service"].(pbp.PostServiceClient)
}
func (g *GrpcClient) CategoryService() pbp.CategoryServiceClient {
	return g.connections["category_service"].(pbp.CategoryServiceClient)
}

func (g *GrpcClient) LikeService() pbp.LikeServiceClient {
	return g.connections["like_service"].(pbp.LikeServiceClient)
}
func (g *GrpcClient) CommentService() pbp.CommentServiceClient {
	return g.connections["comment_service"].(pbp.CommentServiceClient)
}
