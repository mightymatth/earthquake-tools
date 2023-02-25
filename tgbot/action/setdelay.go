package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tgbot/entity"
	"github.com/mightymatth/earthquake-tools/tgbot/storage"
	"log"
	"strconv"
)

type SetDelayAction struct {
	Action
}

const SetDelay Cmd = "SET_DELAY"

func NewSetDelayAction(subID string) SetDelayAction {
	return SetDelayAction{Action{Cmd: SetDelay, Params: Params{P1: subID}}}
}

func (a SetDelayAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
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

func (a SetDelayAction) text() string {
	return `
<b>Subscription ᐅ Delay</b>

<b>Delay</b> is a time period between the earthquake occasion and the moment when the report is received. 
It is expected behavior as data may arrive from various sources.

Enter how many <u>minutes</u> you will tolerate.

e.g.: <code>2</code>, <code>60</code>
`
}

func (a SetDelayAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionAction(sub.SubID, ResetInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(cancel),
	)
	return &kb
}

func ShowSetDelay(chatID int64, subID string, bot *tgbotapi.BotAPI, s storage.Service) {
	sub, err := s.GetSubscription(subID)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	setDelayAction := NewSetDelayAction(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setDelayAction.inlineButtons(sub),
		},
		Text: setDelayAction.text(),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}

func (a SetDelayAction) ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service) {
	delay, err := a.processInputValue(update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, a.WrongInput())
		_, _ = bot.Send(msg)
		ShowSetDelay(update.Message.Chat.ID, a.Params.P1, bot, s)
		return
	}

	delayUpdate := entity.SubscriptionUpdate{Delay: delay}
	_, err = storage.Service.UpdateSubscription(s, a.Params.P1, &delayUpdate)
	if err != nil {
		log.Printf("cannot set delay to subscription: %v", err)
		return
	}

	_ = ResetAwaitInput(ResetInput, update.Message.Chat.ID, s)
	ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
}

func (a SetDelayAction) processInputValue(text string) (float64, error) {
	delay, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse text to number")
	}

	switch {
	case delay < 1, delay > 5e6:
		return 0, fmt.Errorf("invalid range")
	default:
		return delay, nil
	}
}

func (a SetDelayAction) WrongInput() string {
	return "Wrong input. A whole or decimal number expected; in range [1, 5e6]"
}
