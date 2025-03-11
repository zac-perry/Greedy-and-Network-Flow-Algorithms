package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// TODO: network flow
/*
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
*/

type Graph struct {
	edges map[int][]*Edge
}

type Edge struct {
	to       int
	capacity int
	flow     int
}

func NewGraph(vertexCount int) *Graph {
  graph := &Graph{
		edges: make(map[int][]*Edge),
	}

  for i := range vertexCount {
   graph.edges[i] = []*Edge{} 
  }

  return graph
}

// AddEdge() will add an edge from v1 to v2 with the specified weight as the capacity to the graph.
// Additionally, it also adds a backwards edge (with 0 capacity).
func (graph *Graph) AddEdge(v1, v2, weight int) {
	// Forward edge (v1 -> v2)
	forward := &Edge{
		to:       v2,
		capacity: weight,
		flow:     0,
	}

	// Backward edge (v2 -> v1)
	backward := &Edge{
		to:       v1,
		capacity: 0,
		flow:     0,
	}

	graph.edges[v1] = append(graph.edges[v1], forward)
	graph.edges[v2] = append(graph.edges[v2], backward)
}

func (graph *Graph) BFS() {}

func (graph *Graph) FordFulkerson() {}



// readFile will read a file containing a flow network in DIMACS format.
// While reading this, it will go ahead and create the
func readFile(fileName string) (*Graph, int, int) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening the file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		log.Fatal("readFile(): Issue scanning the file -- may be empty..")
	}

	// first row (2 values) -> vertex count and then the arc count
	firstLine := strings.Fields(scanner.Text())
	if len(firstLine) < 2 || len(firstLine) > 2 {
		log.Fatal(
			"readFile(): Error parsing the first line -- either too few arguments or too many. Should only include vertex count and arc count",
		)
	}

	vertexCount, err := strconv.Atoi(firstLine[0])
	if err != nil {
		log.Fatal("readFile(): Error reading in the vertexCount")
	}

	arcCount, err := strconv.Atoi(firstLine[1])
	if err != nil {
		log.Fatal("readFile(): Error reading in the arcCount")
	}
	//fmt.Printf("VERTEX COUNT: %6d\nARCCOUNT: %10d\n", vertexCount, arcCount)

	// Remaining rows: vertex 1, vertex 2, and arc weight
	// - 0 is always the source vertex
	// - vertex with the largest id will always be the sink
	// for example -> 0 is the source, 4 is the sink, 3rd value are all weights
	graph := NewGraph(vertexCount)
	source := 0
	sink := 0
	arcReadCount := 1
	for scanner.Scan() {
		// if there's any extra rows that are not included in the arcCount defined, break.
		if arcReadCount > arcCount {
			break
		}

		line := strings.Fields(scanner.Text())
		if len(line) < 3 || len(line) > 3 {
			// NOTE: maybe break here when one of the lines is fucked?
			log.Fatal("readFile(): vertex and weight row either has too few or too many arguments")
		}
		//fmt.Println(line)

		// decode each value, keep track of the source and sink?
		v1, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal("readFile(): Error reading in v1")
		}

		v2, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal("readFile(): Error reading in the v2")
		}

		weight, err := strconv.Atoi(line[2])
		if err != nil {
			log.Fatal("readFile(): Error reading in the weight")
		}

		// fmt.Printf("  v1: %10d\n  v2: %10d\n  weight: %6d\n", v1, v2, weight)

		// Track the sink (max v2 found).
		if v2 > sink {
			sink = v2
		}

		// Create and add edge to the graph (both forward and backward)
		graph.AddEdge(v1, v2, weight)
		arcReadCount++
	}

	return graph, source, sink
}

func printData(graph *Graph, source, sink int) {
	for k, v := range graph.edges {
		fmt.Println("Key: ", k)
		for _, k := range v {
			fmt.Println(" going to vertex: ", k.to)
			fmt.Println(" capacity       : ", k.capacity)
		}
	}
	fmt.Println("source: ", source)
	fmt.Println("sink: ", sink)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./bin/problem_1 filename")
	}

	// Read in the input file, get file Stats.
	graph, source, sink := readFile(os.Args[1])
	printData(graph, source, sink)

	// now, can do ford fulkerson and bfs
}
