// reverse index package (storage and functionality)
// here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted inde

// implementation "Storage" interface for webapp package (thread unsafe)

package index

import (
	"encoding/json"
	"errors"
	"search_engine/pkg/crawler"
)

func (s *Service) GetIndexDescription() (string, error) { //TODO: reDo for normal json-view w/o \"
	jsonData, err := json.Marshal(s.index)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (s *Service) GetDocsDescription() (*[]crawler.Document, error) {
	return &s.links, nil
}

func (s *Service) GetDoc(i uint32) (*crawler.Document, error) {
	return s.getDoc(i)
}

func (s *Service) PostDoc(doc *crawler.Document) error {
	return s.AddDocumentToStorage(doc)
}

func (s *Service) PostDocs(docs *[]crawler.Document) error {
	for i := 0; i < len(*docs); i++ {
		_ = s.AddDocumentToStorage(&(*docs)[i])
	}
	return nil
}

func (s *Service) PutDoc(doc *crawler.Document) error { //TODO: add mutex ????

	i, ok := s.binarySearchLink(doc.ID)
	if !ok {
		return errors.New("bad index")
	}
	s.links[i] = *doc
	return nil

}

func (s *Service) DeleteDoc(i uint32) error {
	return s.deleteDocumentFromStorage(i)
}

func (s *Service) FindDocs(word string) (*[]crawler.Document, error) {
	idx := s.index[word]
	docs, err := s.findDocsByIndexes(idx)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
