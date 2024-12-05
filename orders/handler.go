package main

import (
	"context"
	pb "github.com/bercivarga/commons/api"
	"google.golang.org/grpc"
	"log"
)

type GrpcHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewGrpcHandler(grpcServer *grpc.Server) {
	handler := &GrpcHandler{}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *GrpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("CreateOrder %v", p)

	items := p.Items

	o := &pb.Order{
		Id:    p.CustomerId,
		Items: items,
	}
	return o, nil
}
