package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/screen"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
)

func botHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, s storage.Service) {
	if update.CallbackQuery != nil {
		scr, err := screen.New(update.CallbackQuery.Data)
		if err != nil {
			log.Printf("cannot create screen: %s", err)
			return
		}

		scr.TakeAction(bot, update.CallbackQuery.Message, s)

		bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
		return
	}

	if update.Message == nil { // ignore any non-Message Updates
		return
	}

	switch update.Message.Text {
	case "/start":
		screen.ShowHome(bot, update.Message.Chat.ID)
		return
	case "/list":
		screen.ShowSubscriptions(update.Message.Chat.ID, bot, s)
		return
	}

	chatState := s.GetChatState(update.Message.Chat.ID)

	screener, err := screen.New(chatState.AwaitInput)
	if err != nil {
		log.Printf("cannot create screen: %s", err)
		return
	}

	scr := screener.Type()
	switch scr.Cmd {
	case screen.CreateSub:
		_, err := s.CreateSubscription(update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			log.Printf("cannot create subscription: %v", err)
		}
		screen.ShowSubscriptions(update.Message.Chat.ID, bot, s)
	default:
		screen.ShowUnknownCommand(bot, update.Message.Chat.ID)
	}
}
