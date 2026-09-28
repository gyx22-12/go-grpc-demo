package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/gyx22-12/go-grpc-demo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 一元调用
	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: "world", Age: 25})
	if err != nil {
		log.Fatalf("SayHello: %v", err)
	}
	log.Printf("unary: %s", resp.GetMessage())

	// 服务端流式
	stream, err := c.SayHelloStream(ctx, &pb.HelloRequest{Name: "stream"})
	if err != nil {
		log.Fatalf("SayHelloStream: %v", err)
	}
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("stream recv: %v", err)
		}
		log.Printf("stream: %s", r.GetMessage())
	}

	// 客户端流式
	stream2, err := c.SayHelloToMany(ctx)
	if err != nil {
		log.Fatalf("SayHelloToMany: %v", err)
	}
	for _, name := range []string{"Alice", "Bob", "Carol"} {
		if err := stream2.Send(&pb.HelloRequest{Name: name}); err != nil {
			log.Fatalf("send: %v", err)
		}
	}
	reply, err := stream2.CloseAndRecv()
	if err != nil {
		log.Fatalf("close: %v", err)
	}
	log.Printf("client-stream: %s", reply.GetMessage())
}
