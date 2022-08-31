package core

import (
	"TelegramBot/config"
	"TelegramBot/core/database"
	handlers2 "TelegramBot/core/handlers"
	"TelegramBot/core/stack"
	"fmt"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	API *telebot.Bot
}

func NewBot(cfg *config.Config, db *database.Database) (*Bot, error) {
	api, err := telebot.NewBot(telebot.Settings{
		Token:   cfg.Token,
		Updates: 60,
		Poller:  nil,
	})
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	botStack := stack.New()

	handlersList := createHandlers(botStack, db, api, cfg)

	useHandle(api, botStack, handlersList)

	if err := setCommands(api); err != nil {
		return nil, fmt.Errorf("set commands: %w", err)
	}

	return &Bot{
		API: api,
	}, nil
}

func (bot *Bot) Run() {
	bot.API.Start()
}

func createHandlers(stack *stack.Stack, db *database.Database, api *telebot.Bot, cfg *config.Config) []*handlers2.HandlerBase {
	return []*handlers2.HandlerBase{
		handlers2.NewCity(stack, handlers2.VlgName),
		handlers2.NewCity(stack, handlers2.MskName),
		handlers2.NewCity(stack, handlers2.KrdName),

		handlers2.NewCity(stack, handlers2.KrdCommand),
		handlers2.NewCity(stack, handlers2.VlgCommand),
		handlers2.NewCity(stack, handlers2.MskCommand),

		handlers2.NewPlace(stack, handlers2.AdminName),
		handlers2.NewPlace(stack, handlers2.SecurityName),
		handlers2.NewPlace(stack, handlers2.ArtName),
		handlers2.NewPlace(stack, handlers2.BarName),

		handlers2.NewBack(stack),
		handlers2.NewPhotoSend(stack, db, api, cfg),
		handlers2.NewTextSend(stack, db, api, cfg),
		handlers2.NewReset(stack),
	}
}

func setCommands(api *telebot.Bot) error {
	teleCommands := []telebot.Command{
		{
			Text:        "start",
			Description: "Начать сначала",
		},
		{
			Text:        handlers2.VlgCommand[1:],
			Description: "Волгоград",
		},
		{
			Text:        handlers2.MskCommand[1:],
			Description: "Москва",
		},
		{
			Text:        handlers2.KrdCommand[1:],
			Description: "Краснодар",
		},
	}

	return api.SetCommands(teleCommands)
}

func useHandle(api *telebot.Bot, stack *stack.Stack, handlersList []*handlers2.HandlerBase) {
	for _, handle := range handlersList {
		api.Handle(handle.Command(), handle.Handle)
	}

	start := handlers2.NewStart(stack)

	api.Handle(start.Command(), start.Handle)
}
