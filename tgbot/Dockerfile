FROM golang:latest as builder

LABEL maintainer="Matija Pevec <pevec@live.com>"

WORKDIR /app

COPY .. .
RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags '-w -s' -o /out/main ./tgbot

FROM alpine:latest
WORKDIR /root
COPY --from=builder /out/main .

CMD ./main
