// Main function for search_engine
package main

import (
	"os"

	"search_engine/pkg/crawler/index"
	"search_engine/pkg/crawler/spider"
	"search_engine/pkg/netsrv"
	"search_engine/pkg/webapp"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	// _ "net/http/pprof"
)

func main() {
	// Set up zerolog to output logs to console
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	urls := []string{"https://go.dev", "https://golang.org"}

	// Create a new Service object using the constructor
	robotService := spider.Service{}

	// Create a new index service
	indexService := index.New()

	// Try to get data from storage
	err := index.ReadDataFromStorage(indexService)
	if err != nil {
		// Log the error reading from storage
		log.Error().Err(err).Msg("Error reading data from storage")
		log.Info().Msg("Attempting to fetch data from the net")

		for _, url := range urls {
			// Scan the URL for documents
			docs, err := robotService.Scan(url, 2)
			if err != nil {
				// Log error parsing the URL
				log.Error().Err(err).Str("url", url).Msg("Error while parsing URL")
				continue
			}

			// Add documents to storage
			indexService.AddDocumentsToStorage(docs)
		}

		// Write data back to storage
		err = index.WriteDataToStorage(indexService)
		if err != nil {
			log.Error().Err(err).Msg("Error writing data to storage")
		} else {
			log.Info().Msg("Data successfully written to storage")
		}
	} else {
		log.Info().Msg("Data successfully read from storage")
	}

	// Start web server
	go webapp.StartServer(indexService)

	// Start network server
	err = netsrv.StartListen(indexService)
	if err != nil {
		log.Error().Err(err).Msg("Error starting network server")
	}

}
