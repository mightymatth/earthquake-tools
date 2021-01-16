package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscription struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name,omitempty"`
	ChatID     int64              `bson:"chat_id,omitempty"`
	MinMag     float64            `bson:"min_mag,omitempty"`
	EqLocation string             `bson:"eq_location,omitempty"`
	MyLocation string             `bson:"my_location,omitempty"`
	Radius     float64            `bson:"radius,omitempty"`
	OffsetSec  int32              `bson:"offset_sec,omitempty"`
}
