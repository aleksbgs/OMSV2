package main

import (
	"context"
	common "github.com/aleksbgs/commons"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:2000")
)

func main() {

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)

	if err != nil {
		log.Fatal("failed to connect to grpc server: ", err.Error())
	}
	defer l.Close()

	store := NewStore()

	svc := NewService(store)

	NewGrpcHandler(grpcServer)

	svc.CreateOrder(context.Background())

	log.Println("starting grpc server at", grpcAddr)

	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err.Error())
	}

}
