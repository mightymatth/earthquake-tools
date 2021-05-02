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

const GeofonID ID = "GEOFON"

type Geofon struct {
	FdsnWs
}

func NewGeofon() Geofon {
	return Geofon{
		NewFdsnWs("GEOFON", "https://geofon.gfz-potsdam.de/fdsnws/event/1/query", GeofonID),
	}
}

func (s Geofon) Transform(r io.Reader) ([]EarthquakeData, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read from buffer: %v", err)
	}

	var eventsRes GeofonResponse
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
			Mag:      feature.Magnitude.Mag.Value,
			MagType:  strings.ToLower(feature.Magnitude.Type),
			Depth:    math.Abs(feature.Origin.Depth.Value / 1000),
			Time:     time.Time(feature.Origin.Time.Value).UTC(),
			Lat:      feature.Origin.Latitude.Value,
			Lon:      feature.Origin.Longitude.Value,
			Location: feature.Description.Text,
			DetailsURL: fmt.Sprintf("https://geofon.gfz-potsdam.de/eqinfo/event.php?id=%s",
				geofonEventID(feature.PublicID)),
			SourceID: s.SourceID,
			EventID:  geofonEventID(feature.PublicID),
		}

		events = append(events, data)
	}

	return events, nil
}

func geofonEventID(publicIDPath string) string {
	r, _ := regexp.Compile("geofon/(.+)$")

	parts := r.FindStringSubmatch(publicIDPath)
	if len(parts) != 2 {
		return ""
	}

	return parts[1]
}

type GeofonResponse struct {
	XMLName         xml.Name `xml:"quakeml"`
	Text            string   `xml:",chardata"`
	Xmlns           string   `xml:"xmlns,attr"`
	Q               string   `xml:"q,attr"`
	EventParameters struct {
		Text     string `xml:",chardata"`
		PublicID string `xml:"publicID,attr"`
		Event    []struct {
			Text        string `xml:",chardata"`
			PublicID    string `xml:"publicID,attr"`
			Description struct {
				Chardata string `xml:",chardata"`
				Text     string `xml:"text"`
				Type     string `xml:"type"`
			} `xml:"description"`
			CreationInfo struct {
				Text         string `xml:",chardata"`
				AgencyID     string `xml:"agencyID"`
				CreationTime string `xml:"creationTime"`
			} `xml:"creationInfo"`
			Magnitude struct {
				Text         string `xml:",chardata"`
				PublicID     string `xml:"publicID,attr"`
				StationCount string `xml:"stationCount"`
				CreationInfo struct {
					Text         string `xml:",chardata"`
					AgencyID     string `xml:"agencyID"`
					CreationTime string `xml:"creationTime"`
				} `xml:"creationInfo"`
				Mag struct {
					Text        string  `xml:",chardata"`
					Value       float64 `xml:"value"`
					Uncertainty string  `xml:"uncertainty"`
				} `xml:"mag"`
				Type     string `xml:"type"`
				OriginID string `xml:"originID"`
				MethodID string `xml:"methodID"`
			} `xml:"magnitude"`
			Origin struct {
				Text     string `xml:",chardata"`
				PublicID string `xml:"publicID,attr"`
				Time     struct {
					Text        string     `xml:",chardata"`
					Value       geofonTime `xml:"value"`
					Uncertainty string     `xml:"uncertainty"`
				} `xml:"time"`
				Longitude struct {
					Text        string  `xml:",chardata"`
					Value       float64 `xml:"value"`
					Uncertainty string  `xml:"uncertainty"`
				} `xml:"longitude"`
				Latitude struct {
					Text        string  `xml:",chardata"`
					Value       float64 `xml:"value"`
					Uncertainty string  `xml:"uncertainty"`
				} `xml:"latitude"`
				Quality struct {
					Text                   string `xml:",chardata"`
					AssociatedPhaseCount   string `xml:"associatedPhaseCount"`
					UsedPhaseCount         string `xml:"usedPhaseCount"`
					AssociatedStationCount string `xml:"associatedStationCount"`
					UsedStationCount       string `xml:"usedStationCount"`
					StandardError          string `xml:"standardError"`
					AzimuthalGap           string `xml:"azimuthalGap"`
					MaximumDistance        string `xml:"maximumDistance"`
					MinimumDistance        string `xml:"minimumDistance"`
					MedianDistance         string `xml:"medianDistance"`
					DepthPhaseCount        string `xml:"depthPhaseCount"`
				} `xml:"quality"`
				EvaluationMode string `xml:"evaluationMode"`
				CreationInfo   struct {
					Text         string `xml:",chardata"`
					AgencyID     string `xml:"agencyID"`
					CreationTime string `xml:"creationTime"`
				} `xml:"creationInfo"`
				Depth struct {
					Text        string  `xml:",chardata"`
					Value       float64 `xml:"value"`
					Uncertainty string  `xml:"uncertainty"`
				} `xml:"depth"`
				MethodID          string `xml:"methodID"`
				EarthModelID      string `xml:"earthModelID"`
				DepthType         string `xml:"depthType"`
				OriginUncertainty struct {
					Text                            string `xml:",chardata"`
					MinHorizontalUncertainty        string `xml:"minHorizontalUncertainty"`
					MaxHorizontalUncertainty        string `xml:"maxHorizontalUncertainty"`
					AzimuthMaxHorizontalUncertainty string `xml:"azimuthMaxHorizontalUncertainty"`
					PreferredDescription            string `xml:"preferredDescription"`
				} `xml:"originUncertainty"`
				EvaluationStatus string `xml:"evaluationStatus"`
			} `xml:"origin"`
			PreferredOriginID    string `xml:"preferredOriginID"`
			PreferredMagnitudeID string `xml:"preferredMagnitudeID"`
		} `xml:"event"`
	} `xml:"eventParameters"`
}

type geofonTime time.Time

func (c *geofonTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	parse, err := time.Parse("2006-01-02T15:04:05.999999Z", v)
	if err != nil {
		return err
	}

	*c = geofonTime(parse)
	return nil
}
