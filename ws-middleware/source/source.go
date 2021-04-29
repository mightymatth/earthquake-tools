package source

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"time"
)

var webhook = flag.String("webhook",
	"http://localhost:3300",
	"webhook address for events")

func init() {
	flag.Parse()
}

type source struct {
	Name     string
	Url      string
	Method   Method
	SourceID ID
}

func (s source) Listen(ctx context.Context, t Transformer) {
	switch s.Method {
	case WEBSOCKET:
		s.handleWS(ctx, t)
	case REST:
		s.handleREST(ctx)
	}
}

type Transformer interface {
	Transform(r io.Reader) ([]EarthquakeData, error)
}

func (s source) handleWS(ctx context.Context, t Transformer) {
	restart := make(chan struct{})
	done := make(chan struct{})

startConnWS:
	log.Printf("[WS][%s] connecting to %s", s.Name, s.Url)

	conn, _, err := websocket.DefaultDialer.Dial(s.Url, nil)
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
			log.Printf("[WS][%s] received: %s", s.Name, message)

			events, err := t.Transform(bytes.NewBuffer(message))
			if err != nil {
				log.Printf("[WS][%s] transform data failed: %s", s.Name, err)
				continue
			}

			log.Printf("events: %v, length: %v", events, len(events))

			for _, data := range events {
				err = sendToWebhook(data)
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

func sendToWebhook(data EarthquakeData) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("cannot marshal earthquake data to JSON: %s", err)
	}

	_, err = http.Post(*webhook, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}

	return nil
}

func (s source) handleREST(ctx context.Context) {

}

//EarthquakeData represents a narrow set of attributes that expresses a single
// earthquake.
type EarthquakeData struct {
	Mag        float64   `json:"mag"`
	MagType    string    `json:"magtype"`
	Time       time.Time `json:"time"`
	Lat        float64   `json:"lat"`
	Lon        float64   `json:"lon"`
	Location   string    `json:"location"`
	DetailsURL string    `json:"details_url"`
	SourceID   ID        `json:"source"`
}

type Method string

const (
	REST      Method = "REST"
	WEBSOCKET Method = "WEBSOCKET"
)

type ID string
