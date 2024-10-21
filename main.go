package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/joho/godotenv"

	"github.com/coderavels/upworkassessment/server/client"
	"github.com/coderavels/upworkassessment/server/handler"
)

type Handler interface {
	ListBooks(w http.ResponseWriter, r *http.Request)
	GetBookCollection(w http.ResponseWriter, r *http.Request)
}

func main() {
	log.Println("Starting our bookshelf backend server")
	err := godotenv.Load()

	port := os.Getenv("PORT")
	portNum := fmt.Sprintf(":%s", port)

	ac := initializeAssessClient()
	h := initializeHandler(ac)

	// ui paths
	// index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/index.html")
	})
	// server static assets from "ui/static" directory
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static"))))

	// api paths
	http.HandleFunc("/api/v1/books", h.ListBooks)
	http.HandleFunc("/api/v1/collection/{bookISBN}", h.GetBookCollection)

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
