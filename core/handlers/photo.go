package handlers

import (
	"TelegramBot/config"
	"TelegramBot/core/database"
	"TelegramBot/core/stack"
	"TelegramBot/tgerror"
	"fmt"
	scribble "github.com/nanobox-io/golang-scribble"
	"gopkg.in/telebot.v3"
	"log"
	"strconv"
)

const (
	errGetCity  = tgerror.TelegramError("city not found in config")
	errGetPlace = tgerror.TelegramError("place not found in config")

	fmtMsgErr = "Произошла ошибка (%s), начните сначала"
)

func NewPhotoSend(stack *stack.Stack, db *scribble.Driver, api *telebot.Bot, cfg *config.Config) *HandlerBase {
	return &HandlerBase{
		name:     telebot.OnPhoto,
		callback: sendCallback(stack, db, api, cfg),
	}
}

func NewTextSend(stack *stack.Stack, db *scribble.Driver, api *telebot.Bot, cfg *config.Config) *HandlerBase {
	return &HandlerBase{
		name:     telebot.OnText,
		callback: sendCallback(stack, db, api, cfg),
	}
}

func sendCallback(stack *stack.Stack, db *scribble.Driver, api *telebot.Bot, cfg *config.Config) func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		query, ok := stack.Get(ctx.Message().Chat.ID)

		if !ok || len(query.City) == 0 {
			stack.Delete(ctx.Message().Chat.ID)

			return ctx.Send(chooseCity)
		}

		if len(query.Place) == 0 {
			return ctx.Send("Выберите один из вариантов", &PlaceMenu)
		}

		chatID, err := getChatID(cfg, query)
		if err != nil {
			log.Println("get chat id error:", err)
			return reset(ctx, fmt.Sprintf(fmtMsgErr, err.Error()), stack)
		}

		_, err = api.Forward(chatID, ctx.Message())
		if err != nil {
			log.Println("send picture error:", err)
			return reset(ctx, fmt.Sprintf(fmtMsgErr, err.Error()), stack)
		}

		var fileID string

		if ctx.Message().Media() != nil {
			fileID = ctx.Message().Media().MediaFile().FileID
		}

		err = db.Write(query.City, fmt.Sprintf("%s-%d", query.Place, ctx.Message().ID), database.Message{
			Timestamp: ctx.Message().Time(),
			Text:      ctx.Message().Text,
			Media:     fileID,
			Place:     query.Place,
			From:      ctx.Message().Sender.Username,
		})
		if err != nil {
			log.Println("save message error:", err)
			return reset(ctx, fmt.Sprintf(fmtMsgErr, err.Error()), stack)
		}

		return ctx.Send("Отправлено")
	}
}

func getChatID(cfg *config.Config, query *stack.Query) (telebot.ChatID, error) {
	places, ok := cfg.Channels[query.City]
	if !ok {
		return 0, errGetCity
	}

	textChatID, ok := places[query.Place]
	if !ok {
		return 0, errGetPlace
	}

	id, err := strconv.Atoi(textChatID.ID)
	if err != nil {
		return 0, err
	}

	return telebot.ChatID(int64(id)), nil
}
