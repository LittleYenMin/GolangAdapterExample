package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func (mu *Mux) lineBotCallbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	log.Println("Events: ", events, "Error: ", err)
	if err != nil {
		w.WriteHeader(400)
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			mu.do(*event)
		}
		log.Println("We Got line Event: ", event)
	}
}
func txetMessageHandler(t linebot.Event) {
	message := t.Message.(*linebot.TextMessage)
	_, err := bot.ReplyMessage(t.ReplyToken, linebot.NewTextMessage(message.Text)).Do()
	if err != nil {
		log.Println("Reply Error: ", err)
	}
}

func main() {
	mux := &Mux{
		make(map[reflect.Type]MessageHandler),
	}
	mux.setListener(&linebot.TextMessage{}, txetMessageHandler)
	var err error
	bot, err = linebot.New(os.Getenv("secret"), os.Getenv("token"))
	log.Println("Bot: ", bot, "Err: ", err)
	http.HandleFunc("/callback", mux.lineBotCallbackHandler)
	port := 8080
	fmt.Printf("Server running on port %d\n", port)
	listenPort := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
