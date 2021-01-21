package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type SubscriptionScreen struct {
	Screen
}

const Sub Cmd = "SUB"

func NewSubscriptionScreen(subID string, reset ResetInputType) SubscriptionScreen {
	return SubscriptionScreen{Screen{
		Cmd:    Sub,
		Params: Params{P1: subID, P2: string(reset)},
	}}
}

func (scr SubscriptionScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(s), scr.inlineButtons())
	bot.Send(message)
}

func (scr SubscriptionScreen) text(s storage.Service) string {
	sub, err := s.GetSubscription(scr.Params.P1)
	if err != nil {
		fmt.Printf("cannot get subscription: %v", err)
		return ""
	}

	return fmt.Sprintf(`
Current subscription settings:
Name: %s
Minimum magnitude: %.1f
Earthquake location: %s
My location: %s
Radius: %.1f km
Time offset: %d s
`, sub.Name, sub.MinMag, sub.EqLocations,
		sub.MyLocation, sub.Radius, sub.OffsetSec)
}

func (scr SubscriptionScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	magnitude := tgbotapi.NewInlineKeyboardButtonData("Magnitude", " ")
	delay := tgbotapi.NewInlineKeyboardButtonData("Delay", " ")
	home := tgbotapi.NewInlineKeyboardButtonData("« Subscriptions",
		NewSubscriptionsScreen("").Encode())
	deleteSub := tgbotapi.NewInlineKeyboardButtonData("Delete",
		NewDeleteSubscriptionScreen(scr.Params.P1, "").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(magnitude, delay),
		tgbotapi.NewInlineKeyboardRow(home, deleteSub),
	)
	return &kb
}
