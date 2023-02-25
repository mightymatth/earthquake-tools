package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tgbot/storage"
	"log"
)

type CreateSubscriptionAction struct {
	Action
}

const CreateSub Cmd = "CREATE_SUB"

func NewCreateSubscriptionAction() CreateSubscriptionAction {
	return CreateSubscriptionAction{Action{Cmd: CreateSub}}
}

func (a CreateSubscriptionAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	err := s.SetAwaitUserInput(msg.Chat.ID, a.Encode())
	if err != nil {
		log.Printf("cannot set chat state: %v", err)
		return
	}

	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(), a.inlineButtons())
	_, _ = bot.Send(message)
}

func (a CreateSubscriptionAction) text() string {
	return `
<b>Subscriptions ᐅ New</b>
Enter short subscription name (e.g. city, region or state you want to observe).
`
}

func (a CreateSubscriptionAction) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionsAction(ResetInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(cancel),
	)
	return &kb
}

func (a CreateSubscriptionAction) ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service) {
	if err := a.validateInput(update.Message.Text); err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, a.WrongInput())
		_, _ = bot.Send(msg)
		ShowCreateSubscription(update.Message.Chat.ID, bot)
		return
	}

	sub, err := storage.Service.CreateSubscription(s, update.Message.Chat.ID, update.Message.Text)
	if err != nil {
		log.Printf("cannot create subscription: %v", err)
		return
	}

	_ = ResetAwaitInput(ResetInput, update.Message.Chat.ID, s)
	ShowSubscription(update.Message.Chat.ID, sub.SubID, bot, s)
}

func ShowCreateSubscription(chatID int64, bot *tgbotapi.BotAPI) {
	setRadiusAction := NewCreateSubscriptionAction()
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setRadiusAction.inlineButtons(),
		},
		Text: setRadiusAction.text(),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}

func (a CreateSubscriptionAction) validateInput(text string) error {
	switch {
	case text == "", len(text) > 64:
		return fmt.Errorf("invalid")
	default:
		return nil
	}
}

func (a CreateSubscriptionAction) WrongInput() string {
	return "Wrong input. A simple, text value expected; no more than 64 characters."
}
