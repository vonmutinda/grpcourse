package cmd

import (
	"fmt"
	"grpcourse/cmd/server"
	"grpcourse/data/protos/blog"
	"grpcourse/data/protos/greet"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
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

// Blog -
type Blog struct {
	Logger *logrus.Logger
}

// RunServer - run grpc server
func serve() {

	fmt.Println("Starting gRPC server ...")

	// 1. Server instance
	svr := server.NewServer()

	tls := true

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

	// 2. instantiate grpc server
	gs := grpc.NewServer(opts...)

	// reflection allows inspection of gRPC API endpoints
	reflection.Register(gs)

	// 3. Register services - you could create a seperate server for the blog
	greet.RegisterGreetServiceServer(gs, svr)

	blog.RegisterBlogServiceServer(gs, svr)

	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("could not get listener : %v", err)
	}

	go func() {

		if err := gs.Serve(lis); err != nil {
			log.Fatalf("gRPC server couldn't serve : %v", err)
		}

	}()

	// Gracefully shut down the grpc server
	fmt.Println("[ EXIT ] Press CTRL + C ...")

	kill := make(chan os.Signal)
	signal.Notify(kill, os.Interrupt)

	<-kill
}
