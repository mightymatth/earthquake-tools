# Earthquake Tools

[![Go Report Card](https://goreportcard.com/badge/github.com/mightymatth/earthquake-tools)](https://goreportcard.com/report/github.com/mightymatth/earthquake-tools)

Set of utilities that help to track and notify users about recent earthquakes.

The current utilities are:

* [Telegram Bot](/tg-bot)
    * the bot is available [here](https://t.me/EarthquakeEventsBot).
* [Earthquake aggregator](/eq-aggregator)
    * used to subscribe to multiple earthquake sources and notify Telegram Bot (or anyone else) about the events.

# Deployment

To start bot and aggregator together, setup `.env` file with valid secrets and use Docker Compose:
```shell
docker-compose up --build -d
```
