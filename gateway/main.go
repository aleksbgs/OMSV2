package main

import (
	common "github.com/aleksbgs/commons"
	pb "github.com/aleksbgs/commons/api"
	_ "github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

var (
	httpAddr           = common.EnvString("HTTP_ADDR", ":8080")
	ordersServicesAddr = "localhost:2000"
)

func main() {
	conn, err := grpc.Dial(ordersServicesAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	log.Println("Dialing Orders service at: ", ordersServicesAddr)
	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Println("starting server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server", err.Error())
	}

}
