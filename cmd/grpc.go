package cmd

import (
	"context"
	"fmt"
	"grpcourse/data/protos/greet"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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

	gs := grpc.NewServer()

	greet.RegisterGreetServiceServer(gs, &Server{logrus.New()})

	lis, err := net.Listen("tcp", ":9001")

	if err != nil {
		log.Fatalf("could not get listener : %v", err)
	}

	if err := gs.Serve(lis); err != nil {
		log.Fatalf("gRPC server couldn't serve : %v", err)
	}
}

// Greet -
func (g *Server) Greet(ctx context.Context, req *greet.GreetRequest) (*greet.GreetResponse, error) {

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

	var a, b int64

	a = req.GetA()
	b = req.GetB()

	g.Logger.Infof("Sum of A : %v and B : %v = %v", a, b, a+b)

	return &greet.SumResponse{Sum: a + b}, nil
}

// GreetAlot - server streaming
func (g *Server) GreetAlot(req *greet.GreetRequest, stream greet.GreetService_GreetAlotServer) error {

	g.Logger.Infof("Now streaming ... ")

	name := req.Greeting.GetFirstName() + " " + req.Greeting.GetSecondName()

	for i := 0; i <= 10; i++ {

		msg := name + " " + strconv.Itoa(i)

		res := &greet.GreetResponse{Response: msg}

		// stream greeting to client
		if err := stream.Send(res); err != nil {
			return err
		}

		// sleep for a second
		time.Sleep(1 * time.Second)
	}

	return nil
}


// PrimeNumberDecomposition -  
func (g *Server)PrimeNumberDecomposition(req *greet.PMRequest, stream greet.GreetService_PrimeNumberDecompositionServer) error {

	number := req.GetNumber()

	g.Logger.Infof("Streaming prime numbers from : %v", number)

	N := number
	var n int64 = 2

	for N > 1 { 

		if N % n == 0 {

			res := &greet.PMResponse{PrimeFactor: n}

			if err := stream.Send(res); err != nil {
				g.Logger.Errorf("cannot send data to stream : %v", err)
			}

			N = N/n

		} else {

			n ++
		}

	}

	return nil
}