package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
)

var db *sql.DB
var mu sync.Mutex

type RequestPayload struct {
	State         string  `json:"state"`
	Amount        float64 `json:"amount,string"`
	TransactionId string  `json:"transactionId"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	sourceType := r.Header.Get("Source-Type")
	if sourceType != "game" && sourceType != "server" && sourceType != "payment" {
		http.Error(w, "Invalid Source-Type", http.StatusBadRequest)
		return
	}

	// Process the request in a thread-safe manner
	mu.Lock()
	defer mu.Unlock()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var processed bool
	err = tx.QueryRow("SELECT processed FROM transactions WHERE transaction_id = $1", payload.TransactionId).Scan(&processed)
	if err == nil && processed {
		http.Error(w, "Transaction already processed", http.StatusBadRequest)
		return
	}

	var balance float64
	err = tx.QueryRow("SELECT balance FROM user_accounts WHERE id = 1").Scan(&balance)
	if err != nil {
		http.Error(w, "User account not found", http.StatusInternalServerError)
		return
	}

	newBalance := balance
	if payload.State == "win" {
		newBalance += payload.Amount
	} else if payload.State == "lost" && balance >= payload.Amount {
		newBalance -= payload.Amount
	} else {
		http.Error(w, "Invalid transaction", http.StatusBadRequest)
		return
	}

	_, err = tx.Exec("UPDATE user_accounts SET balance = $1 WHERE id = 1", newBalance)
	if err != nil {
		http.Error(w, "Balance update failed", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT INTO transactions (transaction_id, state, amount, source_type, user_account_id, processed) VALUES ($1, $2, $3, $4, $5, $6)",
		payload.TransactionId, payload.State, payload.Amount, sourceType, 1, true)
	if err != nil {
		http.Error(w, "Transaction recording failed", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	_, _ = fmt.Fprintf(w, "Transaction processed successfully")
}

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=mysecretpassword dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	http.HandleFunc("/your_url", handleRequest)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
