package main

import (
	"TelegramBot/config"
	"TelegramBot/core"
	"fmt"
	"log"
	"os"

	"github.com/nanobox-io/golang-scribble"
)

func main() {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}

	files, err = os.ReadDir("./go")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}

	files, err = os.ReadDir("/")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}

	cfg, err := config.New("./config.yml")
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
