# Telegram Bot

## Local development

```shell
go run *.go -tg-bot-token "<token>" \
            -mongo-uri "mongodb+srv://<user>:<pasword>@<cluster>/?retryWrites=true&w=majority" \
            -db-namespace "dev" -tg-bot-debug "true"
```

### Using dotenv

```shell
# copy .env.example to .env
cp .env.example .env
# edit .env file with the correct values

# run the bot
go run *.go
```

## Deployment

### Manual

```shell
# build docker image
docker build -t eq-tg-bot .

# run docker container
docker run -d --log-opt max-size=10m --log-opt max-file=5 \
            -e TELEGRAM_BOT_TOKEN="14...w" \
            -e MONGO_URI="mongodb+srv://.../?retryWrites=true&w=majority" \
            -e DB_NAMESPACE="dev" \
            -e TG_BOT_DEBUG="true" \
            --restart=always eq-tg-bot:latest
```

### Fly.io

Install [flyctl](https://fly.io/docs/flyctl/installing/).

```shell
# log in with your account
flyctl auth login

# launch the app
flyctl launch

# set app secrets 
flyctl secrets set TELEGRAM_BOT_TOKEN=<tg-bot-token> \
                   BOT_ENV=prod \
                   MONGO_URI="<mongo-db-uri>"
```

Check [the official documentation](https://fly.io/docs/flyctl/).


## Bot Father configuration

#### Name

Earthquake Events ⚠️

#### Description

It makes you able to receive notifications of recent earthquakes by making subscriptions configured with parameters such as magnitude, location, and observing radius.

#### Commands

```
start - Starts the home screen
list - Lists current subscriptions
```
