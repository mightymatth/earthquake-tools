package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage/mongo"
	"log"
	"net/http"
	"time"
)

func main() {
	_ = loadDotEnv()

	botAPI, err := tgbotapi.NewBotAPI(getEnv("TELEGRAM_BOT_TOKEN", ""))
	if err != nil {
		log.Panic("cannot initialize bot api:", err)
	}

	storageImpl, err := mongo.NewStorage("", getEnv("DB_NAMESPACE", ""))
	if err != nil {
		log.Panicf("cannot create storage: %v", err)
	}
	defer storageImpl.Client.Disconnect(context.Background())

	service := storage.NewService(storageImpl)

	http.HandleFunc("/", EqEventHandler(botAPI, service))
	go func() {
		log.Panic(http.ListenAndServe(":3300", nil))
	}()

	log.Panic(TgBotServer(botAPI, service))
}

func TgBotServer(bot *tgbotapi.BotAPI, s storage.Service) error {
	if getEnv("BOT_ENV", "dev") != "prod" {
		bot.Debug = true
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("cannot get updates: %v", err)
	}

	// Clearing old messages
	time.Sleep(500 * time.Millisecond)
	updates.Clear()

	for update := range updates {
		botHandler(update, bot, s)
	}

	return fmt.Errorf("finished getting updates")
}
