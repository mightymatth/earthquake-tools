package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
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
<b>Delay</b> is a time period between the earthquake occasion and the moment when the report is received. 
It is expected behavior as data may arrive from various sources.

Enter how many <u>minutes</u> you will tolerate.

e.g.: <code>2</code>, <code>60</code>
`
}

func (a SetDelayAction) WrongInput() string {
	return "Wrong input. Integer or decimal number expected."
}

func (a SetDelayAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscription(sub.SubID, ResetInput).Encode())

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
