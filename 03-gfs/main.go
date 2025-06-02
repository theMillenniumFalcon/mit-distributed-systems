package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	ChunkSize         = 64 * 1024 * 1024 // 64MB
	ReplicationFactor = 3
)

// Chunk represents a chunk of a file
type Chunk struct {
	Handle   string    `json:"handle"`
	Servers  []string  `json:"servers"`
	Version  int       `json:"version"`
	Size     int64     `json:"size"`
	Primary  string    `json:"primary"`
	LeaseEnd time.Time `json:"lease_end"`
}

// FileInfo represents metadata about a file
type FileInfo struct {
	Name   string   `json:"name"`
	Chunks []string `json:"chunks"`
	Size   int64    `json:"size"`
}

// Master server handles metadata and chunk allocation
type Master struct {
	files     map[string]*FileInfo
	chunks    map[string]*Chunk
	servers   []string
	nextChunk int
	mu        sync.RWMutex
	port      int
}

// Chunkserver stores actual chunk data
type Chunkserver struct {
	address string
	chunks  map[string][]byte
	master  string
	mu      sync.RWMutex
	dataDir string
}

// Client provides interface to GFS
type Client struct {
	master string
}

// NewMaster creates a new master server
func NewMaster(port int) *Master {
	return &Master{
		files:     make(map[string]*FileInfo),
		chunks:    make(map[string]*Chunk),
		servers:   make([]string, 0),
		nextChunk: 1,
		port:      port,
	}
}

// NewChunkserver creates a new chunkserver
func NewChunkserver(address, master string) *Chunkserver {
	dataDir := fmt.Sprintf("chunkserver_%s", strings.ReplaceAll(address, ":", "_"))
	os.MkdirAll(dataDir, 0755)

	return &Chunkserver{
		address: address,
		chunks:  make(map[string][]byte),
		master:  master,
		dataDir: dataDir,
	}
}

// NewClient creates a new GFS client
func NewClient(master string) *Client {
	return &Client{
		master: master,
	}
}

// Master server HTTP handlers

