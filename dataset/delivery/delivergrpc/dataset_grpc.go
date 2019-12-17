package delivergrpc

import (
	"context"
)

// DataSetServiceHandler the GRPC Handler
type DataSetServiceHandler struct {
	Port string
}

func (s *DataSetServiceHandler) Add(ctx context.Context, request *Request) (*Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a + b

	return &Response{Result: result}, nil
}

// Multiply of calculator
func (s *DataSetServiceHandler) Multiply(ctx context.Context, request *Request) (*Response, error) {
	a, b := request.GetA(), request.GetB()

	result := a * b

	return &Response{Result: result}, nil
}
