package main

import "log"

func main() {
	webService := new(WebService)
	listenPort := 9977
	if err := webService.Start(listenPort); err != nil {
		log.Println(err)
	}
}