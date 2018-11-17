package main

import (
	"reflect"

	"github.com/line/line-bot-sdk-go/linebot"
)

type (
	MessageHandler func(t linebot.Event)
	Mux            struct {
		handlerMap map[reflect.Type]MessageHandler
	}
)

func (mu *Mux) setListener(msg linebot.Message, callback MessageHandler) {
	msgT := reflect.TypeOf(msg)
	mu.handlerMap[msgT] = callback
}

func (mu *Mux) do(event linebot.Event) {
	msgT := reflect.TypeOf(event.Message)
	handler, ok := mu.handlerMap[msgT]
	if ok {
		handler(event)
	}
}
