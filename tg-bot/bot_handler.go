package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mightymatth/earthquake-tools/tg-bot/action"
	"github.com/mightymatth/earthquake-tools/tg-bot/entity"
	"github.com/mightymatth/earthquake-tools/tg-bot/storage"
	"log"
	"strconv"
)

func botHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI, s storage.Service) {
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

	chatState := storage.Service.GetChatState(s, update.Message.Chat.ID)

	if chatState.AwaitInput == "" {
		action.ShowUnknownCommand(bot, update.Message.Chat.ID)
		return
	}

	actionable, err := action.New(chatState.AwaitInput)
	if err != nil {
		log.Printf("cannot perform action: %s", err)
		return
	}

	a := actionable.ToAction()
	switch a.Cmd {
	case action.CreateSub:
		if update.Message.Text == "" {
			createSubAction := action.CreateSubscriptionAction{Action: a}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, createSubAction.WrongInput())
			_, _ = bot.Send(msg)
			action.ShowCreateSubscription(update.Message.Chat.ID, bot)
			return
		}

		sub, err := storage.Service.CreateSubscription(s, update.Message.Chat.ID, update.Message.Text)
		if err != nil {
			log.Printf("cannot create subscription: %v", err)
			return
		}

		_ = action.ResetAwaitInput(action.ResetInput, update.Message.Chat.ID, s)
		action.ShowSubscription(update.Message.Chat.ID, sub.SubID, bot, s)
	case action.SetMagnitude:
		mag, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			setMagAction := action.SetMagnitudeAction{Action: a}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setMagAction.WrongInput())
			_, _ = bot.Send(msg)
			action.ShowSetMagnitude(update.Message.Chat.ID, setMagAction.Params.P1, bot, s)
			return
		}

		magUpdate := entity.SubscriptionUpdate{MinMag: mag}
		_, err = storage.Service.UpdateSubscription(s, a.Params.P1, &magUpdate)
		if err != nil {
			log.Printf("cannot set magnitude to subscription: %v", err)
			return
		}

		_ = action.ResetAwaitInput(action.ResetInput, update.Message.Chat.ID, s)
		action.ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
	case action.SetDelay:
		delay, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			setDelayAction := action.SetDelayAction{Action: a}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setDelayAction.WrongInput())
			_, _ = bot.Send(msg)
			action.ShowSetDelay(update.Message.Chat.ID, setDelayAction.Params.P1, bot, s)
			return
		}

		delayUpdate := entity.SubscriptionUpdate{Delay: delay}
		_, err = storage.Service.UpdateSubscription(s, a.Params.P1, &delayUpdate)
		if err != nil {
			log.Printf("cannot set delay to subscription: %v", err)
			return
		}

		_ = action.ResetAwaitInput(action.ResetInput, update.Message.Chat.ID, s)
		action.ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
	case action.SetLocation:
		inputLoc := update.Message.Location

		if inputLoc == nil {
			setLocationAction := action.SetLocationAction{Action: a}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setLocationAction.WrongInput())
			_, _ = bot.Send(msg)
			action.ShowSetLocation(update.Message.Chat.ID, setLocationAction.Params.P1, bot, s)
			return
		}

		location := entity.Location{
			Lat: inputLoc.Latitude,
			Lng: inputLoc.Longitude,
		}

		locationUpdate := entity.SubscriptionUpdate{Location: &location}
		_, err = storage.Service.UpdateSubscription(s, a.Params.P1, &locationUpdate)
		if err != nil {
			log.Printf("cannot set location to subscription: %v", err)
			return
		}

		_ = action.ResetAwaitInput(action.ResetInput, update.Message.Chat.ID, s)
		action.ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
	case action.SetRadius:
		radius, err := strconv.ParseFloat(update.Message.Text, 64)
		if err != nil {
			setRadiusAction := action.SetRadiusAction{Action: a}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, setRadiusAction.WrongInput())
			_, _ = bot.Send(msg)
			action.ShowSetRadius(update.Message.Chat.ID, setRadiusAction.Params.P1, bot, s)
			return
		}

		radiusUpdate := entity.SubscriptionUpdate{Radius: radius}
		_, err = storage.Service.UpdateSubscription(s, a.Params.P1, &radiusUpdate)
		if err != nil {
			log.Printf("cannot set radius to subscription: %v", err)
			return
		}

		_ = action.ResetAwaitInput(action.ResetInput, update.Message.Chat.ID, s)
		action.ShowSubscription(update.Message.Chat.ID, a.Params.P1, bot, s)
	default:
		action.ShowUnknownCommand(bot, update.Message.Chat.ID)
	}
}
