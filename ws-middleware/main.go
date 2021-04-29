package main

import (
	"context"
	"github.com/mightymatth/earthquake-tools/ws-middleware/source"
	"os"
	"os/signal"
	"sync"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sources := []ListenTransformer{
		source.NewEmscWsSource("EMSC WS",
			"wss://www.seismicportal.eu/standing_order/websocket"),
		source.NewUsgsSource("USGS",
			"https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson"),
	}

	var wg sync.WaitGroup
	for _, src := range sources {
		wg.Add(1)
		go func(src ListenTransformer) {
			defer wg.Done()
			src.Listen(ctx, src)
		}(src)
	}

	c := make(chan struct{})
	go func() {
		wg.Wait()
		c <- struct{}{}
	}()

	select {
	case <-c:
	case <-interrupt:
	}
}

type Listener interface {
	Listen(ctx context.Context, t source.Transformer)
}

type ListenTransformer interface {
	Listener
	source.Transformer
}
