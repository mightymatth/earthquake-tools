package mongo

import (
	"context"
	"fmt"
	"github.com/mightymatth/earthquake-tools/tgbot/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"math"
	"time"
)

const (
	databaseName           = "earthquakesTgBot"
	chatStateCollection    = "chatStates"
	subscriptionCollection = "subscriptions"
)

type Storage struct {
	Client        *mongo.Client
	database      *mongo.Database
	chatStates    *mongo.Collection
	subscriptions *mongo.Collection
}

func NewStorage(uri, namespace string) (*Storage, error) {

	log.Printf("[mongo] Connecting to '%s'", uri)

	s := new(Storage)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("unable to ping mongo db: %s", err)
	}
	log.Printf("[mongo] Successfully connected!")

	if namespace != "" {
		log.Printf("[mongo] Namespace set to: '%s'", namespace)
	}

	s.Client = client
	s.database = s.Client.Database(namespaced(namespace, databaseName))
	s.chatStates = s.database.Collection(namespaced(namespace, chatStateCollection))
	s.subscriptions = s.database.Collection(namespaced(namespace, subscriptionCollection))

	return s, nil
}

func (s *Storage) GetChatState(chatID int64) *entity.ChatState {
	var chatStateDB ChatState

	filter := bson.M{"chat_id": chatID}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.chatStates.FindOne(ctx, filter).Decode(&chatStateDB); err != nil {
		log.Printf("err %v", err)
		chatStateDB = ChatState{
			ChatID: chatID,
		}
	}

	chatState := entity.ChatState{
		ChatID:       chatStateDB.ChatID,
		AwaitInput:   chatStateDB.AwaitInput,
		DisableInput: chatStateDB.DisableInput,
	}

	return &chatState
}

func (s *Storage) SetChatState(
	chatID int64, update *entity.ChatStateUpdate,
) (*entity.ChatState, error) {
	var newChatStateDB ChatState

	updateDB := ChatStateUpdate{
		AwaitInput:   update.AwaitInput,
		DisableInput: update.DisableInput,
	}

	filter := bson.M{"chat_id": chatID}
	upsert := true
	returnDocument := options.After
	findOptions := options.FindOneAndUpdateOptions{
		ReturnDocument: &returnDocument,
		Upsert:         &upsert,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.chatStates.FindOneAndUpdate(
		ctx, filter, bson.M{"$set": updateDB}, &findOptions,
	).Decode(&newChatStateDB); err != nil {
		return nil, err
	}

	newChatState := entity.ChatState{
		ChatID:       newChatStateDB.ChatID,
		AwaitInput:   newChatStateDB.AwaitInput,
		DisableInput: newChatStateDB.DisableInput,
	}

	return &newChatState, nil
}

func (s *Storage) GetSubscription(subHexID string) (*entity.Subscription, error) {
	var subDB Subscription

	subID, err := primitive.ObjectIDFromHex(subHexID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": subID}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := s.subscriptions.FindOne(ctx, filter).Decode(&subDB); err != nil {
		return nil, err
	}

	sub := entity.Subscription{
		ChatID:   subDB.ChatID,
		SubID:    subDB.ID.Hex(),
		Name:     subDB.Name,
		MinMag:   subDB.MinMag,
		Delay:    subDB.Delay,
		Location: subDB.Location.toLocation(),
		Radius:   subDB.Radius,
		Sources:  entity.ToSources(subDB.Sources),
	}

	return &sub, nil
}

func (s *Storage) CreateSubscription(chatID int64, name string) (*entity.Subscription, error) {
	subCreate := Subscription{
		ChatID:  chatID,
		Name:    name,
		MinMag:  1.5,
		Delay:   15,
		Radius:  140,
		Sources: []string{string(entity.EMSC)},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := s.subscriptions.InsertOne(ctx, &subCreate)
	if err != nil {
		return nil, err
	}

	var newSubDB Subscription
	filter := bson.M{"_id": res.InsertedID}
	err = s.subscriptions.FindOne(ctx, filter).Decode(&newSubDB)
	if err != nil {
		return nil, err
	}

	newSub := entity.Subscription{
		SubID:    newSubDB.ID.Hex(),
		ChatID:   newSubDB.ChatID,
		MinMag:   newSubDB.MinMag,
		Delay:    newSubDB.Delay,
		Location: newSubDB.Location.toLocation(),
		Radius:   newSubDB.Radius,
		Sources:  entity.ToSources(newSubDB.Sources),
	}

	return &newSub, nil
}

func (s *Storage) UpdateSubscription(
	subHexID string, subUpdate *entity.SubscriptionUpdate,
) (*entity.Subscription, error) {
	subID, err := primitive.ObjectIDFromHex(subHexID)
	if err != nil {
		return nil, err
	}

	subUpdateDB := SubscriptionUpdate{
		Name:     subUpdate.Name,
		MinMag:   subUpdate.MinMag,
		Delay:    subUpdate.Delay,
		Location: toPoint(subUpdate.Location),
		Radius:   subUpdate.Radius,
		Sources:  entity.SourceIDs(subUpdate.Sources).ToStrings(),
	}

	if subUpdateDB.Location != nil || subUpdateDB.Radius > 0 {
		err = setObserveArea(subID, &subUpdateDB, s)
		if err != nil {
			return nil, fmt.Errorf("cannot get subscription to set observe area: %v", err)
		}
	}

	var newSubDB Subscription

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.M{"_id": subID}
	update := bson.M{"$set": subUpdateDB}

	if err := s.subscriptions.FindOneAndUpdate(
		ctx, filter, update,
	).Decode(&newSubDB); err != nil {
		return nil, err
	}

	newSub := entity.Subscription{
		SubID:    newSubDB.ID.Hex(),
		ChatID:   newSubDB.ChatID,
		MinMag:   newSubDB.MinMag,
		Delay:    newSubDB.Delay,
		Location: newSubDB.Location.toLocation(),
		Radius:   newSubDB.Radius,
		Sources:  entity.ToSources(newSubDB.Sources),
	}

	return &newSub, nil
}

func (s *Storage) DeleteSubscription(subHexID string) error {
	subID, err := primitive.ObjectIDFromHex(subHexID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": subID}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if _, err := s.subscriptions.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetSubscriptions(chatID int64) (subs []entity.Subscription) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	filter := bson.M{"chat_id": chatID}
	cursor, err := s.subscriptions.Find(ctx, filter)
	if err != nil {
		return subs
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var subDB Subscription
		err := cursor.Decode(&subDB)
		if err != nil {
			continue
		}

		sub := entity.Subscription{
			ChatID:   subDB.ChatID,
			SubID:    subDB.ID.Hex(),
			Name:     subDB.Name,
			MinMag:   subDB.MinMag,
			Delay:    subDB.Delay,
			Location: subDB.Location.toLocation(),
			Radius:   subDB.Radius,
			Sources:  entity.ToSources(subDB.Sources),
		}

		subs = append(subs, sub)
	}

	return subs
}

func (s *Storage) GetEventSubscribers(eventData entity.EventData) (chatIDs []int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "min_mag", Value: bson.M{"$lte": eventData.Magnitude}},
		{Key: "delay", Value: bson.M{"$gte": eventData.Delay}},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "radius", Value: 0}},
			bson.D{{Key: "observe_area", Value: bson.M{"$geoIntersects": bson.M{"$geometry": bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: toPoint(&eventData.Location).ToArray()},
			}}}}}},
		},
		{Key: "sources", Value: eventData.Source},
	}

	projection := &options.FindOptions{
		Projection: bson.M{"chat_id": true},
	}

	cursor, err := s.subscriptions.Find(ctx, filter, projection)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var subDB SubscriptionChatID
		err := cursor.Decode(&subDB)
		if err != nil {
			continue
		}

		chatIDs = append(chatIDs, subDB.ChatID)
	}

	return unique(chatIDs), nil
}

