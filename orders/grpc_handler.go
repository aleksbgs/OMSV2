package main

import (
	"context"
	pb "github.com/aleksbgs/commons/api"
	"google.golang.org/grpc"
	"log"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrdersService
}

func NewGrpcHandler(grpcServer *grpc.Server, service OrdersService) {

	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (g *grpcHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order received! Order %v", payload)

	o := &pb.Order{
		ID: "42",
	}
	return o, nil
}
