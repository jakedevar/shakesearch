package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	// "log"
	"net/http"
	// "os"
  "regexp"
  "strconv"
  "strings"
)

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

type Params struct {
  SearchTerm string
  CaseSensitive string
  PageNumber int
  Quantity int
}

type Response struct {
  Results []string
  TotalResults int
}

func populateQueryParams(w http.ResponseWriter, r *http.Request) Params {
  searchTerm, ok := r.URL.Query()["searchTerm"]
  if !ok || len(searchTerm[0]) < 1 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte("missing search searchTerm in URL params"))
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
    SearchTerm: searchTerm[0],
    CaseSensitive: caseSensitive[0],
    PageNumber: parsedPageNumber,
    Quantity: parsedQuantity,
  }
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    params := populateQueryParams(w, r)
    results := searcher.Search(params.SearchTerm, params.CaseSensitive, params.PageNumber, params.Quantity)
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

func (s *Searcher) Search(searchTerm string, caseSensitive string, pageNumber int, quantity int) Response {
  idxs := s.determineSearch(searchTerm, caseSensitive)
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
  resultsLength := len(results)

  var endPageNumber int
  if startPageNumber + quantity > resultsLength {
    endPageNumber = resultsLength
  } else {
    endPageNumber = startPageNumber + quantity
  }

  paginatedResults := results[startPageNumber:endPageNumber]
  return Response {
    Results: paginatedResults,
    TotalResults: resultsLength,
  }
}

func (s *Searcher) determineSearch(searchTerm string, caseSensitive string) [][]int {
  // if caseSensitive == "true" {
  //   searchTerm = "[^a-zA-Z]" + searchTerm + "[^a-zA-Z]"
  //   re := regexp.MustCompile(searchTerm)
  //   return s.SuffixArray.FindAllIndex(re, -1)
  // } else {
    return fuzzySearchResults(s, searchTerm, caseSensitive)
  // }
}

func fuzzySearchResults(s *Searcher, searchTerm string, caseSensitive string) [][]int {
  results := [][]int{}
  splitWorks := strings.Split(s.CompleteWorks, " ")
  var fuzzySearchResults []FuzzyResult
  fuzzySearchResults = fuzzySearch(searchTerm, splitWorks, caseSensitive)
  println(fuzzySearchResults[0].Value, fuzzySearchResults[1].Value, fuzzySearchResults[2].Value, fuzzySearchResults[3].Value, fuzzySearchResults[4].Value)
  for _, item := range fuzzySearchResults {
    fuzzySearchTerm := "[^a-zA-Z]" + item.Value + "[^a-zA-Z]"
    if caseSensitive == "true" {
      fuzzySearchTerm = "[^a-zA-Z]" + fuzzySearchTerm + "[^a-zA-Z]"
    } else {
      fuzzySearchTerm = "(?i)[^a-zA-Z]" + fuzzySearchTerm + "[^a-zA-Z]"
    }
    println(fuzzySearchTerm)
    re := regexp.MustCompile(fuzzySearchTerm)
    results = append(results, s.SuffixArray.FindAllIndex(re, -1)...)
  }
  return results
}

func enableCORS(w *http.ResponseWriter) {
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
