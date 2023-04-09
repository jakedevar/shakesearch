package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
  "regexp"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

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

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
    println(r.URL.Query())
		caseInsensitive, ok := r.URL.Query()["caseInsensitive"]
		if !ok || len(caseInsensitive[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing caseInsensitive in URL params"))
			return
		}
    
		results := searcher.Search(query, caseInsensitive)
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
    enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *Searcher) Search(query []string, caseInsensitive []string) []string {
	idxs := s.determineSearch(query[0], caseInsensitive[0])
	results := []string{}
	for _, idxPair := range idxs {
    start, end := idxPair[0], idxPair[1]
    if start > 250 {
      start = start - 250
    } else {
      start = 0
    }

    if end > len(s.CompleteWorks) - 250 {
      end = len(s.CompleteWorks) 
    } else {
      end = end + 250
    }
		results = append(results, s.CompleteWorks[start:end])
	}
  // println(len(results))
	return results
}

func (s *Searcher) determineSearch(query string, caseInsensitive string) [][]int {
  if caseInsensitive == "true" {
    re := regexp.MustCompile("(?i)" + query)
    return s.SuffixArray.FindAllIndex(re, -1)
  } else {
    re := regexp.MustCompile(query)
    return s.SuffixArray.FindAllIndex(re, -1)
  }
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
