package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"sort"
)

// KeyValue represents a key-value pair used throughout MapReduce
type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// MapFunction is the interface that user-defined map functions must implement
// It takes a filename and its contents, and returns a slice of KeyValue pairs
type MapFunction func(filename string, contents string) []KeyValue

// ReduceFunction is the interface that user-defined reduce functions must implement
// It takes a key and a slice of values for that key, and returns a single value
type ReduceFunction func(key string, values []string) string

// MapReduce represents our MapReduce coordinator
type MapReduce struct {
	mapFunc    MapFunction
	reduceFunc ReduceFunction
	nReduce    int // number of reduce tasks
	inputFiles []string
}

// NewMapReduce creates a new MapReduce instance
func NewMapReduce(mapFunc MapFunction, reduceFunc ReduceFunction, nReduce int, inputFiles []string) *MapReduce {
	return &MapReduce{
		mapFunc:    mapFunc,
		reduceFunc: reduceFunc,
		nReduce:    nReduce,
		inputFiles: inputFiles,
	}
}

// hash function to determine which reduce task should handle a key
func ihash(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32())
}

// RunMapPhase executes the map phase
// For each input file, it runs the map function and partitions the output
func (mr *MapReduce) RunMapPhase() {
	fmt.Println("=== Starting Map Phase ===")

	for i, filename := range mr.inputFiles {
		fmt.Printf("Processing file %d: %s\n", i, filename)

		// Read the input file
		content, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", filename, err)
		}

		// Run the map function
		keyValues := mr.mapFunc(filename, string(content))
		fmt.Printf("  Map produced %d key-value pairs\n", len(keyValues))

		// Partition the output into intermediate files for each reduce task
		buckets := make([][]KeyValue, mr.nReduce)
		for _, kv := range keyValues {
			// Use hash function to determine which reduce task gets this key
			bucket := ihash(kv.Key) % mr.nReduce
			buckets[bucket] = append(buckets[bucket], kv)
		}

		// Write each bucket to an intermediate file
		for r := 0; r < mr.nReduce; r++ {
			filename := fmt.Sprintf("mr-%d-%d", i, r)
			file, err := os.Create(filename)
			if err != nil {
				log.Fatalf("Error creating intermediate file %s: %v", filename, err)
			}

			enc := json.NewEncoder(file)
			for _, kv := range buckets[r] {
				if err := enc.Encode(&kv); err != nil {
					log.Fatalf("Error encoding to intermediate file: %v", err)
				}
			}
			file.Close()
			fmt.Printf("  Created intermediate file: %s (%d pairs)\n", filename, len(buckets[r]))
		}
	}
	fmt.Println("=== Map Phase Complete ===")
}

// RunReducePhase executes the reduce phase
// For each reduce task, it collects all intermediate files and runs the reduce function
func (mr *MapReduce) RunReducePhase() {
	fmt.Println("=== Starting Reduce Phase ===")

	for r := 0; r < mr.nReduce; r++ {
		fmt.Printf("Running reduce task %d\n", r)

		// Collect all intermediate files for this reduce task
		var keyValues []KeyValue
		for m := 0; m < len(mr.inputFiles); m++ {
			filename := fmt.Sprintf("mr-%d-%d", m, r)

			// Check if file exists (some might be empty)
			if _, err := os.Stat(filename); os.IsNotExist(err) {
				continue
			}

			file, err := os.Open(filename)
			if err != nil {
				log.Fatalf("Error opening intermediate file %s: %v", filename, err)
			}

			dec := json.NewDecoder(file)
			for {
				var kv KeyValue
				if err := dec.Decode(&kv); err != nil {
					break // End of file
				}
				keyValues = append(keyValues, kv)
			}
			file.Close()
		}

		fmt.Printf("  Collected %d key-value pairs\n", len(keyValues))

		// Group by key
		keyGroups := make(map[string][]string)
		for _, kv := range keyValues {
			keyGroups[kv.Key] = append(keyGroups[kv.Key], kv.Value)
		}

		// Sort keys for consistent output
		var keys []string
		for key := range keyGroups {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		// Run reduce function for each key and write output
		outputFilename := fmt.Sprintf("mr-out-%d", r)
		file, err := os.Create(outputFilename)
		if err != nil {
			log.Fatalf("Error creating output file %s: %v", outputFilename, err)
		}

		for _, key := range keys {
			values := keyGroups[key]
			result := mr.reduceFunc(key, values)
			fmt.Fprintf(file, "%v %v\n", key, result)
		}
		file.Close()

		fmt.Printf("  Created output file: %s (%d unique keys)\n", outputFilename, len(keys))
	}
	fmt.Println("=== Reduce Phase Complete ===")
}

// Cleanup removes intermediate files
func (mr *MapReduce) Cleanup() {
	fmt.Println("=== Cleaning up intermediate files ===")
	for m := 0; m < len(mr.inputFiles); m++ {
		for r := 0; r < mr.nReduce; r++ {
			filename := fmt.Sprintf("mr-%d-%d", m, r)
			if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
				fmt.Printf("Warning: could not remove %s: %v\n", filename, err)
			}
		}
	}
}

// Run executes the complete MapReduce job
func (mr *MapReduce) Run() {
	fmt.Println("🚀 Starting MapReduce Job")
	fmt.Printf("Input files: %v\n", mr.inputFiles)
	fmt.Printf("Number of reduce tasks: %d\n\n", mr.nReduce)

	mr.RunMapPhase()
	mr.RunReducePhase()
	mr.Cleanup()

	fmt.Println("✅ MapReduce Job Complete!")
}
