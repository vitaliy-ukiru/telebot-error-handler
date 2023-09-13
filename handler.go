package telerr

import tele "gopkg.in/telebot.v3"

type TelebotErrorHandler = func(err error, ctx tele.Context)

type Handler TelebotErrorHandler
