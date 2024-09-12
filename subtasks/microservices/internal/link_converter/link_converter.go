package link_converter

import (
	"context"
	"encoding/json"
	"log"
	"maleykovich/package/api"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	sync.RWMutex
	links map[string]string
	dp    api.Dispatcher
}

func New(dp api.Dispatcher) *Service {
	s := Service{}
	s.links = make(map[string]string, 1024)
	s.dp = dp
	return &s
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func keyExists(m map[string]string, key string) bool {
	_, exists := m[key]
	return exists
}

func (s *Service) generateShortLink() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func (s *Service) AddLink(w http.ResponseWriter, r *http.Request) {

	var request struct {
		OriginalLink string `json:"original_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var shortLink string

	for {
		shortLink = s.generateShortLink()
		s.Lock()
		fl := keyExists(s.links, shortLink)
		if fl {
			s.Unlock()
			continue
		}
		s.links[shortLink] = request.OriginalLink
		s.Unlock()
		break
	}

	response := struct {
		ShortLink string `json:"short_link"`
	}{
		ShortLink: shortLink,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	err := s.dp.Write(context.Background(), request.OriginalLink) //wrie kafka
	if err != nil {
		log.Println(err)
	}
}

func (s *Service) GetOriginalLink(w http.ResponseWriter, r *http.Request) {

	var req struct {
		ShortLink string `json:"short_link"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	shortLink := req.ShortLink

	s.RLock()
	originalLink, exists := s.links[shortLink]
	s.RUnlock()

	if !exists {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	response := struct {
		OriginalLink string `json:"original_link"`
	}{
		OriginalLink: originalLink,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
