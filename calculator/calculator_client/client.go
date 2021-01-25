package main

import (
	"context"
	"f/workspace/play/grpc-go-course/calculator/calculatorpb"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Calculator Client\n")
	cc, err := grpc.Dial("localhost:5051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doSumUnary(c)
	doServerStream(c)
}

func doSumUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do sum Unary RPC...\n")
	req := &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			A: 10,
			B: 20,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	log.Printf("Response from Sum: %v", res.Result)
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Starting to do PrimeNumberDecomposition Server Streaming...\n")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 234234234242,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happpened: %v", err)
		} else {
			fmt.Println(res.GetPrimeFactor())
		}
	}

}
