package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"github.com/pkg/errors"
	"strings"
)

type Screen struct {
	Cmd    Cmd
	Params Params
}

type Cmd string

type Params struct {
	P1 string
	P2 string
}

type ResetInputType string

const ResetInput ResetInputType = "+"

type Screener interface {
	TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service)
}

func New(data string) (Screener, error) {
	s, err := Decode(data)
	if err != nil {
		return nil, fmt.Errorf("unknown screen '%s'", data)
	}

	return s, nil
}

func (s Screen) Encode() string {
	return fmt.Sprintf("%s:%s:%s", s.Cmd, s.Params.P1, s.Params.P2)
}

func Decode(data string) (Screener, error) {
	parts := strings.Split(data, ":")
	cmd, p1, p2 := parts[0], parts[1], parts[2]

	switch Cmd(cmd) {
	case Home:
		return NewHomeScreen(), nil
	case Subs:
		return NewSubscriptionsScreen(ResetInputType(p1)), nil
	case Sub:
		return NewSubscriptionScreen(p1, ResetInputType(p2)), nil
	case CreateSub:
		return NewCreateSubscriptionScreen(), nil
	default:
		return nil, errors.Errorf("screen '%s' doesnt exist", cmd)
	}
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
