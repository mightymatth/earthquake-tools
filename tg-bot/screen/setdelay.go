package screen

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetDelayScreen struct {
	Screen
}

const SetDelay Cmd = "SET_DELAY"

func NewSetDelayScreen(subID string) SetDelayScreen {
	return SetDelayScreen{Screen{Cmd: SetDelay, Params: Params{P1: subID}}}
}

func (scr SetDelayScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
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

func (scr SetDelayScreen) text() string {
	return `
<b>Delay</b> is a time period between earthquake time and time when the report is received.
Enter how many minutes you will tolerate.
e.g.: <code>2</code>, <code>60</code>
`
}

func (scr SetDelayScreen) WrongInput() string {
	return "Wrong input. Integer or decimal number expected."
}

func (scr SetDelayScreen) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionScreen(sub.SubID, ResetInput).Encode())

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

	setDelayScreen := NewSetDelayScreen(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setDelayScreen.inlineButtons(sub),
		},
		Text: setDelayScreen.text(),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	bot.Send(msg)
}
