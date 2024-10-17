package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Book struct {
	Title string `json:"title"`
	ISBN  string `json:"isbn"`
}

type BookDetails struct {
	Title      string `json:"title"`
	ISBN       string `json:"isbn"`
	Publisher  string `json:"publisher"`
	Height     string `json:"height"`
	Published  int    `json:"published"`
	Author     string `json:"author"`
	Related    string `json:"related"`
	Collection string `json:"collection"`
	Width      string `json:"width"`
}

type AssessClientParams struct {
	BaseURL  string
	Username string
	Password string
}

type AssessClient struct {
	baseurl   string
	authToken string
}

func NewAssessClient(params AssessClientParams) AssessClient {
	return AssessClient{
		baseurl:   params.BaseURL,
		authToken: basicAuth(params.Username, params.Password),
	}
}

const (
	getBooksPath       string = "/books"
	getBookDetailsPath string = "/book/%s" // bookISBN
)

func (ac AssessClient) GetBooks() ([]Book, error) {
	requestURL := fmt.Sprintf("%s%s", ac.baseurl, getBooksPath)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", ac.authToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var books []Book
	err = json.Unmarshal(resBody, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (ac AssessClient) GetBook(bookISBN string) (BookDetails, error) {
	requestURL := fmt.Sprintf("%s%s", ac.baseurl, fmt.Sprintf(getBookDetailsPath, bookISBN))
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return BookDetails{}, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", ac.authToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return BookDetails{}, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return BookDetails{}, err
	}

	var book BookDetails
	err = json.Unmarshal(resBody, &book)
	if err != nil {
		return BookDetails{}, err
	}

	return book, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
