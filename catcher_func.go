package telerr

import tele "gopkg.in/telebot.v3"

type CatcherFunc func(err error, ctx tele.Context) (matched bool)

func (c CatcherFunc) Catch(err error, ctx tele.Context) bool {
	return c(err, ctx)
}

func Ignore(_ error, _ tele.Context) {}

func Catch(h Handler) Catcher {
	return CatcherFunc(func(err error, ctx tele.Context) bool {
		h(err, ctx)
		return true
	})
}

func Handleable(catcher Catcher) Handler {
	return func(err error, ctx tele.Context) {
		catcher.Catch(err, ctx)
	}
}
