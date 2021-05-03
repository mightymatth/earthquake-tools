# Earthquake aggregator 

A program used to subscribe to multiple earthquake event sources and pass them to our own webhook. Their source of information is represented by JSON message in a standard format. 

## Development

```shell
# install dependencies
go get 

# run the middleware service
go run . 
```

By default, the earthquake aggregator service will pass events in JSON format to webhook URL `http://localhost:3300`.

However, custom parameters can be set:
```shell
go run . -webhook "https://custom.webhook.url"
```

## Deployment

```shell
# build docker image
docker build -t eq-aggregator .

# run docker container
docker run -d --restart=always eq-aggregator:latest
```

The container will have the same source and webhook defaults. They can be changed by setting environment variables:
````shell
docker run -d --restart=always -e WEBHOOK=https://custom.websocket.url eq-aggregator:latest
````

If webhook is set on host machine, the container needs to be bind to host network. For macOS use `host.docker.internal` instead of `localhost`.

````shell
docker run -d --restart=always --network host -e WEBHOOK=http://host.docker.internal:3300 eq-aggregator:latest
````
