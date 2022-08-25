package handlers

import (
	"TelegramBot/core/stack"
	"gopkg.in/telebot.v3"
)

const (
	startName = "/start"

	backName  = "Назад"
	resetName = "Выйти из отправки отчёта"

	chooseCity = "Выберите город"
)

func NewStart(stack *stack.Stack) *HandlerBase {
	return &HandlerBase{
		name: startName,
		callback: func(ctx telebot.Context) error {
			return reset(ctx, chooseCity, stack)
		},
	}
}

func NewBack(stack *stack.Stack) *HandlerBase {
	return &HandlerBase{
		name: backName,
		callback: func(ctx telebot.Context) error {
			query, ok := stack.Get(ctx.Message().Chat.ID)

			if !ok || len(query.Place) == 0 {
				return reset(ctx, chooseCity, stack)
			}

			return cityCallback(ctx, stack, query.City)
		},
	}
}

func NewReset(stack *stack.Stack) *HandlerBase {
	return &HandlerBase{
		name: resetName,
		callback: func(ctx telebot.Context) error {
			return reset(ctx, chooseCity, stack)
		},
	}
}

func reset(ctx telebot.Context, msg string, stack *stack.Stack) error {
	stack.Delete(ctx.Message().Chat.ID)

	return ctx.Send(msg, telebot.RemoveKeyboard)
}
