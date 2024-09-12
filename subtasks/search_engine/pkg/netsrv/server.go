// server for working with console application (client.go)
package netsrv

import (
	"bufio"
	"encoding/json"
	"net"
	"strings"
	"time"

	"search_engine/pkg/crawler"

	"github.com/rs/zerolog/log"
)

// net address
// service will listen to all requests on all IP-addresses of the computer on port 12345
const addr = "0.0.0.0:12345"

// network service protocol
const proto = "tcp4"

const timeout = 12 // seconds

// FinderDocs is an interface that expects a method FindDocs
// which takes a word as input and returns a slice of documents and an error
type FinderDocs interface {
	FindDocs(word string) (*[]crawler.Document, error)
}

// handleConn handles the incoming connections
func handleConn(conn net.Conn, fr FinderDocs, f *bool) {

	defer conn.Close()                                      // Ensure the connection is closed when the function exits
	conn.SetDeadline(time.Now().Add(time.Second * timeout)) // Set initial deadline for the connection

	reader := bufio.NewReader(conn) // Create a buffered reader to read data from the connection
	for {
		b, err := reader.ReadBytes('\n') // Read data until newline
		if err != nil {
			log.Error().Err(err).Msg("Error reading data from connection")
			return
		}
		msg := strings.TrimSuffix(string(b), "\n") // Remove trailing newline
		msg = strings.TrimSuffix(msg, "\r")        // Remove trailing carriage return if present

		if msg == "exit" {
			_, _ = conn.Write([]byte(msg)) // Echo back the 'exit' message
			*f = true                      // Set flag to true, signifying termination
			log.Info().Msg("Received exit command, closing connection")
			return
		}
		docs, err := fr.FindDocs(msg)
		if err != nil {
			_, err = conn.Write([]byte("-- internal server error -- \n")) // Notify client of server error
			if err != nil {
				log.Error().Err(err).Msg("Error writing to connection")
				return
			}
		} else {
			// Write the found documents to the connection as JSON
			if err := json.NewEncoder(conn).Encode(*docs); err != nil {
				log.Error().Err(err).Msg("Error encoding documents to JSON")
				return
			}
		}
		conn.SetDeadline(time.Now().Add(time.Second * timeout)) // Reset deadline after every successful operation
	}
}

// StartListen starts listening for incoming network connections
func StartListen(t FinderDocs) error {

	// Start the network service
	listener, err := net.Listen(proto, addr)
	if err != nil {
		return err
	}
	defer listener.Close() // Ensure the listener is closed when the function exits
	var f bool             // Flag to control the termination of the service
	for {
		if f {
			break
		}
		conn, err := listener.Accept() // Accept new connections
		if err != nil {
			return err
		}
		// Handle the connection in a separate goroutine
		go handleConn(conn, t, &f)
	}
	return nil
}
