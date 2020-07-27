package server

import (
	"context"
	"fmt"
	blogpb "grpcourse/data/protos/blog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Body     string             `bson:"body"`
}

// CreateBlog -  server handler for creating a new blog
func (b *Server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {

	b.Logger.Infof("CreateBlog endpoint invoked")

	blog := req.GetBlog()

	data := blogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Body:     blog.GetBody(),
	}

	// pass nil to default to context.Background()
	res, err := b.DB.Collection("blog").InsertOne(ctx, data)

	if err != nil {
		b.Logger.Errorf("couldn't create a new blog : %v", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("internal error : %v", err))
	}

	b.Logger.Infof("New blog created successfully")

	// Typecast insertedID to ObjectID
	oid := res.InsertedID.(primitive.ObjectID)

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Body:     blog.GetBody(),
		},
	}, nil
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
