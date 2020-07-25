package cmd

import (
	"context"
	"fmt"
	"grpcourse/data/protos/greet"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
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

	conn, err := grpc.Dial(":9001", grpc.WithInsecure())

	if err != nil {
		panic("could not establish connection")
	}

	defer conn.Close()

	client := greet.NewGreetServiceClient(conn)

	// 1. Unary implementation 
	doUnary(client)

}

// doUnary - 
func doUnary(c greet.GreetServiceClient) {

	req := &greet.GreetRequest{FirstName: "von", SecondName: "Doe"}

	greeting, _ := c.Greet(context.Background(), req)

	fmt.Println("Greetings :\n ", greeting.Response)
}
