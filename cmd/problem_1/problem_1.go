package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

// TODO: All of the huffman encoding shit
/*
Greedy Algorithms: Write an algorithm in the language of your choice to provide
a Huffman encoding from the contents of a .txt file (just the 26 lower case letters
of the alphabet in some arbitrary combination and length) and apply it to
compress the file. Your program should output 2 things, the frequency table for
the characters, and the size of the file (in bytes) before and after the encoding.
Note: You will build the frequency table based on the content of the file.
*/

// General Approach
/*
1. Read in the file and count frequency of each char
2. Then, build a Huffman tree for each char based on their frequency (look up algo for this)
  - NOTE: this will be the hard part
3. Then, traverse this tree, creating the encoding for each letter
4. Store the encoding, (letter as the key, encoded value as the val
5. When printing, traverse the original file again, and look up each letter to get the encoding
*/

func buildHuffmanTree() {
	// create leaf nodes for all chars with their freq as the weight
	// place nodes in a min heap ordered by frequency
	// repeat:
	// extract 2 lowest freq nodes
	// create new internal node with these two as children
	// new node freq is sum of children freq
	// insert new node back into the heap
	// continue
}

func generateHuffmanCodes() {

	// traverse the tree and generate the generate codes
	// left = add a 0, right = add a 1 until we reach the right code
	// store this somewhere i guess
}

func encodeFile() {}

// printFreqTable will just print out the chars and their corresponding frequency of
// occurence within the file
func printFreqTable(charFreqMap map[byte]int) {
	chars := make([]byte, 0)
	for c := range charFreqMap {
		chars = append(chars, c)
	}
	slices.Sort(chars)

	fmt.Println("Length of freq: ", len(charFreqMap))
	fmt.Println("CHAR | FREQ")
	fmt.Println("============")

	for _, c := range chars {
		fmt.Printf("%-4s | %4d\n", string(c), charFreqMap[c])
	}
}

// readFile just takes the corresponding file name and opens it
// It will then read in all chars within the file (lowercase)
// Each char will be stored in a map with the corresponding frequency of occurence
func readFile(fileName string) map[byte]int {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening the file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	charFreqMap := make(map[byte]int)
	for scanner.Scan() {
		char := byte(scanner.Text()[0])

		if char >= 97 && char <= 122 {
			charFreqMap[char]++
		}
	}

	return charFreqMap
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./bin/problem_1 filename")
	}

	// read and print
	freq := readFile(os.Args[1])
	printFreqTable(freq)

	// file size before compression
	fileInfo, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Size of the file before compression: ", fileInfo.Size(), " bytes")
}
