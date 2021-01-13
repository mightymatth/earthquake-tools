package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatState struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ChatID        string             `bson:"chat_id,omitempty"`
	State         string             `bson:"state,omitempty"`
	MinMag        float64            `bson:"min_mag,omitempty"`
	Locations     []string           `bson:"locations,omitempty"`
	MyLocation    string             `bson:"my_location,omitempty"`
	Radius        float64            `bson:"radius,omitempty"`
	TimeOffsetSec int32              `bson:"time_offset_sec,omitempty"`
}
