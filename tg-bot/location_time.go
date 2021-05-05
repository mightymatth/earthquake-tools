package main

import (
	"fmt"
	"gopkg.in/ugjka/go-tz.v2/tz"
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
