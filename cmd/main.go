package main

import (
	"errors"
	"log"

	te "github.com/vitaliy-ukiru/telebot-error-handler"
	tele "gopkg.in/telebot.v3"
)

func Logger(err error, ctx tele.Context) {
	if ctx == nil {
		log.Printf("error from handler: %v", err)
		return
	}

	log.Printf("error from handler: chat_id=%d, user_id=%d %v", ctx.Chat().ID, ctx.Sender().ID, err)
}

func main() {
	ec := te.New(
		te.Is(tele.ErrBlockedByUser, func(err error, ctx tele.Context) {
			log.Printf("user(%d) block bot", ctx.Sender().ID)
		}),
		te.IsForbiddenError(func(err error, ctx tele.Context) {
			log.Printf("bot cannot send message to user")
		}),
		te.Default(Logger), // as middleware
	)

	bot, _ := tele.NewBot(tele.Settings{
		OnError:     ec.OnError,
		Synchronous: true,
		Offline:     true,
	})

	u := tele.Update{
		ID: 12,
		Message: &tele.Message{
			ID:     2,
			Sender: &tele.User{ID: 34},
			Chat:   &tele.Chat{ID: 34, Type: "private"},
		},
	}

	bot.OnError(errors.New("custom error"), bot.NewContext(u))
}
