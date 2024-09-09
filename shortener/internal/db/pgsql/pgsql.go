// DB API implementation for PostgreSQL database
package pgsql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// DB represents the database connection pool and configuration
type DB struct {
	pool         *pgxpool.Pool
	connStr      string
	connStrStart string
	ctx          context.Context
	CtxCancel    context.CancelFunc // Use defer ctxCancel() in main function
}

// New initializes a new DB instance with connection strings
func New() *DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DBUser, DBPassword, DBHost, DBPort, DBName)
	connStrStart := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", DBUser, DBPassword, DBHost, DBPort, DBNameStart)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(CtxTimeoutSec)*time.Second) //TODO: Implement context cancellation
	ctx := context.Background()

	return &DB{
		connStr:      connStr,
		connStrStart: connStrStart,
		ctx:          ctx,
		CtxCancel:    nil, // TODO: Cancel context on main function exit
	}
}

// Open opens a connection to the PostgreSQL database
func (db_ *DB) Open() error {
	log.Info().Msg("Opening connection to PostgreSQL database")

	var err error
	db_.pool, err = pgxpool.New(db_.ctx, db_.connStr)
	if err != nil {
		log.Error().Err(err).Msg("Failed to open a connection pool")
		return fmt.Errorf("failed to open a connection pool: %w", err) // Error wrapping
	}

	log.Info().Msg("Connected to PostgreSQL database")
	return nil
}

// Close the connection pool to the PostgreSQL database
func (db_ *DB) Close() {
	log.Info().Msg("Closing connection to PostgreSQL database")
	db_.pool.Close()
}
