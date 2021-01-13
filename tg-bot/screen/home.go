package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HomeScreen Screen

const Home HomeScreen = "HOME"

func (s HomeScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, s.text(), s.inlineButtons())
	bot.Send(message)
}

func (s HomeScreen) text() string {
	return `
Welcome and welcome,

You are the best!
`
}

func (s HomeScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	settings := tgbotapi.NewInlineKeyboardButtonData("Settings", fmt.Sprintf("%s", Settings))

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(settings),
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
	msg.ReplyMarkup = Home.inlineButtons()

	_, _ = bot.Send(msg)
}

func UnknownCommandText() string {
	return `
Unknown command.
`
}

func ShowHome(bot *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Text:                  Home.text(),
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}
	msg.ReplyMarkup = Home.inlineButtons()

	_, _ = bot.Send(msg)
}
