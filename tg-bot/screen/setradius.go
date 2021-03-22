package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetRadiusScreen struct {
	Screen
}

const SetRadius Cmd = "SET_RADIUS"

func NewSetRadiusScreen(subID string) SetRadiusScreen {
	return SetRadiusScreen{Screen{Cmd: SetRadius, Params: Params{P1: subID}}}
}

func (scr SetRadiusScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	err := s.SetAwaitUserInput(msg.Chat.ID, scr.Encode())
	if err != nil {
		log.Printf("cannot set chat state: %v", err)
		return
	}

	sub, err := s.GetSubscription(scr.Params.P1)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(), scr.inlineButtons(sub))
	bot.Send(message)
}

func (scr SetRadiusScreen) text() string {
	return `
Enter maximum earthquake location radius to receive an alert.
The unit is kilometer (km).
e.g.: <code>100.5</code>, <code>350</code>
`
}

func (scr SetRadiusScreen) WrongInput() string {
	return "Wrong input. Integer or decimal number expected."
}

func (scr SetRadiusScreen) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionScreen(sub.SubID, ResetInput).Encode())

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

	setMagScreen := NewSetMagnitudeScreen(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setMagScreen.inlineButtons(sub),
		},
		Text: setMagScreen.text(),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	bot.Send(msg)
}
