// initialisation of global variables for DB-configuration from config.txt
package pgsql

import (
	"shortener/internal/config"

	"github.com/rs/zerolog/log"
)

// Global variables for configuration
var (
	DBUser        string
	DBPassword    string
	DBHost        string
	DBPort        string
	DBName        string
	DBNameStart   string
	CtxTimeoutSec int
)

// Init initializes global variables by reading from the configuration file
func init() {
	cfg := config.NewConfig() // Create a new instance of Config

	// Parse the configuration file
	if err := cfg.ParseConfig("../../config.txt"); err != nil {
		log.Error().Err(err).Msg("Error parsing config")
	}

	// Assign the values from the config to global variables
	DBUser, _ = cfg.GetString("DBUser")
	DBPassword, _ = cfg.GetString("DBPassword")
	DBHost, _ = cfg.GetString("DBHost")
	DBPort, _ = cfg.GetString("DBPort")
	DBName, _ = cfg.GetString("DBName")
	DBNameStart, _ = cfg.GetString("DBNameStart")

	if timeout, exists := cfg.GetInt("CtxTimeoutSec"); exists {
		CtxTimeoutSec = (timeout)
	}
}
