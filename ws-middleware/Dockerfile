FROM golang:latest as builder

LABEL maintainer="Matija Pevec <pevec@live.com>"

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /out/main .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=builder /out/main .

ENV SOURCE "wss://www.seismicportal.eu/standing_order/websocket"
ENV WEBHOOK "http://localhost:3300"

CMD ./main -source $SOURCE -webhook $WEBHOOK
