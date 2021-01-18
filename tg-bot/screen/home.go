package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type HomeScreen struct {
	Screen
}

const Home Cmd = "HOME"

func NewHomeScreen() HomeScreen {
	return HomeScreen{Screen{Cmd: Home}}
}

func (scr HomeScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(), scr.inlineButtons())
	bot.Send(message)
}

func (scr HomeScreen) text() string {
	return `
Welcome and welcome,

You are the best!
`
}

func (scr HomeScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	subs := tgbotapi.NewInlineKeyboardButtonData("Subscriptions",
		NewSubscriptionsScreen("").Encode())

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
	msg.ReplyMarkup = HomeScreen{}.inlineButtons()

	_, _ = bot.Send(msg)
}

func UnknownCommandText() string {
	return `
Unknown command.
`
}

func ShowHome(bot *tgbotapi.BotAPI, chatID int64) {
	home := HomeScreen{}
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
