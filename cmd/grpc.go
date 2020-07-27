package cmd

import (
	"context"
	"fmt"
	"grpcourse/data/protos/greet"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var grpcServer = &cobra.Command{
	Use:     "Server",
	Aliases: []string{"server", "gs"},
	Short:   "grPC server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(grpcServer)
}

// Server -
type Server struct {
	Logger *logrus.Logger
}

// RunServer - run grpc server
func serve() {

	fmt.Println("gRPC server running...")

	// Server instance
	svr := &Server{logrus.New()}

	tls := false

	opts := []grpc.ServerOption{}

	if tls {

		// cert and pem - pass full path
		cert := "ssl/server.crt"
		pem := "ssl/server.pem"

		cred, err := credentials.NewServerTLSFromFile(cert, pem)

		if err != nil {
			svr.Logger.Errorf("cannot creat tls cert from file : %v", err)
		}

		// gRPC server with opts
		opts = append(opts, grpc.Creds(cred))
	}

	gs := grpc.NewServer(opts...)

	// reflection allows inspection of gRPC API endpoints
	reflection.Register(gs)

	greet.RegisterGreetServiceServer(gs, svr)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("could not get listener : %v", err)
	}

	if err := gs.Serve(lis); err != nil {
		log.Fatalf("gRPC server couldn't serve : %v", err)
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
