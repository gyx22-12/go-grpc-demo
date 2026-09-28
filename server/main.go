package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	pb "github.com/gyx22-12/go-grpc-demo/pb"
	"google.golang.org/grpc"
)

// server 实现 Greeter 服务。
type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, " + req.GetName() + ", age " + fmt.Sprint(req.GetAge())}, nil
}

func (s *server) SayHelloStream(req *pb.HelloRequest, stream pb.Greeter_SayHelloStreamServer) error {
	for i := 0; i < 3; i++ {
		if err := stream.Send(&pb.HelloReply{Message: "Hello, " + req.GetName()}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) SayHelloToMany(stream pb.Greeter_SayHelloToManyServer) error {
	var names []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// 客户端发完了，返回汇总结果
			return stream.SendAndClose(&pb.HelloReply{
				Message: "Hello, " + strings.Join(names, ", "),
			})
		}
		if err != nil {
			return err
		}
		names = append(names, req.GetName())
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("gRPC server listening on %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
