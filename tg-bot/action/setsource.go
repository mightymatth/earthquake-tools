package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

type SetSourceAction struct {
	Action
}

const SetSource Cmd = "SET_SOURCE"

const ToggleSetSource = "+"

func NewSetSourceAction(subID, sourceID, confirm string) SetSourceAction {
	return SetSourceAction{Action{Cmd: SetSource, Params: Params{P1: subID, P2: sourceID, P3: confirm}}}
}

func (a SetSourceAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	sub, err := s.GetSubscription(a.Params.P1)
	if err != nil {
		log.Printf("cannot get subscription: %v", err)
		return
	}

	subSrcs := make(map[string]Source)
	for _, srcID := range sub.Sources {
		if src, found := sourcesM[string(srcID)]; found {
			subSrcs[string(srcID)] = src
		}
	}

	src, found := sourcesM[a.Params.P2]
	if !found {
		log.Printf("invalid source")
		return
	}

	var srcActive bool
	if _, found := subSrcs[string(src.SourceID)]; found {
		srcActive = true
	}

	switch a.Params.P3 {
	case "":
		message := editedMessageConfig(msg.Chat.ID, msg.MessageID,
			a.text(src, srcActive), a.inlineButtons(sub.Sources, srcActive))
		_, _ = bot.Send(message)
	case ToggleSetSource:
		if srcActive {
			delete(subSrcs, string(src.SourceID))
		} else {
			subSrcs[string(src.SourceID)] = src
		}

		newIDs := make([]entity.SourceID, 0, len(subSrcs))
		for k := range subSrcs {
			newIDs = append(newIDs, entity.SourceID(k))
		}

		update := entity.SubscriptionUpdate{Sources: newIDs}

		_, err := s.UpdateSubscription(a.Params.P1, &update)
		if err != nil {
			log.Printf("cannot update subscription: %v", err)
			return
		}

		NewListSourcesAction(a.Params.P1).Perform(bot, msg, s)
	default:
		log.Print("unknown set source parameter")
		return
	}
}

func (a SetSourceAction) text(src Source, srcActive bool) string {
	activeStatus := "INACTIVE ❌"
	if srcActive {
		activeStatus = "ACTIVE ✅"
	}

	return fmt.Sprintf(`
<b>%s</b>
%s

This source is currently %s.
`, src.Name, src.Description, activeStatus)
}

func (a SetSourceAction) inlineButtons(srcIDs []entity.SourceID, srcActive bool) *tgbotapi.InlineKeyboardMarkup {
	toggleButtonText := "Activate ✅"
	if srcActive {
		toggleButtonText = "Deactivate ❌"
	}

	row := make([]tgbotapi.InlineKeyboardButton, 0, 2)

	back := tgbotapi.NewInlineKeyboardButtonData("« Sources",
		NewListSourcesAction(a.Action.Params.P1).Encode())
	row = append(row, back)

	if !srcActive || len(srcIDs) > 1 {
		toggle := tgbotapi.NewInlineKeyboardButtonData(toggleButtonText,
			NewSetSourceAction(a.Action.Params.P1, a.Action.Params.P2, ToggleSetSource).Encode())
		row = append(row, toggle)
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(row...),
	)

	return &kb
}
