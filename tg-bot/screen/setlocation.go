package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetLocationScreen struct {
	Screen
}

const SetLocation Cmd = "SET_LOCATION"

func NewSetLocationScreen(subID string) SetLocationScreen {
	return SetLocationScreen{Screen{Cmd: SetLocation, Params: Params{P1: subID}}}
}

func (scr SetLocationScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
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

func (scr SetLocationScreen) text() string {
	return `
Send the location that will mark the center of your wanted observation area.
To send the location, click <b>Send attachment</b> icon and click on <b>Send location</b>.
`
}

func (scr SetLocationScreen) WrongInput() string {
	return "Wrong input. Location is expected."
}

func (scr SetLocationScreen) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionScreen(sub.SubID, ResetInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(cancel),
	)
	return &kb
}

func ShowSetLocation(chatID int64, subID string, bot *tgbotapi.BotAPI, s storage.Service) {
	sub, err := s.GetSubscription(subID)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	setLocScreen := NewSetLocationScreen(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setLocScreen.inlineButtons(sub),
		},
		Text:                  setLocScreen.text(),
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}
