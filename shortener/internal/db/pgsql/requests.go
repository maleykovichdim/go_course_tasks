package pgsql

import (
	"errors"
	"fmt"
	api "shortener/internal/api"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

func (db_ *DB) PutEmptyShortLinks(shortLinks []string) ([]*api.Link, error) {

	log.Info().Msg("PutEmptyShortLink: Start to insert new short links in the postgress database")

	ctx := db_.ctx

	tx, err := db_.pool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	// cancel transaction in case of error
	defer tx.Rollback(ctx)
	// batch request
	var batch = &pgx.Batch{}

	for _, link := range shortLinks {
		batch.Queue(
			"INSERT INTO links (short_code) VALUES ($1) ON CONFLICT (short_code) DO NOTHING", link,
		)
	}

	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		log.Error().Err(err).Msg("Error while closing the transaction batch")
		return nil, fmt.Errorf("failed to close transaction batch: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		log.Error().Err(err).Msg("Error committing the transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	log.Info().Msg("New empty short links inserted in the postgress database")

	query := `
		SELECT *
		FROM links l
		WHERE l.long_url = 'EMPTY'
	`
	rows, err := db_.pool.Query(ctx, query)
	defer rows.Close()
	if err != nil {
		log.Error().Err(err).Msg("Error executing the query to retrieve empty links")
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var result []*api.Link

	for rows.Next() {
		var link api.Link
		err = rows.Scan(&link.ID, &link.ShortCode, &link.LongUrl, &link.CreatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, &link)
	}

	if err = rows.Err(); err != nil { // Check for errors encountered during iteration
		log.Error().Err(err).Msg("Error during rows iteration")
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	log.Info().Msg("Successfully retrieved empty short links from the PostgreSQL database")
	return result, nil

}

func (db_ *DB) UpdateLink(link *api.Link) error {

	log.Info().Msg("UpdateLink: Start to insert new short links in the postgres database")

	ctx := db_.ctx

	if link.ID < 1 {
		log.Error().Msg("id must be greater than 0")
		return errors.New("id must be greater than 0")
	}

	_, err := db_.pool.Exec(ctx,
		`UPDATE links
		 SET  short_code=$1, long_url=$2, created_at=$3
		 WHERE id=$4`,
		link.ShortCode,
		link.LongUrl,
		link.CreatedAt,
		link.ID,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error updating link")
		return fmt.Errorf("error updating link: %w", err)
	}

	return nil
}

func (db_ *DB) GetLongUrl(shortCode string) (*api.Link, error) {

	log.Info().Msg("GetLongUrl: Start to get long url from the postgress database")

	ctx := db_.ctx

	query := `
	SELECT *
	FROM links 
	WHERE short_code = $1
 	`
	row, err := db_.pool.Query(ctx, query, shortCode)
	if err != nil {
		log.Error().Msg("Error to get link: " + err.Error())
		return nil, err
	}
	var link api.Link
	err = row.Scan(&link.ID, &link.ShortCode, &link.LongUrl, &link.CreatedAt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to scan link")
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("short URL not found: %s", shortCode)
		}
		return nil, fmt.Errorf("error getting link: %w", err)
	}

	log.Info().Str("short_code", shortCode).Str("long_url", link.LongUrl).Msg("Retrieved long URL")
	return &link, nil
}

func (db_ *DB) DeleteLink(linkId uint32) error {
	log.Info().Msg("DeleteLink: Start deleting a link from the PostgreSQL database")

	ctx := db_.ctx

	_, err := db_.pool.Exec(ctx, `DELETE FROM links WHERE id=$1`, linkId)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting link")
		return fmt.Errorf("error deleting link: %w", err)
	}

	log.Info().Msgf("Deleted link with ID %d from the PostgreSQL database", linkId)
	return nil
}