func (m *Master) handleRegisterServer(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	server := r.URL.Query().Get("server")
	if server == "" {
		http.Error(w, "Server address required", http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Add server if not already present
	for _, s := range m.servers {
		if s == server {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	m.servers = append(m.servers, server)
	log.Printf("Registered chunkserver: %s", server)
	w.WriteHeader(http.StatusOK)
}

func (m *Master) handleCreateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.files[filename]; exists {
		http.Error(w, "File already exists", http.StatusConflict)
		return
	}

	m.files[filename] = &FileInfo{
		Name:   filename,
		Chunks: make([]string, 0),
		Size:   0,
	}

	log.Printf("Created file: %s", filename)
	w.WriteHeader(http.StatusOK)
}

func (m *Master) handleGetChunks(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	file, exists := m.files[filename]
	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	chunks := make([]*Chunk, len(file.Chunks))
	for i, chunkHandle := range file.Chunks {
		chunks[i] = m.chunks[chunkHandle]
	}

	json.NewEncoder(w).Encode(chunks)
}

func (m *Master) handleAllocateChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	file, exists := m.files[filename]
	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Create new chunk
	chunkHandle := fmt.Sprintf("chunk_%d", m.nextChunk)
	m.nextChunk++

	// Select servers for replication
	servers := m.selectServers(ReplicationFactor)
	if len(servers) == 0 {
		http.Error(w, "No available servers", http.StatusServiceUnavailable)
		return
	}

	chunk := &Chunk{
		Handle:   chunkHandle,
		Servers:  servers,
		Version:  1,
		Size:     0,
		Primary:  servers[0], // First server is primary
		LeaseEnd: time.Now().Add(60 * time.Second),
	}

	m.chunks[chunkHandle] = chunk
	file.Chunks = append(file.Chunks, chunkHandle)

	log.Printf("Allocated chunk %s for file %s on servers %v", chunkHandle, filename, servers)
	json.NewEncoder(w).Encode(chunk)
}

func (m *Master) selectServers(count int) []string {
	if len(m.servers) < count {
		return m.servers
	}

	// Simple selection - in production would consider load, location, etc.
	selected := make([]string, count)
	for i := 0; i < count && i < len(m.servers); i++ {
		selected[i] = m.servers[i]
	}
	return selected
}

func (m *Master) start() {
	http.HandleFunc("/register", m.handleRegisterServer)
	http.HandleFunc("/create", m.handleCreateFile)
	http.HandleFunc("/chunks", m.handleGetChunks)
	http.HandleFunc("/allocate", m.handleAllocateChunk)

	log.Printf("Master server starting on port %d", m.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", m.port), nil))
}

// Chunkserver HTTP handlers

func (cs *Chunkserver) handleWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	chunkHandle := r.URL.Query().Get("chunk")
	if chunkHandle == "" {
		http.Error(w, "Chunk handle required", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read data", http.StatusBadRequest)
		return
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// Store chunk data in memory and on disk
	cs.chunks[chunkHandle] = data

	// Write to disk
	chunkPath := filepath.Join(cs.dataDir, chunkHandle)
	err = os.WriteFile(chunkPath, data, 0644)
	if err != nil {
		log.Printf("Failed to write chunk to disk: %v", err)
	}

	log.Printf("Stored chunk %s (%d bytes) on %s", chunkHandle, len(data), cs.address)
	w.WriteHeader(http.StatusOK)
}

func (cs *Chunkserver) handleRead(w http.ResponseWriter, r *http.Request) {
	chunkHandle := r.URL.Query().Get("chunk")
	if chunkHandle == "" {
		http.Error(w, "Chunk handle required", http.StatusBadRequest)
		return
	}

	cs.mu.RLock()
	defer cs.mu.RUnlock()

	// Try memory first, then disk
	data, exists := cs.chunks[chunkHandle]
	if !exists {
		// Try reading from disk
		chunkPath := filepath.Join(cs.dataDir, chunkHandle)
		diskData, err := os.ReadFile(chunkPath)
		if err != nil {
			http.Error(w, "Chunk not found", http.StatusNotFound)
			return
		}
		data = diskData
		// Cache in memory
		cs.chunks[chunkHandle] = data
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}

func (cs *Chunkserver) registerWithMaster() error {
	url := fmt.Sprintf("http://%s/register?server=%s", cs.master, cs.address)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("registration failed with status %d", resp.StatusCode)
	}

	log.Printf("Registered with master at %s", cs.master)
	return nil
}

func (cs *Chunkserver) start() {
	// Register with master
	for i := 0; i < 5; i++ {
		if err := cs.registerWithMaster(); err != nil {
			log.Printf("Failed to register with master (attempt %d): %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	http.HandleFunc("/write", cs.handleWrite)
	http.HandleFunc("/read", cs.handleRead)

	port := strings.Split(cs.address, ":")[1]
	log.Printf("Chunkserver starting on %s", cs.address)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Client operations

func (c *Client) createFile(filename string) error {
	url := fmt.Sprintf("http://%s/create?file=%s", c.master, filename)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) writeFile(filename, data string) error {
	// Create file if it doesn't exist
	c.createFile(filename)

	// Allocate chunk
	url := fmt.Sprintf("http://%s/allocate?file=%s", c.master, filename)
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chunk allocation failed with status %d", resp.StatusCode)
	}

	var chunk Chunk
	if err := json.NewDecoder(resp.Body).Decode(&chunk); err != nil {
		return err
	}

	// Write to all replicas
	for _, server := range chunk.Servers {
		writeURL := fmt.Sprintf("http://%s/write?chunk=%s", server, chunk.Handle)
		writeResp, err := http.Post(writeURL, "application/octet-stream", strings.NewReader(data))
		if err != nil {
			log.Printf("Failed to write to server %s: %v", server, err)
			continue
		}
		writeResp.Body.Close()
	}

	log.Printf("Wrote file %s (%d bytes)", filename, len(data))
	return nil
}

func (c *Client) readFile(filename string) (string, error) {
	// Get chunk locations
	url := fmt.Sprintf("http://%s/chunks?file=%s", c.master, filename)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get chunks with status %d", resp.StatusCode)
	}

	var chunks []*Chunk
	if err := json.NewDecoder(resp.Body).Decode(&chunks); err != nil {
		return "", err
	}

	if len(chunks) == 0 {
		return "", nil
	}

	// Read from first chunk, first server
	chunk := chunks[0]
	if len(chunk.Servers) == 0 {
		return "", fmt.Errorf("no servers available for chunk")
	}

	readURL := fmt.Sprintf("http://%s/read?chunk=%s", chunk.Servers[0], chunk.Handle)
	readResp, err := http.Get(readURL)
	if err != nil {
		return "", err
	}
	defer readResp.Body.Close()

	data, err := io.ReadAll(readResp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func main() {
	var (
		mode      = flag.String("mode", "master", "Mode: master, chunkserver, or client")
		port      = flag.Int("port", 8080, "Port to listen on")
		master    = flag.String("master", "localhost:8080", "Master server address")
		operation = flag.String("operation", "read", "Client operation: read or write")
		file      = flag.String("file", "", "File path")
		data      = flag.String("data", "", "Data to write")
	)
	flag.Parse()

	switch *mode {
	case "master":
		m := NewMaster(*port)
		m.start()

	case "chunkserver":
		address := fmt.Sprintf("localhost:%d", *port)
		cs := NewChunkserver(address, *master)
		cs.start()

	case "client":
		client := NewClient(*master)

		switch *operation {
		case "write":
			if *file == "" || *data == "" {
				log.Fatal("File and data required for write operation")
			}
			if err := client.writeFile(*file, *data); err != nil {
				log.Fatalf("Write failed: %v", err)
			}
			fmt.Printf("Successfully wrote to %s\n", *file)

		case "read":
			if *file == "" {
				log.Fatal("File required for read operation")
			}
			content, err := client.readFile(*file)
			if err != nil {
				log.Fatalf("Read failed: %v", err)
			}
			fmt.Printf("Content of %s: %s\n", *file, content)

		default:
			log.Fatal("Unknown operation. Use 'read' or 'write'")
		}

	default:
		log.Fatal("Unknown mode. Use 'master', 'chunkserver', or 'client'")
	}
}
