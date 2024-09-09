package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"shortener/internal/config"
	"shortener/internal/core"
	"shortener/internal/db/pgsql"
	"shortener/internal/memcache_client"
	"shortener/internal/urls"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Define Prometheus metrics
var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "url_shortener_requests_total",
			Help: "Total number of requests processed by the URL shortener",
		},
		[]string{"method"}, // Method label for differentiating between types of requests
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "url_shortener_request_duration_seconds",
			Help:    "Histogram for the duration of requests in the shortener",
			Buckets: prometheus.DefBuckets, // Default buckets for request durations
		},
		[]string{"method"},
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
	hostShortCode, _ = cfg.GetString("hostShortCode")
	ListenServePort, _ = cfg.GetString("ListenServePort")
}

var c *core.Core
var hostShortCode string
var ListenServePort string

func main() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	c = core.NewCore(pgsql.New(), memcache_client.New(), urls.New())
	c.CreateDB()
	err := c.OpenDB()
	if err != nil {
		log.Error().Msg("Failed to create DB connection")
		os.Exit(1)
	}
	defer c.CloseDB()

	router := mux.NewRouter()

	// Metrics endpoint
	router.Handle("/metrics", promhttp.Handler()) // Register the Prometheus metrics handler

	// Create short URL endpoint
	router.HandleFunc("/", handlePost).Methods(http.MethodPost)
	router.HandleFunc("/{shortCode}", handleGet).Methods(http.MethodGet)

	log.Info().Msg("Starting URL shortener server on :8080")

	// Start the HTTP server
	if err := http.ListenAndServe(":"+ListenServePort, router); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}

}

// handlePost processes POST requests at "/"
func handlePost(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("POST request received")
	start := time.Now()

	var data struct {
		Destination string `json:"destination"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil || data.Destination == "" {
		log.Error().Err(err).Msg("Failed to decode JSON or missing destination")
		http.Error(w, "Invalid JSON or missing destination", http.StatusBadRequest)
		recordMetrics("POST", start)
		return
	}

	shortCode, err := c.GetShortCode(data.Destination)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get short code fron Core: ")
		http.Error(w, "Failed to get short code fron Core", http.StatusBadRequest)
		recordMetrics("POST", start)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", hostShortCode, shortCode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ShortURL    string `json:"shortUrl"`
		Destination string `json:"destination"`
	}{
		ShortURL:    shortURL,
		Destination: data.Destination,
	})
	log.Info().Str("destination", data.Destination).Str("shortUrl", shortURL).Msg("POST request handled")
	recordMetrics("POST", start)
}

// handleGet processes GET requests
func handleGet(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GET request received")
	start := time.Now()
	vars := mux.Vars(r)
	shortCode := vars["shortCode"] // Extract the short code from the URL path

	url, err := c.GetLongCode(shortCode)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get Long URL fron Core ")
		http.Error(w, "Failed to get Long URL fron Core", http.StatusBadRequest)
		recordMetrics("GET", start)
		return
	}

	log.Info().Msg("GET request handled. Redirecting to: " + url)
	http.Redirect(w, r, url, http.StatusFound)
	recordMetrics("GET", start)
}

// recordMetrics updates the Prometheus metrics with request details
func recordMetrics(method string, start time.Time) {
	duration := time.Since(start).Seconds()
	requestCounter.WithLabelValues(method).Inc()
	requestDuration.WithLabelValues(method).Observe(duration)
}
