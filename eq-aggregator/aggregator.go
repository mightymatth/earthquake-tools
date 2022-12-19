package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mightymatth/earthquake-tools/eq-aggregator/source"
	"log"
	"net/http"
	"sync"
)

func start(ctx context.Context) <-chan source.EarthquakeData {
	sources := []ListenTransformer{
		source.NewEmscWs(),
		source.NewUsgs(),
		source.NewEmscRest(),
		source.NewIris(),
		source.NewUspBr(),
		source.NewGeofon(),
	}

	output := make(chan source.EarthquakeData)

	var wg sync.WaitGroup
	for _, src := range sources {
		wg.Add(1)
		go func(src ListenTransformer) {
			defer wg.Done()
			src.Listen(ctx, src, output)
		}(src)
	}

	go func() {
		wg.Wait()
		close(output)
	}()

	return output
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

type Listener interface {
	Listen(
		ctx context.Context,
		lt source.LocateTransformer,
		output chan<- source.EarthquakeData,
	)
}

type ListenTransformer interface {
	source.LocateTransformer
	Listener
}
