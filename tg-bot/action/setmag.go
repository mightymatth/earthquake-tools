package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
	"strconv"
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

func (a SetMagnitudeAction) ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service) {
	mag, err := a.processInputValue(update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, a.WrongInput())
		_, _ = bot.Send(msg)
		ShowSetMagnitude(update.Message.Chat.ID, a.Params.P1, bot, s)
		return
	}

	magUpdate := entity.SubscriptionUpdate{MinMag: mag}
	_, err = storage.Service.UpdateSubscription(s, a.Params.P1, &magUpdate)
	if err != nil {
		log.Printf("cannot set magnitude to subscription: %v", err)
		return
	}

	_ = ResetAwaitInput(ResetInput, update.Message.Chat.ID, s)
	ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
}

func (a SetMagnitudeAction) processInputValue(text string) (float64, error) {
	mag, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse text to number")
	}

	switch {
	case mag < 0.1, mag > 10:
		return 0, fmt.Errorf("invalid range")
	default:
		return mag, nil
	}
}

func (a SetMagnitudeAction) WrongInput() string {
	return "Wrong input. A whole or decimal number expected; in range [0.1, 10]"
}
