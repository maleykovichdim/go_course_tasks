// reverse index package (storage and functionality)
// here the data of parsing of sites (links) is stored, in the form of a slice of these numbered documents and a map of the inverted inde

package index

import (
	"search_engine/pkg/crawler"
	"testing"
)

func Test_AddDocumentsToStorage(t *testing.T) {

	s := New()

	data := []crawler.Document{
		{
			ID:    0,
			URL:   "https://yandex.ru",
			Title: "Яндекс",
		},
		{
			ID:    1,
			URL:   "https://google.ru",
			Title: "Google",
		},
	}

	s.AddDocumentsToStorage(data)
	s.PrintLinks()
	s.PrintIndex()
}
