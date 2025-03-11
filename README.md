# COSC581-HW5
Contains code for homework 5 problems

## Folder Structure
```
- cmd/problem_1/problem_1.go: contains code for problem 1
- cmd/problem_2/problem_2.go: contains code for problem 2
```

## How to run 
Compile using the Makefile: `make`

Clean & remove bin files: `make clean`

To run problem_1: `./bin/problem_1`

To run problem_2: `./bin/problem_2`

## Problem 1 
Greedy Algorithms: Write an algorithm in the language of your choice to provide
a Huffman encoding from the contents of a .txt file (just the 26 lower case letters
of the alphabet in some arbitrary combination and length) and apply it to
compress the file. Your program should output 2 things, the frequency table for
the characters, and the size of the file (in bytes) before and after the encoding.
Note: You will build the frequency table based on the content of the file.

### Solution:
MinHeap -- TODO: discuss go implementation for this just to make sure it's covered
Then, go through the algorithm itself
TODO: mention the padding on the last byte in the encoding, also explain this part (concating everything, etc)

## Problem 2
2. Network Flow: Write an algorithm in a language of your choice that takes as input
a flow network in modified DIMACS format and prints all flow augmenting paths.
Ensure infinite looping isn’t possible. I will try to break your code.
Ex. File (tab separated values)

5 6
0 1 4
0 2 3
0 3 5
1 4 7
2 4 18
3 4 6

The first row is the vertex count and then the arc count
The remaining rows are the arcs: vertex 1, vertex 2, and arc weight
0 will always be the source vertex, and the vertex with the largest id will always
be the sink. In the above example that makes 0 the source and 4 the sink.

## TODO
- part 1 
    - comments, clean up
    - explain in the readme
- part 2 
    - actually do it
    - comments
    - explain in the readme
