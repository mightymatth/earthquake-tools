package screen

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Screen string

type Screener interface {
	TakeAction(bot *tgbotapi.BotAPI, msg *tgbotapi.Message)
}

func New(data string) (Screener, error) {
	var s Screener

	switch data {
	case string(Home):
		s = HomeScreen(data)
	case string(Settings):
		s = SettingsScreen(data)
	case string(Magnitude):
		s = MagnitudeScreen(data)
	case string(EditMagnitude):
		s = EditMagnitudeScreen(data)
	case string(Delay):
		s = DelayScreen(data)
	case string(EditDelay):
		s = EditDelayScreen(data)
	default:
		return nil, fmt.Errorf("unknown screen '%s'", data)
	}

	return s, nil
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
