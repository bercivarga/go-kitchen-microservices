package main

import (
	"context"
	common "github.com/bercivarga/commons"
	"google.golang.org/grpc"
	"log"
	"net"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:8081")
)

func main() {
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer l.Close()

	store := NewStore()
	svc := NewService(store)

	NewGrpcHandler(grpcServer)

	log.Println("Server is starting on port " + grpcAddr)

	coErr := svc.CreateOrder(context.Background())

	if coErr != nil {
		log.Println(coErr.Error())
		return
	}

	// start the gRPC server
	if err := grpcServer.Serve(l); err != nil {
		// handle error
		log.Println(err.Error())
		return
	}
}
