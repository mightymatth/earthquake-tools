package entity

type Subscription struct {
	ChatID   int64
	SubID    string
	Name     string
	MinMag   float64
	Delay    float64
	Location *Location
	Radius   float64
	Sources  []SourceID
}

type SubscriptionUpdate struct {
	Name     string
	MinMag   float64
	Delay    float64
	Location *Location
	Radius   float64
	Sources  []SourceID
}

type Location struct {
	Lat, Lng float64
}

type EventData struct {
	Magnitude float64
	Delay     float64
	Location  Location
	Source    SourceID
}

type SourceID string

const (
	EMSC   SourceID = "EMSC"
	EMSCWS SourceID = "EMSCWS"
	USGS   SourceID = "USGS"
	IRIS   SourceID = "IRIS"
	USPBR  SourceID = "USPBR"
	GEOFON SourceID = "GEOFON"
)

func ToSources(s []string) []SourceID {
	srcs := make([]SourceID, len(s))
	for i, v := range s {
		srcs[i] = SourceID(v)
	}
	return srcs
}

type SourceIDs []SourceID

func (srcs SourceIDs) ToStrings() []string {
	s := make([]string, len(srcs))
	for i, v := range srcs {
		s[i] = string(v)
	}
	return s
}
