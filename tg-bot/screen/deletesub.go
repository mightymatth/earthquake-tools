package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
)

type DeleteSubscriptionScreen struct {
	Screen
}

const DeleteSub Cmd = "DELETE_SUB"

const ConfirmDeleteSub = "+"

func NewDeleteSubscriptionScreen(subID, confirm string) DeleteSubscriptionScreen {
	return DeleteSubscriptionScreen{Screen{Cmd: DeleteSub, Params: Params{P1: subID, P2: confirm}}}
}

func (scr DeleteSubscriptionScreen) TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	sub, err := s.GetSubscription(scr.Params.P1)
	if err != nil {
		fmt.Printf("cannot get subscription: %v", err)
		return
	}

	switch scr.Params.P2 {
	case "":
		message := editedMessageConfig(msg.Chat.ID, msg.MessageID, scr.text(sub), scr.inlineButtons())
		bot.Send(message)
	case ConfirmDeleteSub:
		err := s.DeleteSubscription(scr.Params.P1)
		if err != nil {
			fmt.Printf("cannot delete subscription: %v", err)
			return
		}
		NewSubscriptionsScreen("").TakeAction(bot, msg, s)
	default:
		fmt.Print("unknown delete subscription parameter")
		return
	}
}

func (scr DeleteSubscriptionScreen) text(sub *entity.Subscription) string {
	return fmt.Sprintf(`
Do you really want to delete subscription '%s'
`, sub.Name)
}

func (scr DeleteSubscriptionScreen) inlineButtons() *tgbotapi.InlineKeyboardMarkup {
	yes := tgbotapi.NewInlineKeyboardButtonData("Yes",
		NewDeleteSubscriptionScreen(scr.Screen.Params.P1, ConfirmDeleteSub).Encode())
	no := tgbotapi.NewInlineKeyboardButtonData("No",
		NewSubscriptionScreen(scr.Screen.Params.P1, "").Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(yes, no),
	)
	return &kb
}
