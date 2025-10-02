package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/rob0t7/pact-kafka-plugin/proto"
	"google.golang.org/grpc"
)

const (
	defaultPort = "50051"
)

func main() {
	port := os.Getenv("PLUGIN_PORT")
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPactPluginServer(grpcServer, &pactPluginServer{})

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	log.Printf("pact-kafka-plugin gRPC server listening on port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
