# MIT 6.824 Channel-Based Web Crawler

This implementation follows the exact approach taught in **MIT 6.824 Distributed Systems Lecture 2** by Professor Frans Kaashoek and Robert Morris.

## 🎓 Educational Context

This crawler demonstrates Go's channel-based concurrency model as taught in MIT's distributed systems course. It's designed to teach fundamental concepts that apply to distributed computing.

## 🏗️ Architecture: Master-Worker Pattern

### Master Goroutine
- **Single Point of Control**: Maintains the `fetched` map
- **Coordination**: Decides which URLs to crawl
- **Termination Detection**: Counts active workers to know when done
- **No Shared Memory**: Only the master touches the fetched map

### Worker Goroutines
- **Simple Task**: Fetch one URL and send results back
- **No State**: Workers are stateless
- **Communication**: Send results via channel
- **No Recursion**: Unlike mutex version, workers don't call themselves

## 📡 Channel Communication

```go
ch := make(chan []string)  // Unbuffered channel for URL slices
```

### Data Flow
1. **Initial**: Master sends starting URL to channel
2. **Workers**: Each worker fetches a URL and sends found URLs back
3. **Coordination**: Master receives URL slices and spawns new workers
4. **Termination**: When all workers complete (n == 0), crawling stops

## 🔄 Algorithm Walkthrough

### Initialization
```go
n := 1  // Start with 1 pending worker (for initial URL)
fetched := make(map[string]bool)  // Only master accesses this
```

### Main Loop
```go
for urls := range ch {  // Master reads from channel
    for _, u := range urls {
        if fetched[u] == false {
            fetched[u] = true  // Mark as seen
            n += 1            // Increment worker count
            go worker(u, ch, fetcher)  // Spawn worker
        }
    }
    n -= 1  // Current batch processed
    if n == 0 {
        break  // All workers finished
    }
}
```

## 🔍 Key Differences from Mutex Approach

| Aspect | Mutex Approach | Channel Approach |
|--------|----------------|------------------|
| **Shared Memory** | Yes (with locks) | No |
| **Synchronization** | Explicit mutex | Channel provides sync |
| **Architecture** | Recursive calls | Master-worker pattern |
| **Race Conditions** | Possible if bugs | Eliminated by design |
| **Scalability** | Limited | Better for distributed |

## 🚀 Why This Matters for Distributed Systems

### 1. Message Passing Paradigm
- **Preparation**: This pattern appears in MapReduce, Raft, etc.
- **Network Ready**: Easy to adapt for network communication
- **Fault Tolerance**: Workers can fail independently

### 2. Go's Philosophy
> "Don't communicate by sharing memory; share memory by communicating."

### 3. Distributed Patterns
- **Actor Model**: Each goroutine is like an actor
- **Event-Driven**: Master reacts to worker messages
- **Asynchronous**: Non-blocking communication

## 🏃‍♂️ Running the Crawler

```bash
# Run the crawler
go run main.go

# Compare with race detection
go run -race main.go
```

## 📊 Performance Characteristics

- **Concurrency**: High (many workers run simultaneously)
- **Memory**: Lower contention (no shared mutex)
- **Scalability**: Excellent (can distribute across machines)
- **Debugging**: Easier (no deadlock possibilities)

## 🎯 Learning Objectives

After studying this implementation, you should understand:

1. **Channel-based coordination** in Go
2. **Master-worker patterns** for distributed computing
3. **Message passing** as an alternative to shared memory
4. **Termination detection** in concurrent systems
5. **How Go principles** apply to distributed systems

## 🔗 Connection to MIT 6.824 Labs

This crawler teaches concepts directly applicable to:
- **Lab 1 (MapReduce)**: Master-worker coordination
- **Lab 2 (Raft)**: Leader-follower communication
- **Lab 3 (KV Store)**: Client-server protocols

## 📚 Further Reading

- [MIT 6.824 Lecture 2 Notes](https://pdos.csail.mit.edu/6.824/notes/l-rpc.txt)
- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide)
- [Effective Go - Channels](https://golang.org/doc/effective_go.html#channels) 