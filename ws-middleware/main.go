package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var source = flag.String("source",
	"wss://www.seismicportal.eu/standing_order/websocket",
	"source websocket address")

var webhook = flag.String("webhook",
	"http://localhost:3300",
	"webhook address for events")

func main() {
	flag.Parse()
	log.SetFlags(log.LUTC)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", *source)

	c, _, err := websocket.DefaultDialer.Dial(*source, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	log.Print("connected!")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("recv: %s", message)

			_, err = http.Post(*webhook, "application/json", bytes.NewBuffer(message))
			if err != nil {
				log.Println("sending to webhook:", err)
				continue
			}

			log.Print("sent to webhook")
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Println("write close:", err)
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
