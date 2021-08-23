package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetLocationAction struct {
	Action
}

const SetLocation Cmd = "SET_LOCATION"

func NewSetLocationAction(subID string) SetLocationAction {
	return SetLocationAction{Action{Cmd: SetLocation, Params: Params{P1: subID}}}
}

func (a SetLocationAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
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

func (a SetLocationAction) text() string {
	return `
<b>Subscription ᐅ Location</b>

Send the location that will mark the center of your wanted observation area.
To send the location, click <b>Send attachment</b> icon and click on <b>Send location</b>.
`
}

func (a SetLocationAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	cancel := tgbotapi.NewInlineKeyboardButtonData("❌ Cancel",
		NewSubscriptionAction(sub.SubID, ResetInput).Encode())

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

	setLocAction := NewSetLocationAction(subID)
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: setLocAction.inlineButtons(sub),
		},
		Text:                  setLocAction.text(),
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}

func (a SetLocationAction) ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service) {
	inputLoc := update.Message.Location
	if inputLoc == nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, a.WrongInput())
		_, _ = bot.Send(msg)
		ShowSetLocation(update.Message.Chat.ID, a.Params.P1, bot, s)
		return
	}

	location := entity.Location{
		Lat: inputLoc.Latitude,
		Lng: inputLoc.Longitude,
	}

	locationUpdate := entity.SubscriptionUpdate{Location: &location}
	_, err := storage.Service.UpdateSubscription(s, a.Params.P1, &locationUpdate)
	if err != nil {
		log.Printf("cannot set location to subscription: %v", err)
		return
	}

	_ = ResetAwaitInput(ResetInput, update.Message.Chat.ID, s)
	ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
}

func (a SetLocationAction) WrongInput() string {
	return "Wrong input. Location is expected."
}
