package source

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const IrisID ID = "IRIS"

type Iris struct {
	source
}

func NewIris() Iris {
	return Iris{source{
		Name: "IRIS", Url: "https://service.iris.edu/fdsnws/event/1/query",
		Method: REST, SourceID: IrisID,
	}}
}

func (s Iris) Locate() *url.URL {
	lURL, err := url.Parse(s.Url)
	if err != nil {
		log.Fatalf("incorrect URL (%v) from source '%s': %v",
			s.Url, s.Name, err)
	}

	q := lURL.Query()
	q.Set("starttime", time.Now().Add(-24*time.Hour).Format("2006-01-02T15:04:05Z"))
	q.Set("limit", fmt.Sprintf("%d", 10))
	lURL.RawQuery = q.Encode()

	return lURL
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

	features := eventsRes.EventParameters.Event[:10]
	events := make([]EarthquakeData, 0, len(features))
	for _, feature := range features {
		data := EarthquakeData{
			Mag:        feature.Magnitude.Mag.Value,
			MagType:    strings.ToLower(feature.Magnitude.Type),
			Depth:      feature.Origin.Depth.Value/1000,
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
					Value customTime `xml:"value"`
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

type customTime time.Time

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	parse, err := time.Parse("2006-01-02T15:04:05", v)
	if err != nil {
		return err
	}

	*c = customTime(parse)
	return nil
}
