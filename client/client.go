package client

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

func GetDollarFromServer() {
	http.HandleFunc("/", handleDollarFromServer)
	http.ListenAndServe(":8080", nil)
}

func handleDollarFromServer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}
