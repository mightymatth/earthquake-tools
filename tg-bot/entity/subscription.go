package entity

type Subscription struct {
	ChatID     int64
	SubID      string
	Name       string
	MinMag     float64
	Delay      float64
	MyLocation string
	Radius     float64
}

type SubscriptionUpdate struct {
	Name       string
	MinMag     float64
	Delay      float64
	MyLocation string
	Radius     float64
}
