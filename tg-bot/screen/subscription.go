package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type SubscriptionScreen struct {
	Screen
}

const Sub Cmd = "SUB"

func NewSubscriptionScreen(subID string, reset ResetInputType) SubscriptionScreen {
	return SubscriptionScreen{Screen{
		Cmd:    Subs,
		Params: Params{P1: subID, P2: string(reset)},
	}}
}

func (scr SubscriptionScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(), scr.inlineButtons())
	bot.Send(message)
}

func (scr SubscriptionScreen) text() string {
	return `
Here are the settings for modifying subscription for earthquake events.

You can filter out earthquakes by properties such as minimum magnitude, your location/range, etc.
`
}

func (scr SubscriptionScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	magnitude := tgbotapi.NewInlineKeyboardButtonData("Magnitude", " ")
	delay := tgbotapi.NewInlineKeyboardButtonData("Delay", " ")
	home := tgbotapi.NewInlineKeyboardButtonData("« Subscriptions",
		NewSubscriptionsScreen("").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(magnitude, delay),
		tgbotapi.NewInlineKeyboardRow(home),
	)
	return &kb
}
