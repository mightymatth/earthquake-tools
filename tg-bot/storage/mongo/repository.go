package mongo

import (
	"context"
	"flag"
	"fmt"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

var flagMongoURI = flag.String("mongo-uri",
	"mongodb://localhost:27017",
	"MongoDB URI")

const (
	databaseName         = "earthquakesTgBot"
	chatStatesCollection = "chatStates"
)

type Storage struct {
	Client     *mongo.Client
	DefaultCtx context.Context
	database   *mongo.Database
	chatStates *mongo.Collection
}

func NewStorage() (*Storage, error) {
	var mongoURI string
	if mongoURI = os.Getenv("MONGO_URI"); mongoURI == "" {
		mongoURI = *flagMongoURI
	}

	log.Printf("[mongo] Connecting to '%s'", mongoURI)

	s := new(Storage)

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("unable to ping mongo db: %s", err)
	}
	log.Printf("[mongo] Successfully connected!")

	s.Client = client
	s.DefaultCtx = ctx
	s.database = s.Client.Database(databaseName)
	s.chatStates = s.database.Collection(chatStatesCollection)

	return s, nil
}

func (s *Storage) GetChatState(chatID string) (*entity.ChatState, error) {
	var dbChatState ChatState

	upsert := true
	returnDocument := options.After
	findOptions := options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDocument,
		Upsert:         &upsert,
	}
	freshChatState := ChatState{
		ChatID: chatID,
		State:  "",
	}
	filter := bson.D{{"chat_id", bson.D{{"$eq", chatID}}}}

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	if err := s.chatStates.FindOneAndUpdate(
		ctx, filter, bson.M{"$set": freshChatState}, &findOptions,
	).Decode(&dbChatState); err != nil {
		return nil, err
	}

	chatState := entity.ChatState{
		ChatID: dbChatState.ChatID,
		State:  dbChatState.State,
	}

	return &chatState, nil
}
