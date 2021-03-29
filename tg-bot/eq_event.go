package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

func ParseEvent(body io.Reader) (event EarthquakeEvent, err error) {
	err = json.NewDecoder(body).Decode(&event)
	if err != nil {
		return event, fmt.Errorf("cannot decode earthquake event: %v", err)
	}

	return event, err
}

type EarthquakeEvent struct {
	Action string `json:"action"`
	Data   struct {
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Type       string `json:"type"`
		ID         string `json:"id"`
		Properties struct {
			LastUpdate    time.Time `json:"lastupdate"`
			MagType       string    `json:"magtype"`
			EvType        string    `json:"evtype"`
			Lon           float64   `json:"lon"`
			Auth          string    `json:"auth"`
			Lat           float64   `json:"lat"`
			Depth         float64   `json:"depth"`
			UnID          string    `json:"unid"`
			Mag           float64   `json:"mag"`
			Time          time.Time `json:"time"`
			SourceID      string    `json:"source_id"`
			SourceCatalog string    `json:"source_catalog"`
			FlynnRegion   string    `json:"flynn_region"`
		} `json:"properties"`
	} `json:"data"`
}
