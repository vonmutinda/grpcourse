package server

import (
	"context"
	"fmt"
	"grpcourse/data/protos/greet"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server -
type Server struct {
	Logger *logrus.Logger
}

// NewServer - returns Server
func NewServer() *Server {
	return &Server{
		Logger: logrus.New(),
	}
}

// Greet -
func (g *Server) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {

	g.Logger.Infof("Greet func invoked")

	var fName, lName string

	fName = req.Greeting.GetFirstName()
	lName = req.Greeting.GetSecondName()

	// log
	g.Logger.Infof("Greetings to FirstName : %v LastName : %v", fName, lName)

	greeting := &greet.GreetResponse{
		Response: fmt.Sprintf("Hello %s %s", fName, lName),
	}
	return greeting, nil
}

// Sum -
func (g *Server) Sum(ctx context.Context, req *greet.SumRequest) (*greet.SumResponse, error) {

	g.Logger.Infof("Sum func invoked")

	if ctx.Err() == context.Canceled {
		g.Logger.Infof("Request cancelled by client")
		return nil, status.Errorf(codes.DeadlineExceeded, "client cancelled the request : %v", ctx.Err())
	}

	var a, b int64

	a = req.GetA()
	b = req.GetB()

	g.Logger.Infof("Sum of A : %v and B : %v = %v", a, b, a+b)

	time.Sleep(4 * time.Second)

	return &greet.SumResponse{Sum: a + b}, nil
}

// GreetAlot - server streaming
func (g *Server) GreetAlot(req *greet.GreetRequest, stream greet.GreetService_GreetAlotServer) error {

	g.Logger.Infof("GreetAlot func invoked")

	name := req.Greeting.GetFirstName() + " " + req.Greeting.GetSecondName()

	for i := 0; i <= 10; i++ {

		msg := name + " " + strconv.Itoa(i)

		res := &greet.GreetResponse{Response: msg}

		// stream greeting to client
		if err := stream.Send(res); err != nil {
			return err
		}

		// sleep for a second
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// PrimeNumberDecomposition -
func (g *Server) PrimeNumberDecomposition(req *greet.PMRequest, stream greet.GreetService_PrimeNumberDecompositionServer) error {

	g.Logger.Infof("PrimeNumberDecomposition func invoked")

	number := req.GetNumber()

	g.Logger.Infof("Streaming prime numbers from : %v", number)

	N := number
	var n int64 = 2

	for N > 1 {

		if N%n == 0 {

			res := &greet.PMResponse{PrimeFactor: n}

			if err := stream.Send(res); err != nil {
				g.Logger.Fatalf("cannot send data to stream : %v", err)
			}

			N = N / n

		} else {

			n++
		}

	}

	return nil
}

// LongGreet -
func (g *Server) LongGreet(stream greet.GreetService_LongGreetServer) error {

	g.Logger.Infof("LongGreet func invoked")

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			g.Logger.Infof("done streaming : %v", err)

			res := &greet.GreetResponse{
				Response: "finished streaming greeting",
			}

			return stream.SendAndClose(res)
		}

		if err != nil {
			g.Logger.Fatalf("cannot read from stream : %v", err)
			return err
		}

		greeting := "Hallo " + req.Greeting.GetSecondName()

		fmt.Println("greeting : ", greeting)

	}
}

// ComputeAverage -
func (g *Server) ComputeAverage(stream greet.GreetService_ComputeAverageServer) error {

	g.Logger.Infof("ComputeAverage func invoked")

	var (
		sum     int64 = 0
		counter int64 = 0
	)

	for {

		req, err := stream.Recv()

		sum += req.GetNumber()

		if err == io.EOF {
			fmt.Printf("sum : %v counter : %v\n", sum, counter)
			res := &greet.AverageResponse{Average: float64(sum / counter)}
			return stream.SendAndClose(res)
		}

		if err != nil {
			g.Logger.Fatalf("couldn't get values from stream : %v", err)
		}

		counter++
	}

}

// GreetEveryone -
func (g *Server) GreetEveryone(stream greet.GreetService_GreetEveryoneServer) error {

	g.Logger.Infof("GreetEveryone func invoked")

	for {

		req, err := stream.Recv()

		if err == io.EOF {
			g.Logger.Infof("GreetEveryone stream is done for")
			return nil
		}

		if err != nil {
			return fmt.Errorf("could not read greetings from client stream : %v", err)
		}

		res := &greet.GreetResponse{Response: "Hi " + req.Greeting.FirstName}

		if err := stream.Send(res); err != nil {
			g.Logger.Fatalf("cannot send greet everyone res : %v", err)

			return err
		}
	}

}

// SquareRoot -
func (g *Server) SquareRoot(ctx context.Context, req *greet.SquareRootRequest) (*greet.SquareRootResponse, error) {

	g.Logger.Infof("SquareRoot func invoked")

	num := req.GetNumber()

	if num < 0 {

		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("received a negative number : %v", num))
	}

	res := &greet.SquareRootResponse{Square: math.Sqrt(float64(num))}

	return res, nil
}
