package source

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const EmscWsID ID = "EMSC_WS"

type EmscWs struct {
	source
}

func NewEmscWs() EmscWs {
	return EmscWs{source{
		Name: "EMSC WS", Url: "wss://www.seismicportal.eu/standing_order/websocket",
		Method: WEBSOCKET, SourceID: EmscWsID,
	}}
}

func (s EmscWs) Transform(r io.Reader) ([]EarthquakeData, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read from buffer: %v", err)
	}

	var event EmscEvent
	err = json.Unmarshal(buf.Bytes(), &event)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal event data: %v", err)
	}

	if event.Action != "create" {
		return []EarthquakeData{}, nil
	}

	data := EarthquakeData{
		Mag:      event.Data.Properties.Mag,
		MagType:  event.Data.Properties.MagType,
		Depth:    event.Data.Properties.Depth,
		Time:     event.Data.Properties.Time,
		Lat:      event.Data.Properties.Lat,
		Lon:      event.Data.Properties.Lon,
		Location: event.Data.Properties.FlynnRegion,
		DetailsURL: fmt.Sprintf(`https://www.emsc-csem.org/Earthquake/earthquake.php?id=%s`,
			event.Data.Properties.SourceID),
		SourceID: s.SourceID,
		EventID:  event.Data.Properties.SourceID,
	}

	return []EarthquakeData{data}, nil
}

type EmscEvent struct {
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
