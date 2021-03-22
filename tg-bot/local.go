package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage/mongo"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	_ = godotenv.Load()

	bot, err := tgbotapi.NewBotAPI(getEnv("TELEGRAM_BOT_TOKEN", ""))
	if err != nil {
		log.Panic(err)
	}

	storageImpl, err := mongo.NewStorage("")
	if err != nil {
		log.Panicf("cannot create storage: %v", err)
	}
	defer storageImpl.Client.Disconnect(context.Background())

	service := storage.NewService(storageImpl)

	http.HandleFunc("/", EqEventHandler(bot, service))
	go func() {
		log.Panic(http.ListenAndServe(":3300", nil))
	}()

	log.Panic(TgBotServer(bot, service))
}

func getLocationTime(timeUTC time.Time, lat, lon float64) string {
	localTime, err := LocationTime(timeUTC, lat, lon)
	if err != nil {
		log.Printf("cannot get location time: %v", err)
		return "unknown"
	}

	return localTime.Format("Mon, 2 Jan 2006 15:04:05 MST")
}

func TgBotServer(bot *tgbotapi.BotAPI, s storage.Service) error {
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("cannot get updates: %v", err)
	}

	for update := range updates {
		botHandler(update, bot, s)
	}

	return fmt.Errorf("finished getting updates")
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
