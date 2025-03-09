// COSC 581 HW5
// Zachary Perry
// March 11, 2025
// Problem 1: Greedy Algorithms -- Huffman Encoding

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
)

type node struct {
	char  byte
	freq  int
	left  *node
	right *node
}

// returns the root of the tree
func buildHuffmanTree(charFreqMap map[byte]int) *node {
	// loop through map, make nodes for them all
	chars := make([]byte, 0)
	nodes := make([]*node, 0)
	for c := range charFreqMap {
		chars = append(chars, c)
	}

	for _, c := range chars {
		newNode := &node{
			char: c,
			freq: charFreqMap[c],
		}
		nodes = append(nodes, newNode)
	}

	for len(nodes) > 1 {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].freq < nodes[j].freq
		})

		left := nodes[0]
		right := nodes[1]
		nodes = nodes[2:]

		newNode := &node{
			char:  0,
			freq:  left.freq + right.freq,
			left:  left,
			right: right,
		}

		nodes = append(nodes, newNode)
	}

	// return the root
	return nodes[0]
}

func generateHuffmanCodes(node *node, encodings map[byte]string, encoding string) {
	// traverse the tree and generate the generate codes
	// left = add a 0, right = add a 1 until we reach the right code
	// store this somewhere i guess
	if node == nil {
		return
	}

	// traverse each side, adding either 1 or zero to the code?
	if node.char != 0 {
		encodings[node.char] = encoding
	}

	generateHuffmanCodes(node.left, encodings, encoding+"0")
	generateHuffmanCodes(node.right, encodings, encoding+"1")
}

func convertStringtoBytes(fullEncodingString string) []byte {
	// Calculate the number of bytes needed (rounding up)
	numBytes := (len(fullEncodingString) + 7) / 8
	bytes := make([]byte, numBytes)

	// Process each byte
	for i := range numBytes {
		// Calculate start and end indices for this byte
		startIdx := i * 8
		endIdx := startIdx + 8
		if endIdx > len(fullEncodingString) {
			endIdx = len(fullEncodingString)
		}

		// Extract the bit string for this byte
		currBitString := fullEncodingString[startIdx:endIdx]

		// Pad with zeros if we don't have 8 bits
		if len(currBitString) < 8 {
			padding := make([]byte, 8-len(currBitString))
			for i := range padding {
				padding[i] = '0'
			}
			currBitString = currBitString + string(padding)
		}

		// Convert to a byte
		b, _ := strconv.ParseUint(currBitString, 2, 8)
		bytes[i] = byte(b)
	}

	return bytes
}

func writeEncodeFile(fileContents []byte, encodings map[byte]string) []byte {
	outputFile, err := os.Create("output/output_1.bin")
	if err != nil {
		log.Fatal("Error opening the output file")
	}

	defer outputFile.Close()

	stringBits := ""
	for _, c := range fileContents {
		stringBits += encodings[c]
	}

	bits2 := convertStringtoBytes(stringBits)

	outputFile.Write(bits2)

	return bits2
}

// printFreqTable will just print out the chars and their corresponding frequency of
// occurence within the file
func printFreqTable(charFreqMap map[byte]int) {
	chars := make([]byte, 0)
	for c := range charFreqMap {
		chars = append(chars, c)
	}

	slices.Sort(chars)

	fmt.Printf("\n")
	fmt.Printf("      ====================== \n")
	fmt.Printf("      |   CHAR  |  FREQ    |\n")
	fmt.Printf("      ====================== \n")

	for _, c := range chars {
		fmt.Printf("      |%6s   | %5d    |\n", string(c), charFreqMap[c])
	}
	fmt.Printf("       -------------------- \n")
	fmt.Printf("\n")
}

// readFile just takes the corresponding file name and opens it.
// It will then read in all chars within the file (lowercase, a-z).
// Each char will be stored in a map with the corresponding frequency of occurence.
// Additionally, everything read from the file is stored and returned in a byte slice.
func readFile(fileName string) (map[byte]int, []byte) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Error opening the file")
	}

	defer file.Close()

	// Scanner to read char by char.
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Building the freq table.
	charFreqMap := make(map[byte]int)
	fileContents := make([]byte, 0)
	for scanner.Scan() {
		char := byte(scanner.Text()[0])
		fileContents = append(fileContents, char)

		// If the current char is a lowercase letter, a-z, add/increment count in freq table.
		if char >= 97 && char <= 122 {
			charFreqMap[char]++
		}
	}

	return charFreqMap, fileContents
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: ./bin/problem_1 filename")
	}

	// Read in the input file, get file Stats.
	freq, fileContents := readFile(os.Args[1])
	inputFileInfo, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	printFreqTable(freq)
	fmt.Printf("Size of the file BEFORE compression: %5d bytes\n", inputFileInfo.Size())

	// Build out the Huffman tree and generate the codes, starting at the root.
	encodings := make(map[byte]string)
	root := buildHuffmanTree(freq)
	generateHuffmanCodes(root, encodings, "")

	// With the original file contents & encodings, encode the file.
	writeEncodeFile(fileContents, encodings)
	outputFileInfo, err := os.Stat("output/output_1.bin")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Size of the file AFTER compression : %5d bytes\n", outputFileInfo.Size())
}
