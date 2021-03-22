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

## Deployment

```shell
# build docker image
docker build -t eq-tg-bot .

# run docker container
docker run -d --network host --log-opt max-size=10m --log-opt max-file=5 -e MONGO_URI="mongodb+srv://.../?retryWrites=true&w=majority" -e TELEGRAM_BOT_TOKEN="14...w" -e TIMEZONEDB_API_KEY="V...H" -e BOT_ENV=prod --restart=always eq-tg-bot:latest
```

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
