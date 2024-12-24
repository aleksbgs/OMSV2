package gateway

import (
	"context"
	pb "github.com/aleksbgs/commons/api"
	"github.com/aleksbgs/commons/discovery"
	"log"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry: registry}
}

func (g *gateway) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)

	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	c := pb.NewOrderServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: p.CustomerID,
		Items:      p.Items,
	})
}
