package handlers

import (
	"TelegramBot/core/stack"
	"gopkg.in/telebot.v3"
)

const (
	AdminName    = "Админ"
	ArtName      = "Арт"
	BarName      = "Бар"
	SecurityName = "Охрана"
)

func NewPlace(stack *stack.Stack, name string) *HandlerBase {
	return &HandlerBase{
		name: name,
		callback: func(ctx telebot.Context) error {
			return placeCallback(ctx, stack, name)
		},
	}
}

func placeCallback(ctx telebot.Context, stack *stack.Stack, place string) error {
	query, ok := stack.Get(ctx.Message().Chat.ID)

	if !ok || len(query.City) == 0 {
		stack.Delete(ctx.Message().Chat.ID)

		return ctx.Send("Выберите город", telebot.RemoveKeyboard)
	}

	stack.Replace(ctx.Message().Chat.ID, query.City, place)

	return ctx.Send("Отправьте фото", &SendMenu)
}
