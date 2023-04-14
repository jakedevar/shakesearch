# ShakeSearch

A working version of the application can be found [here](https://jakedevar-shakesearch.onrender.com/)

## New Features
These are some new features added to the ShakeSearch application.

### Backend
* Added fuzzy search feature that can interpret misspelled words/sentences
* Preserved the ability to search for exact matches within the text
* Added case-sensitive searches to both fuzzy search and exact match
* Implemented caching using SQLite3 for search terms
* Paginated search results using quantity and page number
* Fixed out-of-bounds error when delivering search results within the 250-character range of the found index
* Trimmed whitespace of search term
* Moved search code and fuzzy search code to separate local packages for better namespacing and readability

### Front End
* Made it possible for search bar to search automatically after user has finished typing (0.5 second delay)
* Added checkboxes to enable case-sensitive searches and exact match searches (or a combination of the two)
* Added a quantity selector enabling user to specify how many search results to display
* Added page buttons at the top and bottom for the pagination of data and convenience
* Implemented error handling for search results
* Added documentation to provide a brief description of features
* Changed the heading of the browser tab to ShakeSearch
* Created a title bar
* Added logos

# Design Decisions & Implementation Details
### Fuzzy Search Notable Details
The fuzzy search feature uses the Levenshtein distance algorithm to take in both the search term and the item from the text to be compared against.
The function returns an integer representing how many letters need to be changed to create the search term.
There is a tolerance of up to 10 characters that could be incorrect from the text in order for a given item from the complete works to return as a match.
If a match is returned, it is placed in a struct with the Levenshtein distance as well as the string used and pushed into a results slice.
The results slice is then sorted for items with the shortest distance, and the first 5 items are returned to the search function to be searched for in the text using the arraysuffix.FindAllIndex() method and returned to the front end.

At it's base, the fuzzy search feature is a somewhat slow and computationaly expensive process. To combat this I:

* Used Go routines to split the complete works up into smaller pieces and run the fuzzy search partial function in parallel.
* Used a map to store any term within the text that has already been searched in order to save time on future iterations.
* Created a string via slicing the complete works by the length of the search term split by spaces - This ensures that I can search by whole words instead of going by character index.
* Lastly, the length of the returned string is compared against the length of the search term using a margin of error of 2 characters longer or shorter than the given search term. If the item does not meet these length requirements, the iteration is skipped.

### Caching
To avoid expensive search operations, I have implemented a form of caching using SQLite3, which is an in-filesystem database to store any terms that have been searched before. This allows for quick page loads, as the entirety of the results from the first search has been stored in the SQLite3 cache.


# Future Changes
There are several things that I would like to address if I had more time:

### Search by Character & Display Character Lines
Searching by character would be a nice feature to have for someone who is trying to look up lines spoken by a specific character. Additionally, I would like the search results on the frontend to be displayed with new lines between a character speaking lines for easier readability.

### Search by Play & Scene
In addition to the regular search function, the ability to search for a term in a particular play would be useful. Preferably, I would implement this in a way that would also make the search feature function even faster by passing a query parameter to the backend to search for text within a given scene from a specific play.

### Search Results
I noticed during testing that there were a few phrases that did not show up when searching, most notably "We are such stuff As dreams are made on." If given more time, I would like to find a solution to this problem.

When you split the previous search into either "We are such stuff" or "As dreams are made on," search results are found. However, if you search the full term, nothing is found. The fuzzy search function does recognize that this sentence exists within the text. Though, searches for it, as well as other longer sentences, are not found, either from within the suffix array or using a regexp method to search the complete works as a string.

### Refactoring
I would like to refactor the code structure to improve readability and efficiency in both the frontend and backend components.

### Testing
I would like to add tests to both frontend and backend components.

### Sytem Resources
If given the option, I would increase the CPU capability of the application server to take advantage of the Go routine functionality in the fuzzy search code.
