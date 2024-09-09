package urls

import (
	"crypto/sha256"
	"encoding/base64"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Service provides methods for URL shortening and validation
type Service struct{}

// New initializes a new Service instance
func New() *Service {
	return &Service{}
}

// ShortenSimple generates a random string of the specified length
func (s *Service) ShortenSimple(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// Shorten generates a short URL for a given long URL by hashing it with a salt
// and encoding the hash using base64, trimmed to a desired length
func (s *Service) Shorten(longURL string) string {
	const salt = "tornado_salt_key_for_shortening_urls_in_this_environment_variable"

	if len(longURL) == 0 {
		// Handle empty input and decide on a strategy, possibly returning an error or log
		return ""
	}

	// Salt the URL and compute its SHA-256 hash
	saltedURL := longURL + salt
	hash := sha256.Sum256([]byte(saltedURL))

	// Encode the hash into base64 and trim the resulting string
	base64Hash := base64.URLEncoding.EncodeToString(hash[:])
	shortURL := strings.TrimRight(base64Hash[:10], "=")

	// Return with an additional random suffix
	return shortURL + s.ShortenSimple(3)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// IsValidURL verifies whether the provided string is a valid URL
func IsValidURL(url string) bool {
	// Pre-compiled regex improved URL validation accuracy and speed by avoiding recompilation
	var urlRegex = regexp.MustCompile(`^(http|https):\/\/[^\s$.?#].[^\s]*$`)
	return urlRegex.MatchString(url)
}
