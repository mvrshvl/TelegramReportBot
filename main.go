package main

import (
	"context"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"TelegramBot/config"
	"TelegramBot/core"
	"TelegramBot/core/database"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	db := database.New(cfg)

	err = db.Connect(ctx)
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
