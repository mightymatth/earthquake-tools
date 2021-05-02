package source_test

import (
	"bytes"
	"github.com/mightymatth/earthquake-tools/ws-middleware/source"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmscWs_Transform(t *testing.T) {
	s := source.NewEmscWs()

	events, err := s.Transform(bytes.NewReader([]byte(emscWsResponse)))
	if err != nil {
		t.Errorf("cannot transform events: %v", err)
		return
	}

	assert.Equal(t, 1, len(events))
	eventsValid(t, events)
}

var emscWsResponse = `
{
   "action":"create",
   "data":{
      "geometry":{
         "type":"Point",
         "coordinates":[
            -71.53,
            -30.07,
            -38.0
         ]
      },
      "type":"Feature",
      "id":"20210502_0000100",
      "properties":{
         "lastupdate":"2021-05-02T13:01:00.0Z",
         "magtype":"ml",
         "evtype":"ke",
         "lon":-71.53,
         "auth":"GUC",
         "lat":-30.07,
         "depth":38.0,
         "unid":"20210502_0000100",
         "mag":2.7,
         "time":"2021-05-02T12:46:04.0Z",
         "source_id":"978886",
         "source_catalog":"EMSC-RTS",
         "flynn_region":"OFFSHORE COQUIMBO, CHILE"
      }
   }
}
`
