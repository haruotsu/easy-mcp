package main

import (
	"context"
	"log"
	"time"

	pb "github.com/haruotsu/easy-mcp/proto"
	"google.golang.org/grpc"
	"fmt"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMCPServiceClient(conn)

	req := &pb.Request{
        Id:     "123",
        Method: "Greet",
        Params: map[string]string{"message": "Hello from client"},
    }

	// サーバーにリクエストを送信
	// ctxはコンテキスト. withtimeoutはタイムアウトを設定
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel() // コンテキストを解放

	res, err := client.SendRequest(ctx, req)
    if err != nil {
        log.Fatalf("Error calling SendRequest: %v", err)
    }
    fmt.Printf("Response from server: ID=%s, Result=%s, Error=%s\n", res.Id, res.Result, res.Error)
}
