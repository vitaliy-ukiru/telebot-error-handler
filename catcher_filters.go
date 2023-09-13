package telerr

import (
	"errors"

	tele "gopkg.in/telebot.v3"
)

type Filter func(err error) bool

func Case(filter Filter, handler Handler) Catcher {
	return CatcherFunc(func(err error, ctx tele.Context) bool {
		if !filter(err) {
			return false
		}

		handler(err, ctx)
		return true
	})
}

func Is(target error, handler Handler) Catcher {
	return CatcherFunc(func(err error, ctx tele.Context) bool {
		if !errors.Is(err, target) {
			return false
		}

		handler(err, ctx)
		return true
	})
}

func As[E error](handler func(err E, ctx tele.Context)) Catcher {
	return CatcherFunc(func(err error, ctx tele.Context) bool {
		var e E
		if errors.As(err, &e) {
			handler(e, ctx)
			return true
		}
		return false
	})
}

type IsForbiddenError Handler

func (i IsForbiddenError) Catch(err error, ctx tele.Context) bool {
	var e *tele.Error
	if !errors.As(err, &e) {
		return false
	}

	if e.Code != 403 {
		return false
	}

	i(err, ctx)
	return true
}

type OnMessage struct {
	Next Catcher
}

func (om OnMessage) Catch(err error, ctx tele.Context) bool {
	return ctx.Update().Message != nil && om.Next.Catch(err, ctx)
}
