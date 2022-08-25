package handlers

import (
	"gopkg.in/telebot.v3"
)

type HandlerBase struct {
	name     string
	callback func(ctx telebot.Context) error
}

func (h *HandlerBase) Command() string {
	return h.name
}

func (h *HandlerBase) Handle(ctx telebot.Context) error {
	return h.callback(ctx)
}

var CityMenu = telebot.ReplyMarkup{
	OneTimeKeyboard: false,
	RemoveKeyboard:  true,
	//ReplyKeyboard: [][]telebot.ReplyButton{
	//	{
	//		telebot.ReplyButton{Text: "Волгоград"},
	//		telebot.ReplyButton{Text: "Москва"},
	//		telebot.ReplyButton{Text: "Краснодар"},
	//	},
	//},
}

var PlaceMenu = telebot.ReplyMarkup{
	OneTimeKeyboard: false,
	RemoveKeyboard:  true,
	ReplyKeyboard: [][]telebot.ReplyButton{
		{
			telebot.ReplyButton{Text: AdminName},
			telebot.ReplyButton{Text: SecurityName},
			telebot.ReplyButton{Text: BarName},
			telebot.ReplyButton{Text: ArtName},
		},
		{
			telebot.ReplyButton{Text: backName},
		},
	},
}

var SendMenu = telebot.ReplyMarkup{
	OneTimeKeyboard: false,
	ReplyKeyboard: [][]telebot.ReplyButton{
		{
			telebot.ReplyButton{Text: backName},
			telebot.ReplyButton{Text: resetName},
		},
	},
}
