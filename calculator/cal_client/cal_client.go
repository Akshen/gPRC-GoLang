package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/gRPC-GoLang/calculator/calpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Hello, I'm the Client....\n")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := calpb.NewCalServiceClient(cc)
	//fmt.Printf("Created Client: %f", c)
	//doUnary(c)

	doServerStreaming(c)

}

func doUnary(c calpb.CalServiceClient) {
	fmt.Println("Starting to do a Unary RPC........")
	req := &calpb.CalRequest{
		FirstNum:  22,
		SecondNum: 8,
	}
	res, err := c.Calculation(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling  RPC %v", err)
	}
	log.Printf("Response from cal: %v", res.Result)
}

func doServerStreaming(c calpb.CalServiceClient) {
	fmt.Println("Client service started.......")
	req := &calpb.PrimeNoDecompositionRequest{
		Number: 120,
	}
	resStream, err := c.PrimeNoDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Error")
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v:", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}
