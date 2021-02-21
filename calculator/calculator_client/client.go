package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"play/grpc-go-course/calculator/calculatorpb"

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
	// doServerStream(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)
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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a FindMaximum doBiDistreaming RPC")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{4, 7, 19, 4, 6, 32}
		for _, number := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Rec()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server stream: %v", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Recieved a new maximum of....: %v", maximum)
		}
		close(waitc)
	}()
	<-waitc
}
