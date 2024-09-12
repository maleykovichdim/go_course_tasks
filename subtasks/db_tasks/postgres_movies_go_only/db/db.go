package db

import (
	"context"
	"time"
)

type Interface interface {
	Movies(ctx context.Context, studioId uint32) (*[]Movie, error) // if studioId == 0 return All movies
	AddMovies(ctx context.Context, movies *[]Movie) error
	ChangeMovie(ctx context.Context, movie *Movie) error
	DeleteMovie(ctx context.Context, movie *Movie) error
}

// Studio represents the "Studios" table
type Studio struct {
	ID   uint32 `json:"id"`   // Unique identifier for the studio
	Name string `json:"name"` // Name of the studio
}

// Movie represents the "Movies" table
type Movie struct {
	ID          uint32  `json:"id"`           // Unique identifier for the movie
	Name        string  `json:"name"`         // Name of the movie
	ReleaseYear int     `json:"release_year"` // Year the movie was released (must be >= 1800)
	StudioID    uint32  `json:"studio_id"`    // Reference to the studio (foreign key)
	BoxOffice   float64 `json:"box_office"`   // Box office earnings
	Rating      string  `json:"rating"`       // Rating (PG-10, PG-13, PG-18)
}

// Actor represents the "Actors" table
type Actor struct {
	ID       uint32    `json:"id"`        // Unique identifier for the actor
	Name     string    `json:"name"`      // Name of the actor
	BirthDay time.Time `json:"birth_day"` // Birth date of the actor
}

// Director represents the "Directors" table
type Director struct {
	ID       uint32    `json:"id"`        // Unique identifier for the director
	Name     string    `json:"name"`      // Name of the director
	BirthDay time.Time `json:"birth_day"` // Birth date of the director
}
