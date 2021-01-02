# [SeismicPortal](https://www.seismicportal.eu/) Websocket Middleware

A program used to subscribe to SeismicPortal events about recent earthquakes and pass them to our own webhook. Their source of information is represented by JSON message delivered via [Websocket service](https://www.seismicportal.eu/realtime.html). 

## Development

```shell
# install dependencies
go get 

# run the middleware service
go run main.go
```

By default, the middleware service will subscribe to Websocket source URL `wss://www.seismicportal.eu/standing_order/websocket` and pass events in JSON format to webhook URL `http://localhost:3300`.

However, custom parameters can be set:
```shell
go run main.go -source "wss://custom.websocket.url" -webhook "https://custom.webhook.url"
```

## Deployment

```shell
# build docker image
docker build -t eq-ws-middleware .

# run docker container
docker run -d --restart=always eq-ws-middleware:latest
```

The container will have the same source and webhook defaults. They can be changed by setting environment variables:
````shell
docker run -d --restart=always -e SOURCE=wss://custom.websocket.url -e WEBHOOK=https:/custom.websocket.url eq-ws-middleware:latest
````

If webhook is set on host machine, the container needs to be bind to host network. If using macOS, use `host.docker.internal` instead of `localhost`.

````shell
docker run -d --restart=always --network host -e WEBHOOK=http://host.docker.internal:3300 eq-ws-middleware:latest
````
