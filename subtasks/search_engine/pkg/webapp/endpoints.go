// package webapp contains the implementation of endpoints for web applications.
package webapp

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"search_engine/pkg/crawler"
)

// getDocs retrieves a description of all documents.
func (api *API) getDocs(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request to fetch all documents")

	docs, err := api.d.GetDocsDescription()
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve documents description")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set response header
	if err := json.NewEncoder(w).Encode(docs); err != nil {
		log.Error().Err(err).Msg("Failed to encode documents response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Successfully fetched all documents")
}

// getIndex retrieves the description of the index.
func (api *API) getIndex(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request to fetch index description")

	indexDescription, err := api.d.GetIndexDescription()
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve index description")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set response header
	w.Write([]byte(indexDescription))
	log.Info().Msg("Successfully fetched index description")
}

// getDoc retrieves a document by its ID.
func (api *API) getDoc(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request to fetch document by ID")

	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("Invalid document ID format")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	doc, err := api.d.GetDoc(uint32(num))
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve document by ID")
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set response header
	if err := json.NewEncoder(w).Encode(doc); err != nil {
		log.Error().Err(err).Msg("Failed to encode document response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Successfully fetched document by ID")
}

// newDoc adds a new document.
func (api *API) newDoc(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request to add a new document")

	var doc crawler.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		log.Error().Err(err).Msg("Invalid document input")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := api.d.PostDoc(&doc); err != nil {
		log.Error().Err(err).Msg("Failed to save new document")
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set response header
	if err := json.NewEncoder(w).Encode(doc); err != nil {
		log.Error().Err(err).Msg("Failed to encode response with new document")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Successfully added a new document")
}

// newDocs adds multiple new documents
func (api *API) newDocs(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request to add new documents")

	var docs []crawler.Document
	if err := json.NewDecoder(r.Body).Decode(&docs); err != nil {
		log.Error().Err(err).Msg("Invalid input for new documents")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := api.d.PostDocs(&docs); err != nil {
		log.Error().Err(err).Msg("Failed to save new documents")
		http.Error(w, "Failed to save documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set response header
	if err := json.NewEncoder(w).Encode(docs); err != nil {
		log.Error().Err(err).Msg("Failed to encode response with new documents")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Successfully added new documents")
}

// updateDoc updates a document by its ID
func (api *API) updateDoc(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request to update document")

	var doc crawler.Document
	if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
		log.Error().Err(err).Msg("Invalid input for document update")
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := api.d.PutDoc(&doc); err != nil {
		log.Error().Err(err).Msg("Failed to update document")
		http.Error(w, "Failed to update document", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Successfully updated document")
	w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
}

// deleteDoc deletes a document by its ID
func (api *API) deleteDoc(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("Request to delete document")

	vars := mux.Vars(r)
	id := vars["id"]
	num, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Error().Err(err).Msg("Invalid document ID format for deletion")
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := api.d.DeleteDoc(uint32(num)); err != nil {
		log.Error().Err(err).Msg("Failed to delete document")
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Successfully deleted document")
	w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
}

// getDocsByKeyword finds documents by keyword in document titles
func (api *API) getDocsByKeyword(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	log.Info().Str("keyword", word).Msg("Request to find documents by keyword")

	docs, err := api.d.FindDocs(word)
	if err != nil {
		log.Error().Err(err).Msg("Failed to find documents by keyword")
		http.Error(w, "Failed to find documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set response header
	if err := json.NewEncoder(w).Encode(docs); err != nil {
		log.Error().Err(err).Msg("Failed to encode documents response")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Info().Msg("Successfully found documents by keyword")
}
