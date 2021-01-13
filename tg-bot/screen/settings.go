package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SettingsScreen Screen

const Settings SettingsScreen = "SETTINGS"

func (s SettingsScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, s.text(), s.inlineButtons())
	bot.Send(message)
}

func (s SettingsScreen) text() string {
	return `
Here are the settings for modifying subscription for earthquake events.

You can filter out earthquakes by properties such as minimum magnitude, your location/range, etc.
`
}

func (s SettingsScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	magnitude := tgbotapi.NewInlineKeyboardButtonData("Magnitude", string(Magnitude))
	delay := tgbotapi.NewInlineKeyboardButtonData("Delay", string(Delay))
	home := tgbotapi.NewInlineKeyboardButtonData("« Home", string(Home))

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(magnitude, delay),
		tgbotapi.NewInlineKeyboardRow(home),
	)
	return &kb
}
