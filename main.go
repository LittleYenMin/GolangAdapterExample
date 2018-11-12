package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

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

func (mu *Mux) test(w http.ResponseWriter, r *http.Request) {
	input, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}

	var msg TestMessage
	if err := json.Unmarshal(input, &msg); err != nil {
		log.Printf("[JSON] %s\n", err)
	}
	fmt.Println(msg)
	mu.do(msg)
	response, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[JSON] %s\n", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func someMessageHandle(t Message) {
	fmt.Printf("%s is under control now!\n", t.(TestMessage).Name)
}

func main() {
	mux := &Mux{
		make(map[reflect.Type]MessageHandler),
	}
	mux.setListener(TestMessage{}, someMessageHandle)
	http.HandleFunc("/test", mux.test)
	port := 8080
	fmt.Printf("Server running on port %d\n", port)
	listenPort := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
