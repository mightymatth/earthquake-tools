package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/screen"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
	"strconv"
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

	if chatState.AwaitInput == "" {
		screen.ShowUnknownCommand(bot, update.Message.Chat.ID)
		return
	}

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
			return
		}

		screen.ResetAwaitInput(screen.ResetInput, update.Message.Chat.ID, s)
		screen.ShowSubscriptions(update.Message.Chat.ID, bot, s)
	case screen.SetMagnitude:
		mag, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			setMagScreen := screen.SetMagnitudeScreen{Screen: scr}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setMagScreen.WrongInput())
			bot.Send(msg)
			screen.ShowSetMagnitude(update.Message.Chat.ID, setMagScreen.Params.P1, bot, s)
			return
		}

		magUpdate := entity.SubscriptionUpdate{MinMag: mag}
		_, err = s.UpdateSubscription(scr.Params.P1, &magUpdate)
		if err != nil {
			log.Printf("cannot set magnitude to subscription: %v", err)
			return
		}

		screen.ResetAwaitInput(screen.ResetInput, update.Message.Chat.ID, s)
		screen.ShowSubscription(update.Message.Chat.ID, scr.Params.P1, bot, s)
	default:
		screen.ShowUnknownCommand(bot, update.Message.Chat.ID)
	}
}
