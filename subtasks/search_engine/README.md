# Search Engine Development
(This project is implemented, but only for the requested technical tasks. Since it was a training project. The main requests are implemented in  the POSTMAN file)

## Task 1: Development of the First Iteration of the Search Engine

The first version of the search engine will utilize the `crawler` package to scan websites and enable users to search for data through the command line.

### Objectives:

#### Task #1: Project Structure
- Create an application with the following standard directory structure:
  - `cmd`
  - `pkg`
- Copy the `crawler` package from the `GoSearch` directory in the course repository into your application.

#### Task #2: Executable Package
- Develop an executable package that uses the `crawler` package to scan the websites `go.dev` and `golang.org` upon execution.
- Combine the results of the scanning from the two websites.

---

## Task 2: Development of a Fast Search Index

To facilitate quick word searches in documents, an inverted index will be implemented (Learn more: [Inverted Index](https://en.wikipedia.org/wiki/Inverted_index)). Given that the index is sorted, we can leverage binary search to enhance search result retrieval.

### Objectives:

#### Task #1: Create Inverted Index
- Create an inverted index for the scanned documents.
- Store the index in the `index` package.
- The key of the index should be each word from the link descriptions, while the value should be the document number.
- All discovered links from the website scans should be formatted in the following structure:

```go
// Document represents a document, a web page obtained by the crawler.
type Document struct {
    ID    int    // Unique identifier for the document
    URL   string // URL of the document
    Title string // Title of the document
}

    Indexed documents should be stored in an array of documents as well as in an object from the index package.
    Each document must include its corresponding number as a field in the data structure, along with the URL and Title.
    The array of documents should be sorted by their numbers, utilizing the sorting functionality from the standard library.

Task #2: Refactor Search Method

    Modify the search method to utilize the inverted index.
    The application should import the search index from the index package as a dependency and perform searches through the index instead of the document array.
    The search results from the index should return the document numbers (the field within the document structure).

Task #3: Utilize Binary Search

    After retrieving the document numbers from the index, implement binary search on the previously sorted array of documents (by document number).
