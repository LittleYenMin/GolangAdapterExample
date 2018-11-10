package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Message interface {
	message()
}

type TestMessage struct {
	Name string `json:"name"`
	Age  *int64 `json:"age"`
}

func (t TestMessage) message() {}

type MessageHandler struct {
	name    string
	handler Handler
}

type Handler interface {
	handle(t Message)
}

type HandlerAdapt func(t Message)

func (h HandlerAdapt) handle(t Message) {
	h(t)
}

type Mux struct {
	m map[string]MessageHandler
}

func (mu *Mux) setListener(msg Message, callback Handler) {
	name := fmt.Sprintf("%T", msg)
	fmt.Println("ListenerName is", name)
	mu.m[name] = MessageHandler{
		name,
		callback,
	}
}

func (mu *Mux) do(event Message) {
	name := fmt.Sprintf("%T", event)
	fmt.Println("doName is", name)
	handler, ok := mu.m[name]
	if ok {
		handler.handler.handle(event)
	}
}

func someMessageHandle(t Message) {
	fmt.Println(t.(TestMessage).Name, "is under control now!")
}

func test(w http.ResponseWriter, r *http.Request) {

	bytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	var msg TestMessage
	error := json.Unmarshal(bytes, &msg)
	if error != nil {
		log.Fatalln("JSON PARSE ERROR")
	}
	fmt.Println(msg)
	mux.do(msg)
	outputBytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalln("CONVERT JSON ERROR")
	}
	w.Header().Set("content-type", "application/json")
	w.Write(outputBytes)
}

var mux Mux

func main() {
	mux = Mux{
		make(map[string]MessageHandler),
	}
	mux.setListener(TestMessage{}, HandlerAdapt(someMessageHandle))
	http.HandleFunc("/test", test)
	port := 8080
	fmt.Println("Server running on port", port)
	listenPort := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
