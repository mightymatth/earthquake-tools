package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type DelayScreen Screen
type EditDelayScreen Screen

const (
	Delay     DelayScreen     = "DELAY"
	EditDelay EditDelayScreen = "EDIT_DELAY"
)

func (s DelayScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, s.text(), s.inlineButtons())
	bot.Send(message)
}

func (s DelayScreen) text() string {
	return `
Current delay set to: X

Data from earthquake data sources may arrive with significant delays; sometimes for a few hours. 
If you set a delay to 5 minutes, you will only receive the events that arrived late up to 5 minutes.
`
}

func (s DelayScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	delay := tgbotapi.NewInlineKeyboardButtonData("Edit Delay", string(EditDelay))
	settings := tgbotapi.NewInlineKeyboardButtonData("« Settings", string(Settings))

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(delay),
		tgbotapi.NewInlineKeyboardRow(settings),
	)
	return &kb
}

func (s EditDelayScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	// TODO: edit message text and set keyboard (not inline)
	//message := editedMessageConfig(msg.Chat.ID, msg.MessageID, m.text(), m.inlineButtons())
	//bot.Send(message)
}
