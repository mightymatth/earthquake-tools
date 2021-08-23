package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type SubscriptionsAction struct {
	Action
}

const Subs Cmd = "SUBS"

func NewSubscriptionsAction(reset ResetInputType) SubscriptionsAction {
	return SubscriptionsAction{Action{
		Cmd:    Subs,
		Params: Params{P1: string(reset)},
	}}
}

func (a SubscriptionsAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	_ = ResetAwaitInput(ResetInputType(a.Params.P1), msg.Chat.ID, s)

	subs := s.GetSubscriptions(msg.Chat.ID)
	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(), a.inlineButtons(subs))
	_, _ = bot.Send(message)
}

func (a SubscriptionsAction) text() string {
	return `
<b>Subscriptions</b>

Here is the list of your active subscriptions.
Click on any subscription to edit it or create a new one.
`
}

var backToHomeButton = tgbotapi.NewInlineKeyboardButtonData("« Home", NewHomeAction().Encode())
var newSubscriptionButton = tgbotapi.NewInlineKeyboardButtonData("＋ New", NewCreateSubscriptionAction().Encode())

func (a SubscriptionsAction) inlineButtons(
	subs []entity.Subscription,
) *tgbotapi.InlineKeyboardMarkup {
	var staticRow = make([]tgbotapi.InlineKeyboardButton, 0, 2)
	staticRow = append(staticRow, backToHomeButton)

	if len(subs) < 10 {
		staticRow = append(staticRow, newSubscriptionButton)
	}

	rows := append(
		a.subscriptionRows(subs, 2),
		tgbotapi.NewInlineKeyboardRow(staticRow...),
	)
	kb := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &kb
}

func (a SubscriptionsAction) subscriptionRows(
	subs []entity.Subscription, iPerRow int,
) (rows [][]tgbotapi.InlineKeyboardButton) {
	tmpRow := make([]tgbotapi.InlineKeyboardButton, 0, iPerRow)
	for i := 0; i < len(subs); i++ {
		btn := tgbotapi.NewInlineKeyboardButtonData(subs[i].Name,
			NewSubscriptionAction(subs[i].SubID, "").Encode())
		tmpRow = append(tmpRow, btn)
		if len(tmpRow) == iPerRow {
			rows = append(rows, tmpRow)
			tmpRow = make([]tgbotapi.InlineKeyboardButton, 0, iPerRow)
		}
	}

	if len(tmpRow) > 0 {
		rows = append(rows, tmpRow)
	}

	return rows
}

func ShowSubscriptions(chatID int64, bot *tgbotapi.BotAPI, s storage.Service) {
	subs := s.GetSubscriptions(chatID)
	subsAction := SubscriptionsAction{}

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: subsAction.inlineButtons(subs),
		},
		Text: subsAction.text(),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}
