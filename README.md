# COSC581-HW5
Contains code for homework 5 problems

## Folder Structure
```
- cmd/problem_1/problem_1.go: contains code for problem 1
- cmd/problem_2/problem_2.go: contains code for problem 2
- input/: contains sample input files for each problem
    - files starting with 1_ represent inputs for problem 1
    - files starting with 2_ represent inputs for problem 2
```

## How to run 
- _Compile using the Makefile_: `make`
- _Clean & remove bin files_: `make clean`
- _To run problem_1_:
```
./bin/problem_1 filename
```
***NOTE: This will output a binary file to the output/ directory***


- _To run problem_2_:
```
./bin/problem_2 filename
```

## Problem #1 
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

## Problem #2
2. Network Flow: Write an algorithm in a language of your choice that takes as input
a flow network in modified DIMACS format and prints all flow augmenting paths.
Ensure infinite looping isn’t possible. I will try to break your code.
Ex. File (tab separated values)

```
5 6
0 1 4
0 2 3
0 3 5
1 4 7
2 4 18
3 4 6
```

The first row is the vertex count and then the arc count
The remaining rows are the arcs: vertex 1, vertex 2, and arc weight
0 will always be the source vertex, and the vertex with the largest id will always
be the sink. In the above example that makes 0 the source and 4 the sink.

### Solution:
To accomplish this, I started with reading in any input files (accounting for the DIMACS formatting) and creating the graph. My approach to creating the graph involves using two specific structs: 
```Go
// Graph struct will represent all edges.
// Key in the map       -> vertex value.
// Val for each vertex  -> list of edges.
type Graph struct {
	edges map[int][]*Edge
}

// Edge struct represents an individual edge.
// to       -> destination vertex, where this edge goes to.
// capacity -> weight of the edge, capacity that can flow through it.
// flow     -> what is currently flowing through the edge.
type Edge struct {
	to       int
	capacity int
	flow     int
}
```
Once the graph is created, with all edges stored in the map, keyed on their vertex value, I then am able to find all of the flow augmenting paths via the Edmonds-Karp algorithm. Once all are found, I just output the corresponding paths, their flows, and the total flow (answer to given input shown below):
```
Augmenting Path #0  0 -> 1 -> 4 | Flow = 4
Augmenting Path #1  0 -> 2 -> 4 | Flow = 3
Augmenting Path #2  0 -> 3 -> 4 | Flow = 5

Total Flow: 12
```

## TODO
- part 1 
    - comments, clean up
    - explain in the readme
