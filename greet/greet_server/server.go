package main

import (
	"com.grpc.tleu/greet/greetpb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	var result int
	cnt := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			var r float64
			// we have finished reading the client stream
			r = float64(result) / float64(cnt)
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: r,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		number := int(req.Greeting.GetNumber())
		result += number
		cnt++
	}
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
