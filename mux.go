package main

import "reflect"

type Message interface {
	message()
}

type TestMessage struct {
	Name string `json:"name"`
	Age  *int64 `json:"age"`
}

func (t TestMessage) message() {}

type (
	MessageHandler func(t Message)
	Mux            struct {
		handlerMap map[reflect.Type]MessageHandler
	}
)

func (mu *Mux) setListener(msg Message, callback MessageHandler) {
	msgT := reflect.TypeOf(msg)
	mu.handlerMap[msgT] = callback
}

func (mu *Mux) do(event Message) {
	msgT := reflect.TypeOf(event)
	handler, ok := mu.handlerMap[msgT]
	if ok {
		handler(event)
	}
}
