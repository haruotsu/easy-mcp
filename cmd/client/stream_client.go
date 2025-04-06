package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "time"

    "google.golang.org/grpc"
    pb "github.com/haruotsu/easy-mcp/proto"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewMCPServiceClient(conn)

    // 双方向ストリーミングの呼び出し
    stream, err := client.StreamMessages(context.Background())
    if err != nil {
        log.Fatalf("Failed to establish stream: %v", err)
    }

    // 非同期でレスポンスを受け取るゴルーチン
    go func() {
        for {
            res, err := stream.Recv()
            if err == io.EOF {
                log.Println("Stream closed by server")
                return
            }
            if err != nil {
                log.Fatalf("Error receiving response: %v", err)
            }
            fmt.Printf("Received from server: %s - %s\n", res.Id, res.Result)
        }
    }()

    // クライアントからメッセージを送信
    for i := 0; i < 5; i++ {
        req := &pb.Request{
            Id:     fmt.Sprintf("%d", i),
            Method: fmt.Sprintf("Command %d", i),
        }
        if err := stream.Send(req); err != nil {
            log.Fatalf("Error sending request: %v", err)
        }
        time.Sleep(1 * time.Second) // ダミーの処理間隔
    }

    // ストリームを終了
    stream.CloseSend()
    time.Sleep(2 * time.Second) // サーバーからのレスポンスを待機
}
