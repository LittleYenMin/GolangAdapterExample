package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

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
	msg := t.(TestMessage)
	fmt.Printf("%s is under control now!\n", msg.Name)
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
