package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Fetcher interface defines the behavior for fetching URLs
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	// Simulate network delay to make concurrency effects visible
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	if res, ok := f[url]; ok {
		fmt.Printf("Fetched: %s\n", url)
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
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
// 1. SERIAL CRAWLER (for comparison)
// ====================
func Serial(url string, fetcher Fetcher, fetched map[string]bool) {
	if fetched[url] {
		return
	}
	fetched[url] = true
	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	for _, u := range urls {
		Serial(u, fetcher, fetched)
	}
}

// ====================
// 2. CONCURRENT CRAWLER WITH MUTEX (for comparison)
// ====================
type fetchState struct {
	mu      sync.Mutex
	fetched map[string]bool
}

func makeState() *fetchState {
	f := &fetchState{}
	f.fetched = make(map[string]bool)
	return f
}

func ConcurrentMutex(url string, fetcher Fetcher, f *fetchState) {
	f.mu.Lock()
	already := f.fetched[url]
	f.fetched[url] = true
	f.mu.Unlock()

	if already {
		return
	}

	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		return
	}
	var done sync.WaitGroup
	for _, u := range urls {
		done.Add(1)
		go func(u string) {
			defer done.Done()
			ConcurrentMutex(u, fetcher, f)
		}(u)
	}
	done.Wait()
}

// ====================
// 3. MIT-STYLE CHANNEL-BASED CRAWLER
// ====================

// ConcurrentChannel implements MIT's approach to channel-based web crawling
// This is the main entry point for the channel-based crawler
func ConcurrentChannel(url string, fetcher Fetcher) {
	ch := make(chan []string)
	go func() {
		ch <- []string{url} // Send initial URL to channel
	}()
	master(ch, fetcher)
}

// master coordinates the crawling using channels
// This function runs in the main goroutine and:
// 1. Maintains the fetched map (no shared memory!)
// 2. Coordinates all workers through the channel
// 3. Decides when crawling is complete
func master(ch chan []string, fetcher Fetcher) {
	n := 1                           // number of pending workers
	fetched := make(map[string]bool) // Only master touches this map!

	for urls := range ch {
		for _, u := range urls {
			if !fetched[u] {
				fetched[u] = true
				n += 1
				go worker(u, ch, fetcher) // Start worker for this URL
			}
		}
		n -= 1 // Current worker finished
		if n == 0 {
			break // All workers done
		}
	}
}

// worker fetches a URL and sends results back through channel
// Each worker:
// 1. Fetches one URL
// 2. Sends found URLs back to master via channel
// 3. Terminates (no recursion!)
func worker(url string, ch chan []string, fetcher Fetcher) {
	_, urls, err := fetcher.Fetch(url)
	if err != nil {
		ch <- []string{} // Send empty slice on error
	} else {
		ch <- urls // Send found URLs to master
	}
}

// ====================
// UTILITY FUNCTIONS
// ====================

func benchmark(name string, crawlFunc func()) {
	fmt.Printf("\n%s\n", "=================================")
	fmt.Printf("Running %s\n", name)
	fmt.Printf("%s\n", "=================================")
	start := time.Now()
	crawlFunc()
	duration := time.Since(start)
	fmt.Printf("\n%s completed in %v\n", name, duration)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("MIT 6.824 Style Channel-Based Web Crawler")
	fmt.Println("==========================================")
	fmt.Println("Demonstrates MIT's approach to concurrent web crawling using channels")
	fmt.Println("Key features:")
	fmt.Println("- No shared memory between goroutines")
	fmt.Println("- Master-worker pattern using channels")
	fmt.Println("- Clean separation of concerns")
	fmt.Println("- No mutex needed!")
	fmt.Println()

	startURL := "https://golang.org/"

	// Compare all three approaches
	benchmark("1. Serial Crawler (Baseline)", func() {
		Serial(startURL, fetcher, make(map[string]bool))
	})

	benchmark("2. Concurrent with Mutex (Shared Memory)", func() {
		state := makeState()
		ConcurrentMutex(startURL, fetcher, state)
	})

	benchmark("3. MIT Channel-Based Crawler (Message Passing)", func() {
		ConcurrentChannel(startURL, fetcher)
	})

	// Explain the MIT approach
	explainMITApproach()
}

func explainMITApproach() {
	fmt.Printf("\n%s\n", "=================================")
	fmt.Println("MIT 6.824 CHANNEL-BASED APPROACH EXPLAINED:")
	fmt.Printf("%s\n", "=================================")
	fmt.Println()
	fmt.Println("🎓 KEY CONCEPTS FROM MIT LECTURE 2:")
	fmt.Println()
	fmt.Println("1. MASTER-WORKER PATTERN:")
	fmt.Println("   - Master: Coordinates everything, maintains fetched map")
	fmt.Println("   - Workers: Fetch URLs, send results back via channel")
	fmt.Println()
	fmt.Println("2. NO SHARED MEMORY:")
	fmt.Println("   - Only master touches the fetched map")
	fmt.Println("   - No mutexes needed!")
	fmt.Println("   - Communication through channels only")
	fmt.Println()
	fmt.Println("3. CHANNEL COMMUNICATION:")
	fmt.Println("   - Workers send URL slices to master")
	fmt.Println("   - Master decides which URLs to crawl")
	fmt.Println("   - Channel provides both communication AND synchronization")
	fmt.Println()
	fmt.Println("4. TERMINATION DETECTION:")
	fmt.Println("   - Master counts active workers (n variable)")
	fmt.Println("   - When n == 0, all work is done")
	fmt.Println("   - Simple and elegant!")
	fmt.Println()
	fmt.Println("🔍 WHY THIS APPROACH WORKS:")
	fmt.Println("- Eliminates race conditions by design")
	fmt.Println("- Follows Go's principle: 'Don't communicate by sharing memory;")
	fmt.Println("  share memory by communicating'")
	fmt.Println("- Scales well to distributed systems")
	fmt.Println()
	fmt.Println("🚀 DISTRIBUTED SYSTEMS CONNECTION:")
	fmt.Println("- This pattern appears in MapReduce, Raft, and other algorithms")
	fmt.Println("- Message passing is fundamental to distributed computing")
	fmt.Println("- Prepares you for thinking about network communication")
}
