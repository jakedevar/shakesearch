package main

import (
	"strings"
  "sort"
  "regexp"
  "unicode"
  "sync"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
	}

	for i := 1; i <= len(s1); i++ {
		matrix[i][0] = i
	}
	for j := 1; j <= len(s2); j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 0
			if s1[i-1] == s2[j-1] {
				cost = 0
			} else {
				cost = 1
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,
				min(
					matrix[i][j-1]+1,
					matrix[i-1][j-1]+cost,
				),
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}


type FuzzyResult struct {
	Value    string
	Distance int
}

type FuzzyResultsSlice []FuzzyResult

func (f FuzzyResultsSlice) Len() int {
  return len(f)
}
func (f FuzzyResultsSlice) Less(i, j int) bool {
  return f[i].Distance < f[j].Distance
}
func (f FuzzyResultsSlice) Swap(i, j int) {
  f[i], f[j] = f[j], f[i]
}

func fuzzySearch(searchTerm string, dataset []string, caseSensitive string) []FuzzyResult {
  lengthOfDataset := len(dataset)

  numGoRoutines := 10
  chunkSize := lengthOfDataset / numGoRoutines

  channelResults := make(chan []FuzzyResult)
  var wg sync.WaitGroup

  for i := 0; i < numGoRoutines; i++ {
    start := i * chunkSize
    end := (i + 1) * chunkSize

    if i == numGoRoutines-1 {
      end = lengthOfDataset
    }

    wg.Add(1)
    go func(start, end int) {
      defer wg.Done()
      partialResults := fuzzySearchPartial(searchTerm, dataset[start:end], caseSensitive)
      channelResults <- partialResults
    }(start, end)
  }

  go func() {
    wg.Wait()
    close(channelResults)
  }()

  results := []FuzzyResult{}
  for partialResults := range channelResults {
    results = append(results, partialResults...)
  }

  sort.Sort(FuzzyResultsSlice(results))

  resultsSliceLength := 5
  var end int
  if len(results) > resultsSliceLength {
    end = resultsSliceLength
  } else {
    end = len(results)
  }
  return results[:end]
}

func fuzzySearchPartial(searchTerm string, dataset []string, caseSensitive string) []FuzzyResult {
  sanitizedSearchTerm := sanitizeSearchTerm(searchTerm)
  lengthOfSanitizedSearchTerm := len(sanitizedSearchTerm)
  splitSearchTerm := strings.Split(searchTerm, " ")
  lengthOfSplitSearchTerm := len(splitSearchTerm)
  lengthOfDataset := len(dataset)
  threshold := 10 
  results := []FuzzyResult{}
  seenItems := make(map[string]bool)
  var distance int
  for i := lengthOfSplitSearchTerm; i < lengthOfDataset ; i++ {
    item := returnStringFromSlice(dataset[i-lengthOfSplitSearchTerm:i])
    sanitizedItem := sanitizeSearchTerm(item)
    lengthOfItem := len(sanitizedItem)
    if lengthOfSanitizedSearchTerm > lengthOfItem + 2 || lengthOfSanitizedSearchTerm < lengthOfItem - 2  || seenItems[item] {
      continue
    } else {
      seenItems[item] = true
    }

    searchTermFirstChar := rune(sanitizedSearchTerm[0])
    if caseSensitive == "true" {
      item = filterString(sanitizedItem, searchTermFirstChar)
      if item == "" {
        continue
      }
      distance = levenshteinDistance(sanitizedSearchTerm, sanitizedItem)
    } else {
      distance = levenshteinDistance(strings.ToLower(sanitizedSearchTerm), strings.ToLower(sanitizedItem))
    }
    if distance <= threshold {
      results = append(results, FuzzyResult{Value: item, Distance: distance})
    }
  }
  sort.Sort(FuzzyResultsSlice(results))

  resultsLength := 5
  var end int
  if len(results) > resultsLength {
    end = resultsLength
  } else {
    end = len(results)
  }
  return results[:end]
}

func filterString(item string, searchTermFirstChar rune) string {
  if len(item) == 0 {
    return ""
  }
  firstCharOfItem := rune(item[0])
  if unicode.IsUpper(searchTermFirstChar) == true {
    if unicode.IsUpper(firstCharOfItem) == false {
      return ""
    }
  } else {
    if unicode.IsUpper(firstCharOfItem) == true {
      return ""
    }
  }
  return item
}

func returnStringFromSlice(slices []string) string {
  slicesLength := len(slices)
  var result string
  for i, slice := range slices {
    if i == slicesLength - 1 {
      result += strings.TrimSpace(slice)
    } else {
      result += strings.TrimSpace(slice) + " "
    }
  }
  return result
}

func sanitizeSearchTerm(searchTerm string) string {
  reg := regexp.MustCompile(`[\W_]`)
  sanitizedSearchTerm := reg.ReplaceAllString(searchTerm, " ")
  sanitizedSearchTerm = strings.Replace(sanitizedSearchTerm, "'", "’", -1)
  return strings.TrimSpace(sanitizedSearchTerm)
}
