package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type SubscriptionsScreen struct {
	Screen
}

const Subs Cmd = "SUBS"

func NewSubscriptionsScreen(reset ResetInputType) SubscriptionsScreen {
	return SubscriptionsScreen{Screen{
		Cmd:    Subs,
		Params: Params{P1: string(reset)},
	}}
}

func (scr SubscriptionsScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	ResetAwaitInput(ResetInputType(scr.Params.P1), msg.Chat.ID, s)

	subs := s.GetSubscriptions(msg.Chat.ID)
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(), scr.inlineButtons(subs))
	bot.Send(message)
}

func (scr SubscriptionsScreen) text() string {
	return `
Here is the list of your active subscriptions.
Click on any subscription to edit it or create a new one.
`
}

func (scr SubscriptionsScreen) inlineButtons(
	subs []entity.Subscription,
) *tgbotapi.InlineKeyboardMarkup {
	home := tgbotapi.NewInlineKeyboardButtonData("« Home", NewHomeScreen().Encode())
	newSub := tgbotapi.NewInlineKeyboardButtonData("+ New", NewCreateSubscriptionScreen().Encode())

	rows := append(
		scr.subscriptionRows(subs, 2),
		tgbotapi.NewInlineKeyboardRow(home, newSub),
	)
	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &kb
}

func (scr SubscriptionsScreen) subscriptionRows(
	subs []entity.Subscription, iPerRow int,
) (rows [][]tgbotapi.InlineKeyboardButton) {
	tmpRow := make([]tgbotapi.InlineKeyboardButton, 0, iPerRow)
	for i := 0; i < len(subs); i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(subs[i].Name, NewHomeScreen().Encode())
		tmpRow = append(tmpRow, btn)
		if len(tmpRow) == 2 {
			rows = append(rows, tmpRow)
			tmpRow = make([]tgbotapi.InlineKeyboardButton, 0, iPerRow)
		}
	}

	if len(tmpRow) > 0 {
		rows = append(rows, tmpRow)
	}

	return rows
}
