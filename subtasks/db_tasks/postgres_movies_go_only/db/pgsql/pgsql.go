// DB API implementation for postgress db
package pgsql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"maleykovich.db/db"
)

const (
	user     = "postgres"
	password = "postgress"
)

type DB struct {
	pool *pgxpool.Pool
}

// ?
func New() *DB {
	return &DB{}
}

// ?
func (db_ *DB) Open(ctx context.Context, dbName string) error {

	var err error
	connStr := `postgres://` + user + `:` + password + `@localhost:5432/` + dbName + `?sslmode=disable`
	db_.pool, err = pgxpool.New(ctx, connStr)
	return err
}

func (db_ *DB) Close() { //ctx ?
	db_.pool.Close()
}

// get movies from DB
func (db_ *DB) Movies(ctx context.Context, studioId uint32) (*[]db.Movie, error) {

	var query string
	var rows pgx.Rows
	var err error
	var result []db.Movie

	// Selectively query based on whether a studio ID is specified
	if studioId == 0 {
		query = `
        SELECT m.*
        FROM Movies m
        `
		rows, err = db_.pool.Query(context.Background(), query)
	} else {
		query = `
		   SELECT m.*
           FROM Movies m
           JOIN Studios s ON m.studio_id = s.id
	       WHERE s.id = $1
        `
		rows, err = db_.pool.Query(context.Background(), query, studioId)
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var movie db.Movie
		err = rows.Scan(&movie.ID, &movie.Name, &movie.ReleaseYear, &movie.StudioID, &movie.BoxOffice, &movie.Rating)
		if err != nil {
			return nil, err
		}
		result = append(result, movie)
	}

	// Check for errors after iterating
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &result, nil
}

// add one movie
func (db_ *DB) AddMovies(ctx context.Context, movies *[]db.Movie) error {

	var numStudios uint32
	db_.pool.QueryRow(ctx, "SELECT count(*) FROM Studios").Scan(&numStudios)

	tx, err := db_.pool.Begin(ctx)
	if err != nil {
		return err
	}
	// cancel transaction in case of error
	defer tx.Rollback(ctx)
	// batch request
	var batch = &pgx.Batch{}

	for _, movie := range *movies {
		if movie.StudioID > numStudios || movie.StudioID < 1 {
			return errors.New("wrong studio name id")
		}
		batch.Queue(
			"INSERT INTO Movies (name, release_year, studio_id, box_office, rating) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (name, release_year) DO NOTHING",
			movie.Name,
			movie.ReleaseYear,
			movie.StudioID,
			movie.BoxOffice,
			movie.Rating,
		)
	}

	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}

	// approving transaction
	return tx.Commit(ctx)
}

func (db_ *DB) ChangeMovie(ctx context.Context, movie *db.Movie) error {

	_, err := db_.pool.Exec(ctx,
		`UPDATE movies
		 SET  name=$1, release_year=$2, studio_id=$3, box_office=$4, rating=$5
		 WHERE id=$6`,
		movie.Name,
		movie.ReleaseYear,
		movie.StudioID,
		movie.BoxOffice,
		movie.Rating,
		movie.ID,
	)
	return err
}

func (db_ *DB) DeleteMovie(ctx context.Context, movieId uint32) error {
	_, err := db_.pool.Exec(ctx,
		`DELETE FROM movies  WHERE id=$1`,
		movieId,
	)
	return err
}
