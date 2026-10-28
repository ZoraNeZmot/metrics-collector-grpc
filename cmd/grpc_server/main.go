package main

import (
	"grpc/inner/server"
	"grpc/pkg/config"
	"log"
)

func main() {
	log.Printf("Server info: %v", config.Conf)
	s := server.New(config.Conf.Host, config.Conf.Port)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
