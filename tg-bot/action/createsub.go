package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
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

func (a CreateSubscriptionAction) WrongInput() string {
	return "Wrong input. Text value expected."
}
