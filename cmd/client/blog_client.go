package client

import (
	"context"
	"fmt"
	blogpb "grpcourse/data/protos/blog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DoCreateBlog -
func DoCreateBlog(client blogpb.BlogServiceClient) {

	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "1001",
			Title:    "Introduction to gRPC",
			Body: `In this section we will be looking at \n
			the minutiea of gRPC....
		`,
		},
	}

	res, err := client.CreateBlog(context.Background(), req)

	if err != nil {

		resErr, ok := status.FromError(err)

		if !ok {
			fmt.Printf("cannot save new blog : %v", err)
		}
		if resErr.Code() == codes.Internal {
			fmt.Printf("mongodb err : %v", err)
		}

		return
	}

	fmt.Printf("Blog created :\n %+v\n", res.GetBlog())
}

// DoReadBlog -
func DoReadBlog(client blogpb.BlogServiceClient) {

	req := &blogpb.ReadBlogRequest{
		Id: "5f1f2df4e60566953662d73f",
	}

	res, err := client.ReadBlog(context.Background(), req)

	if err != nil {

		resErr, ok := status.FromError(err)

		if !ok {
			fmt.Printf("could not fetch blog : %v", err)
		}

		if resErr.Code() == codes.NotFound {
			fmt.Printf("blog not found : %v", err)
		}

		return
	}

	fmt.Printf("blog found : %+v\n", res.GetBlog())
}
