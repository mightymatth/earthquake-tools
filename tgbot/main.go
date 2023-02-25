package main

import (
	"context"
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/aggregator"
	"github.com/mightymatth/earthquake-tools/tgbot/storage"
	"github.com/mightymatth/earthquake-tools/tgbot/storage/mongo"
	"log"
	"time"
)

// TODO: Use cobra for CLI params parsing
var (
	tgBotToken = flag.String("tg-bot-token", "", "telegram bot token, required")
	mongoURI   = flag.String("mongo-uri", "", "mongo URI, required")

	dbNamespace = flag.String("db-namespace", "", "database namespace, optional")
	tgBotDebug  = flag.String("tg-bot-debug", "", "run telegram bot in debug mode (set true), optional")
)

func main() {
	flag.Parse()
	_ = loadDotEnv()

	ctx := context.Background()
	botAPI, err := tgbotapi.NewBotAPI(getParam(*tgBotToken, "TELEGRAM_BOT_TOKEN", ""))
	if err != nil {
		log.Panic("cannot initialize bot api:", err)
	}

	storageImpl, err := mongo.NewStorage(
		getParam(*mongoURI, "MONGO_URI", "mongodb://localhost:27017"),
		getParam(*dbNamespace, "DB_NAMESPACE", ""),
	)
	if err != nil {
		log.Panicf("cannot create storage: %v", err)
	}
	defer storageImpl.Client.Disconnect(ctx)

	service := storage.NewService(storageImpl)

	// TODO: set this behind a flag
	//http.HandleFunc("/", EqEventHandler(botAPI, service))
	//go func() {
	//	log.Panic(http.ListenAndServe(":3300", nil))
	//}()

	events := aggregator.Start(ctx)
	go EqEventSender(events, botAPI, service)

	log.Panic(TgBotServer(botAPI, service))
}

func TgBotServer(bot *tgbotapi.BotAPI, s storage.Service) error {
	if getParam(*tgBotDebug, "TG_BOT_DEBUG", "") != "" {
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
