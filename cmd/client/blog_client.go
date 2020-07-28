package client

import (
	"bufio"
	"context"
	"fmt"
	blogpb "grpcourse/data/protos/blog"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoCreateBlog -
func DoCreateBlog(client blogpb.BlogServiceClient) {

	// 1. prepare blog
	req := &blogpb.CreateBlogRequest{
		Data: &blogpb.CreateBlogRequest_Blog{
			Blog: &blogpb.Blog{
				AuthorId:  "1001",
				ImagePath: "data/temp/cyber_pirate.jpg",
				Title:     "Introduction to gRPC",
				Body: `In this section we will be looking at
				the minutiea of gRPC....
			`,
			},
		},
	}

	// 2. instantiate createblog stream - time out after 5 sec without response
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	stream, err := client.CreateBlog(ctx)

	if err != nil {
		log.Fatalf("cannot create createblog service : %v", err)
	}

	// 3. send blog -
	if err := stream.Send(req); err != nil {
		log.Fatalf("cannot send data to create blog service : %v", err)
		return
	}

	// 4. upload image in chunks
	file, err := os.Open(req.GetBlog().GetImagePath())
	defer file.Close()

	if err != nil {
		log.Fatalf("cannot open image file : %v", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024) // 1KB buffer

	fmt.Println("uploading image chunk ..... ")

	for {

		n, err := reader.Read(buffer)

		if err == io.EOF {
			fmt.Println("image upload complete")
			break
		}

		if err != nil {
			log.Fatalf("cannot read image chunk : %v", err)
		}

		req := &blogpb.CreateBlogRequest{
			Data: &blogpb.CreateBlogRequest_Image{
				Image: buffer[:n],
			},
		}

		if err := stream.Send(req); err != nil {

			resErr, ok := status.FromError(err)

			if resErr.Code() == codes.Unknown && ok {
				// log something
			}

			log.Fatalf("cannot send stream of image chunk to server : %v", err)
		}

	}

	// 5. stop streaming and receive response from server
	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("cannot receive response from server : %v", err)
		return
	}

	log.Printf("Blog created : %+v\n", res.GetBlog())

}

// DoReadBlog -
func DoReadBlog(client blogpb.BlogServiceClient) {

	req := &blogpb.ReadBlogRequest{
		Id: "5f2011c0f7bc9e1a387c2a1e",
	}

	res, err := client.ReadBlog(context.Background(), req)

	if err != nil {

		resErr, ok := status.FromError(err)

		if resErr.Code() == codes.NotFound && ok {
			fmt.Printf("blog not found : %v\n", err)
			return
		}

		fmt.Printf("could not fetch blog : %v\n", err)

		return
	}

	fmt.Printf("blog found : %+v\n", res.GetBlog())
}

// DoUpdateBlog - Update blog content and image as well,
// rather than streaming the image to server, we'll send raw bytes at a go
func DoUpdateBlog(client blogpb.BlogServiceClient) {

	// 1. Prepare image into bytes slice
	file, err := os.Open("data/temp/mojave.jpg")

	if err != nil {
		fmt.Printf("cannot open image file : %v", err)
	}

	// get file size and create buffer
	fi, _ := file.Stat()
	var buffer = make([]byte, fi.Size())

	_, err = bufio.NewReader(file).Read(buffer)

	if err != nil {
		log.Fatalf("cannot read file to buffer : %v", err)
	}

	// 2. create updated body
	req := &blogpb.UpdateBlogRequest{
		Blog: &blogpb.Blog{
			Id:       "5f2011c0f7bc9e1a387c2a1e",
			AuthorId: "2002",
			Title:    "This is a new title",
			Body:     "And the body is not as long as before",
		},
		Image: buffer,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 3. update blog
	res, err := client.UpdateBlog(ctx, req)

	if err != nil {
		resErr, ok := status.FromError(err)

		if ok && resErr.Code() == codes.Internal {
			fmt.Printf("internal server errror : %v\n", err)
			return
		}

		fmt.Printf("cannot update blog : %v\n", err)
		return
	}

	fmt.Printf("Blog updated : %+v\n", res)

}
