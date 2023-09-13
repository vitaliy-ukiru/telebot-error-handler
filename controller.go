package telerr

import (
	tele "gopkg.in/telebot.v3"
)

type ErrorController struct {
	cases       []Catcher
	defaultCase Handler
}

func New(catchers ...Catcher) *ErrorController {
	ec := new(ErrorController)
	ec.setupCases(catchers)
	return ec
}

// Catcher filters and handling error.
//
// Catch method must return `true` if error is handled.
// If returns false controller will check others catchers.
type Catcher interface {
	Catch(err error, ctx tele.Context) (matched bool)
}

func (ec *ErrorController) Default(defaultCase Handler) {
	ec.defaultCase = defaultCase
}

type Default Handler

func (d Default) Catch(err error, ctx tele.Context) bool { d(err, ctx); return true }
func (d Default) defaultCase() Handler                   { return Handler(d) }

type defaultCaseCatcher interface {
	Catcher
	defaultCase() Handler
}

func (ec *ErrorController) setupCases(catchers []Catcher) bool {
	// capacity equals because in bad case (default case)
	result := make([]Catcher, 0, len(catchers))
	var existsDefault bool

	for _, catcher := range catchers {
		defaultCase, ok := catcher.(defaultCaseCatcher)
		if !ok {
			result = append(result, catcher)
			continue
		}
		if existsDefault {
			panic("default case already exists")
		}

		ec.defaultCase = defaultCase.defaultCase()
		existsDefault = true
	}
	ec.cases = result
	return existsDefault
}

// OnError method have signature like telebot.Bot.OnError
// You need assign this method as error handler in bot settings
//
//	tele.NewBot(tele.Settings{
//		// ...
//		OnError: errorController.OnError
//		// ...
//	})
func (ec *ErrorController) OnError(err error, ctx tele.Context) {
	ec.process(err, ctx, false)
}

// Catch method implements Catcher interface for flexible.
// You can use one ErrorController as component for other controller.
//
// For details see example.
func (ec *ErrorController) Catch(err error, ctx tele.Context) bool {
	return ec.process(err, ctx, true)
}

func (ec *ErrorController) process(err error, ctx tele.Context, catcherCall bool) bool {
	for _, catcher := range ec.cases {
		if catcher.Catch(err, ctx) {
			return true
		}
	}

	// if it calls as catcher should be ignored default case
	// because it breaks chain for others handlers.
	if !catcherCall && ec.defaultCase != nil {
		ec.defaultCase(err, ctx)
		return true
	}

	return false
}
