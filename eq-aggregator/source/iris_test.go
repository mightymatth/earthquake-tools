package source_test

import (
	"bytes"
	"github.com/mightymatth/earthquake-tools/eq-aggregator/source"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestIris_Transform(t *testing.T) {
	s := source.NewIris()

	events, err := s.Transform(bytes.NewReader([]byte(irisResponse)))
	if err != nil {
		t.Errorf("cannot transform events: %v", err)
		return
	}

	assert.Equal(t, 1, len(events))
	eventsValid(t, events)
}

func eventsValid(t *testing.T, events []source.EarthquakeData) {
	for _, event := range events {
		assert.GreaterOrEqual(t, event.Mag, 0.0)
		assert.Less(t, event.Mag, 13.0)
		assert.NotEmpty(t, event.MagType)
		assert.GreaterOrEqual(t, event.Depth, 0.0)
		assert.True(t, time.Now().After(event.Time))
		// If the time is incorrectly parsed, it will be set to
		// Mon Jan 2 15:04:05 MST 2006 (Unix time 1136239445),
		// so we check the event time to be after that.
		assert.True(t, event.Time.After(time.Unix(1136239445, 0)))
		assert.LessOrEqual(t, event.Lat, 90.0)
		assert.GreaterOrEqual(t, event.Lat, -90.0)
		assert.LessOrEqual(t, event.Lon, 180.0)
		assert.GreaterOrEqual(t, event.Lat, -180.0)
		assert.NotEmpty(t, event.Location)

		_, err := url.Parse(event.DetailsURL)
		if err != nil {
			t.Errorf("unable to parse URL %v: %v", event.DetailsURL, err)
		}

		assert.NotEmpty(t, event.SourceID)
		assert.NotEmpty(t, event.EventID)

		// Uncomment for debugging
		//t.Logf("%+v", event)
	}
}

var irisResponse = `
This XML file does not appear to have any style information associated with it. The document tree is shown below.
<q:quakeml xmlns:q="http://quakeml.org/xmlns/quakeml/1.2" xmlns:iris="http://service.iris.edu/fdsnws/event/1/" xmlns="http://quakeml.org/xmlns/bed/1.2" xmlns:xsi="http://www.w3.org/2000/10/XMLSchema-instance" xsi:schemaLocation="http://quakeml.org/schema/xsd http://quakeml.org/schema/xsd/QuakeML-1.2.xsd">
<script/>
<eventParameters publicID="smi:service.iris.edu/fdsnws/event/1/query">
<event publicID="smi:service.iris.edu/fdsnws/event/1/query?eventid=11409051">
<type>earthquake</type>
<description xmlns:iris="http://service.iris.edu/fdsnws/event/1/" iris:FEcode="613">
<type>Flinn-Engdahl region</type>
<text>HAWAII</text>
</description>
<preferredMagnitudeID>smi:service.iris.edu/fdsnws/event/1/query?magnitudeid=206153890</preferredMagnitudeID>
<preferredOriginID>smi:service.iris.edu/fdsnws/event/1/query?originid=45059393</preferredOriginID>
<origin xmlns:iris="http://service.iris.edu/fdsnws/event/1/" publicID="smi:service.iris.edu/fdsnws/event/1/query?originid=45059393" iris:contributorOriginId="hv72450627" iris:contributor="hv" iris:contributorEventId="hv72450627" iris:catalog="NEIC PDE">
<time>
<value>2021-05-01T17:22:03.430</value>
</time>
<creationInfo>
<author>hv</author>
</creationInfo>
<latitude>
<value>19.193666</value>
</latitude>
<longitude>
<value>-155.414169</value>
</longitude>
<depth>
<value>31860.001</value>
</depth>
</origin>
<magnitude publicID="smi:service.iris.edu/fdsnws/event/1/query?magnitudeid=206153890">
<mag>
<value>1.91999996</value>
</mag>
<type>Md</type>
<creationInfo>
<author>HV</author>
</creationInfo>
</magnitude>
</event>
</eventParameters>
</q:quakeml>
`
