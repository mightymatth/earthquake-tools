package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const Home Screen = "HOME"

func ShowHome(bot *tgbotapi.BotAPI, chatID int64) {
	text := `
Welcome and welcome,

You are the best!
`

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: chatID,
		},
		Text:                  text,
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}
	msg.ReplyMarkup = HomeButtons()

	_, _ = bot.Send(msg)
}

func HomeButtons() tgbotapi.InlineKeyboardMarkup {
	settings := tgbotapi.NewInlineKeyboardButtonData("Settings", fmt.Sprintf("%s", Settings))

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(settings),
	)
}
