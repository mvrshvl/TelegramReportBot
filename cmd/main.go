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

	tgbot.Run()
}
