// Reverse index package (storage and functionality)
// Here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted index

package index

import (
	"search_engine/pkg/crawler"
	"sync"

	"github.com/rs/zerolog/log"
)

// Handler of service.
type Service struct {
	counter uint32              // counter for Document.ID
	links   []crawler.Document  // crawler.Documents storage (links to sites with titles)
	index   map[string][]uint32 // inverted index storage  [word in title] => []Document.ID
	mu      sync.RWMutex        // mutex for thread safety
}

// New creates and returns a new Service instance
func New() *Service {
	s := Service{}
	s.links = make([]crawler.Document, 0, 1024)
	s.index = make(map[string][]uint32)
	log.Info().Msg("Service initialized.")
	return &s
}

// Clear resets the storage of documents and the index
func (s *Service) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Info().Msg("Clearing the document storage and index.")
	s.links = make([]crawler.Document, 0, 1024)
	s.index = make(map[string][]uint32)
	s.counter = 0
	log.Info().Msg("Document storage and index cleared.")
}
