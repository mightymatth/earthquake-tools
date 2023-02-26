package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tgbot/action"
	"github.com/mightymatth/earthquake-tools/tgbot/storage"
	"log"
)

func botHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, s storage.Service) {
	if update.ChannelPost != nil {
		update.Message = update.ChannelPost
	}

	if update.CallbackQuery != nil {
		a, err := action.New(update.CallbackQuery.Data)
		if err != nil {
			log.Printf("cannot perform action: %s", err)
			return
		}

		a.Perform(bot, update.CallbackQuery.Message, s)

		_, _ = bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
		return
	}

	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	switch update.Message.Text {
	case "/start":
		action.ShowHome(bot, update.Message.Chat.ID)
		return
	case "/list":
		action.ShowSubscriptions(update.Message.Chat.ID, bot, s)
		return
	}

	chatState := s.GetChatState(update.Message.Chat.ID)

	if !update.Message.Chat.IsPrivate() && chatState.DisableInput {
		return
	}

	if chatState.AwaitInput == "" {
		action.ShowUnknownCommand(bot, update.Message.Chat.ID)
		return
	}

	actionable, err := action.New(chatState.AwaitInput)
	if err != nil {
		log.Printf("cannot perform action: %s", err)
		return
	}

	actionable.ProcessUserInput(bot, &update, s)
}
