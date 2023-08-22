package main

import (
	"net/http"

	"github.com/lgustavopalmieri/challenge-Client-Server-API/server"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", server.HandleCurrentDollar)
	http.ListenAndServe(":8080", mux)
}
