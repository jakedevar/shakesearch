notes: so i need to return an array of obects so that the hightlight search function can work
that means i need to create a search object type and then change the arrays and return object of the functions to match as well

i can make a function that maps a given array of sentances and returns an array of searhResult objects with the string and the match term inside
1. create SearchResult Type
2. create function that takes in an array of index pairs and transforms them into an aryay of SearchResult objects/structs with the search term and sentence

1. type 
  a. has SearchTerm prop and Line prop
  b. change Result type to have an array of SearchResult or SearchResults

2. function takes in array of strings and the search term
  a. take logic from search function and implement fucntion logic with it
  b. call array .map and map objects into an array. or build a new array?
  c. return array and deliver to frontend

