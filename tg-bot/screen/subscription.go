package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SubscriptionScreen struct {
	Screen
}

const Sub Cmd = "SUB"

func NewSubscriptionScreen(subID string, reset ResetInputType) SubscriptionScreen {
	return SubscriptionScreen{Screen{
		Cmd:    Sub,
		Params: Params{P1: subID, P2: string(reset)},
	}}
}

func (scr SubscriptionScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	_ = ResetAwaitInput(ResetInputType(scr.Params.P2), msg.Chat.ID, bot, s)

	sub, err := s.GetSubscription(scr.Params.P1)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(sub), scr.inlineButtons(sub))
	bot.Send(message)
}

func (scr SubscriptionScreen) text(sub *entity.Subscription) string {
	return fmt.Sprintf(`
<i>Subscription Settings for</i> <b>%s</b>

📶 Magnitude: ≥ %.1f
⏳ Delay: ≤ %.0f min
📍 Location: %s
⭕️ Radius: %.1f km
`, sub.Name, sub.MinMag, sub.Delay, LocationToHTMLString(sub.Location), sub.Radius)
}

func (scr SubscriptionScreen) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	magnitude := tgbotapi.NewInlineKeyboardButtonData("📶 Magnitude", NewSetMagnitudeScreen(sub.SubID).Encode())
	delay := tgbotapi.NewInlineKeyboardButtonData("⏳ Delay", NewSetDelayScreen(sub.SubID).Encode())
	location := tgbotapi.NewInlineKeyboardButtonData("📍️ Location", NewSetLocationScreen(sub.SubID).Encode())
	radius := tgbotapi.NewInlineKeyboardButtonData("⭕️ Radius", NewSetRadiusScreen(sub.SubID).Encode())
	home := tgbotapi.NewInlineKeyboardButtonData("« Subscriptions",
		NewSubscriptionsScreen("").Encode())
	deleteSub := tgbotapi.NewInlineKeyboardButtonData("🗑 Delete",
		NewDeleteSubscriptionScreen(scr.Params.P1, "").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(magnitude, delay),
		tgbotapi.NewInlineKeyboardRow(location, radius),
		tgbotapi.NewInlineKeyboardRow(home, deleteSub),
	)
	return &kb
}

func ShowSubscription(chatID int64, subID string, bot *tgbotapi.BotAPI, s storage.Service) {
	sub, err := s.GetSubscription(subID)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	subScreen := NewSubscriptionScreen(subID, "")
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: subScreen.inlineButtons(sub),
		},
		Text: subScreen.text(sub),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	bot.Send(msg)
}

func LocationToHTMLString(loc *entity.Location) string {
	if loc == nil {
		return "not set"
	}

	return fmt.Sprintf("<code>%f,%f</code> (<a href=\"http://www.google.com/maps/place/%f,%f\">maps 🌍</a>)",
		loc.Lat, loc.Lng, loc.Lat, loc.Lng)
}
