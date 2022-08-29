package main

import (
	"TelegramBot/config"
	"TelegramBot/core"
	"github.com/nanobox-io/golang-scribble"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db, err := scribble.New("./telegram", nil)
	if err != nil {
		log.Fatal(err)
	}

	tgbot, err := core.NewBot(cfg, db)
	if err != nil {
		log.Fatal(err)
	}

	server := core.New(8090, db)
	go server.Run(cfg)

	tgbot.Run()
}
