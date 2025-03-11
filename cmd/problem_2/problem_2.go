// COSC 581 HW5
// Zachary Perry
// March 11, 2025
// Problem 2: Network Flow -- Augmenting Paths
package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

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

// AugPath struct represents a found augmented path.
// pathVertices -> list containing all vertex values in the augmented path.
// flow         -> flow the augmented path can support.
type AugPath struct {
	pathVertices []int
	flow         int
}

// NewGraph() just creates a new graph instance.
// The map will hold vertex values as the keys and their outgoing edges in a slice as the val.
// NOTE: the edges slice for each vertex will also contain any backwards edges.
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
	// Forward edge (v1 -> v2).
	forward := &Edge{
		to:       v2,
		capacity: weight,
		flow:     0,
	}

	// Backward edge (v2 -> v1).
	backward := &Edge{
		to:       v1,
		capacity: 0,
		flow:     0,
	}

	// Append to respective vertex Edge slice.
	graph.edges[v1] = append(graph.edges[v1], forward)
	graph.edges[v2] = append(graph.edges[v2], backward)
}

// BFS will find the augmented path from the source to the sink.
// For each vertex, it will look at all outgoing edges and follow those that have a valid capacity.
// It will mark verticies as visited to avoid any infinite loops + cycles.
// Once it has reached the sink, it will return the augmented path.
func (graph *Graph) BFS(source, sink int) ([]int, bool) {
	// visited tracks visited nodes, saved path will save the current path from source -> sink if one is found.
	visited := make(map[int]bool)
	savedPath := make(map[int]int)

	queue := make([]int, 0)
	queue = append(queue, source)
	visited[source] = true

	for len(queue) > 0 {
		// Pop off the queue.
		currVertex := queue[0]
		queue = queue[1:]

		for _, edge := range graph.edges[currVertex] {
			next := edge.to
			// If we have not visited this edge yet and the capacity is > flow.
			// Set as visited, save into the path map, and append to the queue.
			if !visited[next] && edge.capacity > edge.flow {
				visited[next] = true
				savedPath[next] = currVertex
				queue = append(queue, next)

				// When a path is found, create the final path found by looping backwards from sink to source using the saved path.
				if next == sink {
					finalPath := []int{sink}
					for v := sink; v != source; v = savedPath[v] {
						finalPath = append([]int{savedPath[v]}, finalPath...)
					}

					return finalPath, true
				}
			}
		}
	}

	return nil, false
}

// AugmentingPaths uses the Edmonds-Karp algorithm.
// It will use BFS to help find the shortest augmenting paths in the given network graph.
func (graph *Graph) AugmentingPaths(source, sink int) []*AugPath {
	allAugmentingPaths := make([]*AugPath, 0)

	for {
		// Call BFS, look for a path.
		path, found := graph.BFS(source, sink)
		if !found {
			break
		}

		// Find the maximum amount of flow that can be sent through the path (minimum capacity).
		// Find edge from v1 to v2.
		// If found, calculate the difference in capacity and flow, update minimumCapacity.
		minimumCapacity := math.MaxInt
		for i := range len(path) - 1 {
			v1, v2 := path[i], path[i+1]
			for _, edge := range graph.edges[v1] {
				if edge.to == v2 {
					diff := edge.capacity - edge.flow
					if diff < minimumCapacity {
						minimumCapacity = diff
					}
					break
				}
			}
		}

		// Update the flow values.
		// Increase the flow on the forward edge (v1->v2) by the found minimumCapacity.
		// Decrease the flow on the backward edge (v2->v1) by the found minimumCapacity.
		for i := range len(path) - 1 {
			v1, v2 := path[i], path[i+1]

			for _, edge := range graph.edges[v1] {
				if edge.to == v2 {
					edge.flow += minimumCapacity
					break
				}
			}

			for _, edge := range graph.edges[v2] {
				if edge.to == v1 {
					edge.flow -= minimumCapacity
					break
				}
			}
		}

		// Save the augmented path, its vertices, and flow.
		newAugmentedPath := &AugPath{
			pathVertices: path,
			flow:         minimumCapacity,
		}

		allAugmentingPaths = append(allAugmentingPaths, newAugmentedPath)
	}

	return allAugmentingPaths
}

// printAugmentingPaths() just prints out the found augmenting paths.
// Will print the path (vertices), along with the flow.
// Also outputs the total flow.
func printAugmentingPaths(paths []*AugPath) {
	totalFlow := 0

	for i, p := range paths {
		currPath := p.pathVertices
		currFlow := p.flow

		fmt.Printf("Augmenting Path #%d ", i)
		for j, v := range currPath {
			if j < len(currPath)-1 {
				fmt.Printf(" %d ->", v)
			} else {
				fmt.Printf(" %d ", v)
			}
		}
		fmt.Printf(" | Flow = %d\n", currFlow)
		totalFlow += currFlow
	}

	fmt.Printf("\nTotal Flow: %d\n", totalFlow)
}

// readFile will read a file containing a flow network in DIMACS format.
// Creates the graph, adding the defined edges and vertices.
func readFile(fileName string) (*Graph, int, int) {
	// Reading in the file contents.
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening the file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		log.Fatal("readFile(): Issue scanning the file -- may be empty..")
	}

	// First row (2 values) -> vertex count and then the arc count.
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

	graph := NewGraph(vertexCount)
	source := 0
	sink := 0
	arcReadCount := 1

	// Read in the remaining rows containing the edges (v1, v2, arc weight).
	for scanner.Scan() {

		// If there's any extra rows that are not included in the arcCount defined, break.
		if arcReadCount > arcCount {
			break
		}

		line := strings.Fields(scanner.Text())
		if len(line) < 3 || len(line) > 3 {
			log.Fatal("readFile(): vertex and weight row either has too few or too many arguments")
		}

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

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./bin/problem_1 filename")
	}

	// Read the file, get the graph.
	graph, source, sink := readFile(os.Args[1])

	// Get the augmented paths and print them out.
	augmentedPaths := graph.AugmentingPaths(source, sink)
	printAugmentingPaths(augmentedPaths)
}
