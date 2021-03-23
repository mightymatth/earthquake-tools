package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetRadiusAction struct {
	Action
}

const SetRadius Cmd = "SET_RADIUS"

func NewSetRadiusAction(subID string) SetRadiusAction {
	return SetRadiusAction{Action{Cmd: SetRadius, Params: Params{P1: subID}}}
}

func (a SetRadiusAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	err := s.SetAwaitUserInput(msg.Chat.ID, a.Encode())
	if err != nil {
		log.Printf("cannot set chat state: %v", err)
		return
	}

	sub, err := s.GetSubscription(a.Params.P1)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(), a.inlineButtons(sub))
	bot.Send(message)
}

func (a SetRadiusAction) text() string {
	return `
Enter maximum earthquake location radius to receive an alert.
The unit is kilometer (km).
e.g.: <code>100.5</code>, <code>350</code>
`
}

func (a SetRadiusAction) WrongInput() string {
	return "Wrong input. Integer or decimal number expected."
}

func (a SetRadiusAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscription(sub.SubID, ResetInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(cancel),
	)
	return &kb
}

func ShowSetRadius(chatID int64, subID string, bot *tgbotapi.BotAPI, s storage.Service) {
	sub, err := s.GetSubscription(subID)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	setRadiusAction := NewSetRadiusAction(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setRadiusAction.inlineButtons(sub),
		},
		Text: setRadiusAction.text(),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}
