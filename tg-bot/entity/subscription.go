package entity

type Subscription struct {
	ChatID     int64
	SubID      string
	Name       string
	MinMag     float64
	EqLocation string
	MyLocation string
	Radius     float64
	OffsetSec  int32
}

type SubscriptionUpdate struct {
	Name       string
	MinMag     float64
	EqLocation string
	MyLocation string
	Radius     float64
	OffsetSec  int32
}
