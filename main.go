package main

import (
	"log"
	"net/http"

	"github.com/lgustavopalmieri/challenge-Client-Server-API/client"
	"github.com/lgustavopalmieri/challenge-Client-Server-API/server"
)

func main() {
	muxServer := http.NewServeMux()
	muxServer.HandleFunc("/cotacao", server.HandleServer)
	go func() {
		log.Println("server.go started at :8080")
		errServer := http.ListenAndServe(":8080", muxServer)
		if errServer != nil {
			log.Fatalf("error starting server.go: %s", errServer)
		}
	}()

	muxClient := http.NewServeMux()
	muxClient.HandleFunc("/", client.HandleClient)
	log.Println("client.go started at :8081")
	errClient := http.ListenAndServe(":8081", muxClient)
	if errClient != nil {
		log.Fatalf("error starting client.go: %s", errClient)
	}
}
