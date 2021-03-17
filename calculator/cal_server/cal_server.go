package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

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

func (*server) PrimeNoDecomposition(req *calpb.PrimeNoDecompositionRequest, stream calpb.CalService_PrimeNoDecompositionServer) error {
	fmt.Println("PrimeNumberDecomposition function was invoked.....")
	number := req.GetNumber()
	var k int32 = 2
	s := ""
	for number > 1 {
		if number%k == 0 { // if k evenly divides into N
			s += strconv.Itoa(int(k)) + " "
			number /= k
		} else {
			k += 1
		}
	}
	res := &calpb.PrimeNoDecompositionResponse{
		Result: s,
	}
	stream.Send(res)
	return nil
}

func (*server) CalAverageofNumbers(stream calpb.CalService_CalAverageofNumbersServer) error {
	fmt.Printf("Calculate Average function was invoked.....\n")
	idx := 0
	var result int32 = 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result /= int32(idx)
			return stream.SendAndClose(&calpb.CalAverageofNumbersResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		number := req.GetNumber()
		idx++
		result += number
	}
}

func (*server) FindMaximum(stream calpb.CalService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request\n")
	var maxi int32 = -1
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading the client stream...%v", err)
			return err
		}
		number := req.GetNumber()
		if number > maxi {
			maxi = number
		}
		sendErr := stream.Send(&calpb.FindMaximumResponse{
			MaxNumber: maxi,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client...%v", err)
			return err
		}
	}
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
