FROM golang:latest as builder

LABEL maintainer="Matija Pevec <pevec@live.com>"

WORKDIR /app

COPY ../go.mod .
COPY ../go.sum .
RUN go mod download
COPY . .

RUN go build -ldflags '-w -s' -o /out/main .


FROM alpine:latest
WORKDIR /root
COPY --from=builder /out/main .

CMD ./main
