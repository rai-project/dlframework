package main

import (
	"log"
	"net"

	agent "github.com/rai-project/dlframework/mxnet/agent"
)

const (
	port = ":50051"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := agent.Register()

	log.Println("Serving")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
