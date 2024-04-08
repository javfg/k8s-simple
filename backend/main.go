// package main contains the main function
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
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

var (
	requestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "backend_requests_total",
		Help: "Total number of requests to the backend",
	})

	requestsGet = promauto.NewCounter(prometheus.CounterOpts{
		Name: "backend_requests_get",
		Help: "Total number of GET requests to the backend",
	})

	requestsOptions = promauto.NewCounter(prometheus.CounterOpts{
		Name: "backend_requests_options",
		Help: "Total number of OPTIONS requests to the backend",
	})
)

func main() {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()

	l.Info().Msg("starting server")
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
		handleHi(w, r, c, l)
	})

	http.Handle("/metrics", promhttp.Handler())

	fmt.Printf("Server is running on port 8080, env vars: %+v\n", c)
	l.Fatal().Err(http.ListenAndServe(":8080", nil)).Send()
}

func handleHi(w http.ResponseWriter, r *http.Request, c Config, l zerolog.Logger) {
	l.Info().Msgf("%s - %s: %s", r.Method, r.RemoteAddr, r.URL.Path)

	requestsTotal.Inc()

	// Handle OPTIONS requests
	if r.Method == "OPTIONS" {
		requestsOptions.Inc()
		w.WriteHeader(http.StatusOK)
		l.Info().Msg("OPTIONS request success")
		return
	}

	requestsGet.Inc()
	dbString, err := queryDBString(r.Context(), c)
	if err != nil {
		l.Error().Err(err).Msg("query db fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{
		Message: "OK",
		Info: Info{
			DBString: *dbString,
			Hostname: c.Host,
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		l.Error().Err(err).Msg("marshal response fail")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func queryDBString(ctx context.Context, c Config) (*string, error) {
	conn, err := pgx.Connect(ctx, fmt.Sprintf("postgresql://%s:%s@%s:5432/mydb", c.DBUsername, c.DBPassword, c.DBHost))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	var dbString *string
	err = conn.QueryRow(context.Background(), "SELECT message FROM mytable LIMIT 1").Scan(dbString)

	return dbString, err
}
