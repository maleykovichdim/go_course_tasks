// put test data into DB

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

func fillDb() {
	// Establish a connection to the PostgreSQL database
	connStr := fmt.Sprintf("user=%s password=%s dbname=movies sslmode=disable", user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed when done

	// Populate the database with test data
	populateDatabase(db)
}

// Function to fill the database with test data
func populateDatabase(db *sql.DB) {
	// Create some sample studios
	studioNames := []string{"Studio A", "Studio B", "Studio C", "Studio D", "Studio E"}
	for _, name := range studioNames {
		_, err := db.Exec("INSERT INTO Studios (name) VALUES ($1) ", name) // ON CONFLICT (name) DO NOTHING

		//fmt.Println("-------", name)
		if err != nil {
			log.Fatalf("Failed to insert studio: %v", err)
		}
	}

	// Create some sample actors
	actorNames := []string{"Actor 1", "Actor 2", "Actor 3", "Actor 4", "Actor 5",
		"Actor 6", "Actor 7", "Actor 8", "Actor 9", "Actor 10"}
	for _, name := range actorNames {
		_, err := db.Exec("INSERT INTO Actors (name, birth_day) VALUES ($1, $2) ", name, time.Now().AddDate(-30, 0, 0))
		if err != nil {
			log.Fatalf("Failed to insert actor: %v", err)
		}
	}

	// Create some sample directors
	directorNames := []string{"Director 1", "Director 2", "Director 3", "Director 4", "Director 5"}
	for _, name := range directorNames {
		_, err := db.Exec("INSERT INTO Directors (name, birth_day) VALUES ($1, $2) ", name, time.Now().AddDate(-35, 0, 0))
		if err != nil {
			log.Fatalf("Failed to insert director: %v", err)
		}
	}

	// Create sample movies
	movies := []struct {
		Name        string
		ReleaseYear int
		StudioName  string
		BoxOffice   float64
		Rating      string
	}{
		{"Movie 1", 2020, "Studio A", 1000000, "PG-13"},
		{"Movie 2", 2019, "Studio B", 1500000, "PG-13"},
		{"Movie 3", 2021, "Studio C", 2000000, "PG-10"},
		{"Movie 4", 2018, "Studio D", 2500000, "PG-18"},
		{"Movie 5", 2022, "Studio E", 3000000, "PG-13"},
		{"Movie 6", 2020, "Studio A", 1200000, "PG-10"},
		{"Movie 7", 2022, "Studio B", 2200000, "PG-13"},
		{"Movie 8", 2018, "Studio C", 900000, "PG-18"},
		{"Movie 9", 2021, "Studio D", 600000, "PG-10"},
		{"Movie 10", 2017, "Studio E", 5000000, "PG-10"},
	}

	// Insert movies into the database
	for _, movie := range movies {
		// First, get the studio_id from the studio name
		var studioID int
		err := db.QueryRow("SELECT id FROM Studios WHERE name = $1", movie.StudioName).Scan(&studioID)
		if err != nil {
			log.Fatalf("Failed to get studio ID: %v", err)
		}

		// Insert the new movie into the Movies table
		_, err = db.Exec(
			"INSERT INTO Movies (name, release_year, studio_id, box_office, rating) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (name, release_year) DO NOTHING",
			movie.Name,
			movie.ReleaseYear,
			studioID,
			movie.BoxOffice,
			movie.Rating,
		)
		if err != nil {
			log.Fatalf("Failed to insert movie: %v", err)
		}
	}

	// Populate Movies_Actors table
	for i, _ := range movies {
		for j := 1; j <= 3; j++ { // Assume we add 3 actors to each movie
			actorID := (i*3 + j) // Simple example to map IDs from 1 to 10
			if actorID > 10 {
				actorID = actorID % 10
				if actorID == 0 {
					actorID = 10
				}
			}
			// fmt.Println(i+1, " ", actorID)
			_, err := db.Exec(
				"INSERT INTO Movies_Actors (movie_id, actor_id) VALUES ($1, $2)",
				i+1,     // Movie ID (1, 2, ..., 10)
				actorID, // Actor ID (1, 2, ..., 10)
			)
			if err != nil {
				log.Fatalf("Failed to insert movie actor relation: %v", err)
			}
		}
	}

	// Populate Movies_Directors table
	for i, _ := range movies {
		directorID := (i % len(directorNames)) + 1 // Use directors in a round-robin fashion
		_, err := db.Exec(
			"INSERT INTO Movies_Directors (movie_id, director_id) VALUES ($1, $2)",
			i+1,        // Movie ID (1, 2, ..., 10)
			directorID, // Director ID (1, 2, ..., 5)
		)
		if err != nil {
			log.Fatalf("Failed to insert movie director relation: %v", err)
		}
	}
	fmt.Println("Database populated with test data successfully!")
}
