// package webapp contains the implementation of the web server.
package webapp

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log" // Importing the zerolog library for structured logging
)

// Addr defines the address where the server listens for requests.
const Addr = ":8080"

// server struct holds the API and the HTTP server instance.
type server struct {
	api *API         // Pointer to the API instance
	srv *http.Server // Pointer to the HTTP server
}

var s server // Global instance of the server

// StartServer initializes and starts the HTTP server.
func StartServer(d Storage) {
	log.Info().Msg("Starting server...") // Log that the server is starting

	s.api = New()         // Initialize a new API
	s.api.d = d           // Assign the storage to the API (consider moving this inside New)
	s.srv = &http.Server{ // Initialize the HTTP server
		Addr:    Addr,         // Set the address for the server
		Handler: s.api.router, // Set the router for handling requests
	}

	// Start listening for incoming requests
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error().Err(err).Msgf("Could not listen on %s", Addr) // Log error if server fails to start
		return                                                    // Exit the function gracefully
	}

	log.Info().Msg("Server started successfully") // Log success message when the server starts
}

// ShutdownServer gracefully shuts down the HTTP server.
func ShutdownServer() {
	log.Info().Msg("Shutting down server...") // Log that the server shutdown process is starting

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second) // Create a context with a timeout
	defer cancel()                                                          // Ensure that the cancel function is called to release resources

	// Attempt to shutdown the server gracefully
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown") // Log error if the shutdown fails
	}

	log.Info().Msg("Server shutdown gracefully") // Log successful shutdown
}
