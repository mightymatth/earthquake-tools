package source

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var webhook = flag.String("webhook",
	"http://localhost:3300",
	"webhook address for events")

var cache *ristretto.Cache

func init() {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic(err)
	}

	cache = c
}

type source struct {
	Name     string
	Url      string
	Method   Method
	SourceID ID
}

func (s source) Listen(ctx context.Context, lt LocateTransformer) {
	switch s.Method {
	case WEBSOCKET:
		s.handleWS(ctx, lt)
	case REST:
		s.handleREST(ctx, lt)
	}
}

func (s source) Locate() *url.URL {
	lURL, err := url.Parse(s.Url)
	if err != nil {
		log.Fatalf("incorrect URL (%v) from source '%s': %v",
			s.Url, s.Name, err)
	}
	return lURL
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

func (s source) handleWS(ctx context.Context, lt LocateTransformer) {
	restart := make(chan struct{})
	done := make(chan struct{})

startConnWS:
	log.Printf("[WS][%s] connecting...", s.Name)

	sourceURL := lt.Locate().String()
	conn, _, err := websocket.DefaultDialer.Dial(sourceURL, nil)
	if err != nil {
		defer close(done)
		log.Fatalf("[WS][%s] websocket dial error: %s", s.Name, err)
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

			for _, data := range events {
				err = s.sendToWebhook(data)
				if err != nil {
					log.Printf("[WS][%s] sending to webhook failed: %s", s.Name, err)
					continue
				}
			}
		}
	}()

	for {
		select {
		case <-restart:
			log.Printf("[WS][%s] restarting!", s.Name)
			time.Sleep(200 * time.Millisecond) // avoid spam
			goto startConnWS
		case <-done:
			return
		case <-ctx.Done():
			// Closing a websocket connection in a clean way.
			err := conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Printf("[WS][%s] write close failed: %v", s.Name, err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (s source) sendToWebhook(data EarthquakeData) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("cannot marshal earthquake data to JSON: %s", err)
	}

	_, err = http.Post(*webhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}

	log.Printf("[WS][%s] sent to webhook: %+v", s.Name, data)

	return nil
}

func (s source) handleREST(ctx context.Context, lt LocateTransformer) {
	log.Printf("[REST][%s] connecting...", s.Name)

	events, err := s.getEvents(ctx, lt)
	if err != nil {
		log.Fatalf("[REST][%s] initial fetch failed: %s", s.Name, err)
	}

	err = s.setEventsToCache(events)
	if err != nil {
		log.Fatalf("[REST][%s] initial set to cache failed: %s", s.Name, err)
	}

	log.Printf("[REST][%s] connected!", s.Name)
	for {
		time.Sleep(15 * time.Second)

		events, err := s.getEvents(ctx, lt)
		if err != nil {
			log.Printf("[REST][%s] initial fetch failed: %s", s.Name, err)
			continue
		}

		cacheEvents, err := s.getEventsFromCache()
		if err != nil {
			log.Fatalf("[REST][%s] fetch from cache failed: %s", s.Name, err)
		}

		diffEvents := difference(events, cacheEvents)

		if len(diffEvents) > 0 {
			err = s.setEventsToCache(events)
			if err != nil {
				log.Fatalf("[REST][%s] set to cache failed: %s", s.Name, err)
			}
		}

		for _, event := range diffEvents {
			err = s.sendToWebhook(event)
			if err != nil {
				log.Printf("[REST][%s] sending to webhook failed: %s", s.Name, err)
				continue
			}
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
	success := cache.Set(s.getHashKey(), events, 0)
	if !success {
		return fmt.Errorf("set to cache failed")
	}

	return nil
}

func (s source) getEventsFromCache() ([]EarthquakeData, error) {
	value, found := cache.Get(s.getHashKey())
	if !found {
		// wait for value to pass through buffers
		time.Sleep(10 * time.Millisecond)

		value, found = cache.Get(s.getHashKey())
		if !found {
			return nil, fmt.Errorf("events not found for this key")
		}
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

//EarthquakeData represents a narrow set of attributes that expresses a single
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
	WEBSOCKET Method = "WEBSOCKET"
)

type ID string
