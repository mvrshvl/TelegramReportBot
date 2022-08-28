package handlers

import (
	"TelegramBot/core/database"
	"TelegramBot/core/stack"
	"gopkg.in/telebot.v3"
)

type City string

const (
	VlgName = "Волгоград"
	MskName = "Москва"
	KrdName = "Краснодар"

	KrdCommand = "/krasnodar"
	VlgCommand = "/volgograd"
	MskCommand = "/moscow"

	choosePlace = "Выберите один из вариантов"
)

func NewCity(stack *stack.Stack, name string) *HandlerBase {
	return &HandlerBase{
		name: name,
		callback: func(ctx telebot.Context) error {
			return cityCallback(ctx, stack, name)
		},
	}
}

func cityCallback(ctx telebot.Context, stack *stack.Stack, city string) error {
	stack.Replace(ctx.Message().Chat.ID, translate(city), "")

	return ctx.Reply(choosePlace, &PlaceMenu)
}

func translate(city string) string {
	switch city {
	case VlgCommand:
		return database.TableVLG
	case KrdCommand:
		return database.TableKRD
	case MskCommand:
		return database.TableKRD
	default:
		return city
	}
}
