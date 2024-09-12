// Reverse index package (storage and functionality)
// Here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted index

package index

import (
	"search_engine/pkg/crawler"
	"strings"

	"github.com/rs/zerolog/log"
)

// Main function to put Documents into storage (thread safe)
func (s *Service) AddDocumentsToStorage(docs []crawler.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, doc := range docs {
		fl := false
		for _, dAdded := range s.links {
			if strings.TrimSpace(dAdded.Title) == strings.TrimSpace(doc.Title) {
				fl = true
				log.Debug().Str("title", doc.Title).Msg("Document already exists in storage, skipping.")
				break
			}
		}
		if fl {
			continue
		}

		// Set ID into documents
		id := s.counter
		s.counter++
		doc.ID = id
		s.links = append(s.links, doc)
		log.Info().Uint32("id", id).Str("title", doc.Title).Msg("Document added to storage.")

		// Create inverted index data
		words := strings.Fields(doc.Title)
		for _, w := range words {
			w = strings.ToLower(w) // attention: To LOWER CASE
			s.index[w] = append(s.index[w], id)
		}
	}

	s.sort()
	log.Info().Int("doc_count", len(docs)).Msg("Documents processed for storage.")
	return nil
}

// Main function to put one Document into storage (thread safe)
func (s *Service) AddDocumentToStorage(doc *crawler.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, dAdded := range s.links {
		if strings.TrimSpace(dAdded.Title) == strings.TrimSpace(doc.Title) {
			log.Debug().Str("title", doc.Title).Msg("Document already exists in storage, skipping.")
			return nil
		}
	}

	// Set ID for the document
	id := s.counter
	s.counter++
	doc.ID = id
	s.links = append(s.links, *doc)
	log.Info().Uint32("id", id).Str("title", doc.Title).Msg("Document added to storage.")

	// Create inverted index data
	words := strings.Fields(doc.Title)
	for _, w := range words {
		w = strings.ToLower(w)
		s.index[w] = append(s.index[w], id)
	}

	return nil
}
