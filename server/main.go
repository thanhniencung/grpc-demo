package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-demo/pb"
	"io"
	"net"
	"strconv"
	"time"
)

type server struct {

}

func (*server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Println(req)
	firstName := req.GetGreeting().GetFirstName()
	result := "response : " + firstName
	res := &pb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *pb.GreetManyTimesRequest, stream pb.DemoService_GreetManyTimesServer ) (error) {
	firstName := req.GetGreeting().GetFirstName()
	for i:=0; i<10; i++ {
		result := "response " + firstName + strconv.Itoa(i)
		res := &pb.GreetManyTimesResponse{
			Result: result,
		}

		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream pb.DemoService_LongGreetServer) (error) {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			break
		}

		firstName := req.Greeting.GetFirstName()
		result += firstName + "! \n"
	}
	return nil
}

func (*server) GreetEveryone(stream pb.DemoService_GreetEveryoneServer) (error) {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "

		sendErr := stream.Send(&pb.GreetEveryoneResponse{
			Result: result,
		})

		if sendErr != nil {
			return err
		}
	}
}

func main()  {
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()
	pb.RegisterDemoServiceServer(s, &server{})
	fmt.Println("server is running..")
	s.Serve(lis);
}