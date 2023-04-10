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
  "strconv"
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

type Params struct {
  Query string
  CaseSensitive string
  PageNumber int
  Quantity int
}

type Response struct {
  Results []string
  TotalResults int
}

func populateQueryParams(w http.ResponseWriter, r *http.Request) Params {
  query, ok := r.URL.Query()["query"]
  if !ok || len(query[0]) < 1 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("missing search query in URL params"))
    return Params {}
  }

  caseSensitive, ok := r.URL.Query()["caseSensitive"]
  if !ok || len(caseSensitive[0]) < 1 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("missing case sensitive in URL params"))
    return Params {}
  }

  pageNumber, ok := r.URL.Query()["pageNumber"]
  if !ok || len(pageNumber[0]) < 1 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("missing case sensitive in URL params"))
    return Params {}
  }
  parsedPageNumber, err := strconv.Atoi(pageNumber[0])
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("invalid page number"))
    return Params {}
  }

  quantity, ok := r.URL.Query()["quantity"]
  if !ok || len(quantity[0]) < 1 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("missing case sensitive in URL params"))
    return Params {}
  }
  parsedQuantity, err := strconv.Atoi(quantity[0])
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("invalid quantity"))
    return Params {}
  }

  return Params {
    Query: query[0],
    CaseSensitive: caseSensitive[0],
    PageNumber: parsedPageNumber,
    Quantity: parsedQuantity,
  }
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    params := populateQueryParams(w, r)
    results := searcher.Search(params.Query, params.CaseSensitive, params.PageNumber, params.Quantity)
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

func (s *Searcher) Search(query string, caseSensitive string, pageNumber int, quantity int) Response {
  idxs := s.determineSearch(query, caseSensitive)
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
  startPageNumber := (pageNumber - 1) * quantity
  endPageNumber := startPageNumber + quantity
  println(len(results), startPageNumber, endPageNumber)
  paginatedResults := results[startPageNumber:endPageNumber]
  resultsLength := len(results)
  return Response {
    Results: paginatedResults,
    TotalResults: resultsLength,
  }
}

func (s *Searcher) determineSearch(query string, caseSensitive string) [][]int {
  if caseSensitive == "true" {
    re := regexp.MustCompile(query)
    return s.SuffixArray.FindAllIndex(re, -1)
  } else {
    re := regexp.MustCompile("(?i)" + query)
    return s.SuffixArray.FindAllIndex(re, -1)
  }
}

func enableCORS(w *http.ResponseWriter) {
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
