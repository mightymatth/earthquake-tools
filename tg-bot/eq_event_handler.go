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
			log.Printf("processing event failed: %v", err)
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

		if len(chatIDs) == 0 {
			return
		}

		eventReport := eventReport{event, time.Now()}

		broadcast(eventReport, chatIDs, bot)
	}
}

func broadcast(eventReport eventReport, chatIDs []int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.MessageConfig{
		Text:                  eventReport.String(),
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}
	msg.ReplyMarkup = eventButtons(eventReport.EarthquakeEvent)

	// sleepPeriod limits the rate of sending messages. Current Telegram limit
	// is 30 messages per second, so sleepPeriod (which is equal to 1/rate) should be
	// somewhat higher than 50 ms keep handing user interactions smoothly.
	sleepPeriod := 50 * time.Millisecond
	retryAttempts := 3

	for _, chatID := range chatIDs {
		msg.BaseChat.ChatID = chatID

		for i := 1; i <= retryAttempts; i++ {
			_, err := bot.Send(msg)

			if err != nil {
				if i == retryAttempts {
					log.Printf("error sending event after %d retr{y,ies}: %v", retryAttempts, err)
					break
				}
				time.Sleep(sleepPeriod)

				continue
			}

			break
		}

		time.Sleep(sleepPeriod)
	}
}

func eventButtons(event EarthquakeEvent) tgbotapi.InlineKeyboardMarkup {
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

func getLocationTime(timeUTC time.Time, lat, lon float64) string {
	localTime, err := LocationTime(timeUTC, lat, lon)
	if err != nil {
		log.Printf("cannot get location time: %v", err)
		return "unknown"
	}

	return localTime.Format("Mon, 2 Jan 2006 15:04:05 MST")
}

func sourceLinkHTML(sourceType, ID string) string {
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

type eventReport struct {
	EarthquakeEvent
	TimeNow time.Time
}

func (e eventReport) String() string {
	return fmt.Sprintf(
		`
📶 Magnitude: <b>%.1f</b> %s
🔻 Depth: %.0f km
📍 Location: %s
⏳ Relative time: %s
⏱ UTC Time: <code>%s</code>
⏰ Local Time: <code>%s</code>
🏣 Source/ID: %s
			`,
		e.Data.Properties.Mag,
		e.Data.Properties.MagType,
		e.Data.Properties.Depth,
		e.Data.Properties.FlynnRegion,
		humanize.RelTime(e.Data.Properties.Time, e.TimeNow, "ago", "later"),
		e.Data.Properties.Time.Format("Mon, 2 Jan 2006 15:04:05 MST"),
		getLocationTime(
			e.Data.Properties.Time,
			e.Data.Properties.Lat,
			e.Data.Properties.Lon,
		),
		sourceLinkHTML(e.Data.Properties.SourceCatalog, e.Data.Properties.SourceID),
	)
}
