package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

func ParseEvent(body io.Reader) (event EarthquakeEvent, err error) {
	err = json.NewDecoder(body).Decode(&event)
	if err != nil {
		return event, fmt.Errorf("cannot decode earthquake event: %v", err)
	}

	log.Printf("parsed event: %+v", event)
	return event, err
}

type EarthquakeEvent struct {
	Mag        float64   `json:"mag"`
	MagType    string    `json:"magtype"`
	Depth      float64   `json:"depth"`
	Time       time.Time `json:"time"`
	Lat        float64   `json:"lat"`
	Lon        float64   `json:"lon"`
	Location   string    `json:"location"`
	DetailsURL string    `json:"details_url"`
	SourceID   string    `json:"source"`
	EventID    string    `json:"event_id"`
}
