// API for internal
package api

import (
	"time"
)

type Link struct {
	ID        uint      `json:"id"`
	ShortCode string    `json:"short_code"`
	LongUrl   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
}

type DBShortener interface {
	CreateDbStructure() error
	Open() error
	Close()
	PutEmptyShortLinks(shortLinks []string) ([]*Link, error)
	UpdateLink(link *Link) error
	GetLongUrl(shortCode string) (*Link, error)
	DeleteLink(linkId uint32) error
}

type CacheClient interface {
	SaveLink(shortURL, longURL string) error
	GetLink(shortURL string) (string, error)
}

type Shortener interface {
	Shorten(longURL string) string
	ShortenSimple(longURL int) string
}
