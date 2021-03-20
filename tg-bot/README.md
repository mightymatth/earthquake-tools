# Telegram Bot

## Local development

Setup `ws-middleware` that will send events to `localhost:3300` and run the next
script:

```shell
➜  tg-bot git:(master) MONGO_URI="mongodb+srv://..." go run .
```

This command will:

* spin up a server on `localhost:3300` that listens to earthquake events and
  send them to the bot,
* start bot to listen to user events.

### Bot Father configuration

#### Name

EMSC Events ⚠️

#### Description

It makes you able to receive notifications of recent earthquakes by making subscriptions configured with parameters such as magnitude, location, and observing radius.

#### Commands

```
start - Starts the home screen
list - Lists current subscriptions
```
