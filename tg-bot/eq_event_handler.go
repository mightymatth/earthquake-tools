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
			return
		}
		defer r.Body.Close()

		w.WriteHeader(http.StatusAccepted)

		go sendEventReport(event, bot, s)
	}
}

func sendEventReport(event EarthquakeEvent, bot *tgbotapi.BotAPI, s storage.Service) {
	eventData := entity.EventData{
		Magnitude: event.Mag,
		Delay:     time.Now().Sub(event.Time).Minutes(),
		Location: entity.Location{
			Lat: event.Lat,
			Lng: event.Lon,
		},
		Source: event.SourceID,
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
	detailsURL := tgbotapi.NewInlineKeyboardButtonURL("Details 📰", event.DetailsURL)
	mapsURL := tgbotapi.NewInlineKeyboardButtonURL("Location 📍",
		fmt.Sprintf("https://www.google.com/maps/place/%f,%f", event.Lat, event.Lon),
	)

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(detailsURL, mapsURL),
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

type eventReport struct {
	EarthquakeEvent
	TimeNow time.Time
}

func (e eventReport) String() string {
	return fmt.Sprintf(`
💥 <b>NEW EARTHQUAKE</b> 💥
📶 Magnitude: <b>%.1f</b> %s
🔻 Depth: %.0f km
📍 Location: %s
⏳ Relative time: %s
⏱ UTC Time: <code>%s</code>
⏰ Local Time: <code>%s</code>
🏣 Source: <code>%s/%s</code>
			`,
		e.Mag,
		e.MagType,
		e.Depth,
		e.Location,
		humanize.RelTime(e.Time, e.TimeNow, "ago", "later"),
		e.Time.Format("Mon, 2 Jan 2006 15:04:05 MST"),
		getLocationTime(e.Time, e.Lat, e.Lon),
		e.SourceID, e.EventID,
	)
}
