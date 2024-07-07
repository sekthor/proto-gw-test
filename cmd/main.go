package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	greeterv1 "github.com/sekthor/proto-gw-test/api/greeter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	greeterv1.UnimplementedGreeterServer
}

func (s *server) Greet(ctx context.Context, req *greeterv1.GreetingRequest) (*greeterv1.GreetingResponse, error) {
	greeting := greeterv1.GreetingResponse{
		Greeting: "Hello",
	}

	if req.GetName() != "" {
		greeting.Greeting = fmt.Sprintf("Hello, %s", req.GetName())
	}

	return &greeting, nil
}

func main() {

	ctx := context.Background()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer()
	greeterv1.RegisterGreeterServer(s, &server{})

	go func() {
		log.Println("Serving gRPC on 0.0.0.0:8080")
		log.Fatal(s.Serve(lis))
	}()

	conn, err := grpc.NewClient(
		"0.0.0.0:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gw := runtime.NewServeMux()
	greeterv1.RegisterGreeterHandler(ctx, gw, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	gwServer := &http.Server{
		Addr:    ":8090",
		Handler: gw,
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
