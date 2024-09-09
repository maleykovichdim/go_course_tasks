package memcache_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"shortener/internal/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// BaseCacheUrl = "http://localhost:8080"
	// BaseCacheUrl    = "http://memcache_service:8081"
	PathSave        = "/save"
	PathGet         = "/get"
	ContentTypeJSON = "application/json"
)

var BaseCacheUrl string

// LinkPair represents the structure of a link pair
type LinkPair struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

// Client provides methods to interact with the memcache microservice
type Client struct {
	BaseCacheUrl string
}

// Init initializes global variables by reading from the configuration file
func init() {
	cfg := config.NewConfig() // Create a new instance of Config

	// Parse the configuration file
	if err := cfg.ParseConfig("../../config.txt"); err != nil {
		log.Error().Err(err).Msg("Error parsing config")
	}

	// Assign the values from the config to global variables
	BaseCacheUrl, _ = cfg.GetString("BaseCacheUrl")
}

// NewClient creates a new client instance
func New() *Client {
	// Initialize zerolog to output in a human-readable format
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	return &Client{BaseCacheUrl: BaseCacheUrl}
}

// SaveLink sends a POST request to save a short and long URL pair in the service
func (c *Client) SaveLink(shortURL, longURL string) error {
	link := LinkPair{ShortURL: shortURL, LongURL: longURL}
	data, err := json.Marshal(link)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal link pair")
		return fmt.Errorf("failed to marshal link pair: %w", err)
	}

	resp, err := http.Post(c.BaseCacheUrl+PathSave, ContentTypeJSON, bytes.NewBuffer(data))
	if err != nil {
		log.Error().Err(err).Msg("Failed to send request to save link")
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Error().Msgf("Failed to save link: %s", string(body))
		return fmt.Errorf("failed to save link: %s", string(body))
	}

	log.Info().Str("short_url", shortURL).Str("long_url", longURL).Msg("Link saved successfully")
	return nil
}

// GetLink sends a GET request to retrieve the long URL for a short URL
func (c *Client) GetLink(shortURL string) (string, error) {
	resp, err := http.Get(c.BaseCacheUrl + PathGet + "?short_url=" + shortURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send request to get link")
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Error().Str("short_url", shortURL).Msgf("Failed to get link: %s", string(body))
		return "", fmt.Errorf("failed to get link: %s", string(body))
	}

	var link LinkPair
	if err := json.NewDecoder(resp.Body).Decode(&link); err != nil {
		log.Error().Err(err).Msg("Failed to decode response")
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	log.Info().Str("short_url", shortURL).Str("long_url", link.LongURL).Msg("Link retrieved successfully")
	return link.LongURL, nil
}
