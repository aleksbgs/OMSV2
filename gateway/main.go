package main

import (
	"context"
	"github.com/akeksbgs/omsv2-gateway/gateway"
	common "github.com/aleksbgs/commons"
	"github.com/aleksbgs/commons/discovery"
	"time"

	"github.com/aleksbgs/commons/discovery/consul"
	_ "github.com/joho/godotenv"

	"log"
	"net/http"
)

var (
	serviceName = "gateway"
	httpAddr    = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceID, serviceName, httpAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("health check failed", err)
			}
			time.Sleep(time.Second * 1)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	mux := http.NewServeMux()

	ordersGateway := gateway.NewGRPCGateway(registry)

	handler := NewHandler(ordersGateway)
	handler.registerRoutes(mux)

	log.Println("starting server on", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("Failed to start http server", err.Error())
	}

}
