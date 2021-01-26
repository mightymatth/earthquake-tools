package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscription struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	ChatID     int64              `bson:"chat_id,omitempty"`
	MinMag     float64            `bson:"min_mag,omitempty"`
	Delay      float64            `bson:"delay,omitempty"`
	MyLocation string             `bson:"my_location,omitempty"`
	Radius     float64            `bson:"radius,omitempty"`
}

type SubscriptionUpdate struct {
	Name       string  `bson:"name,omitempty"`
	MinMag     float64 `bson:"min_mag,omitempty"`
	Delay      float64 `bson:"delay,omitempty"`
	MyLocation string  `bson:"my_location,omitempty"`
	Radius     float64 `bson:"radius,omitempty"`
}
