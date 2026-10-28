package main

import (
	"grpc/inner/service"
	"log"
)

func main() {
	err := service.CreateAndStartService("GRPs", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to start process %v", err)
	}
}
