package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tgbot/entity"
	"github.com/mightymatth/earthquake-tools/tgbot/storage"
	"log"
)

var sources []Source
var sourcesM map[string]Source

func init() {
	createSources()
}

type ListSourcesAction struct {
	Action
}

const ListSources Cmd = "LIST_SRCS"

func NewListSourcesAction(subID string) ListSourcesAction {
	return ListSourcesAction{Action{
		Cmd:    ListSources,
		Params: Params{P1: subID},
	}}
}

func (a ListSourcesAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	sub, err := s.GetSubscription(a.Params.P1)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(sub), a.inlineButtons(sub))
	_, _ = bot.Send(message)
}

func (a ListSourcesAction) text(sub *entity.Subscription) string {
	return fmt.Sprintf(`
<i>Subscription sources for</i> <b>%s</b>
Active sources on which you are subscribed to are marked with a green checkmark.
`, sub.Name)
}

func (a ListSourcesAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	back := tgbotapi.NewInlineKeyboardButtonData("« Subscription",
		NewSubscriptionAction(a.Action.Params.P1, "").Encode())

	rows := append(
		a.sourcesRows(sub, 2),
		tgbotapi.NewInlineKeyboardRow(back),
	)

	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &kb
}

func (a ListSourcesAction) sourcesRows(
	sub *entity.Subscription, iPerRow int,
) (rows [][]tgbotapi.InlineKeyboardButton) {
	var activeM = make(map[string]Source)
	for _, srcID := range sub.Sources {
		if src, found := sourcesM[string(srcID)]; found {
			activeM[string(srcID)] = src
		}
	}

	tmpRow := make([]tgbotapi.InlineKeyboardButton, 0, iPerRow)

	for _, src := range sources {
		var active string
		if _, found := activeM[string(src.SourceID)]; found {
			active = "✅"
		}

		btn := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", src.Name, active),
			NewSetSourceAction(sub.SubID, string(src.SourceID), "").Encode())
		tmpRow = append(tmpRow, btn)
		if len(tmpRow) == iPerRow {
			rows = append(rows, tmpRow)
			tmpRow = make([]tgbotapi.InlineKeyboardButton, 0, iPerRow)
		}
	}

	if len(tmpRow) > 0 {
		rows = append(rows, tmpRow)
	}

	return rows
}

type Source struct {
	Name        string
	Description string
	SourceID    entity.SourceID
}

func createSources() {
	sources = []Source{
		{
			Name: "EMSC",
			Description: `
<a href="https://www.emsc-csem.org/">EMSC</a> data <a href="https://www.seismicportal.eu/fdsn-wsevent.html">feed</a> based on <a href="https://www.fdsn.org/webservices/">FDSN Web Service</a> event specification.
The service reports almost all earthquakes from Europe and Mediterranean region, and also more significant ones from the entire world.
`,
			SourceID: entity.EMSC,
		},
		{
			Name: "EMSC/ws",
			Description: `
<a href="https://www.emsc-csem.org/">EMSC</a> data feed connected to <a href="https://www.seismicportal.eu/realtime.html">(near) Realtime Notification using Websocket</a>.
The service reports almost all earthquakes from Europe and Mediterranean region, and also more significant ones from the entire world. 
The reports are almost the same as the EMSC service based on FDSN event service. The difference is that it's theoretically a few seconds faster, but it has a little bit less availability under heavy load when bigger earthquakes occur. 
`,
			SourceID: entity.EMSCWS,
		},
		{
			Name: "USGS",
			Description: `
<a href="https://earthquake.usgs.gov/">USGS</a> data <a href="https://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_day.geojson">feed</a>.
The service reports almost all earthquakes from North America. It is reliable and fast source, so is highly recommended to use.
`,
			SourceID: entity.USGS,
		},
		{
			Name: "IRIS",
			Description: `
<a href="https://www.iris.edu/hq/">IRIS</a> data <a href="https://service.iris.edu/fdsnws/">feed</a> based on <a href="https://www.fdsn.org/webservices/">FDSN Web Service</a> event specification.
The service reports almost all earthquakes from North America. It is reliable and fast source, so is highly recommended to use. The reports have very similar content and response time as USGS source.
`,
			SourceID: entity.IRIS,
		},
		{
			Name: "USPBR",
			Description: `
<a href="http://www.moho.iag.usp.br/">Centro de sismologia USP</a> data <a href="http://www.moho.iag.usp.br/rq/">feed</a> based on <a href="https://www.fdsn.org/webservices/">FDSN Web Service</a> event specification.
The service reports earthquakes mainly from South America. 
`,
			SourceID: entity.USPBR,
		},
		{
			Name: "GEOFON",
			Description: `
<a href="https://geofon.gfz-potsdam.de/">GEOFON</a> data <a href="http://geofon.gfz-potsdam.de/fdsnws/">feed</a> based on <a href="https://www.fdsn.org/webservices/">FDSN Web Service</a> event specification.
The service reports only significant earthquakes from the entire world. 
`,
			SourceID: entity.GEOFON,
		},
	}

	sourcesM = make(map[string]Source)
	for _, src := range sources {
		sourcesM[string(src.SourceID)] = src
	}
}
