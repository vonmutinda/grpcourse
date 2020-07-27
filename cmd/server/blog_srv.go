package server

import "go.mongodb.org/mongo-driver/bson/primitive"

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Body     string             `bson:"body"`
}
