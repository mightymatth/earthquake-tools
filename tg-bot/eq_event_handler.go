package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
	"net/http"
	"time"
)

func EqEventHandler(bot *tgbotapi.BotAPI, s storage.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		event, err := ParseEvent(r.Body)
		if err != nil {
			log.Printf("process event: %v\n", err)
		}

		if event.Action != "create" {
			return
		}

		eventData := entity.EventData{
			Magnitude: event.Data.Properties.Mag,
			Delay:     time.Now().Sub(event.Data.Properties.Time).Minutes(),
			Location: entity.Location{
				Lat: event.Data.Properties.Lat,
				Lng: event.Data.Properties.Lon,
			},
		}

		chatIDs, err := s.GetEventSubscribers(eventData)
		if err != nil {
			log.Printf("cannot get event subscribers: %v\n", err)
		}

		if len(chatIDs) < 1 {
			return
		}

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
			Text:                  text,
			ParseMode:             tgbotapi.ModeHTML,
			DisableWebPagePreview: true,
		}
		msg.ReplyMarkup = EventButtons(event)

		for _, chatID := range chatIDs {
			msg.BaseChat.ChatID = chatID
			_, _ = bot.Send(msg)
		}
	}
}

func EventButtons(event EarthquakeEvent) tgbotapi.InlineKeyboardMarkup {
	detailsURL := tgbotapi.NewInlineKeyboardButtonURL("Details & Updates",
		fmt.Sprintf("https://www.seismicportal.eu/eventdetails.html?unid=%s", event.Data.ID),
	)
	mapsURL := tgbotapi.NewInlineKeyboardButtonURL("Location 📍",
		fmt.Sprintf("http://www.google.com/maps/place/%f,%f",
			event.Data.Properties.Lat,
			event.Data.Properties.Lon),
	)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(detailsURL),
		tgbotapi.NewInlineKeyboardRow(mapsURL),
	)
}

func SourceLinkHTML(sourceType, ID string) string {
	switch SourceType(sourceType) {
	case emsc:
		return fmt.Sprintf(
			`<a href="https://www.emsc-csem.org/Earthquake/earthquake.php?id=%s">%s/%s</a>`,
			ID, sourceType, ID)
	default:
		return fmt.Sprintf(`<code>%s/%s</code>`, sourceType, ID)
	}
}

type SourceType string

const (
	emsc SourceType = "EMSC-RTS"
)
