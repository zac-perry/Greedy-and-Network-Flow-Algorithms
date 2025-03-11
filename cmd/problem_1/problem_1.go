// COSC 581 HW5
// Zachary Perry
// March 11, 2025
// Problem 1: Greedy Algorithms -- Huffman Encoding

package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
)

// node struct represents a node in the MinHeap.
// It will hold the 
type node struct {
	char  byte
	freq  int
	left  *node
	right *node
}

// NOTE: All of this heap stuff is required in order to use Go's stdlib implementation of heap
// In this case, defining a type and implementing the proper interface methods (Go's library will then call certain things when needed)
// for example, when calling push, Go stdlib container/heap will call len, less, and swap to make sure it retains its heap property (in this case, min)
type MinHeap []*node

func (minHeap MinHeap) Len() int {
	return len(minHeap)
}

func (minHeap MinHeap) Less(i, j int) bool {
	return minHeap[i].freq < minHeap[j].freq
}

func (minHeap MinHeap) Swap(i, j int) {
	minHeap[i], minHeap[j] = minHeap[j], minHeap[i]
}

func (minHeap *MinHeap) Push(x any) {
	*minHeap = append(*minHeap, x.(*node))
}

func (minHeap *MinHeap) Pop() any {
	placeholder := *minHeap
	length := len(placeholder)
	value := placeholder[length-1]
	*minHeap = placeholder[0 : length-1]

	return value
}

// Builds out the huffman tree using a minheap
func buildHuffmanTree(charFreqMap map[byte]int) *node {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	// loop and push onto the heap
	for k, v := range charFreqMap {
		newNode := &node{
			char: k,
			freq: v,
		}
		heap.Push(minHeap, newNode)
	}

	// NOTE: to access and index, need to do (*minHeap)[i]
	// This is the huffman tree builder man part
	// Makes a new node (freq is left + right)
	// Left and right are the two shmallest atm
	for minHeap.Len() > 1 {
		left := heap.Pop(minHeap).(*node)
		right := heap.Pop(minHeap).(*node)

		newNode := &node{
			char:  0,
			freq:  left.freq + right.freq,
			left:  left,
			right: right,
		}

		heap.Push(minHeap, newNode)
	}

	return heap.Pop(minHeap).(*node)
}

func generateHuffmanCodes(node *node, encodings map[byte]string, encoding string) {
	// traverse the tree and generate the generate codes
	// left = add a 0, right = add a 1 until we reach the right code
	// store this somewhere i guess
	if node == nil {
		return
	}

	// traverse each side, adding either 1 or zero to the code?
	if node.left == nil && node.right == nil {
		// if node.char != 0 {
		encodings[node.char] = encoding
		return
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

	fmt.Printf("Size of the file AFTER compression : %5d bytes\n\n", outputFileInfo.Size())
}
