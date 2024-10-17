package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

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

	server := &http.Server{Addr: portNum}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Waiting for SIGINT (kill -2)
	<-stop

	err = server.Shutdown(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func initializeAssessClient() handler.AssessClient {
	assessClientUsername := os.Getenv("ASSESS_CLIENT_USER")
	assessClientPassword := os.Getenv("ASSESS_CLIENT_PASSWORD")
	assessClientBaseURL := os.Getenv("ASSESS_CLIENT_BASEURL")

	return client.NewAssessClient(client.AssessClientParams{
		BaseURL:  assessClientBaseURL,
		Username: assessClientUsername,
		Password: assessClientPassword,
	})
}

func initializeHandler(ac handler.AssessClient) Handler {
	return handler.NewHandler(handler.HandlerParams{
		AssessClient: ac,
	})
}
