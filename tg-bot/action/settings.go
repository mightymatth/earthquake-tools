package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type SettingsAction struct {
	Action
}

const Settings Cmd = "SETTINGS"

func NewSettingsAction(reset ResetInputType) SettingsAction {
	return SettingsAction{Action{
		Cmd:    Settings,
		Params: Params{P1: string(reset)},
	}}
}

func (a SettingsAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(), a.inlineButtons())
	_, _ = bot.Send(message)
}

func (a SettingsAction) text() string {
	return `
<b>Settings</b>

`
}

func (a SettingsAction) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	disableInput := tgbotapi.NewInlineKeyboardButtonData("🚫 Disable input", NewDisableInputAction("").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(disableInput),
		tgbotapi.NewInlineKeyboardRow(backToHomeButton),
	)
	return &kb
}
