package middleware

import (
	"time"

	telerr "github.com/vitaliy-ukiru/telebot-error-handler"
	tele "gopkg.in/telebot.v3"
)

// ChainMiddleware works as separated Catcher in chain.
// Will be executed by her place in chain.
type ChainMiddleware struct{}

func (ChainMiddleware) Catch(_ error, ctx tele.Context) bool {
	ctx.Set("chain_middleware_start", time.Now())
	return false // execute next catcher in chain
}

// WrapperMiddleware wraps one catcher.
// This middleware calls always before call next. For
// controller WrapperMiddleware+Catcher is one object.
//
// Also, controller implements Catcher, i.g. you can wrap all controller
// in this middleware. It's good for create collections of controllers.
//
// You can manage next function calls.
func WrapperMiddleware(next telerr.Catcher) telerr.Catcher {
	return telerr.CatcherFunc(func(err error, ctx tele.Context) bool {
		// for example handling only errors in channels.
		if ctx.Chat().Type == tele.ChatChannel {
			return next.Catch(err, ctx)
		}
		return false
	})
}
