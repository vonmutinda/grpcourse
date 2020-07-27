package cmd

import (
	"context"
	"fmt"
	"grpcourse/cmd/server"
	blogpb "grpcourse/data/protos/blog"
	"grpcourse/data/protos/greet"
	"net"
	"os"
	"os/signal"

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

	// 3. Register services
	// - you could alternatively create a seperate server for the blog
	greet.RegisterGreetServiceServer(gs, svr)

	blogpb.RegisterBlogServiceServer(gs, svr)

	// 4. create a tcp listener
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		svr.Logger.Fatalf("could not get listener : %v", err)
	}

	go func() {

		if err := gs.Serve(lis); err != nil {
			svr.Logger.Fatalf("gRPC server couldn't listen to port: %v", err)
		}
	}()

	// Gracefully shut down the grpc server
	fmt.Println("[ EXIT ] Press CTRL + C ...")

	kill := make(chan os.Signal)
	signal.Notify(kill, os.Interrupt)

	<-kill

	svr.Logger.Println("Server gracefully shutting down....")

	svr.Logger.Println("Closing mongodb connection....")

	svr.DB.Client().Disconnect(context.TODO())

	lis.Close()

	os.Exit(0)
}
