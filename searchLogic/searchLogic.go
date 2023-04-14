package searchLogic

import (
  "pulley.com/shakesearch/fuzzySearch"
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"net/http"
  "regexp"
  "strconv"
  "strings"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

type Searcher struct {
  CompleteWorks string
  SuffixArray *suffixarray.Index
  SearchCache *sql.DB
}

type Params struct {
  SearchTerm string
  CaseSensitive string
  PageNumber int
  Quantity int
  ExactMatch string
}

type Response struct {
  Results SearchResults
  TotalResults int
}

type SearchResult struct {
  SearchTerm string
  Line string
}

type SearchResults []SearchResult

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

  exactMatch, ok := r.URL.Query()["exactMatch"]
  if !ok || len(exactMatch[0]) < 1 {
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

  trimWhiteSpace := strings.TrimSpace(searchTerm[0]) 
  sanitizedSearchTerm := strings.Replace(trimWhiteSpace, "'", "’", -1) 
  return Params {
    SearchTerm: sanitizedSearchTerm,
    CaseSensitive: caseSensitive[0],
    PageNumber: parsedPageNumber,
    Quantity: parsedQuantity,
    ExactMatch: exactMatch[0],
  }
}

func HandleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    params := populateQueryParams(w, r)
    results := searcher.Search(params.SearchTerm, params.CaseSensitive, params.PageNumber, params.Quantity, params.ExactMatch)
    buf := &bytes.Buffer{}
    enc := json.NewEncoder(buf)
    err := enc.Encode(results)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte("encoding failure"))
      return
    }
    // enableCORS(&w)
    w.Header().Set("Content-Type", "application/json")
    w.Write(buf.Bytes())
  }
}

func (s *Searcher) InitializeSearchCache() {
  db, err := sql.Open("sqlite3", "./searchCache.db")
  if err != nil {
    println("Error opening search cache db")
  }

  tableDeleteStatement, err := db.Prepare("DROP TABLE IF EXISTS searchCache")
  if err != nil {
    println("Error preparing table delete statement")
  }
  tableDeleteStatement.Exec()
  if err != nil {
    println("Error deleting table")
  }

  tableCreateStatment, err := db.Prepare("create table searchCache (id integer not null primary key, storedSearchTerm text, searchResult text)")
  _, err = tableCreateStatment.Exec()
  if err != nil {
    println("Error creating table")
  }
   
  s.SearchCache = db
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

func (s *Searcher) Search(searchTerm string, caseSensitive string, pageNumber int, quantity int, exactMatch string) Response {
  
  rows, _ := s.SearchCache.Query("SELECT * FROM searchCache WHERE storedSearchTerm = ?", searchTerm + caseSensitive + exactMatch)
  defer rows.Close()

  var id int
  var storedSearchTerm string
  var searchResult string
  for rows.Next() {
    rows.Scan(&id, &storedSearchTerm, &searchResult)
  }
  var results SearchResults
  if id > 0 {
    err := json.Unmarshal([]byte(searchResult), &results)
    if err != nil {
      println("Error parsing search results")
    }
  } else {
    results = s.determineSearch(searchTerm, caseSensitive, exactMatch)
    jsonResults, err := json.Marshal(results)
    if err != nil {
      println("Error marshalling results")
    }
    insertStatement := "INSERT INTO searchCache (storedSearchTerm, searchResult) VALUES (?, ?)"
    _, err = s.SearchCache.Exec(insertStatement, searchTerm + caseSensitive + exactMatch, string(jsonResults))
    if err != nil {
      println("Error inserting into search cache")
    }
  }
  startPageNumber := (pageNumber - 1) * quantity
  resultsLength := len(results)
  if resultsLength == 0 {
    startPageNumber = 0
  }

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

func (s *Searcher) determineSearch(searchTerm string, caseSensitive string, exactMatch string) SearchResults {
  if exactMatch == "true" {
    searchTerm = regexp.QuoteMeta(searchTerm)
    return searchExactMatch(s, searchTerm, caseSensitive)
  } else {
    return fuzzySearchResults(s, searchTerm, caseSensitive)
  }
}

func searchExactMatch(s *Searcher, searchTerm string, caseSensitive string) SearchResults {
  var rePattern string
  if caseSensitive == "true" {
    rePattern = searchTerm
  } else {
    rePattern = "(?i)" + searchTerm
  }
  re := regexp.MustCompile(rePattern)
  idxs := s.SuffixArray.FindAllIndex(re, -1)
  return createSearchResultsArray(s, idxs, searchTerm)
}

func fuzzySearchResults(s *Searcher, searchTerm string, caseSensitive string) SearchResults {
  results := []SearchResult{}
  splitWorks := strings.Fields(s.CompleteWorks)
  // splitWorks := strings.Split("Hamlet’s mother, now wife of Claudius. POLONIUS, Lord Chamberlain.", " ")
  var fuzzySearchResults []fuzzySearch.FuzzyResult
  fuzzySearchResults = fuzzySearch.FuzzySearch(searchTerm, splitWorks, caseSensitive)
  for _, item := range fuzzySearchResults {
    fuzzySearchTerm := regexp.QuoteMeta(item.Value)  
    if caseSensitive == "true" {
      fuzzySearchTerm = "[^a-zA-Z]" + fuzzySearchTerm + "[^a-zA-Z]"
    } else {
      fuzzySearchTerm = "(?i)" + fuzzySearchTerm
    }
    re := regexp.MustCompile(fuzzySearchTerm)
    idxs := s.SuffixArray.FindAllIndex(re, -1)
    results = append(results, createSearchResultsArray(s, idxs, item.Value)...)
  }
  return results
}

func createSearchResultsArray(s *Searcher, idxs [][]int, searchTerm string) SearchResults {
  results := []SearchResult{}
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
    searchResult := SearchResult {
      SearchTerm: searchTerm,
      Line: s.CompleteWorks[start:end],
    }
    results = append(results, searchResult)
  }
  return results
}

// func enableCORS(w *http.ResponseWriter) {
//   (*w).Header().Set("Access-Control-Allow-Origin", "*")
//   (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
//   (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// }
