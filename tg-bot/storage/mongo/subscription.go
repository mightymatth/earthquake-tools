package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscription struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	ChatID      int64              `bson:"chat_id,omitempty"`
	MinMag      float64            `bson:"min_mag,omitempty"`
	Delay       float64            `bson:"delay,omitempty"`
	MyLocation  Point              `bson:"my_location,omitempty"`
	Radius      float64            `bson:"radius,omitempty"`
	ObserveArea ObserveArea        `bson:"observe_area,omitempty"`
}

type SubscriptionUpdate struct {
	Name        string      `bson:"name,omitempty"`
	MinMag      float64     `bson:"min_mag,omitempty"`
	Delay       float64     `bson:"delay,omitempty"`
	MyLocation  *Point       `bson:"my_location,omitempty"`
	Radius      float64     `bson:"radius,omitempty"`
	ObserveArea *ObserveArea `bson:"observe_area,omitempty"`
}

//Point represents a GeoJSON type
type Point struct {
	Lat float64 `bson:"lat,omitempty"`
	Lng float64 `bson:"lng,omitempty"`
}

func (p Point) ToArray() (arr PointAsArray) {
	arr[0] = p.Lng
	arr[1] = p.Lat

	return arr
}

// ObserveArea is a GeoJSON type.
type ObserveArea struct {
	Type        string `bson:"type"`
	Coordinates []Path `bson:"coordinates"`
}

type Path []PointAsArray

//PointAsArray represents a GeoJSON Point in array format.
// First element is longitude, second is latitude.
type PointAsArray [2]float64

//NewObserveArea creates a GeoJSON type Polygon with a single path (called exterior ring).
func NewObserveArea(path Path) ObserveArea {
	return ObserveArea{
		"Polygon",
		[]Path{path},
	}
}
