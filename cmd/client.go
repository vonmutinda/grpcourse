package cmd

import (
	"grpcourse/cmd/client"
	blogpb "grpcourse/data/protos/blog"
	"grpcourse/data/protos/greet"
	"log"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var grpcClient = &cobra.Command{
	Use:     "Client",
	Aliases: []string{"client", "gc"},
	Short:   "grPC client",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	rootCmd.AddCommand(grpcClient)
}

func start() {

	// by default
	tls := true

	opts := grpc.WithInsecure()

	if tls {

		cert := "ssl/ca.crt" // Certificate Authority Trust

		creds, err := credentials.NewClientTLSFromFile(cert, "")

		if err != nil {
			log.Fatalf("could not construct TLS credentials : %v", err)
			return
		}

		opts = grpc.WithTransportCredentials(creds)
	}

	// grpc.Dial(target should be explicit rather than the port)
	// eg. "localhost:PORT" rather than ":PORT"
	conn, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		panic("could not establish connection")
	}

	defer conn.Close()

	// 1. Coursework - greet and calculator
	// gclient := greet.NewGreetServiceClient(conn)
	// exercise(gclient) // comment or uncomment

	// 2. blog service
	bclient := blogpb.NewBlogServiceClient(conn)
	client.DoCreateBlog(bclient) // NB: use created id when reading/updating
	client.DoReadBlog(bclient)
	client.DoUpdateBlog(bclient)

}

func exercise(c greet.GreetServiceClient) {

	// Rather than creating a branch for each implementation,
	// all func will be done here.. comment/uncomment appropriately.

	// 1. Unary implementation
	client.DoUnary(c)
	client.DoUnarySum(c)

	// 2. Server streaming implementation
	client.DoGreetStream(c)
	client.DoPMStream(c)

	//3. Client streaming
	client.DoClientStreaming(c)
	client.DoComputeAverage(c)

	// 4. BiDi streaming
	client.DoBiDiStreaming(c)

	// ----------- 5. Advanced concepts
	// (Error Handling)
	client.DoSquareRoot(c)
}
