package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
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

	http.HandleFunc("/", EarthquakeEventServer(bot))
	go http.ListenAndServe(":3300", nil)

	TgBotServer(bot)
}

func EarthquakeEventServer(bot *tgbotapi.BotAPI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := ParseEvent(r.Body)
		if err != nil {
			log.Printf("process event: %v", err)
		}

		if event.Action != "create" {
			return
		}

		chatID := 307010667

		text := fmt.Sprintf(
			`
📶 Magnitude: <b>%.1f</b> %s
🔻 Depth: %.0f km
📍 Location: %s
⏳ Relative time: %s
⏱ UTC Time: <code>%s</code>
⏰ Local Time: <code>%s</code>
🏣 Source/ID: %s
			`,
			event.Data.Properties.Mag,
			event.Data.Properties.MagType,
			event.Data.Properties.Depth,
			event.Data.Properties.FlynnRegion,
			humanize.RelTime(event.Data.Properties.Time, time.Now(), "ago", "later"),
			event.Data.Properties.Time.Format("Mon, 2 Jan 2006 15:04:05 MST"),
			getLocationTime(
				event.Data.Properties.Time,
				event.Data.Properties.Lat,
				event.Data.Properties.Lon,
			),
			SourceLinkHTML(event.Data.Properties.SourceCatalog, event.Data.Properties.SourceID),
		)

		msg := tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID: int64(chatID),
			},
			Text:                  text,
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		}
		msg.ReplyMarkup = EventButtons(event)

		_, _ = bot.Send(msg)
	}
}

func getLocationTime(timeUTC time.Time, lat, lon float64) string {
	localTime, err := LocationTime(timeUTC, lat, lon)
	if err != nil {
		fmt.Printf("cannot get location time: %v", err)
		return "unknown"
	}

	return localTime.Format("Mon, 2 Jan 2006 15:04:05 MST")
}

func EventButtons(event EarthquakeEvent) tgbotapi.InlineKeyboardMarkup {
	detailsURL := tgbotapi.NewInlineKeyboardButtonURL("Details & Updates",
		fmt.Sprintf("https://www.seismicportal.eu/eventdetails.html?unid=%s", event.Data.ID),
	)
	mapsURL := tgbotapi.NewInlineKeyboardButtonURL("Location 📍",
		fmt.Sprintf("http://www.google.com/maps/place/%f,%f",
			event.Data.Properties.Lat,
			event.Data.Properties.Lon,
		),
	)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(detailsURL),
		tgbotapi.NewInlineKeyboardRow(mapsURL),
	)
}

func TgBotServer(bot *tgbotapi.BotAPI) {
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		// TODO: Implement callback queries when needed.
		//if update.CallbackQuery != nil{
		//	fmt.Print(update)
		//	bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID,update.CallbackQuery.Data))
		//
		//	bot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,update.CallbackQuery.Data))
		//}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// TODO: Implement user message responds
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID
		//
		//bot.Send(msg)
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
