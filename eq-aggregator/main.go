package main

import (
	"context"
	"flag"
	"github.com/mightymatth/earthquake-tools/eq-aggregator/source"
	"os"
	"os/signal"
	"sync"
)

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sources := []ListenTransformer{
		source.NewEmscWs(),
		source.NewUsgs(),
		source.NewEmscRest(),
		source.NewIris(),
		source.NewUspBr(),
		source.NewGeofon(),
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
	Listen(ctx context.Context, lt source.LocateTransformer)
}

type ListenTransformer interface {
	source.LocateTransformer
	Listener
}
