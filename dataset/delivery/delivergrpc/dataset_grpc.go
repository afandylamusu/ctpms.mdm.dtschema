package delivergrpc

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	port string
}

func (s *Server) Add(ctx context.Context, request *Request) (*Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a + b

	return &Response{Result: result}, nil
}

func (s *Server) Multiply(ctx context.Context, request *Request) (*Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a * b

	return &Response{Result: result}, nil
}

func RunServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.ConnectionTimeout(time.Second),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: time.Second * 10,
			Timeout:           time.Second * 20,
		}),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime:             time.Second,
				PermitWithoutStream: true,
			}),
		grpc.MaxConcurrentStreams(5),
	)

	RegisterAddServiceServer(s, &Server{port})
	log.Println("Run GRPC AddServiceServer: " + port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
