# Earthquake aggregator 

A program used to subscribe to multiple earthquake event sources and pass them to 
our own webhook. Their source of information is represented by JSON message 
in a standard format. 

## Development

```shell
# run service and pass events to webhook
go run cmd/aggregator/*.go -webhook "http://localhost:3300"
```
