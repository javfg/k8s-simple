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

// Info represents the contents of a response.
type Info struct {
	DBString string `json:"db_string"`
	Hostname string `json:"hostname"`
}

// Response represents an endpoint response.
type Response struct {
	Message string `json:"message"`
	Info    Info   `json:"info"`
}

// Config represents the configuration of the application.
type Config struct {
	Host       string
	DBHost     string
	DBUsername string
	DBPassword string
}

func main() {
	c := Config{
		Host:       os.Getenv("HOSTNAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBUsername: os.Getenv("POSTGRES_USER"),
		DBPassword: os.Getenv("POSTGRES_PASSWORD"),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle the request
		handleHi(w, r, c)
	})

	fmt.Printf("Server is running on port 8080, env vars: %+v\n", c)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHi(w http.ResponseWriter, r *http.Request, c Config) {
	log.Printf("%s - %s: %s\n", r.Method, r.RemoteAddr, r.URL.Path)

	// Handle OPTIONS requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()
	dbString := queryDBString(ctx, c)

	response := Response{
		Message: "OK",
		Info: Info{
			DBString: dbString,
			Hostname: c.Host,
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func queryDBString(ctx context.Context, c Config) string {
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s:5432/mydb", c.DBUsername, c.DBPassword, c.DBHost))

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
