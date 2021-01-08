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
	restart := make(chan struct{})
	done := make(chan struct{})

	conn := listen(restart, done)

	for {
		select {
		case <-restart:
			log.Printf("restart")
			time.Sleep(200 * time.Millisecond) // avoid spam
			restart = make(chan struct{})
			conn = listen(restart, done)
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Closing a websocket connection in a clean way.
			err := conn.WriteMessage(
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

func listen(restart, done chan<- struct{}) *websocket.Conn {
	log.Printf("connecting to %s", *source)

	conn, _, err := websocket.DefaultDialer.Dial(*source, nil)
	if err != nil {
		defer close(done)
		log.Fatal("dial:", err)
	}
	log.Print("connected!")

	go func() {
		defer conn.Close()
		defer close(restart)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			log.Printf("received: %s", message)

			_, err = http.Post(*webhook, "application/json", bytes.NewBuffer(message))
			if err != nil {
				log.Println("error sending to webhook:", err)
				continue
			}

			log.Print("sent to webhook")
		}
	}()

	return conn
}
