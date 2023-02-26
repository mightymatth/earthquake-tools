package source

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/url"
	"time"
)

// FdsnWs is a base for implementing a feed from
// [FDSNWS](https://www.fdsn.org/webservices/) event source.
type FdsnWs struct {
	source
}

func NewFdsnWs(name, url string, sourceID ID) FdsnWs {
	return FdsnWs{source{
		Name: name, Url: url,
		Method: REST, SourceID: sourceID,
	}}
}

func (s FdsnWs) Locate() *url.URL {
	u, err := url.Parse(s.Url)
	if err != nil {
		log.Fatalf("incorrect URL (%v) from source '%s': %v",
			s.Url, s.Name, err)
	}

	q := u.Query()
	q.Set("starttime", time.Now().Add(-48*time.Hour).Format("2006-01-02"))
	q.Set("limit", fmt.Sprintf("%d", 10))
	u.RawQuery = q.Encode()

	return u
}

type fdsnwsTime time.Time

func (c *fdsnwsTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string

	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}

	parse, err := time.Parse("2006-01-02T15:04:05", v)
	if err != nil {
		return err
	}

	*c = fdsnwsTime(parse)
	return nil
}
