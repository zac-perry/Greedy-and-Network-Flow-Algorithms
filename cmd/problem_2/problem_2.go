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

func readFile(fileName string) {
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
	fmt.Printf("VERTEX COUNT: %6d\nARCCOUNT: %10d\n", vertexCount, arcCount)

	// Remaining rows: vertex 1, vertex 2, and arc weight
	// - 0 is always the source vertex
	// - vertex with the largest id will always be the sink
	// for example -> 0 is the source, 4 is the sink, 3rd value are all weights

	// TODO, make sure to only scan for the given number of arc counts
	// keep track of the sink value
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
		fmt.Println(line)

		// decode each value, keep track of the source and sink?
		v1, err := strconv.Atoi(line[0])
		if err != nil {
			// NOTE: maybe end the iteration and read in the rest of the lines?
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

		fmt.Printf("  v1: %10d\n  v2: %10d\n  weight: %6d\n", v1, v2, weight)

		// Track the sink
		if v2 > sink {
			sink = v2
		}
		arcReadCount++
	}

	fmt.Printf("\nSourceVertex: %d\nSinkVertex: %d\n", source, sink)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./bin/problem_1 filename")
	}

	// Read in the input file, get file Stats.
	readFile(os.Args[1])
}
