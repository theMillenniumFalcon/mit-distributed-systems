# Concurrent Web Crawler - Shared Memory Approaches

This implementation demonstrates three different approaches to web crawling, focusing on shared memory concurrency patterns.

## 🎯 Overview

This crawler showcases the evolution from serial to concurrent programming:

1. **Serial Crawler** - Baseline sequential implementation
2. **Concurrent with Mutex** - Shared memory with synchronization
3. **Concurrent with Channels** - Basic channel implementation

## 🔄 Three Implementations

### 1. Serial Crawler
```go
func Serial(url string, fetcher Fetcher, fetched map[string]bool)
```
- **Simplest approach**: One URL at a time
- **Sequential processing**: No concurrency
- **Baseline for comparison**: Shows what we're trying to improve

### 2. Concurrent with Mutex (Shared Memory)
```go
type FetchState struct {
    mu      sync.Mutex
    fetched map[string]bool
}
```
- **Shared state**: Multiple goroutines access same map
- **Synchronization**: Mutex protects the shared map
- **Recursive pattern**: Each goroutine spawns more goroutines
- **WaitGroup**: Ensures all goroutines complete

### 3. Concurrent with Channels (Local Implementation)
```go
func ConcurrentChannel(url string, fetcher Fetcher)
```
- **Message passing**: Communication via channels
- **Master-worker**: Coordinator and workers
- **No shared memory**: Each component owns its data

## 🔒 Synchronization Patterns

### Mutex Approach
```go
f.mu.Lock()
already := f.fetched[url]
f.fetched[url] = true
f.mu.Unlock()
```
**Pros:**
- Direct shared memory access
- Familiar to developers from other languages
- Fine-grained control over locking

**Cons:**
- Risk of race conditions
- Potential for deadlocks
- Requires careful lock management

### WaitGroup Coordination
```go
var done sync.WaitGroup
for _, u := range urls {
    done.Add(1)
    go func(u string) {
        defer done.Done()
        ConcurrentMutex(u, fetcher, f)
    }(u)
}
done.Wait()
```

## ⚠️ Common Concurrency Pitfalls

### 1. Race Conditions
```bash
# Detect race conditions
go run -race main.go
```

### 2. Variable Capture in Closures
```go
// WRONG - captures loop variable
for _, u := range urls {
    go func() {
        ConcurrentMutex(u, fetcher, f) // u changes!
    }()
}

// CORRECT - pass as parameter
for _, u := range urls {
    go func(u string) {
        ConcurrentMutex(u, fetcher, f)
    }(u)
}
```

## 📊 Performance Comparison

Running the benchmark shows:
- **Serial**: Slowest, but simple and safe
- **Concurrent Mutex**: Fast, but requires careful synchronization
- **Concurrent Channels**: Fast and safe by design

## 🔍 Key Learning Points

### When to Use Mutex
- **Shared state**: When multiple goroutines need access to same data
- **Fine-grained control**: When you need precise control over locking
- **Legacy integration**: When working with existing code patterns

### When to Use Channels
- **Communication**: When goroutines need to coordinate
- **Safety**: When you want to eliminate race conditions by design
- **Distributed thinking**: When preparing for network communication

## 🚀 Comparison with MIT Channel Approach

This implementation's channel approach differs from the MIT version:

| Aspect | This Implementation | MIT Approach |
|--------|-------------------|--------------|
| **Architecture** | Still recursive | Pure master-worker |
| **State Management** | Mixed | Master-only state |
| **Complexity** | More complex | Simpler, cleaner |
| **Distributed Ready** | Less ready | More ready |

For the pure MIT approach, see the `../channel-based-crawler/` directory.

## 🏃‍♂️ Running the Examples

```bash
# Run all three implementations
go run main.go

# Check for race conditions
go run -race main.go

# Run specific tests
go test -v
```

## 🎓 Educational Value

This implementation teaches:
1. **Evolution of concurrent programming**
2. **Trade-offs** between different approaches
3. **Common pitfalls** in concurrent programming
4. **Go's concurrency primitives** (goroutines, mutexes, channels)
5. **Race condition detection** and prevention

## 📚 Next Steps

1. Study the **MIT channel-based approach** in `../channel-based-crawler/`
2. Experiment with **different fetcher implementations**
3. Try **adding rate limiting** or **request caching**
4. Explore **distributed crawler** patterns

## 🔗 Related Concepts

- **MapReduce**: Similar master-worker patterns
- **Producer-Consumer**: Channel communication patterns
- **Actor Model**: Message-passing paradigms
- **CSP (Communicating Sequential Processes)**: Go's theoretical foundation