package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	err    error
)

func init() {

	client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	err = client.Connect(ctx)
}

// GetDB - return a MongoDB client ref
func GetDB() *mongo.Database {

	// connect to db -
	db := client.Database("grpcourse")

	if err != nil {
		log.Fatalf("cannot connect to mongodb client : %v", err)
	}

	return db
}
