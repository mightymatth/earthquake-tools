package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetMagnitudeAction struct {
	Action
}

const SetMagnitude Cmd = "SET_MAG"

func NewSetMagnitudeAction(subID string) SetMagnitudeAction {
	return SetMagnitudeAction{Action{Cmd: SetMagnitude, Params: Params{P1: subID}}}
}

func (a SetMagnitudeAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
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

func (a SetMagnitudeAction) text() string {
	return `
Enter minimum magnitude to receive an alert.
e.g.: <code>4.3</code>, <code>5</code>
`
}

func (a SetMagnitudeAction) WrongInput() string {
	return "Wrong input. Integer or decimal number expected."
}

func (a SetMagnitudeAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscription(sub.SubID, ResetInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(cancel),
	)
	return &kb
}

func ShowSetMagnitude(chatID int64, subID string, bot *tgbotapi.BotAPI, s storage.Service) {
	sub, err := s.GetSubscription(subID)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	setMagAction := NewSetMagnitudeAction(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setMagAction.inlineButtons(sub),
		},
		Text: setMagAction.text(),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}
