package main

import (
	"TelegramBot/config"
	"TelegramBot/core"
	"fmt"
	"github.com/nanobox-io/golang-scribble"
	"log"
	"os"
)

const redhatCfg = `token: "5440030440:AAGiPBn9JHb7VW60aRhf5S6HUqw6V8eEVgI"
app_id: 10550974
app_hash: "698997f261777bbfb3bf91db82641034"
server: "149.154.167.50:443"
key: "pub_keys.pem"
login: admin
password: admin
channels:
  volgograd:
      Охрана:
        id: "-1001675961636"
        link: "https://t.me/+IbKT55gmFFY3Mjli"
      Бар:
        id: "-1001749986535"
        link: "https://t.me/+39UtHLL-Io84ODFi"
      Арт:
        id: "-1001709159006"
        link: "https://t.me/+bmt4pt_U-xo0Mzhi"
      Админ:
        id: "-1001699370600"
        link: "https://t.me/+sCjHoRiemWNlOTVi"
  moscow:
      Охрана:
        id: "-1001697572469"
        link: "https://t.me/+lCSB8ONa-_A5MjBi"
      Бар:
        id: "-1001605922666"
        link: "https://t.me/+ZI_4oOBxTO1jMDAy"
      Арт:
        id: "-1001770807023"
        link: "https://t.me/+H0gKILCOoEZhNzI6"
      Админ:
        id: "-1001579885451"
        link: "https://t.me/+csVBWNTtOio1Yzgy"
  krasnodar:
      Охрана:
        id: "-1001633711454"
        link: "https://t.me/+RefSzW2YfyRiM2Qy"
      Бар:
        id: "-1001636988740"
        link: "https://t.me/+6ZIc6zsrSyA2MDEy"
      Арт:
        id: "-1001799331606"
        link: "https://t.me/+NfhguhhKwWlhMmZi"
      Админ:
        id: "-1001678631440"
        link: "https://t.me/+UUluhUnk4tZhYTIy"`

func main() {
	err := os.WriteFile("config.yml", []byte(redhatCfg), os.ModePerm) // only for redhat
	if err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(".")
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
