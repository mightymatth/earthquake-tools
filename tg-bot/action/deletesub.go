package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type DeleteSubscriptionAction struct {
	Action
}

const DeleteSub Cmd = "DELETE_SUB"

const ConfirmDeleteSub = "+"

func NewDeleteSubscriptionAction(subID, confirm string) DeleteSubscriptionAction {
	return DeleteSubscriptionAction{Action{Cmd: DeleteSub, Params: Params{P1: subID, P2: confirm}}}
}

func (a DeleteSubscriptionAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	sub, err := s.GetSubscription(a.Params.P1)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	switch a.Params.P2 {
	case "":
		message := editedMessageConfig(msg.Chat.ID, msg.MessageID, a.text(sub), a.inlineButtons())
		_, _ = bot.Send(message)
	case ConfirmDeleteSub:
		err := s.DeleteSubscription(a.Params.P1)
		if err != nil {
			log.Printf("cannot delete subscription: %v", err)
			return
		}
		NewSubscriptionsAction("").Perform(bot, msg, s)
	default:
		log.Print("unknown delete subscription parameter")
		return
	}
}

func (a DeleteSubscriptionAction) text(sub *entity.Subscription) string {
	return fmt.Sprintf(`
Do you really want to delete subscription <b>%s</b>?
`, sub.Name)
}

func (a DeleteSubscriptionAction) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	yes := tgbotapi.NewInlineKeyboardButtonData("✅ Yes",
		NewDeleteSubscriptionAction(a.Action.Params.P1, ConfirmDeleteSub).Encode())
	no := tgbotapi.NewInlineKeyboardButtonData("❌ No",
		NewSubscriptionAction(a.Action.Params.P1, "").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(yes, no),
	)
	return &kb
}
