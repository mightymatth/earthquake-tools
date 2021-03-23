package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type HomeAction struct {
	Action
}

const Home Cmd = "HOME"

func NewHomeAction() HomeAction {
	return HomeAction{Action{Cmd: Home}}
}

func (a HomeAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(), a.inlineButtons())
	_, _ = bot.Send(message)
}

func (a HomeAction) text() string {
	return `
Welcome to <i>EMSC Events ⚠️</i> Bot

It makes you able to receive notifications of recent earthquakes by making subscriptions configured with parameters such as magnitude, location, and observing radius.

Developed by @mpevec
Source code: <a href="https://github.com/mightymatth/earthquake-tools">GitHub</a>
`
}

func (a HomeAction) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	subs := tgbotapi.NewInlineKeyboardButtonData("Subscriptions",
		NewSubscriptionsAction("").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(subs),
	)
	return &kb
}

func ShowUnknownCommand(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Text:                  UnknownCommandText(),
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}
	msg.ReplyMarkup = HomeAction{}.inlineButtons()

	_, _ = bot.Send(msg)
}

func UnknownCommandText() string {
	return `
Unknown command.
`
}

func ShowHome(bot *tgbotapi.BotAPI, chatID int64) {
	home := HomeAction{}
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Text:                  home.text(),
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}
	msg.ReplyMarkup = home.inlineButtons()

	_, _ = bot.Send(msg)
}
