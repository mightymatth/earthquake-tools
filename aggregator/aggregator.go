package aggregator

import (
	"context"
	"github.com/mightymatth/earthquake-tools/aggregator/source"
	"sync"
)

// Start starts to listen to earthquake sources pushing the data to the output channel.
func Start(ctx context.Context) <-chan source.EarthquakeData {
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
