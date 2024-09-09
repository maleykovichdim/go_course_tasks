package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

// Config holds the parsed configuration parameters in maps
type Config struct {
	stringParams map[string]string
	intParams    map[string]int
}

// NewConfig initializes a new Config instance
func NewConfig() *Config {
	return &Config{
		stringParams: make(map[string]string),
		intParams:    make(map[string]int),
	}
}

// ParseConfig reads a file and populates the Config structure with parameters
func (c *Config) ParseConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Err(err).Msg("could not open config file")
		return fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			log.Error().Str("line", line).Msg("invalid line in config file")
			return fmt.Errorf("invalid line in config file: %s", line)
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		// Check if the value is a quoted string
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			// Remove surrounding quotes from the value
			value = strings.Trim(value, "\"")
			c.stringParams[key] = value // Store as string
			//log.Info().Str("key", key).Str("value", value).Msg("Parsed string parameter")
		} else {
			// Try to convert the value to an integer
			intValue, err := strconv.Atoi(value)
			if err != nil {
				log.Error().Err(err).Str("value", value).Msg("invalid integer value")
				return fmt.Errorf("invalid integer value for key %s: %w", key, err)
			}
			c.intParams[key] = intValue // Store as int if conversion is successful
			//log.Info().Str("key", key).Int("value", intValue).Msg("Parsed integer parameter")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("error reading config file")
		return fmt.Errorf("error reading config file: %w", err)
	}

	return nil
}

// GetString retrieves a string parameter by key
func (c *Config) GetString(key string) (string, bool) {
	value, exists := c.stringParams[key]
	return value, exists
}

// GetInt retrieves an integer parameter by key
func (c *Config) GetInt(key string) (int, bool) {
	value, exists := c.intParams[key]
	return value, exists
}
