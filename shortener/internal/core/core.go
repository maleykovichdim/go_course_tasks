// Business logic for generating and storing short codes. High-level API.
package core

import (
	"shortener/internal/api"
	"shortener/internal/config"
	"sync"

	"github.com/rs/zerolog/log"
)

var pSize int
var pThreshold int
var sizeShortCode int

type Core struct {
	db    api.DBShortener
	cache api.CacheClient
	short api.Shortener

	pregenerated []*api.Link
	dbM          sync.RWMutex
	preM         sync.Mutex
}

// init initializes global variables by reading from the configuration file.
func init() {
	cfg := config.NewConfig() // Create a new instance of Config

	// Parse the configuration file and log errors if any occur
	if err := cfg.ParseConfig("../../config.txt"); err != nil {
		log.Error().Err(err).Msg("Error parsing config file")
	}

	// Assign configuration values to global variables if they exist
	if pSizeValue, exists := cfg.GetInt("pSize"); exists {
		pSize = pSizeValue
	}
	if pThresholdValue, exists := cfg.GetInt("pThreshold"); exists {
		pThreshold = pThresholdValue
	}
	if sizeShortCodeValue, exists := cfg.GetInt("sizeShortCode"); exists {
		sizeShortCode = sizeShortCodeValue
	}
}

// NewCore initializes a new Core instance with provided database, cache, and shortener interfaces.
func NewCore(dbs api.DBShortener, cc api.CacheClient, sh api.Shortener) *Core {
	c := Core{
		db:           dbs,
		cache:        cc,
		short:        sh,
		pregenerated: make([]*api.Link, 0, pSize*2),
	}
	return &c
}

// OpenDB opens the database connection. Logs errors if connection fails.
func (c *Core) OpenDB() error {
	if err := c.db.Open(); err != nil {
		log.Error().Err(err).Msg("Failed to open database")
		return err
	}
	log.Info().Msg("Database connection opened successfully")
	return nil
}

// CloseDB closes the database connection.
func (c *Core) CloseDB() {
	c.db.Close()
	log.Info().Msg("Database connection closed")
}

// service generates new short links and stores them if the pre-generated pool is below threshold.
func (c *Core) service(isZero bool) error {
	if len(c.pregenerated) <= pThreshold {
		fresh := make([]string, 0, pSize)
		for _ = range pSize {
			fresh = append(fresh, c.short.ShortenSimple(sizeShortCode))
		}
		c.dbM.Lock()
		mellow, err := c.db.PutEmptyShortLinks(fresh)
		c.dbM.Unlock()
		if err != nil {
			log.Error().Err(err).Msg("Database error when putting empty links")
			return err
		}

		if !isZero {
			c.preM.Lock()
		}
		c.pregenerated = append(c.pregenerated, mellow...)
		if !isZero {
			c.preM.Unlock()
		}
		log.Info().Msg("Successfully generated and stored new short links")
	}
	return nil
}

// CreateDB creates the necessary database structure for the shortener. Logs errors if creation fails.
func (c *Core) CreateDB() error {
	if err := c.db.CreateDbStructure(); err != nil {
		log.Error().Err(err).Msg("Failed to create database structure")
		return err
	}
	log.Info().Msg("Database structure created successfully")
	return nil
}

// GetShortCode provides a short code for the given destination, generating a new one if necessary.
func (c *Core) GetShortCode(destination string) (string, error) {
	log.Info().Msgf("Request to get short code for destination: %s", destination)

	c.preM.Lock()
	defer c.preM.Unlock()

	// Generate links if the pre-generated pool is empty
	if len(c.pregenerated) == 0 {
		if err := c.service(true); err != nil {
			log.Error().Err(err).Msg("Error during service generation")
			return "", err
		}
	}

	// Trigger background service if the pool size is below threshold
	if len(c.pregenerated) < pThreshold {
		go c.service(false)
	}

	i := len(c.pregenerated) - 1
	a := c.pregenerated[i]
	c.pregenerated = c.pregenerated[:i]

	a.LongUrl = destination

	c.dbM.Lock()
	err := c.db.UpdateLink(a)
	c.dbM.Unlock()
	if err != nil {
		log.Error().Err(err).Msg("Error updating link in database")
		return "", err
	}

	go c.cache.SaveLink(a.ShortCode, a.LongUrl)
	log.Info().Msgf("Short code %s generated for destination %s", a.ShortCode, destination)
	return a.ShortCode, nil
}

// GetLongCode retrieves the full URL for a given short code, sourcing from cache or database.
func (c *Core) GetLongCode(shortCode string) (string, error) {

	log.Info().Msgf("Request to get long code for short code: %s", shortCode)
	destination, err := c.cache.GetLink(shortCode)
	if err != nil {
		log.Warn().Err(err).Msg("Cache lookup failed")
		c.dbM.RLock()
		a, err := c.db.GetLongUrl(shortCode)
		c.dbM.RUnlock()
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving URL from database")
			return "", err
		}
		destination = a.LongUrl
		log.Warn().Msgf("Cache miss for short code: %s. URL retrieved from database.", shortCode)
		go c.cache.SaveLink(shortCode, a.LongUrl)
	}
	return destination, nil
}
