package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mightymatth/earthquake-tools/aggregator"
	"github.com/mightymatth/earthquake-tools/aggregator/source"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var webhook = flag.String(
	"webhook", "", "webhook address for events (e.g. http://localhost:3300)",
)

func main() {
	flag.Parse()

	if *webhook == "" {
		log.Println("WARNING: no webhook set, will only log events")
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output := aggregator.Start(ctx)

	done := make(chan struct{})
	go func() {
		for {
			select {
			case event, ok := <-output:
				if !ok {
					output = nil
					done <- struct{}{}
					return
				}

				log.Println(event)

				if *webhook == "" {
					continue
				}

				err := sendToWebhook(event, *webhook)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()

	select {
	case <-done:
	case <-interrupt:
	case <-ctx.Done():
	}
}

func sendToWebhook(data source.EarthquakeData, webhookURL string) error {
	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("cannot marshal earthquake data to JSON: %s", err)
	}

	_, err = http.Post(webhookURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}

	log.Printf("sent to webhook: %+v", webhookURL)

	return nil
}
