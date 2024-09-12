// create postgress DB and tables
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // Importing PostgreSQL driver
)

const (
	user     = "postgres"
	password = "postgress"
)

func main() {
	// Connecting to PostgreSQL (usually to your base database, e.g., postgres)
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable", user, password)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Dropping the database if it exists (outside of a transaction)
	_, err = db.Exec("DROP DATABASE IF EXISTS movies")
	if err != nil {
		log.Fatalf("Failed to drop database: %v", err)
	}

	// Creating a new database (also outside of a transaction)
	_, err = db.Exec("CREATE DATABASE movies")
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}

	// Closing the current connection
	db.Close()

	// Now opening a new connection to the just created database
	newDBConnStr := fmt.Sprintf("user=%s password=%s dbname=movies sslmode=disable", user, password)
	newDB, err := sql.Open("postgres", newDBConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to the new database: %v", err)
	}
	defer newDB.Close()

	// Executing SQL script to create tables
	sqlScript := `
    DROP TABLE IF EXISTS Movies_Directors;
    DROP TABLE IF EXISTS Movies_Actors;
    DROP TABLE IF EXISTS Movies;
    DROP TABLE IF EXISTS Studios;

    CREATE TABLE Studios (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
    );

    CREATE TABLE Movies (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        release_year INT CHECK (release_year >= 1800),
        studio_id INT REFERENCES Studios(id),
        box_office DECIMAL,
        rating VARCHAR(5) CHECK (rating IN ('PG-10','PG-13','PG-18')),
        UNIQUE (name, release_year)
    );

    DROP TABLE IF EXISTS Actors;

    CREATE TABLE Actors (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        birth_day DATE
    );

    DROP TABLE IF EXISTS Directors;

    CREATE TABLE Directors (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        birth_day DATE
    );

    CREATE TABLE Movies_Actors (
        movie_id INT REFERENCES Movies(id) ON DELETE CASCADE,
        actor_id INT REFERENCES Actors(id) ON DELETE CASCADE,
        PRIMARY KEY (movie_id, actor_id)
    );

    CREATE TABLE Movies_Directors (
        movie_id INT REFERENCES Movies(id) ON DELETE CASCADE,
        director_id INT REFERENCES Directors(id) ON DELETE CASCADE,
        PRIMARY KEY (movie_id, director_id)
    );`

	// Executing SQL script to create tables
	_, err = newDB.Exec(sqlScript)
	if err != nil {
		log.Fatalf("Failed to execute SQL script: %v", err)
	}

	fmt.Println("Database and tables created successfully!")

	fillDb()

	requests()
}
