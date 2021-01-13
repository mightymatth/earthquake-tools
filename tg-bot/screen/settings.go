package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const Settings Screen = "SETTINGS"

func SettingsButtons() tgbotapi.InlineKeyboardMarkup {
	home := tgbotapi.NewInlineKeyboardButtonData("≪ Home", fmt.Sprintf("%s", Home))
	magnitude := tgbotapi.NewInlineKeyboardButtonData("Magnitude", fmt.Sprintf("%s", Magnitude))

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(home, magnitude),
	)
}

const Magnitude Screen = "MAGNITUDE"
