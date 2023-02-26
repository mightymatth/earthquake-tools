package action

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tgbot/entity"
	"github.com/mightymatth/earthquake-tools/tgbot/storage"
	"log"
)

type DisableInputAction struct {
	Action
}

const DisableInput Cmd = "DISABLE_INPUT"

const ToggleDisableInput = "+"

func NewDisableInputAction(confirm string) DisableInputAction {
	return DisableInputAction{Action{Cmd: DisableInput, Params: Params{P1: confirm}}}
}

func (a DisableInputAction) Perform(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, s storage.Service) {
	state := s.GetChatState(msg.Chat.ID)

	switch a.Params.P1 {
	case "":
		message := editedMessageConfig(msg.Chat.ID, msg.MessageID,
			a.text(*state), a.inlineButtons(state.DisableInput))
		_, _ = bot.Send(message)
	case ToggleDisableInput:
		update := entity.ChatStateUpdate{DisableInput: !state.DisableInput}

		_, err := s.SetChatState(msg.Chat.ID, &update)
		if err != nil {
			log.Printf("cannot set chat state: %v", err)
			return
		}

		NewSettingsAction("").Perform(bot, msg, s)
	default:
		log.Print("unknown disable input parameter")
		return
	}
}

func (a DisableInputAction) text(src entity.ChatState) string {
	status := "ENABLED ✅"
	if src.DisableInput {
		status = "DISABLED ❌"
	}

	return fmt.Sprintf(`
<b>Settings ᐅ Disable input</b>

In groups and channels, bot can be triggered with any message.
To disable this behavior, set user input to disabled state.

User input is currently %s.
`, status)
}

func (a DisableInputAction) inlineButtons(inputDisabled bool) *tgbotapi.InlineKeyboardMarkup {
	toggleButtonText := "Disable ❌"
	if inputDisabled {
		toggleButtonText = "Enable ✅"
	}

	back := tgbotapi.NewInlineKeyboardButtonData("« Settings",
		NewSettingsAction("").Encode())
	toggle := tgbotapi.NewInlineKeyboardButtonData(toggleButtonText,
		NewDisableInputAction(ToggleDisableInput).Encode())

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(back, toggle),
	)

	return &kb
}
