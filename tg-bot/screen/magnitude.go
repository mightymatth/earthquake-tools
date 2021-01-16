package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type MagnitudeScreen Screen
type EditMagnitudeScreen Screen


func (scr MagnitudeScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(), scr.inlineButtons())
	bot.Send(message)
}

func (scr MagnitudeScreen) text() string {
	return `
You have set your minimum magnitude level to X.
`
}

func (scr MagnitudeScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	magnitude := tgbotapi.NewInlineKeyboardButtonData(
		"Edit Magnitude", "")
	settings := tgbotapi.NewInlineKeyboardButtonData(
		"« Subscription", "")

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(magnitude),
		tgbotapi.NewInlineKeyboardRow(settings),
	)
	return &kb
}

func (scr EditMagnitudeScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	// TODO: edit message text and set keyboard (not inline)
	//message := editedMessageConfig(msg.Chat.ID, msg.MessageID, m.text(), m.inlineButtons())
	//bot.Send(message)
}
