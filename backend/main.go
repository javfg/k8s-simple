// package main contains the main function
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
)

// Response represents an endpoint response.
type Response struct {
	Message string `json:"message"`
	Info    struct {
		DBString string `json:"message"`
		Hostname string `json:"hostname"`
	} `json:"info"`
}

func main() {
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle the request
		handleHi(w, r)
	})
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHi(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS requests
	if r.Method == "OPTIONS" {
		fmt.Println("OPTIONS /hi")
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	dbString := queryDBString(ctx)
	hostname := os.Getenv("HOSTNAME")

	response := Response{
		Message: "OK",
		Info: struct {
			DBString string `json:"message"`
			Hostname string `json:"hostname"`
		}{
			DBString: dbString,
			Hostname: hostname,
		},
	}

	fmt.Println("GET /hi")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func queryDBString(ctx context.Context) string {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@db:5432/mydb", username, password))

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	var dbString string
	err = conn.QueryRow(context.Background(), "SELECT message FROM mytable LIMIT 1").Scan(&dbString)
	if err != nil {
		log.Fatal(err)
	}

	return dbString
}
