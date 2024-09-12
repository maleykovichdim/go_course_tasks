// API of the web application
package webapp

import (
	"net/http"
	"search_engine/pkg/crawler"

	"github.com/gorilla/mux"
)

type Storage interface {
	GetIndexDescription() (string, error)
	GetDocsDescription() (*[]crawler.Document, error)

	GetDoc(i uint32) (*crawler.Document, error)
	PostDoc(doc *crawler.Document) error //todo pointer ??
	PostDocs(docs *[]crawler.Document) error
	PutDoc(docs *crawler.Document) error
	DeleteDoc(i uint32) error

	FindDocs(word string) (*[]crawler.Document, error)
}

// API of the web application, handling all router and storage interactions.
type API struct {
	router *mux.Router
	d      Storage
}

// New initializes a new API instance, sets up routes, and returns the API pointer.
func New() *API {
	api := API{
		router: mux.NewRouter(),
	}

	api.endpoints()
	return &api
}

// endpoints sets up all the HTTP endpoints and middleware for the API.
func (api *API) endpoints() {
	// Apply middlewares for logging, header management, etc.
	api.router.Use(logMiddleware)
	api.router.Use(headersMiddleware)
	api.router.Use(middleware)

	// Set up route endpoints for the API.

	// Retrieve a document by its ID
	api.router.HandleFunc("/doc/{id}", api.getDoc).Methods(http.MethodGet)

	// Retrieve all documents
	api.router.HandleFunc("/docs", api.getDocs).Methods(http.MethodGet)

	// Retrieve index of documents
	api.router.HandleFunc("/index", api.getIndex).Methods(http.MethodGet)

	// Add a new document
	api.router.HandleFunc("/doc", api.newDoc).Methods(http.MethodPost)

	// Add multiple new documents
	api.router.HandleFunc("/docs", api.newDocs).Methods(http.MethodPost)

	// Update an existing document by ID
	api.router.HandleFunc("/doc", api.updateDoc).Methods(http.MethodPut)

	// Delete a document by its ID
	api.router.HandleFunc("/doc/{id}", api.deleteDoc).Methods(http.MethodDelete)

	// Find documents by keyword in the title
	api.router.HandleFunc("/find", api.getDocsByKeyword).Methods(http.MethodGet) // request: /?word=....
}
