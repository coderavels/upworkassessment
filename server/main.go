package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/coderavels/upworkassessment/server/client"
	"github.com/coderavels/upworkassessment/server/handler"
)

const portNum string = ":8021"

type Handler interface {
	ListBooks(w http.ResponseWriter, r *http.Request)
	GetBookCollection(w http.ResponseWriter, r *http.Request)
}

func main() {
	log.Println("Starting our bookshelf backend server")
	err := godotenv.Load()

	ac := initializeAssessClient()
	h := initializeHandler(ac)

	http.HandleFunc("/books", h.ListBooks)
	http.HandleFunc("/collection/{bookISBN}", h.GetBookCollection)

	log.Println("Started on port", portNum)

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func initializeAssessClient() handler.AssessClient {
	assessClientUsername := os.Getenv("ASSESS_CLIENT_USER")
	assessClientPassword := os.Getenv("ASSESS_CLIENT_PASSWORD")
	assessClientBaseURL := os.Getenv("ASSESS_CLIENT_BASEURL")

	return client.NewAssessClient(AssessClientParams{
		BaseURL:  assessClientBaseURL,
		Username: assessClientUsername,
		Password: assessClientPassword,
	})
}

func initializeHandler(ac handler.AssessClient) Handler {
	return handler.NewHandler(HandlerParams{
		AssessClient: ac,
	})
}
