// Package pgsql implements testing for the Postgres database API.
package pgsql

import (
	"context"
	"fmt"
	"log"
	"testing"

	"maleykovich.db/db"
)

const (
	dbTestName = "movies_test" // Test database name
)

var (
	dbTest  = New()                // Create a new database instance
	ctxTest = context.Background() // Context for database operations
)

// createDB creates the test database.
func createDB() {
	var err error
	_ = dbTest.Open(ctxTest, "postgres") // Open connection to Postgres

	// Terminate user connections to the database for removal
	disconnectUsersSQL := fmt.Sprintf(`SELECT pg_terminate_backend(pg_stat_activity.pid)
        FROM pg_stat_activity
        WHERE pg_stat_activity.datname = '%s'
          AND pid <> pg_backend_pid();`, dbTestName)

	_, err = dbTest.pool.Exec(context.Background(), disconnectUsersSQL)
	if err != nil {
		log.Printf("Could not disconnect users from the database: %v\n", err)
	}

	// Drop the database if it exists
	dropDBSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s;", dbTestName)
	_, err = dbTest.pool.Exec(context.Background(), dropDBSQL)
	if err != nil {
		log.Fatalf("Failed to drop database: %v\n", err)
	}

	// Create a new database
	createDBSQL := fmt.Sprintf("CREATE DATABASE %s;", dbTestName)
	_, err = dbTest.pool.Exec(context.Background(), createDBSQL)
	if err != nil {
		log.Fatalf("Failed to create database: %v\n", err)
	}

	fmt.Println("Database has been dropped and recreated successfully.")
	dbTest.Close() // Close database connection
}

// xTestMain initializes the database and fills it with test data.
func xTestMain() {
	createDB() // Create database
	err := dbTest.Open(ctxTest, dbTestName)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	dbTest.Close()
	fillDb(dbTestName) // Fill the database with test data
}

// TestDB_PG_API tests the Postgres database API.
func TestDB_PG_API(t *testing.T) {
	xTestMain() // Create the database and fill it with data
	fmt.Println("Start...")

	// List of movies to add to the database
	movies := []db.Movie{
		{Name: "The Shawshank Redemption", ReleaseYear: 1994, StudioID: 1, BoxOffice: 58.3, Rating: "PG-10"},
		{Name: "The Godfather", ReleaseYear: 1972, StudioID: 2, BoxOffice: 250.0, Rating: "PG-18"},
		{Name: "The Dark Knight", ReleaseYear: 2008, StudioID: 3, BoxOffice: 1004.6, Rating: "PG-13"},
		{Name: "Forrest Gump", ReleaseYear: 1994, StudioID: 4, BoxOffice: 678.2, Rating: "PG-13"},
		{Name: "Inception", ReleaseYear: 2010, StudioID: 5, BoxOffice: 836.8, Rating: "PG-13"},
	}

	// Open the database for testing
	err := dbTest.Open(ctxTest, dbTestName)
	if err != nil {
		log.Fatalf("Failed to open database: %v\n", err)
	}
	defer dbTest.Close() // Ensure the database is closed after the test

	// Get the current number of movies in the database
	mv, err := dbTest.Movies(ctxTest, 0)
	if err != nil {
		log.Fatalf("Failed to get movies: %v\n", err)
	}
	numBefore := len(*mv)

	// Add movies to the database
	err = dbTest.AddMovies(ctxTest, &movies)
	if err != nil {
		log.Fatalf("Failed to add movies: %v\n", err)
	}

	mv, err = dbTest.Movies(ctxTest, 0)
	if err != nil {
		log.Fatalf("Failed to get movies: %v\n", err)
	}

	// Expected number of movies after addition
	numWaitAfter := numBefore + len(movies)
	if len(*mv) != numWaitAfter {
		log.Fatalf("Wrong number of rows in movies table: have: %d, expected: %d\n", len(*mv), numWaitAfter)
	}

	// Get IDs for testing operations
	n := uint32(numWaitAfter - 2) // ID of the movie to delete
	m := uint32(numWaitAfter - 1) // ID of the movie to update

	// Delete the movie by ID
	err = dbTest.DeleteMovie(ctxTest, n)
	if err != nil {
		log.Fatalf("Failed to delete movie: %v\n", err)
	}
	numWaitAfter-- // Decrease the expected number after deletion

	// Check the current number of movies after deletion
	mv, err = dbTest.Movies(ctxTest, 0)
	if err != nil {
		log.Fatalf("Failed to get movies: %v\n", err)
	}
	fmt.Println(mv)
	if len(*mv) != numWaitAfter {
		log.Fatalf("Wrong number of rows in movies table: have: %d, expected: %d\n", len(*mv), numWaitAfter)
	}

	// Verify that the deleted movie no longer exists in the table
	isMovieDeleted := true
	for _, a := range *mv {
		if a.ID == n {
			isMovieDeleted = false // Movie still exists if this is false
		}
	}
	if !isMovieDeleted {
		log.Fatalf("Failed to delete movie with ID %d: it still exists in the database\n", n)
	}

	// Prepare the movie for updating
	movie := movies[m-uint32(numBefore)]
	movie.ID = m                                     // Set the movie ID
	movie.Name = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // New name for updating

	// Update the movie in the database
	err = dbTest.ChangeMovie(ctxTest, &movie)
	if err != nil {
		log.Fatalf("Failed to update movie with ID %d: %v\n", m, err)
	}

	// Check the updated movie
	mv, err = dbTest.Movies(ctxTest, 0)
	if err != nil {
		log.Fatalf("Failed to get movies: %v\n", err)
	}
	fmt.Println(mv)
	isMovieUpdated := false
	for _, a := range *mv {
		if a.ID == m && a.Name == "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" {
			isMovieUpdated = true // Movie updated successfully if this is true
		}
	}
	if !isMovieUpdated {
		log.Fatalf("Failed to update movie with ID %d: it does not exist (data part)\n", m)
	}

	// Get movies for studioID=1
	mv, err = dbTest.Movies(ctxTest, 1)
	if err != nil {
		log.Fatalf("Failed to get movies: %v\n", err)
	}
	fmt.Println(mv)
	if len(*mv) != 3 {
		log.Fatalf("Failed to get expected number of movies for studioID=1: have: %d, expected: 3\n", len(*mv))
	}

	fmt.Println("Done...")

	// TODO: Remove test database if needed
}
