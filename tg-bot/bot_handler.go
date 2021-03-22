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

		_, _ = bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{CallbackQueryID: update.CallbackQuery.ID})
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

	chatState := storage.Service.GetChatState(s, update.Message.Chat.ID)

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
		if update.Message.Text == "" {
			// TODO: show user a message should be text.
			// This happens in scenario when a user sends sticker or something else...
			return
		}

		_, err := storage.Service.CreateSubscription(s, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			log.Printf("cannot create subscription: %v", err)
			return
		}

		_ = screen.ResetAwaitInput(screen.ResetInput, update.Message.Chat.ID, bot, s)
		screen.ShowSubscriptions(update.Message.Chat.ID, bot, s)
	case screen.SetMagnitude:
		mag, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			setMagScreen := screen.SetMagnitudeScreen{Screen: scr}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setMagScreen.WrongInput())
			_, _ = bot.Send(msg)
			screen.ShowSetMagnitude(update.Message.Chat.ID, setMagScreen.Params.P1, bot, s)
			return
		}

		magUpdate := entity.SubscriptionUpdate{MinMag: mag}
		_, err = storage.Service.UpdateSubscription(s, scr.Params.P1, &magUpdate)
		if err != nil {
			log.Printf("cannot set magnitude to subscription: %v", err)
			return
		}

		_ = screen.ResetAwaitInput(screen.ResetInput, update.Message.Chat.ID, bot, s)
		screen.ShowSubscription(update.Message.Chat.ID, scr.Params.P1, bot, s)
	case screen.SetDelay:
		delay, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			setDelayScreen := screen.SetDelayScreen{Screen: scr}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setDelayScreen.WrongInput())
			_, _ = bot.Send(msg)
			screen.ShowSetDelay(update.Message.Chat.ID, setDelayScreen.Params.P1, bot, s)
			return
		}

		delayUpdate := entity.SubscriptionUpdate{Delay: delay}
		_, err = storage.Service.UpdateSubscription(s, scr.Params.P1, &delayUpdate)
		if err != nil {
			log.Printf("cannot set delay to subscription: %v", err)
			return
		}

		_ = screen.ResetAwaitInput(screen.ResetInput, update.Message.Chat.ID, bot, s)
		screen.ShowSubscription(update.Message.Chat.ID, scr.Params.P1, bot, s)
	case screen.SetLocation:
		inputLoc := update.Message.Location

		if inputLoc == nil {
			setLocationScreen := screen.SetLocationScreen{Screen: scr}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setLocationScreen.WrongInput())
			_, _ = bot.Send(msg)
			screen.ShowSetLocation(update.Message.Chat.ID, setLocationScreen.Params.P1, bot, s)
			return
		}

		location := entity.Location{
			Lat: inputLoc.Latitude,
			Lng: inputLoc.Longitude,
		}

		locationUpdate := entity.SubscriptionUpdate{Location: &location}
		_, err = storage.Service.UpdateSubscription(s, scr.Params.P1, &locationUpdate)
		if err != nil {
			log.Printf("cannot set location to subscription: %v", err)
			return
		}

		_ = screen.ResetAwaitInput(screen.ResetInput, update.Message.Chat.ID, bot, s)
		screen.ShowSubscription(update.Message.Chat.ID, scr.Params.P1, bot, s)
	case screen.SetRadius:
		radius, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			setRadiusScreen := screen.SetRadiusScreen{Screen: scr}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setRadiusScreen.WrongInput())
			_, _ = bot.Send(msg)
			screen.ShowSetRadius(update.Message.Chat.ID, setRadiusScreen.Params.P1, bot, s)
			return
		}

		radiusUpdate := entity.SubscriptionUpdate{Radius: radius}
		_, err = storage.Service.UpdateSubscription(s, scr.Params.P1, &radiusUpdate)
		if err != nil {
			log.Printf("cannot set radius to subscription: %v", err)
			return
		}

		_ = screen.ResetAwaitInput(screen.ResetInput, update.Message.Chat.ID, bot, s)
		screen.ShowSubscription(update.Message.Chat.ID, scr.Params.P1, bot, s)
	default:
		screen.ShowUnknownCommand(bot, update.Message.Chat.ID)
	}
}