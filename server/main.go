package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "tempconv-go/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTempConvServiceServer
}

func (s *server) CelsiusToFahrenheit(ctx context.Context, req *pb.TempRequest) (*pb.TempResponse, error) {
	f := req.Value*9/5 + 32
	return &pb.TempResponse{Value: f}, nil
}

func (s *server) FahrenheitToCelsius(ctx context.Context, req *pb.TempRequest) (*pb.TempResponse, error) {
	c := (req.Value - 32) * 5 / 9
	return &pb.TempResponse{Value: c}, nil
}

func main() {
	// Allow ports to be configured via environment variables for public tunneling
	grpcPort := ":50051"
	httpPort := ":8080"

	// Support common platform env var `PORT` (used by many PaaS providers)
	if p := getenv("PORT"); p != "" {
		httpPort = ":" + p
	}
	if p := getenv("GRPC_PORT"); p != "" {
		grpcPort = ":" + p
	}
	if p := getenv("HTTP_PORT"); p != "" {
		httpPort = ":" + p
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTempConvServiceServer(grpcServer, &server{})

	log.Printf("gRPC server running on %s", grpcPort)

	// Run gRPC server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start HTTP server and pass the gRPC address so the HTTP handler can call it
	// Use localhost for in-container gRPC connection
	go startHTTPServer("127.0.0.1"+grpcPort, httpPort)

	// Keep program alive
	select {}
}

func getenv(key string) string {
	return os.Getenv(key)
}
