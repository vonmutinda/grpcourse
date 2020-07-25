package cmd

import (
	"context"
	"fmt"
	"grpcourse/data/protos/greet"
	"log"
	"net"

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

	fName = req.GetFirstName()
	lName = req.GetSecondName()

	// log
	g.Logger.Infof("Greetings to FirstName : %v LastName : %v", fName, lName)

	greeting := &greet.GreetResponse{
		Response: fmt.Sprintf("Hello %s %s", fName, lName),
	}
	return greeting, nil
}
