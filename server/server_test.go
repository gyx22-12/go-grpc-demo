package main

import (
	"context"
	"net"
	"testing"

	pb "github.com/gyx22-12/go-grpc-demo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestSayHello(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	go func() { _ = s.Serve(lis) }()
	defer s.Stop()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	got, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "test"})
	if err != nil {
		t.Fatalf("SayHello: %v", err)
	}
	if want := "Hello, test"; got.GetMessage() != want {
		t.Fatalf("unexpected message: got %q, want %q", got.GetMessage(), want)
	}
}
