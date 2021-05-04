package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/ugjka/go-tz.v2/tz"
	"net/http"
	"time"
	_ "time/tzdata"
)

func LocationTime(timeUTC time.Time, lat, lon float64) (*time.Time, error) {
	zone, err := tz.GetZone(tz.Point{
		Lat: lat, Lon: lon,
	})
	if err != nil {
		return nil, fmt.Errorf("get zone failed: %v", err)
	}

	loc, err := time.LoadLocation(zone[0])
	if err != nil {
		return nil, fmt.Errorf("load location failed: %v", err)
	}

	localTime := timeUTC.In(loc)

	return &localTime, nil
}

func LocationTimeTimezoneDB(timeUTC time.Time, lat, lon float64) (*time.Time, error) {
	resp, err := http.Get(fmt.Sprintf(
		"http://api.timezonedb.com/v2.1/get-time-zone?by=position&format=json&lat=%f&lng=%f&time=%d&key=%s",
		lat, lon, timeUTC.Unix(), getEnv("TIMEZONEDB_API_KEY", "")))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("timezonedb failed: %v", err)
	}

	var timezoneRes TimezoneDBResponse

	err = json.NewDecoder(resp.Body).Decode(&timezoneRes)
	if err != nil {
		return nil, fmt.Errorf("timezonedb response parse failed: %v", err)
	}
	if timezoneRes.Status != "OK" {
		return nil, fmt.Errorf("timezonedb invalid response: (status: %v, message: %v)",
			timezoneRes.Status, timezoneRes.Message)
	}

	loc, err := time.LoadLocation(timezoneRes.ZoneName)
	if err != nil {
		return nil, fmt.Errorf("load location: %v", err)
	}

	localTime := timeUTC.In(loc)

	return &localTime, nil
}

type TimezoneDBResponse struct {
	Status           string      `json:"status"`
	Message          string      `json:"message"`
	CountryCode      string      `json:"countryCode"`
	CountryName      string      `json:"countryName"`
	ZoneName         string      `json:"zoneName"`
	Abbreviation     string      `json:"abbreviation"`
	GmtOffset        int         `json:"gmtOffset"`
	Dst              string      `json:"dst"`
	ZoneStart        int         `json:"zoneStart"`
	ZoneEnd          interface{} `json:"zoneEnd"`
	NextAbbreviation interface{} `json:"nextAbbreviation"`
	Timestamp        int64       `json:"timestamp"`
	Formatted        string      `json:"formatted"`
}
