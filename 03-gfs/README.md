# Basic Google File System (GFS) Implementation

A complete basic implementation of Google File System (GFS).

## 🎯 What We Built

This implementation demonstrates the core concepts of GFS with:

### Core Components

1. **Master Server** (`main.go:40-219`)
   - Manages file namespace and metadata
   - Handles chunk allocation and location tracking
   - Coordinates chunkserver registration
   - Manages chunk leases for primary/secondary replicas

2. **Chunkserver** (`main.go:221-341`) 
   - Stores actual file chunks on disk and in memory
   - Handles read/write operations for chunks
   - Registers with master server
   - Supports replication across multiple servers

3. **Client** (`main.go:343-436`)
   - Provides file system interface (read/write operations)
   - Communicates with master for metadata
   - Performs direct I/O with chunkservers
   - Handles chunk location caching

## 🏗️ Architecture

```
Client
  |
  v
Master Server (metadata, chunk locations)
  |
  v
Chunkservers (actual data storage)
```

### Key GFS Concepts Implemented

1. **Chunk-based Storage**
   - 64MB chunk size (configurable)
   - Each file split into chunks
   - Chunks identified by unique handles

2. **Replication**
   - 3-way replication (configurable)
   - Primary-secondary replica model
   - Write-through to all replicas

3. **Master-Chunkserver Architecture**
   - Single master for metadata
   - Multiple chunkservers for data
   - Client communicates with both

4. **Consistency Model**
   - Primary replica coordinates writes
   - Lease-based consistency
   - All replicas updated synchronously

## 📁 File Structure

```
03-gfs/
├── main.go              # Complete GFS implementation
├── go.mod               # Go module configuration
├── README.md            # This comprehensive documentation
├── test_gfs.sh          # Manual testing instructions
├── demo.sh              # Automated testing script
└── chunkserver_*/       # Chunkserver data directories (created at runtime)
```

## 🚀 How to Use

### Quick Start - Automated Demo
```bash
cd 03-gfs
./demo.sh  # Runs complete end-to-end test automatically
```

### Manual Testing - Step by Step

#### 1. Start Master Server
```bash
go run main.go -mode=master -port=8080
```

#### 2. Start Chunkservers (in separate terminals)
```bash
# Terminal 2
go run main.go -mode=chunkserver -port=8081 -master=localhost:8080

# Terminal 3  
go run main.go -mode=chunkserver -port=8082 -master=localhost:8080

# Terminal 4
go run main.go -mode=chunkserver -port=8083 -master=localhost:8080
```

#### 3. Client Operations
```bash
# Write a file
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/hello.txt -data="Hello GFS!"

# Read a file
go run main.go -mode=client -master=localhost:8080 -operation=read -file=/hello.txt

# Write larger files
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/test.txt -data="MIT 6.824 Distributed Systems - GFS Implementation"
```

## 🔍 GFS Operations Flow

### Write Operation
1. Client contacts master for chunk allocation
2. Master returns chunk handle + chunkserver locations
3. Client writes to all replica chunkservers
4. Primary chunkserver coordinates the write
5. Data replicated to secondary replicas

### Read Operation  
1. Client contacts master for chunk locations
2. Master returns chunk handle + chunkserver addresses
3. Client reads directly from nearest chunkserver
4. No master involvement in data transfer

### Append Operation
1. Client sends append request to master
2. Master checks if record fits in current chunk
3. If yes, returns primary chunkserver; if no, creates new chunk
4. Client appends data through primary chunkserver

## 🔧 Technical Details

### Data Structures
- **`FileInfo`**: File metadata (name, chunks, size)
- **`Chunk`**: Chunk metadata (handle, servers, version, lease)
- **HTTP-based communication** between all components

### Consistency Guarantees
- **File namespace**: Strongly consistent (single master)
- **Data regions**: Consistent across replicas
- **Concurrent writes**: Undefined but consistent

### Fault Tolerance Features
- Multiple replicas per chunk (default: 3)
- Chunkserver failure detection and automatic re-registration
- Persistent storage on disk with in-memory caching
- Master coordinates all metadata operations

### Implementation Highlights
- **Chunk Size**: 64MB (configurable via constants)
- **Replication Factor**: 3 replicas (configurable)
- **Communication**: HTTP REST APIs
- **Storage**: Both disk and memory for performance
- **Leases**: 60-second chunk leases for consistency

## 🎓 Learning Outcomes

This implementation demonstrates key distributed systems concepts:

- **Distributed System Design**: Clear separation of metadata and data
- **Fault Tolerance**: Multiple replicas and failure handling mechanisms
- **Scalability**: Master handles metadata, chunkservers handle data
- **Consistency**: Primary-backup replication model
- **Performance**: Large chunks reduce metadata overhead

## ⚠️ Limitations

This is a simplified educational implementation and lacks:
- Production-grade fault tolerance (master single point of failure)
- Sophisticated load balancing and server selection
- Advanced consistency guarantees for concurrent operations
- Security features and authentication
- Performance optimizations and caching strategies
- Garbage collection and space reclamation
- Network partition handling

## 🧪 Testing Examples

### Basic File Operations
```bash
# Write and read a simple file
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/simple.txt -data="Hello World"
go run main.go -mode=client -master=localhost:8080 -operation=read -file=/simple.txt

# Test with multiple files
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/file1.txt -data="First file"
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/file2.txt -data="Second file"
go run main.go -mode=client -master=localhost:8080 -operation=read -file=/file1.txt
go run main.go -mode=client -master=localhost:8080 -operation=read -file=/file2.txt
```

### Testing Replication
1. Start master and 3 chunkservers
2. Write a file (gets replicated to all servers)
3. Stop one chunkserver
4. Read the file (should still work from remaining replicas)

## 🏆 Perfect for MIT 6.824

This implementation covers the essential GFS concepts from Lecture 3:
- ✅ Master-chunkserver architecture
- ✅ Large chunk sizes (64MB)  
- ✅ Replication for fault tolerance
- ✅ Weak consistency model
- ✅ Single master design
- ✅ Client-direct data transfer
- ✅ Lease-based consistency

Ready for educational exploration and understanding of distributed file systems! 