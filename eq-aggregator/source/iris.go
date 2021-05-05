package source

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
	"time"
)

const IrisID ID = "IRIS"

type Iris struct {
	FdsnWs
}

func NewIris() Iris {
	return Iris{
		NewFdsnWs("IRIS", "https://service.iris.edu/fdsnws/event/1/query", IrisID),
	}
}

func (s Iris) Transform(r io.Reader) ([]EarthquakeData, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read from buffer: %v", err)
	}

	var eventsRes IrisResponse
	err = xml.Unmarshal(buf.Bytes(), &eventsRes)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal event data: %v", err)
	}

	maxFeatures := 10
	if featuresLen := len(eventsRes.EventParameters.Event); featuresLen < maxFeatures {
		maxFeatures = featuresLen
	}

	features := eventsRes.EventParameters.Event[:maxFeatures]
	events := make([]EarthquakeData, 0, len(features))
	for _, feature := range features {
		data := EarthquakeData{
			Mag:        feature.Magnitude.Mag.Value,
			MagType:    strings.ToLower(feature.Magnitude.Type),
			Depth:      math.Abs(feature.Origin.Depth.Value / 1000),
			Time:       time.Time(feature.Origin.Time.Value).UTC(),
			Lat:        feature.Origin.Latitude.Value,
			Lon:        feature.Origin.Longitude.Value,
			Location:   feature.Description.Text,
			DetailsURL: irisDetailsURL(feature.PublicID),
			SourceID:   s.SourceID,
			EventID:    feature.Origin.ContributorEventId,
		}

		events = append(events, data)
	}

	return events, nil
}

func irisDetailsURL(publicIDPath string) string {
	r, _ := regexp.Compile("eventid=([0-9]+)")

	parts := r.FindStringSubmatch(publicIDPath)
	if len(parts) != 2 {
		return ""
	}

	return fmt.Sprintf("https://ds.iris.edu/ds/nodes/dmc/tools/event/%s", parts[1])
}

type IrisResponse struct {
	XMLName         xml.Name `xml:"quakeml"`
	Text            string   `xml:",chardata"`
	Q               string   `xml:"q,attr"`
	Iris            string   `xml:"iris,attr"`
	Xmlns           string   `xml:"xmlns,attr"`
	Xsi             string   `xml:"xsi,attr"`
	SchemaLocation  string   `xml:"schemaLocation,attr"`
	EventParameters struct {
		Text     string `xml:",chardata"`
		PublicID string `xml:"publicID,attr"`
		Event    []struct {
			Text        string `xml:",chardata"`
			PublicID    string `xml:"publicID,attr"`
			Type        string `xml:"type"`
			Description struct {
				Chardata string `xml:",chardata"`
				Iris     string `xml:"iris,attr"`
				FEcode   string `xml:"FEcode,attr"`
				Type     string `xml:"type"`
				Text     string `xml:"text"`
			} `xml:"description"`
			PreferredMagnitudeID string `xml:"preferredMagnitudeID"`
			PreferredOriginID    string `xml:"preferredOriginID"`
			Origin               struct {
				Text                string `xml:",chardata"`
				Iris                string `xml:"iris,attr"`
				PublicID            string `xml:"publicID,attr"`
				ContributorOriginId string `xml:"contributorOriginId,attr"`
				Contributor         string `xml:"contributor,attr"`
				ContributorEventId  string `xml:"contributorEventId,attr"`
				Catalog             string `xml:"catalog,attr"`
				Time                struct {
					Text  string     `xml:",chardata"`
					Value fdsnwsTime `xml:"value"`
				} `xml:"time"`
				CreationInfo struct {
					Text   string `xml:",chardata"`
					Author string `xml:"author"`
				} `xml:"creationInfo"`
				Latitude struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value"`
				} `xml:"latitude"`
				Longitude struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value"`
				} `xml:"longitude"`
				Depth struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value"`
				} `xml:"depth"`
			} `xml:"origin"`
			Magnitude struct {
				Text     string `xml:",chardata"`
				PublicID string `xml:"publicID,attr"`
				Mag      struct {
					Text  string  `xml:",chardata"`
					Value float64 `xml:"value"`
				} `xml:"mag"`
				Type         string `xml:"type"`
				CreationInfo struct {
					Text   string `xml:",chardata"`
					Author string `xml:"author"`
				} `xml:"creationInfo"`
			} `xml:"magnitude"`
		} `xml:"event"`
	} `xml:"eventParameters"`
}
