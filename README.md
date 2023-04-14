# ShakeSearch

A working version of the application can be found [here](https://jakedevar-shakesearch.onrender.com/)

##New Features
These are some new features added to the ShakeSearch application.

###Backend
* Added fuzzy search feature that can handle misspelled words/sentances
* Preserved the ability to search for the exact match within the text
* Added case sensitive searches to both fuzy search and exact match
* Added caching using SQLite3 for search terms
* Paginated search results using quanity and page number
* Fixed out of bounds error when delivering search results within the 250 character range of found index
* Trimed whitespace of search term

##Front End
* Made the search bar able to search without using a submit button for searches at the speed of thought
* Added buttons to enable case sensitive searches and exact match searches (or a combination of the two) 
* Added a quantity amount to how many results would be delivered to the front end
* Added page buttons at the top and bottom for the pagination of data and convinience
* Added error handeling for search results
* Made the title of the browser tab ShakeSearch
* Added a logo

