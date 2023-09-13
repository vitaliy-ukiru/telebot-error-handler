package telerr_test

import (
	"log"

	_ "github.com/stretchr/testify/assert"
	te "github.com/vitaliy-ukiru/telebot-error-handler"
	"github.com/vitaliy-ukiru/telebot-error-handler/filters"
	tb "gopkg.in/telebot.v3"
)

func ExampleErrorController_Catch() {
	_ = te.New(
		// ignoring errors
		te.New(
			te.Case(
				filters.Any(
					filters.Is(tb.ErrMessageNotModified),
					filters.Is(tb.ErrSameMessageContent),
				),
				te.Ignore,
			),

			te.Default(func(err error, ctx tb.Context) {
				panic("this handler must be ignored and not called")
			}),
		),
		te.Is(tb.ErrKickedFromGroup, func(err error, ctx tb.Context) {

		}),

		te.IsForbiddenError(func(err error, ctx tb.Context) {
			log.Printf("forbidden access error")
		}),

		te.Default(func(err error, ctx tb.Context) {

		}),
	)
}
