package main

import (
	"context"
	"fmt"
	"gitlab.com/tleuzhan13/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()

	result := "Hello " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v \n", req)
	firstName := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{Result: fmt.Sprintf("%d) Hello, %v\n", i, firstName)}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending greet many times responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}
	return nil
}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}
