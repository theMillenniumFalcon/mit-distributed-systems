package main

import (
	"fmt"
	"os"
)

func main() {
	// Check if we have input files
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run *.go <input_file1> [input_file2] ...")
		fmt.Println("Example: go run *.go sample1.txt sample2.txt")
		os.Exit(1)
	}
	
	// Get input files from command line arguments
	inputFiles := os.Args[1:]
	
	// Verify all input files exist
	for _, filename := range inputFiles {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Printf("Error: File %s does not exist\n", filename)
			os.Exit(1)
		}
	}
	
	fmt.Println("MapReduce Word Count Example")
	fmt.Println("============================")
	fmt.Printf("Input files: %v\n\n", inputFiles)
	
	// Create and run the MapReduce job
	// We use 3 reduce tasks to demonstrate partitioning
	nReduce := 3
	mr := NewMapReduce(WordCountMap, WordCountReduce, nReduce, inputFiles)
	
	// Run the job
	mr.Run()
	
	// Show the results
	fmt.Println("\n📊 Results:")
	for i := 0; i < nReduce; i++ {
		outputFile := fmt.Sprintf("mr-out-%d", i)
		if _, err := os.Stat(outputFile); err == nil {
			fmt.Printf("Output file: %s\n", outputFile)
		}
	}
	
	fmt.Println("\nTo see the word counts, check the mr-out-* files!")
	fmt.Println("Example: cat mr-out-0")
} 