package entity

type Subscription struct {
	ChatID   int64
	SubID    string
	Name     string
	MinMag   float64
	Delay    float64
	Location *Location
	Radius   float64
}

type SubscriptionUpdate struct {
	Name     string
	MinMag   float64
	Delay    float64
	Location *Location
	Radius   float64
}

type Location struct {
	Lat, Lng float64
}

type EventData struct {
	Magnitude float64
	Delay     float64
	Location  Location
}
