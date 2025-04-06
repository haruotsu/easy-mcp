package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/haruotsu/easy-mcp/proto"
	"google.golang.org/grpc/reflection"
)

type server struct{
	pb.UnimplementedMCPServiceServer
}

func (s *server) SendRequest(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("Received request: %v", req)
	return &pb.Response{
		Id: req.Id,
		Result: fmt.Sprintf("Hello, %s", req.Method),
		Error: "",
	}, nil
}

func main(){
	log.Println("Starting server...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterMCPServiceServer(grpcServer, &server{})

	// リフレクションを有効化
    reflection.Register(grpcServer)
	
	log.Printf("Server started on port %d", 50051)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
