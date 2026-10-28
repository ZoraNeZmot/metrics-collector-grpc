package main

import (
	"grpc/inner/service"
	"log"
)

const addr = "192.168.0.1"

func main() {
	err := service.CreateAndStartService("gRPC", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to start process %v", err)
	}
}
