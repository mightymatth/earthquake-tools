package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"time"
)

const EmscRestID ID = "EMSC_REST"

type EmscRest struct {
	source
}

func NewEmscRest() EmscRest {
	return EmscRest{source{
		Name: "EMSC", Url: "https://www.seismicportal.eu/fdsnws/event/1/query?limit=10&format=json",
		Method: REST, SourceID: EmscRestID,
	}}
}

func (s EmscRest) Transform(r io.Reader) ([]EarthquakeData, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read from buffer: %v", err)
	}

	var eventsRes EmscRestResponse
	err = json.Unmarshal(buf.Bytes(), &eventsRes)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal event data: %v", err)
	}

	features := eventsRes.Features[:10]
	events := make([]EarthquakeData, 0, len(features))
	for _, feature := range features {
		data := EarthquakeData{
			Mag:      feature.Properties.Mag,
			MagType:  feature.Properties.Magtype,
			Depth:    math.Abs(feature.Geometry.Coordinates[2]),
			Time:     feature.Properties.Time,
			Lat:      feature.Geometry.Coordinates[1],
			Lon:      feature.Geometry.Coordinates[0],
			Location: feature.Properties.FlynnRegion,
			DetailsURL: fmt.Sprintf(`https://www.emsc-csem.org/Earthquake/earthquake.php?id=%s`,
				feature.Properties.SourceID),
			SourceID: s.SourceID,
			EventID:  feature.Properties.SourceID,
		}

		events = append(events, data)
	}

	return events, nil
}

type EmscRestResponse struct {
	Type     string `json:"type"`
	Metadata struct {
		Totalcount int `json:"totalCount"`
	} `json:"metadata"`
	Features []struct {
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Type       string `json:"type"`
		ID         string `json:"id"`
		Properties struct {
			Lastupdate    time.Time `json:"lastupdate"`
			Magtype       string    `json:"magtype"`
			Evtype        string    `json:"evtype"`
			Lon           float64   `json:"lon"`
			Auth          string    `json:"auth"`
			Lat           float64   `json:"lat"`
			Depth         float64   `json:"depth"`
			Unid          string    `json:"unid"`
			Mag           float64   `json:"mag"`
			Time          time.Time `json:"time"`
			SourceID      string    `json:"source_id"`
			SourceCatalog string    `json:"source_catalog"`
			FlynnRegion   string    `json:"flynn_region"`
		} `json:"properties"`
	} `json:"features"`
}