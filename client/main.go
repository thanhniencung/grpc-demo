package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-demo/pb"
	"io"
	"time"
)



func main()  {
	cc, err := grpc.Dial("localhost:50051",
		grpc.WithInsecure())

	if err != nil {
		fmt.Println(err)
	}

	defer cc.Close()

	c := pb.NewDemoServiceClient(cc)

	//unary(c);

	//fmt.Println("=============")

	//serverStreaming(c)

	//clientStreaming(c)

	biStreamming(c)
}

func clientStreaming(c pb.DemoServiceClient) {
	reqs := []*pb.LongGreetRequest {
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "A",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "B",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "C",
			},
		},
	}

	stream, _ := c.LongGreet(context.Background())

	// Send data from client
	for _, req := range reqs {
		fmt.Println("Sending req: %v \n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}


	// handle response from server
	res, _ := stream.CloseAndRecv()
	fmt.Println("Response: %v\n", res.Result)
}

func biStreamming(c pb.DemoServiceClient) {
	stream, _ := c.GreetEveryone(context.Background())

	reqs := []*pb.GreetEveryoneRequest {
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "A",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "B",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "C",
			},
		},
	}

	waitc := make(chan struct{})

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending message : %v", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			fmt.Println("Response: %v", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func unary(c pb.DemoServiceClient) {
	req := &pb.GreetRequest{
		Greeting: &pb.Greeting{
			FirstName: "ryan",
			LastName: "nguyen",
		},
	}
	res, err := c.Greet(context.Background(), req)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v\n", res.Result)
}

func serverStreaming(c pb.DemoServiceClient)  {
	req := &pb.GreetManyTimesRequest{
		Greeting: &pb.Greeting{
			FirstName: "ryan",
			LastName: "Nguyen",
		},
	}

	res, _ := c.GreetManyTimes(context.Background(), req)

	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			break
		}

		if err == nil {
			fmt.Println(msg.GetResult())
		}
	}
}