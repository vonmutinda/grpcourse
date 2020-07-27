package db

import (
	"context"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoDB(t *testing.T) {

	var collection = GetDB().Collection("test")

	res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})

	if err != nil {
		t.Fatalf("cannot insert test data : %v", err)
	}

	fmt.Printf("data inserted. id : %v", res.InsertedID)
}
