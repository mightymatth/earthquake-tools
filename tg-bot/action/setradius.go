package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
	"strconv"
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
If <code>0</code> (zero) value is provided, it will subscribe to the entire world.
e.g.: <code>100.5</code>, <code>350</code>
`
}

func (a SetRadiusAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionAction(sub.SubID, ResetInput).Encode())

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

func (a SetRadiusAction) ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service) {
	radius, err := a.processInputValue(update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, a.WrongInput())
		_, _ = bot.Send(msg)
		ShowSetRadius(update.Message.Chat.ID, a.Params.P1, bot, s)
		return
	}

	radiusUpdate := entity.SubscriptionUpdate{Radius: radius}
	_, err = storage.Service.UpdateSubscription(s, a.Params.P1, &radiusUpdate)
	if err != nil {
		log.Printf("cannot set radius to subscription: %v", err)
		return
	}

	_ = ResetAwaitInput(ResetInput, update.Message.Chat.ID, s)
	ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
}

func (a SetRadiusAction) processInputValue(text string) (float64, error) {
	radius, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse text to number")
	}

	if !(radius == 0 || radius >= 1 && radius <= 2000) {
		return 0, fmt.Errorf("invalid range")
	}

	return radius, nil
}

func (a SetRadiusAction) WrongInput() string {
	return "Wrong input. A whole or decimal number expected; in range [0] ∪ [1, 2000]"
}
