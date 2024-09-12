// Console client for interacting with a search engine
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure zerolog to write logs to console
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	for {
		// Attempt to connect to the server at localhost on port 12345
		conn, err := net.Dial("tcp4", "localhost:12345")
		if err != nil {
			// Log the error and retry
			log.Error().Err(err).Msg("Failed to connect. Retrying...")
			continue
		}

		// Reader for capturing user input from the console
		reader := bufio.NewReader(os.Stdin)

		// Buffer for storing response from the server
		buffer := make([]byte, 4096)
		for {
			fmt.Print("Enter message: ")

			// Read user input from the console
			msg, err := reader.ReadString('\n')
			if err != nil {
				// Log error reading input
				log.Error().Err(err).Msg("Failed to read from input")
				continue
			}

			// Send the message to the server
			_, err = conn.Write([]byte(msg))
			if err != nil {
				// Log error sending message and retry connection
				log.Error().Err(err).Msg("Failed to send message. Retrying connection...")
				break
			}

			// Read the server's response
			n, err := conn.Read(buffer)
			if err != nil {
				// Log error receiving message and retry connection
				log.Error().Err(err).Msg("Failed to receive message. Retrying connection...")
				break
			}

			// Display the server's response to the user
			fmt.Printf("Server response: %s\n", string(buffer[:n]))

			// Exit the client if the user inputs "exit"
			if msg == "exit\n" {
				conn.Close()
				return
			}
		}

		// Close the connection before retrying
		conn.Close()
	}
}
