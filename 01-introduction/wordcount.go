package main

import (
	"regexp"
	"strconv"
	"strings"
)

// WordCountMap is the map function for word counting
// It takes a filename and file contents, and emits (word, "1") for each word
func WordCountMap(filename string, contents string) []KeyValue {
	// Split contents into words using regex
	// This regex finds sequences of letters
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	words := wordRegex.FindAllString(contents, -1)

	var keyValues []KeyValue
	for _, word := range words {
		// Convert to lowercase for case-insensitive counting
		word = strings.ToLower(word)
		// Emit (word, "1") - each occurrence counts as 1
		keyValues = append(keyValues, KeyValue{Key: word, Value: "1"})
	}

	return keyValues
}

// WordCountReduce is the reduce function for word counting
// It takes a word and a list of counts (all "1"s), and returns the total count
func WordCountReduce(key string, values []string) string {
	total := 0

	// Sum up all the "1"s for this word
	for _, value := range values {
		count, err := strconv.Atoi(value)
		if err != nil {
			// If we can't parse the number, assume it's 1
			count = 1
		}
		total += count
	}

	return strconv.Itoa(total)
}
