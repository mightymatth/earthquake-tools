package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatState struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ChatID       int64              `bson:"chat_id,omitempty"`
	AwaitInput   string             `bson:"await_input"`
	DisableInput bool               `bson:"disable_input"`
}

type ChatStateUpdate struct {
	AwaitInput   string `bson:"await_input"`
	DisableInput bool   `bson:"disable_input"`
}
