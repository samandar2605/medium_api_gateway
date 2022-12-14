package main


import (
	"log"

	_ "github.com/lib/pq"
	"github.com/samandar2605/medium_api_gateway/api"
	"github.com/samandar2605/medium_api_gateway/config"
	grpcPkg "github.com/samandar2605/medium_api_gateway/pkg/grpc_client"
)

func main() {
	cfg := config.Load(".")

	grpcConn, err := grpcPkg.New(cfg)
	if err != nil {
		log.Fatalf("failed to get grpc connections: %v", err)
	}
	apiServer := api.New(&api.RouterOptions{
		Cfg:        &cfg,
		GrpcClientI: grpcConn,
	})
	err = apiServer.Run(cfg.HttpPort)
	if err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}