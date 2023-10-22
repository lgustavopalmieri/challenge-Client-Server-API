package server

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	currentdollar "github.com/lgustavopalmieri/challenge-Client-Server-API/current-dollar"
)

func HandleServer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	usdbrlData, _ := HandleCurrentDollar(ctx, w, r)
	persistCurrentDollar(ctx, usdbrlData)
}

func HandleCurrentDollar(ctx context.Context, w http.ResponseWriter, r *http.Request) (*currentdollar.USDBRL, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	select {
	case <-ctxWithTimeout.Done():
		log.Println("API response timeout:", ctxWithTimeout.Err())
		http.Error(w, "API response timeout", http.StatusGatewayTimeout)
		return nil, ctxWithTimeout.Err()
	default:
		resp, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}

		var data map[string]currentdollar.USDBRL
		err = json.Unmarshal(responseBody, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil, err
		}

		usdbrlData, ok := data["USDBRL"]
		if !ok {
			http.Error(w, "Invalid response format", http.StatusInternalServerError)
			return nil, err
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usdbrlData)
		return &usdbrlData, nil
	}
}

var db *sql.DB

func persistCurrentDollar(ctx context.Context, usdbrlData *currentdollar.USDBRL) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	var err error
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.ExecContext(ctxWithTimeout, `
        CREATE TABLE IF NOT EXISTS usd_brl_data (
            code TEXT,
            codein TEXT,
            name TEXT,
            high TEXT,
            low TEXT,
            varBid TEXT,
            pctChange TEXT,
            bid TEXT,
            ask TEXT,
            timestamp TEXT,
            create_date TEXT
        )
    `)
	if err != nil {
		log.Fatal(err)
	}

	insertQuery := `
        INSERT INTO usd_brl_data (
            code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `

	_, err = db.ExecContext(ctxWithTimeout, insertQuery,
		usdbrlData.Code, usdbrlData.Codein, usdbrlData.Name, usdbrlData.High, usdbrlData.Low,
		usdbrlData.VarBid, usdbrlData.PctChange, usdbrlData.Bid, usdbrlData.Ask,
		usdbrlData.Timestamp, usdbrlData.CreateDate)
	if err != nil {
		select {
		case <-ctx.Done():
			log.Println("Persistence canceled:", ctx.Err())
		default:
			log.Println("Persistence error:", err)
		}
		return err
	}

	return nil
}
