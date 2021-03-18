package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gRPC-GoLang/calculator/calpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)

	doErrorUnary(c)

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

func doClientStreaming(c calpb.CalServiceClient) {
	fmt.Println("Starting todo a client streaming service....")

	stream, err := c.CalAverageofNumbers(context.Background())
	if err != nil {
		log.Fatalf("Error while calling CalAverageofNumber function %v", err)
	}

	requests := []int32{3, 4, 5, 6, 8}

	for _, req := range requests {
		fmt.Printf("Sending requests...%v\n", req)
		stream.Send(&calpb.CalAverageofNumbersRequest{
			Number: req,
		})
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error encountered while recieving response...%v", err)
	}
	fmt.Printf("Average Response....%v\n", res)
}

func doBiDiStreaming(c calpb.CalServiceClient) {
	fmt.Println("Starting todo a BiDi streaming....")
	//we create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []int32{13, 4, 5, 16, 8, 88, 9, 34}

	waitc := make(chan struct{})
	//we send a bunch of messages to the client (go rountine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(&calpb.FindMaximumRequest{
				Number: req,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		// function to recieve a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while recieving the msg %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetMaxNumber())
		}
		close(waitc)
	}()

	//block until everything is done
	<-waitc
}

func doErrorUnary(c calpb.CalServiceClient) {
	fmt.Println("starting the doError function")
	//correct call
	dosqrtcall(c, 10)
	//error call
	dosqrtcall(c, -10)
}

func dosqrtcall(c calpb.CalServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calpb.SquareRootRequest{
		Number: n,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//actual error from gRPC (user error)
			fmt.Println("Error message from server", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot %v", err)
			return
		}
	}
	fmt.Printf("Results of Square Root of number %v is %v\n", n, res.GetNumberRoot())
}
