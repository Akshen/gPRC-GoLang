package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gRPC-GoLang/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Hello, I'm the Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created Client: %f", c)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC........")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Akshen",
			LastName:  "Doke",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}
