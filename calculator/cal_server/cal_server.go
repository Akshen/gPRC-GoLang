package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/gRPC-GoLang/calculator/calpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculation(ctx context.Context, req *calpb.CalRequest) (*calpb.CalResponse, error) {
	fmt.Printf("Calculate function was invoked with %v \n", req)
	first_num := req.FirstNum
	second_num := req.SecondNum
	result := second_num + first_num
	res := &calpb.CalResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator Started...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calpb.RegisterCalServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
