package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const UsgsID ID = "USGS"

type Usgs struct {
	source
}

func NewUsgs() Usgs {
	return Usgs{source{
		Name: "USGS", Url: "https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson",
		Method: REST, SourceID: UsgsID,
	}}
}

func (s Usgs) Transform(r io.Reader) ([]EarthquakeData, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read from buffer: %v", err)
	}

	var eventsRes UsgsGeoJsonResponse
	err = json.Unmarshal(buf.Bytes(), &eventsRes)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal event data: %v", err)
	}

	features := eventsRes.Features[:10]
	events := make([]EarthquakeData, 0, len(features))
	for _, feature := range features {
		data := EarthquakeData{
			Mag:        feature.Properties.Mag,
			MagType:    feature.Properties.Magtype,
			Depth:      feature.Geometry.Coordinates[2],
			Time:       time.Unix(feature.Properties.Time/1000, 0),
			Lat:        feature.Geometry.Coordinates[1],
			Lon:        feature.Geometry.Coordinates[0],
			Location:   feature.Properties.Place,
			DetailsURL: feature.Properties.URL,
			SourceID:   s.SourceID,
			EventID:    feature.ID,
		}

		events = append(events, data)
	}

	return events, nil
}

type UsgsGeoJsonResponse struct {
	Type     string `json:"type"`
	Metadata struct {
		Generated int64  `json:"generated"`
		URL       string `json:"url"`
		Title     string `json:"title"`
		Status    int    `json:"status"`
		API       string `json:"api"`
		Count     int    `json:"count"`
	} `json:"metadata"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			Mag     float64     `json:"mag"`
			Place   string      `json:"place"`
			Time    int64       `json:"time"`
			Updated int64       `json:"updated"`
			Tz      interface{} `json:"tz"`
			URL     string      `json:"url"`
			Detail  string      `json:"detail"`
			Felt    interface{} `json:"felt"`
			Cdi     interface{} `json:"cdi"`
			Mmi     interface{} `json:"mmi"`
			Alert   interface{} `json:"alert"`
			Status  string      `json:"status"`
			Tsunami int         `json:"tsunami"`
			Sig     int         `json:"sig"`
			Net     string      `json:"net"`
			Code    string      `json:"code"`
			Ids     string      `json:"ids"`
			Sources string      `json:"sources"`
			Types   string      `json:"types"`
			Nst     int         `json:"nst"`
			Dmin    float64     `json:"dmin"`
			Rms     float64     `json:"rms"`
			Gap     float64     `json:"gap"`
			Magtype string      `json:"magType"`
			Type    string      `json:"type"`
			Title   string      `json:"title"`
		} `json:"properties"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		ID string `json:"id"`
	} `json:"features"`
	Bbox []float64 `json:"bbox"`
}
