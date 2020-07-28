package server

import (
	"bytes"
	"context"
	"fmt"
	blogpb "grpcourse/data/protos/blog"
	"io"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blogItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	CoverImage string             `bson:"image"`
	AuthorID   string             `bson:"author_id"`
	Title      string             `bson:"title"`
	Body       string             `bson:"body"`
}

const maxImageSize = 1 << 20

// CreateBlog -  server handler for creating a new blog
func (b *Server) CreateBlog(stream blogpb.BlogService_CreateBlogServer) error {

	b.Logger.Infof("CreateBlog endpoint invoked")

	// 1. receive the blog and image metadata first
	req, err := stream.Recv()
	blog := req.GetBlog()

	// imageData := new(bytes.Buffer)
	imageData := new(bytes.Buffer)

	imgSize := 0

	// 2. accept image chunks
	for {

		// a. stream from client
		ch, err := stream.Recv()

		if err == io.EOF {
			b.Logger.Infof("Done receiving image data : %v", err)
			break
		}

		if err != nil {
			return status.Errorf(codes.Unknown, fmt.Sprintf("cannot receive image data : %v", err))
		}

		chunk := ch.GetImage()

		imgSize += len(chunk)

		if imgSize > maxImageSize {
			imageData.Reset() // clean bytes buffer
			return status.Errorf(codes.InvalidArgument, fmt.Sprintf("file too large!"))
		}

		// b. write bytes to buffer
		_, err = imageData.Write(chunk)

		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("cannot write image data : %v", err))
		}
	}

	// 3. save image to disks - TODO : migrate the logic to a func
	imageName := "new_image.jpg"

	file, err := os.Create("data/images/" + imageName)

	if err != nil {
		b.Logger.Errorf("cannot create image : %v", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("cannot create file : %v", err))
	}

	_, err = imageData.WriteTo(file)

	if err != nil {
		b.Logger.Errorf("cannot save image to disk : %v", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("cannot save image to disk : %v", err))
	}

	// 4. prepare document and save to collection
	data := blogItem{
		CoverImage: file.Name(),
		AuthorID:   blog.GetAuthorId(),
		Title:      blog.GetTitle(),
		Body:       blog.GetBody(),
	}

	// pass context.TODO() or nil to default to context.Background()
	res, err := b.DB.Collection("blog").InsertOne(nil, data)

	if err != nil {
		b.Logger.Errorf("couldn't create a new blog : %v", err)
		return status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))
	}

	b.Logger.Infof("New blog created successfully")

	// Typecast insertedID to ObjectID
	oid := res.InsertedID.(primitive.ObjectID)

	// 5. return response
	response := &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:        oid.Hex(),
			ImagePath: blog.GetImagePath(),
			AuthorId:  blog.GetAuthorId(),
			Title:     blog.GetTitle(),
			Body:      blog.GetBody(),
		},
	}

	return stream.SendAndClose(response)
}

// ReadBlog - server handler for fetching a single blog from collection
func (b *Server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {

	b.Logger.Infof("Readblog endpoint invoked")

	id := req.GetId()

	oid, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		b.Logger.Errorf("could not parse blog id : %v")
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("cannot parse blog id : %v", err))
	}

	data := new(blogItem)

	if err := b.DB.Collection("blog").FindOne(ctx, bson.D{{"_id", oid}}).Decode(data); err != nil {

		if err == mongo.ErrNoDocuments {
			b.Logger.Errorf("document not found : %v", err)
			return nil, status.Errorf(codes.NotFound, fmt.Sprintf("blog not found : %v", err))
		}

		b.Logger.Errorf("could not fetch blog : %v", err)
	}

	b.Logger.Printf("document fetched : %+v\n", data)

	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Title:    data.Title,
			Body:     data.Body,
		},
	}, nil
}

// UpdateBlog -
func (b *Server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {

	b.Logger.Infof("UpdateBlog func invoked")

	// 1. Get blog from request
	blog := req.GetBlog()

	oid, err := primitive.ObjectIDFromHex(blog.GetId())

	if err != nil {
		b.Logger.Errorf("could not parse blog id : %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot parse blog id : %v", err))
	}

	// 1 a. Create image to file - ideally, image metadata could be passed from client
	image := new(bytes.Buffer)
	file, err := os.Create("data/images/updated_new_image.jpg")

	if err != nil {
		b.Logger.Errorf("cannot open image : %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot open image : %v", err))
	}

	image.Write(req.GetImage())  // write bytes to buffer
	_, err = image.WriteTo(file) // write to file

	if err != nil {
		b.Logger.Errorf("cannot save image : %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot save image to disk : %v", err))
	}

	// 2b. alternatively
	datab := &blogItem{
		AuthorID:   blog.GetAuthorId(),
		Title:      blog.GetTitle(),
		Body:       blog.GetBody(),
		CoverImage: file.Name(),
	}

	_, err = b.DB.Collection("blog").ReplaceOne(ctx, bson.D{{"_id", oid}}, datab)

	if err != nil {
		b.Logger.Errorf("cannot update record : %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("cannot update blog : %v", err))
	}

	return &blogpb.UpdateBlogResponse{
		Blog: &blogpb.Blog{
			Id:        oid.Hex(),
			AuthorId:  datab.AuthorID,
			Title:     datab.Title,
			Body:      datab.Body,
			ImagePath: datab.CoverImage,
		},
	}, nil
}
