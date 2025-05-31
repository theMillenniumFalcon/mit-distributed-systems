# MapReduce Implementation in Go

A beginner-friendly implementation of the MapReduce programming model for learning distributed systems concepts.

## 🎯 What is MapReduce?

MapReduce is a programming model designed for processing large datasets across multiple machines. It was originally developed by Google and consists of two main phases:

1. **Map Phase**: Processes input data and produces intermediate key-value pairs
2. **Reduce Phase**: Aggregates the intermediate data to produce final results

## 📁 Files Overview

- `mapreduce.go` - Core MapReduce framework implementation
- `wordcount.go` - Example map and reduce functions for counting words
- `main.go` - Main program that ties everything together
- `sample1.txt`, `sample2.txt` - Sample input files for testing
- `README.md` - This documentation

## 🔧 How It Works

### 1. Map Phase
```
Input: Text files
↓
Map Function: Extracts words and emits (word, "1") pairs
↓
Partitioning: Groups keys by hash for different reduce tasks
↓
Output: Intermediate files (mr-X-Y format)
```

### 2. Reduce Phase
```
Input: Intermediate files grouped by key
↓
Reduce Function: Sums up counts for each word
↓
Output: Final result files (mr-out-X format)
```

### 3. Key Components

**KeyValue Struct**: Represents a key-value pair
```go
type KeyValue struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}
```

**MapFunction**: User-defined function that processes input
```go
type MapFunction func(filename string, contents string) []KeyValue
```

**ReduceFunction**: User-defined function that aggregates values
```go
type ReduceFunction func(key string, values []string) string
```

## 🚀 Running the Example

### Prerequisites
- Go installed on your system
- Basic familiarity with command line

### Steps

1. **Navigate to the directory**:
   ```bash
   cd 01-Introduction
   ```

2. **Run the word count example**:
   ```bash
   go run *.go sample1.txt sample2.txt
   ```

3. **Check the results**:
   ```bash
   cat mr-out-0
   cat mr-out-1
   cat mr-out-2
   ```

### Expected Output
The program will show detailed progress of the MapReduce job:
```
🚀 Starting MapReduce Job
Input files: [sample1.txt sample2.txt]
Number of reduce tasks: 3

=== Starting Map Phase ===
Processing file 0: sample1.txt
  Map produced 35 key-value pairs
  Created intermediate file: mr-0-0 (12 pairs)
  Created intermediate file: mr-0-1 (11 pairs)
  Created intermediate file: mr-0-2 (12 pairs)
...
=== Map Phase Complete ===

=== Starting Reduce Phase ===
...
=== Reduce Phase Complete ===

✅ MapReduce Job Complete!
```

## 🧠 Understanding the Code

### Map Function (WordCountMap)
This function:
1. Takes a filename and its contents
2. Uses regex to extract words
3. Converts words to lowercase
4. Emits (word, "1") for each occurrence

### Reduce Function (WordCountReduce)
This function:
1. Takes a word and list of counts
2. Sums up all the "1"s
3. Returns the total count as a string

### Partitioning
Words are distributed across reduce tasks using a hash function:
```go
bucket := ihash(kv.Key) % mr.nReduce
```

This ensures that all instances of the same word go to the same reduce task.

## 🎓 Learning Concepts

### 1. Parallelism
- Map tasks can run in parallel (each processing different files)
- Reduce tasks can run in parallel (each handling different keys)

### 2. Fault Tolerance
- Intermediate files allow recovery if a reduce task fails
- Tasks can be re-executed on different machines

### 3. Scalability
- Adding more machines allows processing larger datasets
- Work is automatically distributed across available resources

### 4. Data Locality
- Map tasks typically run on machines storing the input data
- Reduces network traffic and improves performance

## 🔄 Extending the Implementation

### Creating Your Own MapReduce Job

1. **Define your map function**:
   ```go
   func MyMapFunction(filename string, contents string) []KeyValue {
       // Your logic here
       return keyValues
   }
   ```

2. **Define your reduce function**:
   ```go
   func MyReduceFunction(key string, values []string) string {
       // Your aggregation logic here
       return result
   }
   ```

3. **Run your job**:
   ```go
   mr := NewMapReduce(MyMapFunction, MyReduceFunction, nReduce, inputFiles)
   mr.Run()
   ```

### Example Ideas
- **Character Count**: Count characters instead of words
- **Line Count**: Count lines in files
- **Word Length**: Average length of words
- **Grep**: Find lines matching a pattern

## 🔍 Comparison with Real MapReduce

### Similarities
- Two-phase processing (Map → Reduce)
- Key-value pair abstraction
- Partitioning by key hash
- Intermediate file storage

### Simplifications
- Single machine (real MapReduce uses clusters)
- No fault tolerance (real systems handle machine failures)
- Sequential execution (real systems run tasks in parallel)
- No optimization (real systems optimize data movement)

## 📚 Further Reading

- [Original MapReduce Paper](https://research.google/pubs/pub62/) by Google
- [Hadoop MapReduce Documentation](https://hadoop.apache.org/docs/current/hadoop-mapreduce-client/hadoop-mapreduce-client-core/MapReduceTutorial.html)
- MIT 6.824 Distributed Systems Course

## 🎯 Next Steps

After understanding this implementation, you can:
1. Study the MIT 6.824 labs that build distributed MapReduce
2. Learn about Hadoop and Spark
3. Explore other distributed computing patterns
4. Build your own distributed systems

Happy learning! 🚀 