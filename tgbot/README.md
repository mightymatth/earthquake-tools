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
# run the following commands from the current directory (tgbot)

# build docker image
docker build -t eq-tg-bot .

# run docker container
docker run \
            -e TELEGRAM_BOT_TOKEN="14...w" \
            -e MONGO_URI="mongodb+srv://.../?retryWrites=true&w=majority" \
            -e DB_NAMESPACE="dev" \
            -e TG_BOT_DEBUG="true" \
            --restart=always eq-tg-bot:latest
```

### Fly.io

Install [flyctl](https://fly.io/docs/flyctl/installing/).

```shell
# run the following commands from the root directory

# log in with your account
flyctl auth login

# set app name
export APP_NAME=eq-tg-bot

# launch the app
fly launch --path tgbot --region ams --copy-config --remote-only --no-deploy --name $APP_NAME  

# set app secrets 
flyctl secrets set TELEGRAM_BOT_TOKEN="<token>" \
                   MONGO_URI="<mongo-db-uri>" -a $APP_NAME 

# deploy
fly deploy --dockerfile tgbot/Dockerfile -a $APP_NAME        

# if deployment fails, try to restart it
fly scale count 0 -a $APP_NAME
fly scale count 1 -a $APP_NAME

# test deploy (run from root directory)
fly deploy --build-only --no-cache --config tgbot/fly.toml
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
