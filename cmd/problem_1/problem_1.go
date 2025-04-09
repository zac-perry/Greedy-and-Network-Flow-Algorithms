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
// char -> represents the character read in from the file for this node.
// freq -> frequency the character appears in the file.
// left -> left node connection in the min heap.
// right -> right node connection in the min heap.
// NOTE: This is also used to create new 'parent' nodes, that point to two 'char' nodes, with the frequency set to the total of both left.freq + right.freq.
type node struct {
	char  byte
	freq  int
	left  *node
	right *node
}

// NOTE: Go provides a Heap implementation. In order to use it, you must implement the interface provided.
// NOTE: Once implemented, Go will call the additional functions needed when calling Push and Pop, using your implementations.

// MinHeap []*node type represents the MinHeap using a slice of nodes.
type MinHeap []*node

// Len() function will return the length of the min heap (required for interface).
func (minHeap MinHeap) Len() int {
	return len(minHeap)
}

// Less() function will compare the frequences of two nodes in the min heap (required for interface).
// In this case, comparing minHeap[i] < minHeap[j] will result in us building a min heap.
func (minHeap MinHeap) Less(i, j int) bool {
	return minHeap[i].freq < minHeap[j].freq
}

// Swap() will swap two elements in the min heap (required for interface).
func (minHeap MinHeap) Swap(i, j int) {
	minHeap[i], minHeap[j] = minHeap[j], minHeap[i]
}

// Push() will add a new node to the min heap (required for interface).
// NOTE: Go will automatically handle all re-orderings, etc if needed.
func (minHeap *MinHeap) Push(x any) {
	*minHeap = append(*minHeap, x.(*node))
}

// Pop() will remove a node from the heap (required for interface).
// NOTE: Go will automatically handle re-ordering elements, etc if needed.
func (minHeap *MinHeap) Pop() any {
	placeholder := *minHeap
	length := len(placeholder)
	value := placeholder[length-1]
	*minHeap = placeholder[0 : length-1]

	return value
}

// buildHuffmanTree() function builds out the Huffman Tree utilizing the min heap.
func buildHuffmanTree(charFreqMap map[byte]int) *node {
	minHeap := &MinHeap{}
	heap.Init(minHeap)

	// Create a node for each char in the charFreqMap.
	// Add the nodes to the Min heap which will order nodes by frequency (smallest frequency = highest priority).
	for k, v := range charFreqMap {
		newNode := &node{
			char: k,
			freq: v,
		}
		heap.Push(minHeap, newNode)
	}

	// Build the Huffman Tree.
	// First, pop off the first two nodes (both have the current smallest frequencies).
	// Create a parent node, set freq to the total frequency, and push onto the Min Heap.
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

	// Return the head node of the Min Heap.
	return heap.Pop(minHeap).(*node)
}

// generateHuffmanCodes() will create the Huffman Code for each node.
// Traverse the tree recursivley until you reach a leaf. Once found, add the encoding to the map.
// Each encoding is build based on the sides of the tree traversed to find it (left will add a "0", right will add a "1").
func generateHuffmanCodes(node *node, encodings map[byte]string, encoding string) {
	if node == nil {
		return
	}

	// If we have found a leaf node, set the encoding and return.
	if node.left == nil && node.right == nil {
		encodings[node.char] = encoding
		return
	}

	// Traverse either side of the tree to find all leaf nodes, building the encoding.
	// Left = encoding + "0".
	// Right = encoding + "1".
	if node.left != nil {
		generateHuffmanCodes(node.left, encodings, encoding+"0")
	}
	if node.right != nil {
		generateHuffmanCodes(node.right, encodings, encoding+"1")
	}
}

// convertStringtoBytes() function will take a concatenated string and convert it to bytes.
// This string contains the Huffman codes for each letter in the file, in the order they appear in the file.
func convertStringtoBytes(fullEncodingString string) []byte {
	numBytes := (len(fullEncodingString) + 7) / 8
	bytes := make([]byte, numBytes)

	// Process each byte and calculate start and end indices for this byte.
	for i := range numBytes {
		startIdx := i * 8
		endIdx := startIdx + 8

		if endIdx > len(fullEncodingString) {
			endIdx = len(fullEncodingString)
		}

		// Extract out the current bit string sequence.
		currBitString := fullEncodingString[startIdx:endIdx]

		// Account for padding of the final element.
		// NOTE: If decoding, must account for the fact that the last element could have padding.
		if len(currBitString) < 8 {
			padding := make([]byte, 8-len(currBitString))
			for i := range padding {
				padding[i] = '0'
			}
			currBitString = currBitString + string(padding)
		}

		// Finish by converting the current bit string sequence to an actual byte, saving it to the bytes slice.
		b, _ := strconv.ParseUint(currBitString, 2, 8)
		bytes[i] = byte(b)
	}

	return bytes
}

// writeEncodeFile() function will write the encoded bytes to the output binary file.
func writeEncodeFile(fileContents []byte, encodings map[byte]string) []byte {
	outputFile, err := os.Create("output/output_1.bin")
	if err != nil {
		log.Fatal("Error opening the output file")
	}

	defer outputFile.Close()

	// Traverse the original file contents, looking up the Huffman code for each letter.
	// Concat the codes to a string.
	encodingString := ""
	for _, c := range fileContents {
		encodingString += encodings[c]
	}

	// Convert the entire string to bytes then write to the binary file.
	encodedBytes := convertStringtoBytes(encodingString)
	outputFile.Write(encodedBytes)

	return encodedBytes
}

// printFreqTable() function will just print out the chars and their corresponding frequency of occurence within the file.
func printFreqTable(charFreqMap map[byte]int) {

	// NOTE: This is needed to sort by key in Go because the map implementation doesn't auto sort (lame).
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

// readFile() function just takes the corresponding file name and opens it.
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

	printFreqTable(freq)
	fmt.Printf("Size of the file BEFORE compression: %5d bytes\n", inputFileInfo.Size())
	fmt.Printf("Size of the file AFTER compression : %5d bytes\n\n", outputFileInfo.Size())
}
