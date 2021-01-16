package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type DelayScreen Screen
type EditDelayScreen Screen

func (scr DelayScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(), scr.inlineButtons())
	bot.Send(message)
}

func (scr DelayScreen) text() string {
	return `
Current delay set to: X

Data from earthquake data sources may arrive with significant delays; sometimes for a few hours. 
If you set a delay to 5 minutes, you will only receive the events that arrived late up to 5 minutes.
`
}

func (scr DelayScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	delay := tgbotapi.NewInlineKeyboardButtonData("Edit Delay", "")
	settings := tgbotapi.NewInlineKeyboardButtonData("« Subscription", "")

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(delay),
		tgbotapi.NewInlineKeyboardRow(settings),
	)
	return &kb
}

func (scr EditDelayScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	// TODO: edit message text and set keyboard (not inline)
	//message := editedMessageConfig(msg.Chat.ID, msg.MessageID, m.text(), m.inlineButtons())
	//bot.Send(message)
}
