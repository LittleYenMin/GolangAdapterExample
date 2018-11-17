package main

import (
	"fmt"
	"log"
	"net/http"
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
	bot, err = linebot.New("004a55e347fd117cd562e75ff78e0721", "oOrido9o+9oKAlmWV7OzDglYNIxwe+mOBjU2N+9aORiVRteq1qBx4jCMR5LZkEf/19ojH0KbwMrjLnKu8tWrvP2Ke27ul0AX+Tr7LHWYDrx8wN7sq5saY0Rod3EpCe17A7ZmO8LCGa5CB/TsIIM5KgdB04t89/1O/w1cDnyilFU=")
	log.Println("Bot: ", bot, "Err: ", err)
	http.HandleFunc("/callback", mux.lineBotCallbackHandler)
	port := 8080
	fmt.Printf("Server running on port %d\n", port)
	listenPort := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
