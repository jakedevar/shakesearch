package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
  "pulley.com/shakesearch/searchLogic"
  "github.com/joho/godotenv"
)

func main() {
	searcher := searchLogic.Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
  if err != nil {
    println("Error loading .env file")
  }

  searcher.InitializeSearchCache()

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", searchLogic.HandleSearch(searcher))

  port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}


