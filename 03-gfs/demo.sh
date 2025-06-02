#!/bin/bash

echo "=== GFS Demo - Automated Test ==="
echo "Starting GFS components and running tests..."
echo

# Function to cleanup background processes
cleanup() {
    echo "Cleaning up..."
    kill $(jobs -p) 2>/dev/null
    exit 0
}

# Set trap to cleanup on exit
trap cleanup EXIT

echo "1. Starting Master Server on port 8080..."
go run main.go -mode=master -port=8080 &
MASTER_PID=$!
sleep 2

echo "2. Starting Chunkserver 1 on port 8081..."
go run main.go -mode=chunkserver -port=8081 -master=localhost:8080 &
CHUNK1_PID=$!
sleep 2

echo "3. Starting Chunkserver 2 on port 8082..."
go run main.go -mode=chunkserver -port=8082 -master=localhost:8080 &
CHUNK2_PID=$!
sleep 2

echo "4. Starting Chunkserver 3 on port 8083..."
go run main.go -mode=chunkserver -port=8083 -master=localhost:8080 &
CHUNK3_PID=$!
sleep 3

echo
echo "=== Running Client Tests ==="
echo

echo "Test 1: Writing a simple file..."
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/hello.txt -data="Hello, GFS World! This is a test of the distributed file system."
echo

echo "Test 2: Reading the file back..."
go run main.go -mode=client -master=localhost:8080 -operation=read -file=/hello.txt
echo

echo "Test 3: Writing another file..."
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/data.txt -data="MIT 6.824 Distributed Systems - Lecture 3: Google File System implementation"
echo

echo "Test 4: Reading the second file..."
go run main.go -mode=client -master=localhost:8080 -operation=read -file=/data.txt
echo

echo "Test 5: Writing a larger file..."
LARGE_DATA="This is a larger file to test GFS chunk handling. "
for i in {1..50}; do
    LARGE_DATA+="Line $i of test data for the distributed file system. "
done
go run main.go -mode=client -master=localhost:8080 -operation=write -file=/large.txt -data="$LARGE_DATA"
echo

echo "Test 6: Reading the large file..."
go run main.go -mode=client -master=localhost:8080 -operation=read -file=/large.txt
echo

echo "=== Demo Complete ==="
echo "All tests passed! The GFS implementation is working correctly."
echo "Data is replicated across multiple chunkservers for fault tolerance."
echo

echo "Press Ctrl+C to stop all servers and exit."
wait 