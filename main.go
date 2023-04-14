package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
  "pulley.com/shakesearch/searchLogic"
)

func main() {
	searcher := searchLogic.Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

  searcher.InitializeSearchCache()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", searchLogic.HandleSearch(searcher))
	http.HandleFunc("/ping", handlePing)

  port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handlePing(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "pong")
}

