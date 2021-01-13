package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type MagnitudeScreen Screen
type EditMagnitudeScreen Screen

const (
	Magnitude     MagnitudeScreen     = "MAGNITUDE"
	EditMagnitude EditMagnitudeScreen = "EDIT_MAGNITUDE"
)

func (s MagnitudeScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, s.text(), s.inlineButtons())
	bot.Send(message)
}

func (s MagnitudeScreen) text() string {
	return `
You have set your minimum magnitude level to X.
`
}

func (s MagnitudeScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	magnitude := tgbotapi.NewInlineKeyboardButtonData("Edit Magnitude", string(EditMagnitude))
	settings := tgbotapi.NewInlineKeyboardButtonData("« Settings", string(Settings))

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(magnitude),
		tgbotapi.NewInlineKeyboardRow(settings),
	)
	return &kb
}

func (s EditMagnitudeScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	// TODO: edit message text and set keyboard (not inline)
	//message := editedMessageConfig(msg.Chat.ID, msg.MessageID, m.text(), m.inlineButtons())
	//bot.Send(message)
}
