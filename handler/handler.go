package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/coderavels/upworkassessment/client"
)

//go:generate go run go.uber.org/mock/mockgen -source=./handler.go -destination=./mocks/mock_handler.go -package=mocks . AssessClient

type HandlerParams struct {
	AssessClient AssessClient
}

type Handler struct {
	assessClient AssessClient
}

func NewHandler(params HandlerParams) Handler {
	return Handler{
		assessClient: params.AssessClient,
	}
}

type AssessClient interface {
	GetBooks() ([]client.Book, error)
	GetBook(bookISBN string) (client.BookDetails, error)
}

type Book struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}

type BookDetails struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
	Width string `json:"width"`
}

type GetBookCollectionResponse struct {
	Shelves      [][]BookDetails `json:"shelves"`
	TotalShelves int             `json:"total_shelves"`
}

func (h Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		books, err := h.assessClient.GetBooks()
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to fetch books, %s", err.Error()), 500)
			return
		}

		var listBooksResponse []Book
		for _, b := range books {
			listBooksResponse = append(listBooksResponse, Book{
				Title: b.Title,
				ISBN:  b.ISBN,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(listBooksResponse)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal books to resp, %s", err.Error()), 500)
			return
		}

		w.Write(resp)
		return
	}

	http.Error(w, "unsupported method", 405)
	return
}

func (h Handler) GetBookCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		widthQVal := r.URL.Query().Get("width")

		width, err := strconv.Atoi(widthQVal)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to parse width query arg %s, %s", widthQVal, err.Error()), 400)
			return
		}

		bookISBN := r.PathValue("bookISBN")

		booksInCollection, err := h.getBookCollectionFromClient(bookISBN)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed while getting book collection for %s, %s", bookISBN, err.Error()), 500)
			return
		}

		organisedCollection, err := organiseCollectionInShelves(booksInCollection, width)
		if err != nil {
			http.Error(w, fmt.Sprintf("error while organising collection in shelves, %s", err.Error()), 500)
			return
		}

		var shelves [][]BookDetails
		for _, oc := range organisedCollection {
			var shelf []BookDetails
			for _, b := range oc {
				shelf = append(shelf, BookDetails{
					Title: b.Title,
					ISBN:  b.ISBN,
					Width: b.Width,
				})
			}
			shelves = append(shelves, shelf)
		}
		response := GetBookCollectionResponse{
			Shelves:      shelves,
			TotalShelves: len(shelves),
		}

		w.Header().Set("Content-Type", "application/json")
		resp, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to marshal collection to resp, %s", err.Error()), 500)
			return
		}

		w.Write(resp)

		return
	}

	http.Error(w, "unsupported method", 405)
	return
}

func (h Handler) getBookCollectionFromClient(bookISBN string) ([]client.BookDetails, error) {
	var booksInCollection []client.BookDetails
	seenBooks := map[string]struct{}{}
	for {
		if bookISBN == "" {
			break
		}
		if _, ok := seenBooks[bookISBN]; ok {
			break
		}

		bookDetails, err := h.assessClient.GetBook(bookISBN)
		if err != nil {
			return nil, fmt.Errorf("error while get book %s, %w", bookISBN, err)
		}

		booksInCollection = append(booksInCollection, bookDetails)
		seenBooks[bookISBN] = struct{}{}
		bookISBN = bookDetails.Related
	}

	return booksInCollection, nil
}

func organiseCollectionInShelves(bookCollection []client.BookDetails, width int) ([][]client.BookDetails, error) {
	booksAlreadyShelved := map[string]struct{}{}
	var organisedCollection [][]client.BookDetails

	var booksLeftToBeShelved []client.BookDetails
	booksLeftToBeShelved = append(booksLeftToBeShelved, bookCollection...)

	for {
		if len(booksLeftToBeShelved) == 0 {
			break
		}
		booksOnShelf, err := fillShelfToMax(booksLeftToBeShelved, width)
		if err != nil {
			return nil, fmt.Errorf("error while filling shelf to maximum, %w", err)
		}

		slices.SortFunc(booksOnShelf, func(a, b client.BookDetails) int {
			return strings.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
		})
		organisedCollection = append(organisedCollection, booksOnShelf)

		for _, b := range booksOnShelf {
			booksAlreadyShelved[b.ISBN] = struct{}{}
		}

		var booksLeftAfterThisIteration []client.BookDetails
		for _, b := range booksLeftToBeShelved {
			if _, ok := booksAlreadyShelved[b.ISBN]; !ok {
				booksLeftAfterThisIteration = append(booksLeftAfterThisIteration, b)
			}
		}

		booksLeftToBeShelved = booksLeftAfterThisIteration
	}

	return organisedCollection, nil
}

func fillShelfToMax(booksToShelve []client.BookDetails, shelfWidth int) ([]client.BookDetails, error) {
	N := len(booksToShelve)
	maxWidthCovered := make([]int, shelfWidth+1)
	prevBookIdx := make([]int, shelfWidth+1)

	for i := range prevBookIdx {
		prevBookIdx[i] = -1
	}

	for i := 0; i < N; i++ {
		bookWidth, err := strconv.Atoi(strings.TrimSuffix(booksToShelve[i].Width, "cm"))
		if err != nil {
			return nil, fmt.Errorf("error while parsing book width %s, %w", booksToShelve[i].Width, err)
		}
		for j := shelfWidth; j >= bookWidth; j-- {
			if maxWidthCovered[j-bookWidth]+bookWidth > maxWidthCovered[j] {
				maxWidthCovered[j] = maxWidthCovered[j-bookWidth] + bookWidth
				prevBookIdx[j] = i
			}
		}
	}

	maxWidthCoveredOnShelf := maxWidthCovered[shelfWidth]

	var shelvedBooks []client.BookDetails
	for i := maxWidthCovered[shelfWidth]; i > 0 && maxWidthCoveredOnShelf > 0; {
		if prevBookIdx[i] != -1 {
			shelvedBook := booksToShelve[prevBookIdx[i]]
			bookWidth, err := strconv.Atoi(strings.TrimSuffix(shelvedBook.Width, "cm"))
			if err != nil {
				return nil, fmt.Errorf("error while parsing book width while shelving %s, %w", shelvedBook.Width, err)
			}
			shelvedBooks = append(shelvedBooks, shelvedBook)
			maxWidthCoveredOnShelf -= bookWidth
			i -= bookWidth
		} else {
			break
		}
	}

	return shelvedBooks, nil
}
