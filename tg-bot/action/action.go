package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"github.com/pkg/errors"
	"strings"
)

type Action struct {
	Cmd    Cmd
	Params Params
}

type Cmd string

type Params struct {
	P1 string
	P2 string
	P3 string
}

type ResetInputType string

const ResetInput ResetInputType = "+"

type Actionable interface {
	Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service)
	ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service)
}

func New(data string) (Actionable, error) {
	s, err := Decode(data)
	if err != nil {
		return nil, fmt.Errorf("unknown action '%s'", data)
	}

	return s, nil
}

func (a Action) ToAction() Action {
	return a
}

func (a Action) Encode() string {
	return fmt.Sprintf("%s:%s:%s:%s", a.Cmd, a.Params.P1, a.Params.P2, a.Params.P3)
}

func Decode(data string) (Actionable, error) {
	parts := strings.Split(data, ":")
	if len(parts) != 4 {
		return nil, errors.Errorf("empty or data format: '%v'", data)
	}

	cmd, p1, p2, p3 := parts[0], parts[1], parts[2], parts[3]

	switch Cmd(cmd) {
	case Home:
		return NewHomeAction(), nil
	case Subs:
		return NewSubscriptionsAction(ResetInputType(p1)), nil
	case Sub:
		return NewSubscriptionAction(p1, ResetInputType(p2)), nil
	case CreateSub:
		return NewCreateSubscriptionAction(), nil
	case DeleteSub:
		return NewDeleteSubscriptionAction(p1, p2), nil
	case SetName:
		return NewSetNameAction(p1), nil
	case ListSources:
		return NewListSourcesAction(p1), nil
	case SetSource:
		return NewSetSourceAction(p1, p2, p3), nil
	case SetMagnitude:
		return NewSetMagnitudeAction(p1), nil
	case SetDelay:
		return NewSetDelayAction(p1), nil
	case SetLocation:
		return NewSetLocationAction(p1), nil
	case SetRadius:
		return NewSetRadiusAction(p1), nil
	case Settings:
		return NewSettingsAction(ResetInputType(p1)), nil
	case DisableInput:
		return NewDisableInputAction(p1), nil
	default:
		return nil, errors.Errorf("action with command '%s' doesn't exist", cmd)
	}
}

func (a Action) ProcessUserInput(bot *tgbotapi.BotAPI, update *tgbotapi.Update, s storage.Service) {
	ShowUnknownCommand(bot, update.Message.Chat.ID)
}

func ResetAwaitInput(resetInput ResetInputType, chatID int64, s storage.Service) error {
	switch resetInput {
	case ResetInput:
		err := s.SetAwaitUserInput(chatID, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func editedMessageConfig(
	chatID int64, msgID int, newText string,
	newReplyMarkup *tgbotapi.InlineKeyboardMarkup,
) tgbotapi.EditMessageTextConfig {
	editedMessage := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      chatID,
			MessageID:   msgID,
			ReplyMarkup: newReplyMarkup,
		},
		Text:                  newText,
		ParseMode:             tgbotapi.ModeHTML,
		DisableWebPagePreview: true,
	}

	return editedMessage
}
