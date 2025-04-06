.PHONY: proto
proto:
	protoc --go_out=. --go-grpc_out=. proto/mcp.proto

server:
	go run server/server.go

client:
	go run client/main.go
