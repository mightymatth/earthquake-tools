package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetMagnitudeScreen struct {
	Screen
}

const SetMagnitude Cmd = "SET_MAG"

func NewSetMagnitudeScreen(subID string) SetMagnitudeScreen {
	return SetMagnitudeScreen{Screen{Cmd: SetMagnitude, Params: Params{P1: subID}}}
}

func (scr SetMagnitudeScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
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

func (scr SetMagnitudeScreen) text() string {
	return `
Enter minimum magnitude to receive an alert.
e.g.: <code>4.3</code>, <code>5</code>
`
}

func (scr SetMagnitudeScreen) WrongInput() string {
	return "Wrong input. Integer or decimal number expected."
}

func (scr SetMagnitudeScreen) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionScreen(sub.SubID, ResetInput).Encode())

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