func unique(int64Slice []int64) []int64 {
	keys := make(map[int64]bool)
	var list []int64
	for _, entry := range int64Slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (p *Point) toLocation() *entity.Location {
	if p == nil {
		return nil
	}

	return &entity.Location{
		Lat: p.Lat,
		Lng: p.Lng,
	}
}

func toPoint(loc *entity.Location) *Point {
	if loc == nil {
		return nil
	}

	return &Point{
		Lat: loc.Lat,
		Lng: loc.Lng,
	}
}

func setObserveArea(subID primitive.ObjectID, sub *SubscriptionUpdate, s *Storage) error {
	dbSubOld, err := s.GetSubscription(subID.Hex())
	if err != nil {
		return err
	}

	var location *Point
	if sub.Location != nil && sub.Location != toPoint(dbSubOld.Location) {
		location = sub.Location
	} else if dbSubOld.Location != nil {
		location = toPoint(dbSubOld.Location)
	} else {
		return nil
	}

	var radius float64
	if sub.Radius > 0 && sub.Radius != dbSubOld.Radius {
		radius = sub.Radius
	} else if dbSubOld.Radius > 0 {
		radius = dbSubOld.Radius
	} else {
		return nil
	}

	degToRad := math.Pi / 180
	radToDeg := 180 / math.Pi
	earthRadius := 6371
	pointsTotal := 13 // used for circle approximation

	latR := (radius / float64(earthRadius)) * radToDeg
	lngR := latR / math.Cos(location.Lat*degToRad)

	points := make([]PointAsArray, pointsTotal+1, pointsTotal+1)

	for i := 0; i < pointsTotal; i++ {
		theta := math.Pi * float64(i) / (float64(pointsTotal) / 2)
		ey := location.Lat + (latR * math.Sin(theta))
		ex := location.Lng + (lngR * math.Cos(theta))
		points[i] = Point{Lat: ey, Lng: ex}.ToArray()
	}
	points[pointsTotal] = points[0]

	sub.ObserveArea = NewObserveArea(points)

	return nil
}

func namespaced(namespace, text string) string {
	if namespace == "" {
		return text
	}

	return fmt.Sprintf("%s-%s", namespace, text)
}
