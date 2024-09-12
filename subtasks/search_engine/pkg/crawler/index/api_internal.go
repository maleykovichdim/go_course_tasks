// Reverse index package (storage and functionality)
// Here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted index

// Internal function of storage

package index

import (
	"errors"
	"search_engine/pkg/crawler"
	"sort"

	"github.com/rs/zerolog/log"
)

// Find index in slice for Document ID
func (s *Service) binarySearchLink(targetID uint32) (uint, bool) {
	items := s.links
	left, right := 0, len(items)-1 // Initialize the left and right boundaries

	for left <= right {
		mid := left + (right-left)/2 // Calculate the middle index (to avoid overflow)

		if items[mid].ID == targetID {
			log.Debug().Uint32("targetID", targetID).Msg("Document found in storage.")
			return uint(mid), true // Item found, return its index
		} else if items[mid].ID < targetID {
			left = mid + 1 // Target is in the right half
		} else {
			right = mid - 1 // Target is in the left half
		}
	}
	log.Debug().Uint32("targetID", targetID).Msg("Document not found in storage.")
	return 0, false // Item not found
}

// Sorting slice of documents
func (s *Service) sort() {
	sort.Slice(s.links, func(i, j int) bool {
		return s.links[i].ID < s.links[j].ID
	})
	log.Info().Msg("Documents sorted by ID.")
}

// Delete Document from storage
func (s *Service) deleteDocumentFromStorage(ind uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Info().Uint32("documentID", ind).Msg("Attempting to delete document from storage.")
	i, ok := s.binarySearchLink(ind)
	if !ok {
		log.Error().Uint32("documentID", ind).Msg("Error: document not found in storage.")
		return errors.New("no object in storage")
	}

	s.links = append(s.links[:i], s.links[i+1:]...)
	log.Info().Uint32("documentID", ind).Msg("Document successfully deleted from storage.")

	// TODO: remove ID from s.index if it is necessary

	return nil
}

// Get Document reference in storage by ID
func (s *Service) getDoc(ind uint32) (*crawler.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Info().Uint32("documentID", ind).Msg("Retrieving document from storage.")
	i, ok := s.binarySearchLink(ind)
	if !ok {
		log.Error().Uint32("documentID", ind).Msg("Error: document not found in storage.")
		return nil, errors.New("no object in storage")
	}
	log.Debug().Uint32("documentID", ind).Msg("Document retrieved from storage.")
	return &s.links[i], nil
}

// Find Documents by slice of Documents.ID
func (s *Service) findDocsByIndexes(docsNumber []uint32) (*[]crawler.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Info().Int("numberOfDocuments", len(docsNumber)).Msg("Finding documents by indexes.")
	var result []crawler.Document
	if len(docsNumber) == 0 {
		log.Warn().Msg("No document IDs provided.")
		return &result, nil
	}

	for _, dn := range docsNumber {
		k, ok := s.binarySearchLink(dn)
		if !ok {
			log.Debug().Uint32("documentID", dn).Msg("Document not found in storage during search.")
			continue
		}
		result = append(result, s.links[k])
		log.Debug().Uint32("documentID", dn).Msg("Document found during search.")
	}

	log.Info().Msg("Document search completed.")
	return &result, nil
}

// Not used now
func removeDuplicates[T comparable](slice []T) []T {
	uniqueValues := make(map[T]bool)
	var result []T
	for _, item := range slice {
		if _, ok := uniqueValues[item]; !ok {
			uniqueValues[item] = true
			result = append(result, item)
		}
	}
	return result
}
