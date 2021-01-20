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
	}

	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

	chatState := s.GetChatState(update.Message.Chat.ID)

	log.Printf("chatState: %v", chatState)

	scr, err := screen.New(chatState.AwaitInput)
	if err != nil {
		log.Printf("cannot create screen: %s", err)
		return
	}

	log.Printf("cmd %v", scr.Type().Cmd)

	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	//msg.ReplyToMessageID = update.Message.MessageID
	//
	//bot.Send(msg)

	screen.ShowUnknownCommand(bot, update.Message.Chat.ID)
}
