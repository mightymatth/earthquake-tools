package source

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

const UspBrID ID = "USPBR"

type UspBr struct {
	FdsnWs
}

func NewUspBr() UspBr {
	return UspBr{
		NewFdsnWs("USPBR", "http://www.moho.iag.usp.br/fdsnws/event/1/query", UspBrID),
	}
}

func (s UspBr) Transform(r io.Reader) ([]EarthquakeData, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("cannot read from buffer: %v", err)
	}

	if buf.Len() <= 0 {
		return []EarthquakeData{}, nil
	}

	var eventsRes UspBrResponse
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
			Depth:      feature.Origin.Depth.Value / 1000,
			Time:       time.Time(feature.Origin.Time.Value).UTC(),
			Lat:        feature.Origin.Latitude.Value,
			Lon:        feature.Origin.Longitude.Value,
			Location:   feature.Description.Text,
			DetailsURL: fmt.Sprintf("http://www.moho.iag.usp.br/eq/event/%s", geofonEventID(feature.PublicID)),
			SourceID:   s.SourceID,
			EventID:    geofonEventID(feature.PublicID),
		}

		events = append(events, data)
	}

	return events, nil
}

type UspBrResponse struct {
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
			TypeCertainty string `xml:"typeCertainty"`
			CreationInfo  struct {
				Text         string `xml:",chardata"`
				AgencyID     string `xml:"agencyID"`
				Author       string `xml:"author"`
				CreationTime string `xml:"creationTime"`
			} `xml:"creationInfo"`
			Magnitude struct {
				Text         string `xml:",chardata"`
				PublicID     string `xml:"publicID,attr"`
				StationCount string `xml:"stationCount"`
				CreationInfo struct {
					Text         string `xml:",chardata"`
					AgencyID     string `xml:"agencyID"`
					Author       string `xml:"author"`
					CreationTime string `xml:"creationTime"`
				} `xml:"creationInfo"`
				Mag struct {
					Text        string  `xml:",chardata"`
					Value       float64 `xml:"value"`
					Uncertainty string  `xml:"uncertainty"`
				} `xml:"mag"`
				Type             string `xml:"type"`
				OriginID         string `xml:"originID"`
				MethodID         string `xml:"methodID"`
				EvaluationStatus string `xml:"evaluationStatus"`
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
					DepthPhaseCount        string `xml:"depthPhaseCount"`
					StandardError          string `xml:"standardError"`
					AzimuthalGap           string `xml:"azimuthalGap"`
					MaximumDistance        string `xml:"maximumDistance"`
					MinimumDistance        string `xml:"minimumDistance"`
					MedianDistance         string `xml:"medianDistance"`
				} `xml:"quality"`
				EvaluationMode string `xml:"evaluationMode"`
				CreationInfo   struct {
					Text         string `xml:",chardata"`
					AgencyID     string `xml:"agencyID"`
					Author       string `xml:"author"`
					CreationTime string `xml:"creationTime"`
				} `xml:"creationInfo"`
				Depth struct {
					Text        string  `xml:",chardata"`
					Value       float64 `xml:"value"`
					Uncertainty string  `xml:"uncertainty"`
				} `xml:"depth"`
				OriginUncertainty struct {
					Text                            string `xml:",chardata"`
					MinHorizontalUncertainty        string `xml:"minHorizontalUncertainty"`
					MaxHorizontalUncertainty        string `xml:"maxHorizontalUncertainty"`
					AzimuthMaxHorizontalUncertainty string `xml:"azimuthMaxHorizontalUncertainty"`
					PreferredDescription            string `xml:"preferredDescription"`
				} `xml:"originUncertainty"`
				MethodID         string `xml:"methodID"`
				EarthModelID     string `xml:"earthModelID"`
				EvaluationStatus string `xml:"evaluationStatus"`
				DepthType        string `xml:"depthType"`
			} `xml:"origin"`
			PreferredOriginID    string `xml:"preferredOriginID"`
			PreferredMagnitudeID string `xml:"preferredMagnitudeID"`
			Type                 string `xml:"type"`
		} `xml:"event"`
	} `xml:"eventParameters"`
}
