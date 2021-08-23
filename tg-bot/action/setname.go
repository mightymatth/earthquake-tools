package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetNameAction struct {
	Action
}

const SetName Cmd = "SET_NAME"

func NewSetNameAction(subID string) SetNameAction {
	return SetNameAction{Action{Cmd: SetName, Params: Params{P1: subID}}}
}

func (a SetNameAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
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
	_, _ = bot.Send(message)
}

func (a SetNameAction) text() string {
	return `
<b>Subscription ᐅ Name</b>

Enter short subscription name (e.g. city, region or state you want to observe).
`
}

func (a SetNameAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionAction(sub.SubID, ResetInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(cancel),
	)
	return &kb
}

func ShowSetName(chatID int64, subID string, bot *tgbotapi.BotAPI, s storage.Service) {
	sub, err := s.GetSubscription(subID)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	setNameAction := NewSetNameAction(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setNameAction.inlineButtons(sub),
		},
		Text:                  setNameAction.text(),
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}

func (a SetNameAction) ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service) {
	if err := a.processInputValue(update.Message.Text); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, a.WrongInput())
		_, _ = bot.Send(msg)
		ShowSetName(update.Message.Chat.ID, a.Params.P1, bot, s)
		return
	}

	nameUpdate := entity.SubscriptionUpdate{Name: update.Message.Text}
	_, err := storage.Service.UpdateSubscription(s, a.Params.P1, &nameUpdate)
	if err != nil {
		log.Printf("cannot set name to subscription: %v", err)
		return
	}

	_ = ResetAwaitInput(ResetInput, update.Message.Chat.ID, s)
	ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
}

func (a SetNameAction) processInputValue(text string) error {
	switch {
	case text == "", len(text) > 64:
		return fmt.Errorf("invalid")
	default:
		return nil
	}
}

func (a SetNameAction) WrongInput() string {
	return "Wrong input. A simple, text value expected; no more than 64 characters."
}
