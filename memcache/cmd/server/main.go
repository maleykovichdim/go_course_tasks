package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"memcache/internal/config"
	rw "memcache/internal/redis"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var ListenServePort string
var redisURL string
var password string
var db int

var (
	// Counter for tracking the total number of API requests
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of API requests",
		},
		[]string{"method"}, // Label to differentiate between the methods
	)

	// Histogram for tracking request durations
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Histogram of response latency (seconds) for API requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"}, // Label to differentiate between the methods
	)
)

func init() {
	// Register metrics with Prometheus
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)

	cfg := config.NewConfig() // Create a new instance of Config
	// Parse the configuration file
	if err := cfg.ParseConfig("../../config.txt"); err != nil {
		log.Error().Err(err).Msg("Error parsing config")
	}

	// Assign the values from the config to global variables
	ListenServePort, _ = cfg.GetString("ListenServePort")
	redisURL, _ = cfg.GetString("redisURL")
	password, _ = cfg.GetString("password")

	if timeout, exists := cfg.GetInt("db"); exists {
		db = (timeout)
	}

}

type Request struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

func main() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Initialize Redis cache
	cache := rw.NewRedisCache(redisURL, password, db)

	// Expose metrics endpoint for Prometheus
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Start the timer
		defer func() {
			// Record metrics
			requestDuration.WithLabelValues("save").Observe(time.Since(start).Seconds())
			requestCounter.WithLabelValues("save").Inc()
		}()

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var request Request

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Error().Err(err).Msg("Failed to decode JSON")
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := cache.Save(request.ShortURL, request.LongURL); err != nil {
			log.Error().Err(err).Msg("Failed to save URL to cache")
			http.Error(w, "Failed to save URL", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(request)
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Start the timer
		defer func() {
			// Record metrics
			requestDuration.WithLabelValues("get").Observe(time.Since(start).Seconds())
			requestCounter.WithLabelValues("get").Inc()
		}()

		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		shortURL := r.URL.Query().Get("short_url")
		longURL, err := cache.Get(shortURL)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get URL from cache")
			http.Error(w, "Failed to retrieve URL", http.StatusInternalServerError)
			return
		}

		if longURL == "" {
			http.Error(w, "Short URL not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(Request{ShortURL: shortURL, LongURL: longURL})
	})

	log.Info().Msg("Starting HTTP server on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal().Err(err).Msg("Failed to start memcache server")
	}
}
