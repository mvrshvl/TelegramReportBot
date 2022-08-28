package core

import (
	"TelegramBot/config"
	handlers2 "TelegramBot/core/handlers"
	"TelegramBot/core/stack"
	"fmt"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	API *telebot.Bot
}

func NewBot(cfg *config.Config) (*Bot, error) {
	api, err := telebot.NewBot(telebot.Settings{
		Token:   cfg.Token,
		Updates: 60,
		Poller:  nil,
	})
	if err != nil {
		return nil, fmt.Errorf("create bot: %w", err)
	}

	botStack := stack.New()

	handlersList := createHandlers(botStack, api, cfg)

	useHandle(api, botStack, handlersList)

	if err := setCommands(api); err != nil {
		return nil, fmt.Errorf("set commands: %w", err)
	}

	//_, err = telegram.NewClient(telegram.ClientConfig{
	//	// where to store session configuration. must be set
	//	SessionFile: "storage/session.json",
	//	// host address of mtproto server. Actually, it can be any mtproxy, not only official
	//	ServerHost: cfg.Server,
	//	// public keys file is path to file with public keys, which you must get from https://my.telegram.org
	//	PublicKeysFile:  cfg.Key,
	//	AppID:           cfg.AppID,   // app id, could be find at https://my.telegram.org
	//	AppHash:         cfg.AppHash, // app hash, could be find at https://my.telegram.org
	//	InitWarnChannel: true,        // if we want to get errors, otherwise, client.Warnings will be set nil
	//})
	//if err != nil {
	//	return nil, err
	//}

	//_, err = client.AuthSendCode("89020976661", int32(cfg.AppID), cfg.AppHash, &telegram.CodeSettings{})
	//if err != nil {
	//	return nil, err
	//}
	//
	//auth, err := client.AuthSignIn(
	//	phoneNumber,
	//	setCode.PhoneCodeHash,
	//	code,
	//)
	//
	//_, err = client.AuthImportBotAuthorization(1, int32(cfg.AppID), cfg.AppHash, cfg.Token)
	//if err != nil {
	//	return nil, err
	//}
	//
	//chat, err := client.GetChatInfoByHashLink(cfg.Channels["Волгоград"]["Охрана"].Link)
	//if err != nil {
	//	return nil, err
	//}
	//
	//channelSimpleData, ok := chat.(*telegram.Channel)
	//if !ok {
	//	return nil, fmt.Errorf("not a chan")
	//}
	//
	//msgs, err := client.MessagesGetHistory(&telegram.MessagesGetHistoryParams{
	//	Peer:       &telegram.InputPeerChannel{ChannelID: channelSimpleData.ID, AccessHash: channelSimpleData.AccessHash},
	//	OffsetID:   0,
	//	OffsetDate: 0,
	//	AddOffset:  0,
	//	Limit:      math.MaxInt32,
	//	MaxID:      math.MaxInt32,
	//	MinID:      0,
	//	Hash:       0,
	//})
	//if err != nil {
	//	return nil, err
	//}

	return &Bot{
		API: api,
	}, nil
}

func (bot *Bot) Run() {
	bot.API.Start()
}

func createHandlers(stack *stack.Stack, api *telebot.Bot, cfg *config.Config) []*handlers2.HandlerBase {
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
		handlers2.NewPhotoSend(stack, api, cfg),
		handlers2.NewTextSend(stack, api, cfg),
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
