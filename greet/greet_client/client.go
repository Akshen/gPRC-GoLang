package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	//doUnary(c)
	//doServerStreaming(c)

	doClientStreaming(c)
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

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC........")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Akshen",
			LastName:  "Doke",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Stream Service GreetManyTimes RPC: %v", err)
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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting todo a client streaming service....")
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Akshen",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bharat",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Abhijit",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Anita",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet %v", err)
	}

	//Iterating over the slice defined above and sending each msg individually.
	for _, req := range requests {
		fmt.Printf("Sending request..%v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error encountered while recieving response %v", err)
	}
	fmt.Printf("LongGreet Response...%v\n", res)
}
