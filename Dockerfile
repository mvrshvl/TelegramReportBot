FROM golang:1.18-alpine

RUN apk add --no-cache musl-dev make linux-headers git gcc bash curl openssl

COPY . ./bot
WORKDIR bot

RUN go mod tidy
RUN go build -o tgbot ./main.go

ENTRYPOINT ["./tgbot"]