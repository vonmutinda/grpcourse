package cmd

import (
	"context"
	"fmt"
	"grpcourse/data/protos/greet"
	"io"
	"log"

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

	// 2. Server streaming implementation
	doGreetStream(client)
	doPMStream(client)

	//3. 

}

// doUnary -
func doUnary(c greet.GreetServiceClient) {

	req1 := &greet.GreetRequest{Greeting: &greet.Greeting{FirstName: "Jon", SecondName: "Snow"}}
	req2 := &greet.SumRequest{A: 10, B: 10}

	greeting, _ := c.Greet(context.Background(), req1)

	sum, _ := c.Sum(context.Background(), req2)

	fmt.Println("Greetings :\n", greeting.Response)

	fmt.Printf("Sum of %v and %v = %v\n", req2.A, req2.B, sum.Sum)

}

// doGreetStream -
func doGreetStream(c greet.GreetServiceClient) {

	req := &greet.GreetRequest{
		Greeting: &greet.Greeting{FirstName: "Jane", SecondName: "Doe"},
	}

	stream, err := c.GreetAlot(context.Background(), req)

	if err != nil {
		log.Fatalf("cannot stream greetings : %v", err)
	}

	for {
		res, err := stream.Recv()

		if err == io.EOF { 
			break // end of stream
		}

		if err != nil {
			log.Fatalf("cannot receive greating from stream : %v", err)
		}

		fmt.Println("messages : ", res.GetResponse())
	}

	fmt.Println("end of streaming. Bye.")

}

// doPMStream -
func doPMStream(c greet.GreetServiceClient) {

	req := &greet.PMRequest{Number: 120}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("cannot stream prime numbers : %v", err)
	}

	for {

		res, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("could not read from stream : %v", err)
		}

		fmt.Println(res.GetPrimeFactor())
	}

	fmt.Println("end of streaming. Bye.")
}
