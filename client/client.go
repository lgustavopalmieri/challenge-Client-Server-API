package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotationResponse struct {
	Bid string `json:"bid"`
}

func HandleClient(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}

	client := &http.Client{
		Timeout: 300 * time.Millisecond,
	}

	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out:", err)
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		} else {
			http.Error(w, "Error during request", http.StatusInternalServerError)
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	var cotation CotationResponse
	err = json.Unmarshal(body, &cotation)
	if err != nil {
		http.Error(w, "Error unmarshaling JSON", http.StatusInternalServerError)
		return
	}

	bid := cotation.Bid

	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", bid)), 0644)
	if err != nil {
		http.Error(w, "Error writing file cotacao.txt", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Dólar: %s", bid)
}
