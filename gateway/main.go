package main

import (
	common "github.com/bercivarga/commons"
	pb "github.com/bercivarga/commons/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	httpPort     = common.EnvString("HTTP_PORT", "8080")
	orderService = common.EnvString("ORDER_SERVICE", "localhost:8081")
)

func main() {
	conn, err := grpc.NewClient(orderService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer func(conn *grpc.ClientConn) {
		grpcErr := conn.Close()
		if grpcErr != nil {
			log.Println(grpcErr.Error())
			return
		}
	}(conn)

	log.Println("Connected to order service at " + orderService)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	wrappedMux := LoggingMiddleware(mux)

	constructedPort := ":" + httpPort

	log.Println("Server is starting on port " + constructedPort)

	if err := http.ListenAndServe(constructedPort, wrappedMux); err != nil {
		// print the error in the console in a readable format
		log.Println(err.Error())

		return
	}
}
