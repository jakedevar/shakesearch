notes: i think handeling misspellings would be a great idea for the app
i can implement the fuzzy search function with the levenshtien distance algorithm and use a sliding scale on the frontend to determine how much error is 
allowed. by default i can try it out with two to three letters allowed to be mispelled. then i can have an exact match flag of some sort to signal to 
the determine search function if i should use the fuzzy search or hard search with regex. 

the next problem i see with this though is that the highlighted text in the frontend will be affected by this decision
i can't use an exact match because it's not returning an exact match. do i also implement the levenshtien algorithm on the frontend code? 
that would easily solve that match problem. but this will greatly reduce efficiency. i think that might be acceptable seeing as max its n will be 50.

1. either way i need to implement the levenshtien distance on the backend.
2. make a way to specify either a sliding scale or a boolean in the query string to determine fuzzy.
3. make sure that fuzzy searches are highlighted.

2. i think based on implementing the search 

1. how to implement fuzzy on the backend
  a. i can implement the levenshtien distance algorithm 
  b. i can split up the complete works that is stored in the Searcher struct so that it makes an array of strings
  c. then i can have this function return an array of results 
  d. i can then sort this results by the distance that was found in ascending order
  e. i can return this list and then do a search using the Suffix array to return lines with the words in it

  within the determine search function i can look for an exact match flag. if it is set to false then the default can just be a fuzzy search in the else block
  then i can double up on the logic inside to handle fuzzy search with case sensitive or not. this can be done later as i determine if it's needed

  what needs to be done is implementing the fuzzy search function in the else branch and passing in the searcTerm, the split works, and a 3
  then returning the sorted list. i would also need to either make the array unique or just make it more efficient with some logic in the fuzzy search algo. 
  this would go through the array and skip an iteration if the word has already been seen. this seems like it really is a must if i am to implement this
  as the levenshtien distance is O(n+m) so if i can skip some iteration this would udoubtably be good. the question is now where to put it? it looks like
  i should put it in the actuall fuzzySearch function as to reduce logic in the searchLogic file. and also i can even reduce how many times the levenshtien
  distance function is called. saving over half of iterations i would imagine as conjunctions would be handled
  1. implement a map in the fuzzySearch function (string, bool)
  2. if the word has been seen then skip iteration

  so now that the above has been implemented then the else logic needs to go into place right now i am spliting the works by a regular expression that splits
  . ,  ' and such. but now that i am thinking about it. what if you want to search with actuall punctuation? then i should just split by white space

  I have implemented the logic for the fuzzy search. it turns out that putting in a distance of 3 makes the search insanly slow. i think the find all index method
  is really slow even if it O(n + logn). the fuzzy search function is actually really fast it turns out. so i think what i need to do is paginate the search terms
  and then add a loading pinwheel. this way I can make the search results return faster  
  1. pass in the query struct to get the quantity and the page number
  2. implement pagination logic (should i extract that to a function?)
    upon further inspection i really just need to do like 5 results at a time. but how would i implement pagination logic so that all the results are returned
    eventually?

  

notes: just jutting down some thoughts here. if i return the actual distance that the word was from the mispelled word i can sort the results so that the ones
with the least amount of distance are sent first.
