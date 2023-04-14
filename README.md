# ShakeSearch

A working version of the application can be found [here](https://jakedevar-shakesearch.onrender.com/)

## New Features
These are some new features added to the ShakeSearch application.

### Backend
* Added fuzzy search feature that can handle misspelled words/sentances
* Preserved the ability to search for the exact match within the text
* Added case sensitive searches to both fuzy search and exact match
* Added caching using SQLite3 for search terms
* Paginated search results using quanity and page number
* Fixed out of bounds error when delivering search results within the 250 character range of found index
* Trimed whitespace of search term
* Moved search code and fuzzy search code to seperate local packages for better name spacing and readability

### Front End
* Made the search bar able to search without using a submit button for searches at the speed of thought
* Added checkboxes to enable case sensitive searches and exact match searches (or a combination of the two) 
* Added a quantity selector to specify how many search results to display
* Added page buttons at the top and bottom for the pagination of data and convinience
* Added error handeling for search results
* Added some documentation to provide a brief description of features
* Made the heading of the browser tab ShakeSearch
* Created title bar
* Added logos

# Design Decisions & Implementation Details
### Fuzzy Search Notable Details
The fuzzy search feature takes in the search term, the complete works text (split by strings.Fields()), and the caseSentive query parameter as arguments.
The split complete works slice is split up by 

The fuzzy search feature uses the Levenshtein distance algorithm to take in both the search term and the item from the text to be compared against. The function returns an integer representing how many letters are needed to be changed in order to create the search term.
The threshold I am using for this is 10, meaning that there is a tolerance of up to 10 characters that could be incorrect from the text in order for a given item from the complete works to return as a match.
If a match is returned it is placed in a struct with the Levenshtein distance as well as the string used and pushed into a results slice. The results slice is then sorted for items with the shortest distance and the first 5 items are returned to the search fuction to be searched for in the text using the arraysuffix.FindAllIndex() method and returned to the front end.

Given that this method of searching for words similar to the search term requires that the entirety of the complete works is searched. As well as the big O of the Levenshtein distance algorithm is O(n * m) this is by nature a very ineficcient process.
To combat this, I used go routines to split the complete works up into smaller peices




# Future Changes

There are several things that I would like to address if I had more time.

