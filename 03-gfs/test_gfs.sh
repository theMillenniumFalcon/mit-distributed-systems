#!/bin/bash

echo "=== Basic GFS Test ==="
echo "This script demonstrates the basic GFS functionality"
echo

echo "1. Starting Master Server..."
echo "   Run: go run main.go -mode=master -port=8080"
echo "   (Start this in a separate terminal)"
echo

echo "2. Starting Chunkservers..."
echo "   Run in separate terminals:"
echo "   go run main.go -mode=chunkserver -port=8081 -master=localhost:8080"
echo "   go run main.go -mode=chunkserver -port=8082 -master=localhost:8080"
echo

echo "3. Client Operations:"
echo "   Write a file:"
echo "   go run main.go -mode=client -master=localhost:8080 -operation=write -file=/test.txt -data=\"Hello, GFS World!\""
echo
echo "   Read the file:"
echo "   go run main.go -mode=client -master=localhost:8080 -operation=read -file=/test.txt"
echo

echo "4. Testing with larger data:"
echo "   go run main.go -mode=client -master=localhost:8080 -operation=write -file=/large.txt -data=\"$(cat /dev/urandom | base64 | head -c 1000)\""
echo "   go run main.go -mode=client -master=localhost:8080 -operation=read -file=/large.txt"
echo

echo "=== Manual Testing Steps ==="
echo "1. Open 4 terminal windows"
echo "2. In terminal 1: Start master server"
echo "3. In terminals 2-3: Start chunkservers"
echo "4. In terminal 4: Run client commands"
echo

echo "Example session:"
echo "Terminal 1: go run main.go -mode=master -port=8080"
echo "Terminal 2: go run main.go -mode=chunkserver -port=8081 -master=localhost:8080"
echo "Terminal 3: go run main.go -mode=chunkserver -port=8082 -master=localhost:8080"
echo "Terminal 4: go run main.go -mode=client -master=localhost:8080 -operation=write -file=/hello.txt -data=\"Hello GFS\""
echo "Terminal 4: go run main.go -mode=client -master=localhost:8080 -operation=read -file=/hello.txt" 