// creation of a database structure with the required table
package pgsql

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// CreateDbStructure initializes the database structure
func (db_ *DB) CreateDbStructure() error {
	log.Info().Msg("Start database creation")

	var err error
	db_.pool, err = pgxpool.New(db_.ctx, db_.connStrStart)
	if err != nil {
		log.Error().Err(err).Msg("Error opening a connection pool")
		return fmt.Errorf("error opening a connection pool: %w", err) // Error wrapping
	}

	// Drop the existing database
	if _, err = db_.pool.Exec(db_.ctx, "DROP DATABASE IF EXISTS "+DBName); err != nil {
		log.Error().Err(err).Msg("Failed to drop database")
		return fmt.Errorf("failed to drop database: %w", err) // Error wrapping
	}

	// Create the new database
	if _, err = db_.pool.Exec(db_.ctx, "CREATE DATABASE "+DBName); err != nil {
		log.Error().Err(err).Msg("Failed to create database")
		return fmt.Errorf("failed to create database: %w", err) // Error wrapping
	}

	// Close the current pool to switch to a new one for the newly created database
	db_.pool.Close()

	// Open a new connection pool for the new database
	db_.pool, err = pgxpool.New(db_.ctx, db_.connStr) // Open a new DB
	if err != nil {
		log.Error().Err(err).Msg("Error opening a connection pool to the new database")
		return fmt.Errorf("error opening a connection pool to the new database: %w", err) // Error wrapping
	}
	defer db_.pool.Close() // Ensure the pool is closed when the function exits

	sqlScript := `
    DROP TABLE IF EXISTS links;

    CREATE TABLE links (
        id SERIAL PRIMARY KEY,
        short_code VARCHAR(13) NOT NULL UNIQUE,
        long_url TEXT NOT NULL DEFAULT 'EMPTY',
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

    CREATE INDEX short_code_idx ON links USING HASH (short_code);
    `

	// Execute the SQL script to create necessary tables
	if _, err = db_.pool.Exec(db_.ctx, sqlScript); err != nil {
		log.Error().Err(err).Msg("Failed to execute SQL script for table creation")
		return fmt.Errorf("failed to execute SQL script: %w", err) // Error wrapping
	}

	log.Info().Msg("Database creation completed successfully")
	return nil
}
