package main

import (
	"TelegramBot/config"
	"TelegramBot/core"
	"log"
)

func main() {
	cfg, err := config.New("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	tgbot, err := core.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}

	server := core.New(8090, tgbot.API)
	go server.Run(cfg)

	tgbot.Run()
}
