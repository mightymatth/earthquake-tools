package source

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/mightymatth/earthquake-tools/eq-aggregator/cache"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var sourceCache cache.Cacher

func init() {
	c := cache.NewInMemory()

	// Change to Ristretto if needed.
	//c, err := cache.NewRistretto()
	//if err != nil {
	//	panic(err)
	//}

	sourceCache = c
}

type source struct {
	Name     string
	Url      string
	Method   Method
	SourceID ID
}

func (s source) Listen(
	ctx context.Context, lt LocateTransformer, output chan<- EarthquakeData,
) {
	switch s.Method {
	case WEBSOCKET:
		s.listenWS(ctx, lt, output)
	case REST:
		s.listenREST(ctx, lt, output)
	}
}

func (s source) Locate() *url.URL {
	u, err := url.Parse(s.Url)
	if err != nil {
		log.Fatalf("incorrect URL (%v) from source '%s': %v",
			s.Url, s.Name, err)
	}

	return u
}

type Transformer interface {
	Transform(r io.Reader) ([]EarthquakeData, error)
}

type Locator interface {
	Locate() *url.URL
}

type LocateTransformer interface {
	Transformer
	Locator
}

func (s source) listenWS(
	ctx context.Context, lt LocateTransformer, output chan<- EarthquakeData,
) {
	restart := make(chan struct{})

start:
	log.Printf("[WS][%s] connecting...", s.Name)

	u := lt.Locate().String()
	conn, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Printf("[WS][%s] websocket dial error: %s", s.Name, err)
		time.Sleep(5 * time.Second)
		goto start
	}
	log.Printf("[WS][%s] connected!", s.Name)

	go func() {
		defer func() {
			_ = conn.Close()
			restart <- struct{}{}
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("[WS][%s] read error: %s", s.Name, err)
				return
			}

			events, err := lt.Transform(bytes.NewBuffer(message))
			if err != nil {
				log.Printf("[WS][%s] transform data failed: %s", s.Name, err)
				continue
			}

			for _, event := range events {
				output <- event
			}
		}
	}()

	for {
		select {
		case <-restart:
			log.Printf("[WS][%s] restarting!", s.Name)
			time.Sleep(200 * time.Millisecond) // avoid spam
			goto start
		case <-ctx.Done():
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Printf("[WS][%s] write close failed: %v", s.Name, err)
				return
			}

			<-time.After(time.Second)
			return
		}
	}
}

func (s source) listenREST(
	ctx context.Context, lt LocateTransformer, output chan<- EarthquakeData,
) {
start:
	log.Printf("[REST][%s] connecting...", s.Name)

	events, err := s.getEvents(ctx, lt)
	if err != nil {
		log.Printf("[REST][%s] initial fetch failed: %s", s.Name, err)
		time.Sleep(5 * time.Second)
		goto start
	}

	err = s.setEventsToCache(events)
	if err != nil {
		log.Printf("[REST][%s] initial set to cache failed: %s", s.Name, err)
		time.Sleep(5 * time.Second)
		goto start
	}

	log.Printf("[REST][%s] connected!", s.Name)
	for {
		time.Sleep(15 * time.Second)

		events, err := s.getEvents(ctx, lt)
		if err != nil {
			log.Printf("[REST][%s] fetch failed: %s", s.Name, err)
			continue
		}

		cacheEvents, err := s.getEventsFromCache()
		if err != nil {
			log.Printf("[REST][%s] fetch from cache failed: %s", s.Name, err)
			continue
		}

		diffEvents := difference(events, cacheEvents)

		if len(diffEvents) > 0 {
			err = s.setEventsToCache(events)
			if err != nil {
				log.Printf("[REST][%s] set to cache failed: %s", s.Name, err)
				continue
			}
		}

		for _, event := range diffEvents {
			output <- event
			//err = s.sendToWebhook(event)
			//if err != nil {
			//	log.Printf("[REST][%s] sending to webhook failed: %s", s.Name, err)
			//	continue
			//}
		}
	}
}

// difference returns the events in `a` that aren't in `b`.
func difference(a, b []EarthquakeData) []EarthquakeData {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x.EventID] = struct{}{}
	}
	var diff []EarthquakeData
	for _, x := range a {
		if _, found := mb[x.EventID]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func (s source) getEvents(ctx context.Context, lt LocateTransformer) ([]EarthquakeData, error) {
	sourceURL := lt.Locate().String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot create GET request: %v", err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch from endpoint: %v", err)
	}
	defer resp.Body.Close()

	events, err := lt.Transform(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot transform data: %v", err)
	}

	return events, nil
}

func (s source) setEventsToCache(events []EarthquakeData) error {
	success := sourceCache.Set(s.getHashKey(), events)
	if !success {
		return fmt.Errorf("set to cache failed")
	}

	return nil
}

func (s source) getEventsFromCache() ([]EarthquakeData, error) {
	value, found := sourceCache.Get(s.getHashKey())
	if !found {
		return nil, fmt.Errorf("events not found for this key")
	}

	events, ok := value.([]EarthquakeData)
	if !ok {
		return nil, fmt.Errorf("cannot cast from cache to exact struct")
	}

	return events, nil
}

func (s source) getHashKey() string {
	key := fmt.Sprintf("%s%s%s", s.Name, s.SourceID, s.Method)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(key)))
}

// EarthquakeData represents a narrow set of attributes that expresses a single
// earthquake.
type EarthquakeData struct {
	Mag        float64   `json:"mag"`
	MagType    string    `json:"magtype"`
	Depth      float64   `json:"depth"`
	Time       time.Time `json:"time"`
	Lat        float64   `json:"lat"`
	Lon        float64   `json:"lon"`
	Location   string    `json:"location"`
	DetailsURL string    `json:"details_url"`
	SourceID   ID        `json:"source"`
	EventID    string    `json:"event_id"`
}

type Method string

const (
	REST      Method = "REST"
	WEBSOCKET Method = "WS"
)

type ID string
