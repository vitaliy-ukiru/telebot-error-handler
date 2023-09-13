package collections

import (
	te "github.com/vitaliy-ukiru/telebot-error-handler"
	tele "gopkg.in/telebot.v3"
)

type Filter struct {
	Filter  te.Filter
	Catcher te.Catcher
}

func NewFilter(filter te.Filter, catcher te.Catcher) *Filter {
	return &Filter{Filter: filter, Catcher: catcher}
}

func (f Filter) Catch(err error, ctx tele.Context) (matched bool) {
	return f.Filter(err) && f.Catcher.Catch(err, ctx)
}

func NewHandlerFilter(filter te.Filter, handler te.Handler) *Filter {
	return &Filter{Filter: filter, Catcher: te.Catch(handler)}
}
