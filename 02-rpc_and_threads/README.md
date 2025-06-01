# Concurrent Web Crawler

This project implements a simple concurrent web crawler in Go, demonstrating three different approaches to concurrency as taught in MIT 6.824 Distributed Systems course (Lecture 2: RPC and Threads).

## Overview

This crawler demonstrates key concepts in concurrent programming:
- **Shared State with Mutex**: Using locks to protect shared data structures
- **Message Passing with Channels**: Using Go channels for communication between goroutines
- **Race Condition Education**: Understanding why synchronization is necessary

## Three Implementations

### 1. Serial Crawler
- Processes URLs sequentially, one at a time
- Simple but slow - no parallelism
- Baseline for comparison

### 2. Concurrent Crawler with Mutex
- Uses goroutines for parallel processing
- Shared memory approach with `sync.Mutex` to protect the visited URLs map
- Demonstrates proper synchronization with locks
- Uses `sync.WaitGroup` to coordinate goroutines

### 3. Concurrent Crawler with Channels
- Uses goroutines for parallel processing  
- Message passing approach with channels
- No shared memory - follows Go's principle: "Don't communicate by sharing memory; share memory by communicating"
- Master-worker pattern with channels for coordination

## Key Learning Points

- **Goroutines**: Lightweight threads managed by Go runtime
- **Mutex**: Mutual exclusion locks to protect shared data
- **Channels**: Type-safe communication between goroutines
- **Race Conditions**: Problems that occur when multiple threads access shared data
- **Synchronization**: Techniques to coordinate concurrent execution

## Running the Code

### Basic Execution
```bash
cd 02-rpc_and_threads/web-crawler
go run main.go
```

### Race Detection
Go provides a built-in race detector to find race conditions:
```bash
go run -race main.go
```

This will detect any race conditions in the code and help you understand why synchronization is important.

### Sample Output
```
Concurrent Web Crawler
=====================================================
This demonstrates three approaches to web crawling:
1. Serial (sequential)
2. Concurrent with Mutex (shared memory)
3. Concurrent with Channels (message passing)

==================================================
Running Serial Crawler
==================================================
Fetched: https://golang.org/
Fetched: https://golang.org/pkg/
Fetched: https://golang.org/cmd/
Fetched: https://golang.org/cmd/go/
Fetched: https://golang.org/cmd/gofmt/
Fetched: https://golang.org/pkg/fmt/
Fetched: https://golang.org/pkg/os/

Serial Crawler completed in 363.786541ms

==================================================
Running Concurrent Crawler (Mutex - Shared Memory)
==================================================
Fetched: https://golang.org/
Fetched: https://golang.org/cmd/
Fetched: https://golang.org/cmd/go/
Fetched: https://golang.org/cmd/gofmt/
Fetched: https://golang.org/pkg/
Fetched: https://golang.org/pkg/os/
Fetched: https://golang.org/pkg/fmt/

Concurrent Crawler (Mutex - Shared Memory) completed in 252.135583ms

==================================================
Running Concurrent Crawler (Channels - Message Passing)
==================================================
Fetched: https://golang.org/
Fetched: https://golang.org/pkg/
Fetched: https://golang.org/pkg/fmt/
Fetched: https://golang.org/pkg/os/
Fetched: https://golang.org/cmd/
Fetched: https://golang.org/cmd/go/
Fetched: https://golang.org/cmd/gofmt/

Concurrent Crawler (Channels - Message Passing) completed in 225.906541ms

==================================================
DEMONSTRATION: Why we need synchronization
==================================================
🎓 RACE CONDITIONS EXPLAINED:
Race conditions occur when multiple goroutines access shared data
without proper synchronization. This can lead to:
- Data corruption
- Inconsistent state
- Program crashes

🔍 TO SEE RACE CONDITIONS IN ACTION:
Run this program with the race detector:
  go run -race main.go

The race detector will show you exactly where race conditions occur!

⚠️  In a real scenario, concurrent map access without synchronization
would cause 'fatal error: concurrent map writes' and crash the program.
This is exactly why we use mutexes or channels for safe concurrent access.

==================================================
KEY LEARNINGS:
==================================================
1. Serial: Simple but slow - no parallelism
2. Mutex: Fast with goroutines, but requires careful locking
3. Channels: Fast with goroutines, uses Go's message passing

🔍 ADVANCED DEBUGGING:
Run with: go run -race main.go
This will detect race conditions in concurrent code!

🎓 DISTRIBUTED SYSTEMS CONNECTION:
- Concurrency is fundamental to distributed systems
- These patterns appear in MapReduce, Raft, and other algorithms
- Understanding Go's concurrency model helps with distributed programming
```

## Code Structure

```
main.go
├── Fetcher Interface          # Defines how to fetch URLs
├── fakeFetcher Implementation  # Mock fetcher for testing
├── Serial Crawler             # Sequential implementation
├── Concurrent Mutex Crawler   # Shared memory with locks
├── Concurrent Channel Crawler # Message passing
└── Educational Demonstrations # Race condition explanations
```

## Educational Features

### 1. Simulated Network Delay
The fake fetcher includes random delays to make concurrency effects visible and realistic.

### 2. Race Condition Education
The code includes educational explanations about race conditions and why synchronization is needed, with safe demonstrations using Go's race detector.

### 3. Performance Comparison
Benchmark functions show the performance differences between approaches.

### 4. Detailed Comments
Extensive comments explain the "why" behind each implementation choice.

## Important Go Concepts Demonstrated

### Variable Capture in Goroutines
```go
for _, u := range urls {
    go func(u string) {  // Pass u as parameter
        // Use u safely here
    }(u)
}
```
This pattern prevents the common mistake of capturing loop variables incorrectly.

### Mutex Usage
```go
f.mu.Lock()
// Critical section - only one goroutine can execute this
already := f.fetched[url]
f.fetched[url] = true
f.mu.Unlock()
```

### Channel Communication
```go
ch := make(chan []string)
go worker(url, ch, fetcher)  // Send work to worker
urls := <-ch                 // Receive results
```

## Performance Observations

From the sample output above, you can see clear performance improvements:
- **Serial**: ~364ms - Sequential processing
- **Mutex**: ~252ms - 30% faster with proper synchronization
- **Channels**: ~226ms - 38% faster with message passing

The concurrent approaches demonstrate significant performance gains while maintaining correctness through proper synchronization.

## Extending the Crawler

This is a foundation that you can extend for learning:

1. **Add real HTTP fetching** instead of the mock fetcher
2. **Implement request limiting** to avoid overwhelming servers
3. **Add depth limiting** to control crawl scope
4. **Store results** in a database or file
5. **Add URL filtering** to crawl only specific domains
6. **Implement politeness delays** between requests

## Related MIT 6.824 Concepts

This crawler demonstrates concepts that appear throughout the distributed systems course:
- **Concurrency**: Foundation for distributed systems
- **State Management**: How to handle shared state across multiple processes
- **Communication Patterns**: Message passing vs shared memory
- **Fault Tolerance**: Handling failures in concurrent systems

## Next Steps

After understanding this crawler, you can proceed to:
- Lab 1: MapReduce (uses similar concurrency patterns)
- Lab 2: Raft (consensus algorithm with complex state management)
- Lab 3: Fault-tolerant Key/Value Service
- Lab 4: Sharded Key/Value Service

## Resources

- [MIT 6.824 Course](https://pdos.csail.mit.edu/6.824/)
- [Go Concurrency Tutorial](https://tour.golang.org/concurrency/1)
- [Effective Go - Concurrency](https://golang.org/doc/effective_go.html#concurrency) 