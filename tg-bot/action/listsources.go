package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
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
			Name:        "EMSC",
			Description: "EMSC data feed.",
			SourceID:    entity.EMSC,
		},
		{
			Name:        "EMSC/ws",
			Description: "EMSC data feed via WebSocket.",
			SourceID:    entity.EMSCWS,
		},
		{
			Name:        "USGS",
			Description: "USGS data feed.",
			SourceID:    entity.USGS,
		},
		{
			Name:        "IRIS",
			Description: "IRIS data feed.",
			SourceID:    entity.IRIS,
		},
		{
			Name:        "USPBR",
			Description: "USPBR data feed.",
			SourceID:    entity.USPBR,
		},
		{
			Name:        "GEOFON",
			Description: "GEOFON data feed.",
			SourceID:    entity.GEOFON,
		},
	}

	sourcesM = make(map[string]Source)
	for _, src := range sources {
		sourcesM[string(src.SourceID)] = src
	}
}
