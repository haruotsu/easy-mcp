package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net"

    "google.golang.org/grpc"
    pb "github.com/haruotsu/easy-mcp/proto"
)

type server struct {
    pb.UnimplementedMCPServiceServer
}

func (s *server) SendRequest(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    log.Printf("Received: %s with method %s", req.Id, req.Method)
    return &pb.Response{
        Id:     req.Id,
        Result: fmt.Sprintf("Hello, %s!", req.Method),
    }, nil
}

func (s *server) StreamMessages(stream pb.MCPService_StreamMessagesServer) error {
    for {
        // クライアントからのリクエストを受け取る
        req, err := stream.Recv()
        if err == io.EOF {
            log.Println("Stream closed by client")
            return nil
        }
        if err != nil {
            log.Printf("Error receiving stream: %v", err)
            return err
        }

        log.Printf("Stream received: %s - %s", req.Id, req.Method)

        // ストリームレスポンスを送信
        res := &pb.Response{
            Id:     req.Id,
            Result: fmt.Sprintf("Processed: %s", req.Method),
        }
        if err := stream.Send(res); err != nil {
            log.Printf("Error sending stream response: %v", err)
            return err
        }
    }
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterMCPServiceServer(grpcServer, &server{})

    log.Println("Streaming gRPC server listening on :50051")
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
