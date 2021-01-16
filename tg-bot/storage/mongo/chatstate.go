package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatState struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ChatID        int64             `bson:"chat_id,omitempty"`
	State         string             `bson:"state,omitempty"`
}
