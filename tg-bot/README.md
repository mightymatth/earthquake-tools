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

The bot let you subscribe to certain rules and send notifications on earthquakes
with given parameters.

#### Commands

```
start - Starts the home screen
list - Lists current subscriptions
```
