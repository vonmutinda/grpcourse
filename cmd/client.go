package cmd

import (
	"context"
	"fmt"
	"grpcourse/data/protos/greet"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
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
	tls := false

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

	client := greet.NewGreetServiceClient(conn)

	// Rather than creating a branch for each implementation,
	// all func will be done here.. comment/uncomment appropriately.

	// 1. Unary implementation
	doUnary(client)
	// doUnarySum(client)

	// 2. Server streaming implementation
	// doGreetStream(client)
	// doPMStream(client)

	//3. Client streaming
	// doClientStreaming(client)
	// doComputeAverage(client)

	// 4. BiDi streaming
	// doBiDiStreaming(client)

	// ----------- 5. Advanced concepts
	// (Error Handling)
	// doSquareRoot(client)

}

// doUnary - Unary RPC implementation for Greet
func doUnary(c greet.GreetServiceClient) {

	req1 := &greet.GreetRequest{Greeting: &greet.Greeting{FirstName: "Jon", SecondName: "Snow"}}

	// Deadlines with gRPC, create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	greeting, err := c.Greet(ctx, req1)

	if err != nil {

		resErr, ok := status.FromError(err)

		if ok {

			if resErr.Code() == codes.DeadlineExceeded {
				log.Printf("Timeout : %v", resErr.Message())
			}
		}

		log.Printf("could not greet : %v", err)
	}

	fmt.Println("Greetings :\n", greeting.Response)
}

func doUnarySum(c greet.GreetServiceClient) {

	req2 := &greet.SumRequest{A: 10, B: 10}

	// Deadlines with gRPC, create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	sum, err := c.Sum(ctx, req2)

	if err != nil {

		resErr, ok := status.FromError(err)

		if !ok {

			log.Printf("could not do sum of %v + %v : %v", req2.GetA(), req2.GetB(), err)
		}

		if resErr.Code() == codes.DeadlineExceeded {
			log.Printf("Timeout : %v", resErr.Message())
		}

		return
	}

	fmt.Printf("Sum of %v and %v = %v\n", req2.A, req2.B, sum.Sum)
}

// doGreetStream - stream data from server
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

// doClientStreaming - stream data to server
func doClientStreaming(c greet.GreetServiceClient) {

	lc, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("cannot create long greeting client : %v", err)
	}

	for i := 0; i <= 10; i++ {
		// go func(n int) {

		req := &greet.GreetRequest{
			Greeting: &greet.Greeting{
				FirstName:  "R",
				SecondName: "Kelly " + strconv.Itoa(i),
			},
		}

		if err := lc.Send(req); err != nil {
			log.Fatalf("cannot create greeting request : %v", err)
		}
		// }(i)
	}

	res, err := lc.CloseAndRecv()

	if err != nil {
		log.Fatalf("could not recieve response from LongGreet : %v", err)
	}

	fmt.Println("Response : ", res.GetResponse())

}

func doComputeAverage(c greet.GreetServiceClient) {

	ca, err := c.ComputeAverage(context.Background())

	var vals []int

	if err != nil {
		log.Fatalf("cannot instantiate compute average client : %v", err)
	}

	// send streamig values
	for i := 1; i <= 20; i += 2 {

		vals = append(vals, i)

		req := &greet.NumberRequest{Number: int64(i)}

		if err := ca.Send(req); err != nil {
			log.Fatalf("cannot stream value to compute average : %v", err)
		}
	}

	res, err := ca.CloseAndRecv()

	if err != nil {
		fmt.Printf("cannot receive and close compute average client : %v\n", err)
	}

	fmt.Printf("(%v)Average of %v = %v\n", len(vals), vals, res.GetAverage())
}

func doBiDiStreaming(c greet.GreetServiceClient) {

	bic, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("cannot create bi-di client stream : %v", err)
	}

	requests := []*greet.GreetRequest{
		{Greeting: &greet.Greeting{FirstName: "Jonas", SecondName: "Khanwald"}},
		{Greeting: &greet.Greeting{FirstName: "Mikkel", SecondName: "Khanwald"}},
		{Greeting: &greet.Greeting{FirstName: "Martha", SecondName: "Nielsen"}},
		{Greeting: &greet.Greeting{FirstName: "Hannah", SecondName: "Khanwald"}},
		{Greeting: &greet.Greeting{FirstName: "Ulrich", SecondName: "Nielsen"}},
		{Greeting: &greet.Greeting{FirstName: "Katharina", SecondName: "Nielsen"}},
		{Greeting: &greet.Greeting{FirstName: "Charlotte", SecondName: "Doppler"}},
	}

	done := make(chan struct{})

	// send stream to server
	go func() {

		for _, req := range requests {

			time.Sleep(300 * time.Millisecond)

			if err := bic.Send(req); err != nil {
				log.Fatalf("cannot send request to greet everyone : %v", err)
			}
		}

		defer bic.CloseSend()
	}()

	// receive stream from server
	go func() {

		for {

			res, err := bic.Recv()

			if err == io.EOF {

				log.Printf("done greeting everyone : %v", err)

				break
			}

			if err != nil {

				log.Fatalf("cannot receive greeting from greet everyone : %v", err)

				break
			}

			fmt.Printf("res : %v\n", res.GetResponse())
		}

		close(done)

	}()

	<-done
}

// doSquareRoot -
func doSquareRoot(c greet.GreetServiceClient) {

	req := &greet.SquareRootRequest{Number: -25}

	sqc, err := c.SquareRoot(context.Background(), req)

	if err != nil {

		resErr, ok := status.FromError(err)

		if !ok {
			log.Fatalf("could not do square root of %v : %v", req.Number, err)
			return
		}

		if resErr.Code() == codes.InvalidArgument {
			// do something
			log.Fatalf("invalid argument : %v", resErr.Message())
			return
		}

	}

	fmt.Printf("SquareRoot of %v : %v\n", req.Number, sqc.GetSquare())
}
