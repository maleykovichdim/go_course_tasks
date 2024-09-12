//sql-requests according the task

package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

func requests() {
	// Establish a connection to the PostgreSQL database
	connStr := fmt.Sprintf("user=%s password=%s dbname=movies sslmode=disable", user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close() // Ensure the database connection is closed when done

	// Call functions to perform queries
	fetchMoviesWithStudios(db)
	fetchMoviesForActor(db, "Actor 1")                                        // Replace with the actor's actual name
	countMoviesForDirector(db, "Director 1")                                  // Replace with the director's actual name
	fetchMoviesForMultipleDirectors(db, []string{"Director 1", "Director 2"}) // Replace with directors' names
	countMoviesForActor(db, "Actor 1")
	fetchActorsAndDirectorsNotInLessThan2Movies(db)
	countMoviesWithBoxOfficeGreaterThan(db, 1000)
	countDirectorsWithMoviesBoxOfficeGreaterThan(db, 1000)
	fetchDistinctActorNames(db)
	countDuplicateMovies(db)
}

// Function to fetch movies with studio names
func fetchMoviesWithStudios(db *sql.DB) {
	rows, err := db.Query(`
        SELECT m.name AS movie_name, s.name AS studio_name
        FROM Studios s 
        JOIN Movies m ON m.studio_id = s.id
		WHERE s.name = $1;
    `, "Studio A")
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var movieName string
		var studioName string
		if err := rows.Scan(&movieName, &studioName); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Movie: %s, Studio: %s\n", movieName, studioName)
	}
}

// Function to fetch movies for a specific actor
func fetchMoviesForActor(db *sql.DB, actorName string) {
	rows, err := db.Query(`
        SELECT m.name AS movie_name
        FROM Movies m
        JOIN Movies_Actors ma ON m.id = ma.movie_id
        JOIN Actors a ON ma.actor_id = a.id
        WHERE a.name = $1;
    `, actorName)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("Movies for actor %s:\n", actorName)
	for rows.Next() {
		var movieName string
		if err := rows.Scan(&movieName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(movieName)
	}
}

// Function to count movies for a specific director
func countMoviesForDirector(db *sql.DB, directorName string) {
	var count int
	err := db.QueryRow(`
        SELECT COUNT(m.id) AS movie_count
        FROM Movies m
        JOIN Movies_Directors md ON m.id = md.movie_id
        JOIN Directors d ON md.director_id = d.id
        WHERE d.name = $1;
    `, directorName).Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("Director %s has directed %d movies.\n", directorName, count)
}

// Function to fetch movies for multiple directors
func fetchMoviesForMultipleDirectors(db *sql.DB, directorNames []string) {
	// Create a placeholder string for the SQL IN clause
	placeholder := ""
	args := make([]interface{}, len(directorNames))

	for i, name := range directorNames {
		if i > 0 {
			placeholder += ", "
		}
		placeholder += fmt.Sprintf("$%d", i+1)
		args[i] = name
	}

	query := fmt.Sprintf(`
        SELECT m.name AS movie_name
        FROM Movies m
        WHERE m.id IN (
            SELECT movie_id
            FROM Movies_Directors
            WHERE director_id IN (
                SELECT id
                FROM Directors
                WHERE name IN (%s)
            )
        );
    `, placeholder)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("Movies for directors %v:\n", directorNames)
	for rows.Next() {
		var movieName string
		if err := rows.Scan(&movieName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(movieName)
	}
}

// Function to count movies for a specific actor
func countMoviesForActor(db *sql.DB, actorName string) {
	var count int
	err := db.QueryRow(`
        SELECT COUNT(m.id) AS movie_count
        FROM Movies m
        JOIN Movies_Actors ma ON m.id = ma.movie_id
        JOIN Actors a ON ma.actor_id = a.id
        WHERE a.name = $1;
    `, actorName).Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("Actor %s has acted in %d movies.\n", actorName, count)
}

// Function to fetch actors and directors who participated in more than 2 movies
func fetchActorsAndDirectorsNotInLessThan2Movies(db *sql.DB) {
	query := `
        SELECT a.name AS name, 'Actor' AS role
        FROM Actors a
        JOIN Movies_Actors ma ON a.id = ma.actor_id
        GROUP BY a.id
        HAVING COUNT(ma.movie_id) > 2

        UNION ALL

        SELECT d.name AS name, 'Director' AS role
        FROM Directors d
        JOIN Movies_Directors md ON d.id = md.director_id
        GROUP BY d.id
        HAVING COUNT(md.movie_id) > 2;
    `

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	fmt.Println("Actors and Directors in more than 2 movies:")
	for rows.Next() {
		var name string
		var role string
		if err := rows.Scan(&name, &role); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n", role, name)
	}
}

// Function to count the number of movies with box office greater than a certain amount
func countMoviesWithBoxOfficeGreaterThan(db *sql.DB, boxOfficeThreshold float64) {
	var count int
	err := db.QueryRow(`
        SELECT COUNT(*) AS movie_count
        FROM Movies
        WHERE box_office > $1;
    `, boxOfficeThreshold).Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("Number of movies with box office greater than %.2f: %d\n", boxOfficeThreshold, count)
}

// Function to count directors whose movies made more than a certain amount
func countDirectorsWithMoviesBoxOfficeGreaterThan(db *sql.DB, boxOfficeThreshold float64) {
	var count int
	err := db.QueryRow(`
        SELECT COUNT(DISTINCT d.id) AS director_count
        FROM Directors d
        JOIN Movies_Directors md ON d.id = md.director_id
        JOIN Movies m ON md.movie_id = m.id
        WHERE m.box_office > $1;
    `, boxOfficeThreshold).Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("Number of directors whose movies made more than %.2f: %d\n", boxOfficeThreshold, count)
}

// Function to fetch distinct actor names
func fetchDistinctActorNames(db *sql.DB) {
	rows, err := db.Query(`
        SELECT DISTINCT a.name AS actor_name
        FROM Actors a;
    `)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	fmt.Println("Distinct actor names:")
	for rows.Next() {
		var actorName string
		if err := rows.Scan(&actorName); err != nil {
			log.Fatal(err)
		}
		fmt.Println(actorName)
	}
}

// Function to count duplicate movie titles
func countDuplicateMovies(db *sql.DB) {
	var count int
	err := db.QueryRow(`
        SELECT COUNT(*) AS duplicate_count
        FROM Movies
        WHERE name IN (
            SELECT name
            FROM Movies
            GROUP BY name
            HAVING COUNT(*) > 1
        );
    `).Scan(&count)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	fmt.Printf("Number of movies with duplicate titles: %d\n", count)
}
