package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"github.com/pkg/errors"
	"strings"
)

type Screen struct {
	Cmd Cmd
	Arg string
}

type Cmd string

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
	return fmt.Sprintf("%s %s", s.Cmd, s.Arg)
}

func Decode(data string) (Screener, error) {
	parts := strings.Split(data, " ")
	cmd, arg := parts[0], parts[1]

	switch Cmd(cmd) {
	case Home:
		return NewHomeScreen(), nil
	case Subs:
		return NewSubscriptionsScreen(), nil
	case Sub:
		return NewSubscriptionScreen(arg), nil
	default:
		return nil, errors.Errorf("screen '%s' doesnt exist", cmd)
	}
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
