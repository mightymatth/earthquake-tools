package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
)

var webhook = flag.String("webhook",
	"http://localhost:3300",
	"webhook address for events")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	output := start(ctx)

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

				err := sendToWebhook(event, *webhook)
				if err != nil {
					log.Println(err)
				}
				log.Println(event)
			}
		}
	}()

	select {
	case <-done:
	case <-interrupt:
	case <-ctx.Done():
	}
}
