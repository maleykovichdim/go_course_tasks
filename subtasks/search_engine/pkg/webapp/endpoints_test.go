// testing endpoints of web application
package webapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"search_engine/pkg/crawler"
	"strings"
	"testing"
)

// MockDescriptor implements the Descriptor interface for testing purposes.
type MockDescriptor struct{}

func (m *MockDescriptor) GetIndexDescription() (string, error) {
	return `{"index": "mock data"}`, nil
}

func (m *MockDescriptor) GetDocsDescription() (*[]crawler.Document, error) {
	docs := []crawler.Document{
		{
			ID:    1,
			URL:   "https://example.com",
			Title: "Example Title 1",
			Body:  "This is the body of the first document.",
		},
		{
			ID:    2,
			URL:   "https://golang.org",
			Title: "Golang",
			Body:  "Here is some information about Go.",
		},
		{
			ID:    3,
			URL:   "https://developer.mozilla.org",
			Title: "MDN Web Docs",
			Body:  "Resources for developers, by developers.",
		},
	}
	return &docs, nil
}

func (m *MockDescriptor) GetDoc(i uint32) (*crawler.Document, error) {
	doc := crawler.Document{
		ID:    1,
		URL:   "https://example.com",
		Title: "Example Title 1",
		Body:  "This is the body of the first document.",
	}
	doc.ID = i
	return &doc, nil
}

func (m *MockDescriptor) PostDoc(doc *crawler.Document) error {
	doc.ID = 99999
	return nil
}

func (m *MockDescriptor) PostDocs(docs *[]crawler.Document) error {
	for i := 0; i < len(*docs); i++ {
		i := uint32(i)
		(*docs)[i].ID = i + 99999
	}
	return nil
}

func (m *MockDescriptor) PutDoc(doc *crawler.Document) error { //TODO: add mutex ????
	return nil
}

func (m *MockDescriptor) DeleteDoc(i uint32) error {
	return nil
}

func (m *MockDescriptor) FindDocs(word string) (*[]crawler.Document, error) {
	if word != "key_word" {
		return nil, errors.New("wrong key word")
	}
	docs := []crawler.Document{
		{
			ID:    1,
			URL:   "https://example.com",
			Title: "Example Title 1",
			Body:  "This is the body of the first document.",
		},
		{
			ID:    2,
			URL:   "https://golang.org",
			Title: "Golang",
			Body:  "Here is some information about Go.",
		},
		{
			ID:    3,
			URL:   "https://developer.mozilla.org",
			Title: "MDN Web Docs",
			Body:  "Resources for developers, by developers.",
		},
	}
	return &docs, nil
}

// Setup function runs before each test
func setup(t *testing.T) {
	t.Log("Setup: Executing before the test")
	s.api = New()
	s.api.d = &MockDescriptor{}
}

func TestAPI_getDocs(t *testing.T) {

	//init
	setup(t)

	//request
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := `[{"id":1,"url":"https://example.com","title":"Example Title 1","body":"This is the body of the first document."},{"id":2,"url":"https://golang.org","title":"Golang","body":"Here is some information about Go."},{"id":3,"url":"https://developer.mozilla.org","title":"MDN Web Docs","body":"Resources for developers, by developers."}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}

func TestAPI_getIndex(t *testing.T) {
	//init
	setup(t)

	//request
	req := httptest.NewRequest(http.MethodGet, "/index", nil)

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := "{\"index\": \"mock data\"}"
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAPI_getDoc(t *testing.T) {
	//init
	setup(t)

	//request
	req := httptest.NewRequest(http.MethodGet, "/doc/1", nil)

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := `{"id":1,"url":"https://example.com","title":"Example Title 1","body":"This is the body of the first document."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAPI_newDoc(t *testing.T) {
	//init
	setup(t)

	//request
	doc := crawler.Document{
		URL:   "https://example.com",
		Title: "Example Title 1",
		Body:  "This is the body of the first document.",
	}

	payload, _ := json.Marshal(doc)
	req := httptest.NewRequest(http.MethodPost, "/doc", bytes.NewBuffer(payload))

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := `{"id":99999,"url":"https://example.com","title":"Example Title 1","body":"This is the body of the first document."}`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAPI_newDocs(t *testing.T) {
	//init
	setup(t)

	//request
	doc := []crawler.Document{
		{
			URL:   "https://example1.com",
			Title: "Example Title 1",
			Body:  "This is the body of the first document1.",
		},
		{
			URL:   "https://example.com2",
			Title: "Example Title 2",
			Body:  "This is the body of the first document2.",
		},
		{
			URL:   "https://example3.com",
			Title: "Example Title 3",
			Body:  "This is the body of the first document3.",
		},
	}

	payload, _ := json.Marshal(doc)
	req := httptest.NewRequest(http.MethodPost, "/docs", bytes.NewBuffer(payload))

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := `[{"id":99999,"url":"https://example1.com","title":"Example Title 1","body":"This is the body of the first document1."},{"id":100000,"url":"https://example.com2","title":"Example Title 2","body":"This is the body of the first document2."},{"id":100001,"url":"https://example3.com","title":"Example Title 3","body":"This is the body of the first document3."}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAPI_updateDoc(t *testing.T) {
	//init
	setup(t)

	//request
	doc := crawler.Document{
		URL:   "https://example.com",
		Title: "Example Title 1",
		Body:  "This is the body of the first document.",
	}

	payload, _ := json.Marshal(doc)
	req := httptest.NewRequest(http.MethodPut, "/doc", bytes.NewBuffer(payload))

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := ``
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAPI_deleteDoc(t *testing.T) {
	//init
	setup(t)

	//request
	req := httptest.NewRequest(http.MethodDelete, "/doc/3", nil)

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := ``
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAPI_getDocsByKeyword(t *testing.T) {

	//init
	setup(t)

	//request
	req := httptest.NewRequest(http.MethodGet, "/find?word=key_word", nil)

	//execute
	rr := httptest.NewRecorder()
	s.api.router.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("Wrong response code: get %d, а want %d", rr.Code, http.StatusOK)
	}

	//check
	expected := `[{"id":1,"url":"https://example.com","title":"Example Title 1","body":"This is the body of the first document."},{"id":2,"url":"https://golang.org","title":"Golang","body":"Here is some information about Go."},{"id":3,"url":"https://developer.mozilla.org","title":"MDN Web Docs","body":"Resources for developers, by developers."}]`
	if strings.TrimSpace(rr.Body.String()) != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

}
