package main

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// Fetcher interface defines the behavior for fetching URLs
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	// Simulate network delay to make concurrency effects visible
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	if res, ok := f[url]; ok {
		fmt.Printf("Fetched: %s\n", url)
		return res.urls, nil
	}
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher - simulates a real website structure
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/cmd/": &fakeResult{
		"Commands",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
			"https://golang.org/cmd/go/",
			"https://golang.org/cmd/gofmt/",
		},
	},
	"https://golang.org/cmd/go/": &fakeResult{
		"Command go",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/cmd/gofmt/": &fakeResult{
		"Command gofmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
		},
	},
}

// ====================
// 1. SERIAL CRAWLER
// ====================
// Processes URLs sequentially, one at a time
func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return
	}
	fetched[url] = true

	urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}

	for _, u := range urls {
		Serial(u, fetcher, fetched)
	}
}

// ====================
// 2. CONCURRENT CRAWLER WITH MUTEX
// ====================
// Uses goroutines with mutex to protect shared state

// FetchState holds shared state for concurrent crawler with mutex
type FetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func makeState() *FetchState {
	return &FetchState{
		fetched: make(map[string]bool),
	}
}

// ConcurrentMutex - concurrent crawler using mutex for shared state
func ConcurrentMutex(url string, fetcher Fetcher, f *FetchState) {
	// Critical section: check and mark URL as fetched
	f.mu.Lock()
	already := f.fetched[url]
	f.fetched[url] = true
	f.mu.Unlock()

	if already {
		return
	}

	urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}

	// Launch goroutines for each found URL
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		go func(u string) {
			defer done.Done()
			ConcurrentMutex(u, fetcher, f)
		}(u) // Important: pass u as parameter to capture the value!
	}
	done.Wait() // Wait for all goroutines to complete
}

// ====================
// 3. CONCURRENT CRAWLER WITH CHANNELS
// ====================
// Uses channels for communication, no shared memory

// ConcurrentChannel - concurrent crawler using channels for communication
func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url}
	}()
	master(ch, fetcher)
}

// master coordinates the crawling using channels
func master(ch chan []string, fetcher Fetcher) {
	n := 1 // number of pending workers
	fetched := make(map[string]bool)

	for urls := range ch {
		for _, u := range urls {
			if !fetched[u] {
				fetched[u] = true
				n++
				go worker(u, ch, fetcher)
			}
		}
		n--
		if n == 0 {
			break
		}
	}
}

// worker fetches URLs and sends results back through channel
func worker(url string, ch chan []string, fetcher Fetcher) {
	urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		ch <- []string{}
	} else {
		ch <- urls
	}
}

// ====================
// UTILITY FUNCTIONS
// ====================

// Benchmark function to measure performance and visualize differences
func benchmark(name string, crawlFunc func()) {
	fmt.Printf("\n%s\n", strings.Repeat("=", 50))
	fmt.Printf("Running %s\n", name)
	fmt.Printf("%s\n", strings.Repeat("=", 50))
	start := time.Now()
	crawlFunc()
	duration := time.Since(start)
	fmt.Printf("\n%s completed in %v\n", name, duration)
}

// Demonstrate race conditions (for educational purposes)
func demonstrateRaceCondition() {
	fmt.Printf("\n%s\n", strings.Repeat("=", 50))
	fmt.Println("DEMONSTRATION: Why we need synchronization")
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	fmt.Println("🎓 RACE CONDITIONS EXPLAINED:")
	fmt.Println("Race conditions occur when multiple goroutines access shared data")
	fmt.Println("without proper synchronization. This can lead to:")
	fmt.Println("- Data corruption")
	fmt.Println("- Inconsistent state")
	fmt.Println("- Program crashes")
	fmt.Println()
	fmt.Println("🔍 TO SEE RACE CONDITIONS IN ACTION:")
	fmt.Println("Run this program with the race detector:")
	fmt.Println("  go run -race main.go")
	fmt.Println()
	fmt.Println("The race detector will show you exactly where race conditions occur!")
	fmt.Println()
	fmt.Println("⚠️  In a real scenario, concurrent map access without synchronization")
	fmt.Println("would cause 'fatal error: concurrent map writes' and crash the program.")
	fmt.Println("This is exactly why we use mutexes or channels for safe concurrent access.")
}

func main() {
	// Set random seed for consistent demonstration
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Concurrent Web Crawler")
	fmt.Println("=====================================================")
	fmt.Println("This demonstrates three approaches to web crawling:")
	fmt.Println("1. Serial (sequential)")
	fmt.Println("2. Concurrent with Mutex (shared memory)")
	fmt.Println("3. Concurrent with Channels (message passing)")
	fmt.Println()

	startURL := "https://golang.org/"

	// 1. Serial Crawler - baseline
	benchmark("Serial Crawler", func() {
		Serial(startURL, fetcher, make(map[string]bool))
	})

	// 2. Concurrent Crawler with Mutex - shared memory approach
	benchmark("Concurrent Crawler (Mutex - Shared Memory)", func() {
		state := makeState()
		ConcurrentMutex(startURL, fetcher, state)
	})

	// 3. Concurrent Crawler with Channels - message passing approach
	benchmark("Concurrent Crawler (Channels - Message Passing)", func() {
		ConcurrentChannel(startURL, fetcher)
	})

	// Educational demonstration of race conditions
	demonstrateRaceCondition()

	// Always show the key learnings, even if race condition demo panicked
	showKeyLearnings()
}

func showKeyLearnings() {
	fmt.Printf("\n%s\n", strings.Repeat("=", 50))
	fmt.Println("KEY LEARNINGS:")
	fmt.Printf("%s\n", strings.Repeat("=", 50))
	fmt.Println("1. Serial: Simple but slow - no parallelism")
	fmt.Println("2. Mutex: Fast with goroutines, but requires careful locking")
	fmt.Println("3. Channels: Fast with goroutines, uses Go's message passing")
	fmt.Println()
	fmt.Println("🔍 ADVANCED DEBUGGING:")
	fmt.Println("Run with: go run -race main.go")
	fmt.Println("This will detect race conditions in concurrent code!")
	fmt.Println()
	fmt.Println("🎓 DISTRIBUTED SYSTEMS CONNECTION:")
	fmt.Println("- Concurrency is fundamental to distributed systems")
	fmt.Println("- These patterns appear in MapReduce, Raft, and other algorithms")
	fmt.Println("- Understanding Go's concurrency model helps with distributed programming")
}
