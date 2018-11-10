package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpToGoogle() {
	resp, err := http.Get("https://google.com/")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.StatusCode)
}
