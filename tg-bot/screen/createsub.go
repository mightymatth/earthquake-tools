package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type CreateSubscriptionScreen struct {
	Screen
}

const CreateSub Cmd = "CREATE_SUB"

func NewCreateSubscriptionScreen() CreateSubscriptionScreen {
	return CreateSubscriptionScreen{Screen{Cmd: CreateSub}}
}

func (scr CreateSubscriptionScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	err := s.SetAwaitUserInput(msg.Chat.ID, entity.CreateSubName)
	if err != nil {
		log.Printf("cannot set chat state: %v", err)
		return
	}

	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(), scr.inlineButtons())
	bot.Send(message)
}

func (scr CreateSubscriptionScreen) text() string {
	return `
Enter short subscription name.
`
}

func (scr CreateSubscriptionScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("Cancel",
		NewSubscriptionsScreen(ResetInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(cancel),
	)
	return &kb
}
