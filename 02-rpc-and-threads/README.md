# Web Crawler Implementations

This directory contains two different implementations of concurrent web crawlers, each demonstrating different approaches to concurrency in Go and distributed systems.

## 📁 Directory Structure

```
02-rpc_and_threads/
├── shared-memory-crawler/     # Original implementation with 3 approaches
│   ├── main.go               # Serial, Mutex, and Basic Channel implementations
│   ├── go.mod                # Go module definition
│   └── README.md             # Detailed explanation of shared memory approaches
│
├── channel-based-crawler/     # MIT 6.824 style implementation
│   ├── main.go               # Pure channel-based master-worker pattern
│   ├── go.mod                # Go module definition
│   └── README.md             # MIT approach explanation and theory
│
└── README.md                 # This file - overview of both approaches
```

## 🎯 Two Different Educational Approaches

### 1. Shared Memory Crawler (`shared-memory-crawler/`)

**Purpose**: Demonstrates the evolution from serial to concurrent programming

**Approaches Covered:**
- **Serial**: Sequential baseline
- **Mutex**: Shared memory with synchronization
- **Channels**: Basic channel usage (still somewhat hybrid)

**Key Learning:**
- How to add concurrency to serial programs
- When and why to use mutexes
- Common pitfalls in concurrent programming
- Race condition detection and prevention

**Best For:**
- Understanding concurrency fundamentals
- Learning Go's concurrency primitives
- Seeing performance improvements step by step

### 2. Channel-Based Crawler (`channel-based-crawler/`)

**Purpose**: Demonstrates MIT 6.824's pure channel-based approach

**Approach:**
- **Pure Master-Worker Pattern**: Clean separation of concerns
- **No Shared Memory**: Eliminates race conditions by design
- **Message Passing**: Follows Go's communication philosophy

**Key Learning:**
- How to design concurrent systems without shared state
- Master-worker coordination patterns
- Channel-based synchronization
- Distributed systems thinking

**Best For:**
- Understanding distributed systems patterns
- Learning message-passing paradigms
- Preparing for MapReduce, Raft, and other distributed algorithms

## 🔄 Architectural Comparison

| Aspect | Shared Memory Approach | MIT Channel Approach |
|--------|----------------------|----------------------|
| **State Management** | Shared map with mutex | Master-only state |
| **Communication** | Shared memory + channels | Pure channels |
| **Architecture** | Recursive goroutines | Master-worker pattern |
| **Synchronization** | Explicit (mutex) | Implicit (channels) |
| **Complexity** | Medium | Low |
| **Distributed Ready** | Needs adaptation | Ready for distribution |
| **Race Conditions** | Possible if bugs exist | Eliminated by design |

## 🚀 When to Use Each Approach

### Use Shared Memory When:
- Working with existing shared-state systems
- Need fine-grained control over locking
- Porting from other languages (C++, Java, etc.)
- Performance-critical sections with minimal contention

### Use Channels When:
- Building new systems from scratch
- Want to eliminate race conditions by design
- Preparing for distributed deployment
- Following Go's idiomatic patterns

## 🎓 Learning Path Recommendation

1. **Start with Shared Memory Crawler**
   - Understand the progression from serial to concurrent
   - Learn about race conditions and synchronization
   - Practice with Go's concurrency primitives

2. **Then Study Channel-Based Crawler**
   - See how to eliminate shared state entirely
   - Understand master-worker patterns
   - Prepare for distributed systems concepts

3. **Compare and Contrast**
   - Run both implementations with `-race` flag
   - Measure performance differences
   - Understand when to use each approach

## 🏃‍♂️ Quick Start

### Run Shared Memory Crawler
```bash
cd shared-memory-crawler
go run main.go
go run -race main.go  # Check for race conditions
```

### Run Channel-Based Crawler
```bash
cd channel-based-crawler
go run main.go
go run -race main.go  # Should show no races by design
```

## 🔗 Connection to MIT 6.824

Both crawlers prepare you for MIT's distributed systems course:

- **Lab 1 (MapReduce)**: Uses master-worker patterns like channel-based crawler
- **Lab 2 (Raft)**: Uses message passing and state management concepts
- **Lab 3 (KV Store)**: Combines both approaches for client-server architecture

## 📚 Educational Value

### Shared Memory Crawler Teaches:
- Concurrency fundamentals
- Mutex usage and pitfalls
- Race condition detection
- Performance optimization

### Channel-Based Crawler Teaches:
- Message passing paradigms
- Distributed systems patterns
- Go's concurrency philosophy
- Clean architectural design

## 🎯 Next Steps

After studying both implementations:

1. **Experiment**: Modify both crawlers to add features
2. **Performance Test**: Compare with different workloads
3. **Extend**: Add real HTTP fetching, rate limiting, etc.
4. **Apply**: Use these patterns in your own projects

## 📖 Further Reading

- [MIT 6.824 Course Materials](https://pdos.csail.mit.edu/6.824/)
- [Go Concurrency Patterns](https://talks.golang.org/2012/concurrency.slide)
- [Effective Go - Concurrency](https://golang.org/doc/effective_go.html#concurrency)
- [The Go Memory Model](https://golang.org/ref/mem) 