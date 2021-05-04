package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SubscriptionAction struct {
	Action
}

const Sub Cmd = "SUB"

func NewSubscriptionAction(subID string, reset ResetInputType) SubscriptionAction {
	return SubscriptionAction{Action{
		Cmd:    Sub,
		Params: Params{P1: subID, P2: string(reset)},
	}}
}

func (a SubscriptionAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	_ = ResetAwaitInput(ResetInputType(a.Params.P2), msg.Chat.ID, s)

	sub, err := s.GetSubscription(a.Params.P1)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(sub), a.inlineButtons(sub))
	_, _ = bot.Send(message)
}

func (a SubscriptionAction) text(sub *entity.Subscription) string {
	return fmt.Sprintf(`
<i>Subscription Settings for</i> <b>%s</b>

📶 Magnitude: ≥ %.1f
⏳ Delay: ≤ %.0f min
📍 Location: %s
⭕️ Radius: %.1f km
`, sub.Name, sub.MinMag, sub.Delay, LocationToHTMLString(sub.Location), sub.Radius)
}

func (a SubscriptionAction) inlineButtons(sub *entity.Subscription) *tgbotapi.InlineKeyboardMarkup {
	name := tgbotapi.NewInlineKeyboardButtonData("📖 Name", NewSetNameAction(sub.SubID).Encode())
	srcs := tgbotapi.NewInlineKeyboardButtonData("ℹ️ Sources", NewListSourcesAction(sub.SubID).Encode())
	magnitude := tgbotapi.NewInlineKeyboardButtonData("📶 Magnitude", NewSetMagnitudeAction(sub.SubID).Encode())
	delay := tgbotapi.NewInlineKeyboardButtonData("⏳ Delay", NewSetDelayAction(sub.SubID).Encode())
	location := tgbotapi.NewInlineKeyboardButtonData("📍️ Location", NewSetLocationAction(sub.SubID).Encode())
	radius := tgbotapi.NewInlineKeyboardButtonData("⭕️ Radius", NewSetRadiusAction(sub.SubID).Encode())
	home := tgbotapi.NewInlineKeyboardButtonData("« Subscriptions",
		NewSubscriptionsAction("").Encode())
	deleteSub := tgbotapi.NewInlineKeyboardButtonData("🗑 Delete",
		NewDeleteSubscriptionAction(a.Params.P1, "").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(name, srcs),
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

	subAction := NewSubscriptionAction(subID, "")
	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: subAction.inlineButtons(sub),
		},
		Text: subAction.text(sub),

		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, _ = bot.Send(msg)
}

func LocationToHTMLString(loc *entity.Location) string {
	if loc == nil {
		return "not set"
	}

	return fmt.Sprintf("<code>%f,%f</code> (<a href=\"http://www.google.com/maps/place/%f,%f\">maps 🌍</a>)",
		loc.Lat, loc.Lng, loc.Lat, loc.Lng)
}
