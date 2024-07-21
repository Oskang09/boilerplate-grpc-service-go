package app

import (
	"service/app/api"
	"service/app/bootstrap"
	"service/app/config"
	"service/package/grpc"

	pb "service/protobuf/go/v1"
)

func Start(port string) {
	bs := bootstrap.New()
	handler := api.New(bs)
	server := grpc.NewServer(config.IsLocal())
	pb.RegisterExampleServiceServer(server, handler)
	if err := server.Start(port); err != nil {
		panic(err)
	}
}
