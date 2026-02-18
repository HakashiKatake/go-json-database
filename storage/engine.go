package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/HakashiKatake/Go-Json-Database/core"
)

// FileStorageEngine implements the StorageEngine interface with thread-safe file operations
type FileStorageEngine struct {
	dataDir string
	mu      sync.RWMutex
	locks   map[string]*os.File // File locks per collection
	locksMu sync.Mutex          // Protects the locks map
}

// CollectionFile represents the structure of a collection file
type CollectionFile struct {
	Metadata  CollectionMetadata         `json:"metadata"`
	Documents map[string]core.Document   `json:"documents"`
}

// CollectionMetadata contains metadata about a collection
type CollectionMetadata struct {
	Collection    string    `json:"collection"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	DocumentCount int       `json:"document_count"`
}

// NewFileStorageEngine creates a new file-based storage engine
func NewFileStorageEngine(dataDir string) (*FileStorageEngine, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	return &FileStorageEngine{
		dataDir: dataDir,
		locks:   make(map[string]*os.File),
	}, nil
}

// getCollectionPath returns the file path for a collection
func (e *FileStorageEngine) getCollectionPath(collection string) string {
	return filepath.Join(e.dataDir, collection+".json")
}

// acquireFileLock acquires an exclusive file lock for a collection
func (e *FileStorageEngine) acquireFileLock(collection string) (*os.File, error) {
	e.locksMu.Lock()
	defer e.locksMu.Unlock()

	// Check if we already have a lock file open
	if lockFile, exists := e.locks[collection]; exists {
		// Try to acquire the lock
		if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
			return nil, fmt.Errorf("failed to acquire file lock: %w", err)
		}
		return lockFile, nil
	}

	// Open lock file
	lockPath := filepath.Join(e.dataDir, collection+".lock")
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	// Acquire exclusive lock
	if err := syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX); err != nil {
		lockFile.Close()
		return nil, fmt.Errorf("failed to acquire file lock: %w", err)
	}

	e.locks[collection] = lockFile
	return lockFile, nil
}

// releaseFileLock releases the file lock for a collection
func (e *FileStorageEngine) releaseFileLock(lockFile *os.File) error {
	if lockFile == nil {
		return nil
	}
	return syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
}

// readCollectionFile reads the entire collection file
func (e *FileStorageEngine) readCollectionFile(collection string) (*CollectionFile, error) {
	path := e.getCollectionPath(collection)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return empty collection
		return &CollectionFile{
			Metadata: CollectionMetadata{
				Collection:    collection,
				Version:       1,
				CreatedAt:     time.Now(),
				DocumentCount: 0,
			},
			Documents: make(map[string]core.Document),
		}, nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection file: %w", err)
	}

	// Parse JSON
	var collFile CollectionFile
	if err := json.Unmarshal(data, &collFile); err != nil {
		return nil, fmt.Errorf("failed to parse collection file: %w", err)
	}

	return &collFile, nil
}

// writeCollectionFileAtomic writes the collection file atomically using temp file + rename
func (e *FileStorageEngine) writeCollectionFileAtomic(collection string, collFile *CollectionFile) error {
	path := e.getCollectionPath(collection)
	
	// Update metadata
	collFile.Metadata.DocumentCount = len(collFile.Documents)

	// Marshal to JSON
	data, err := json.MarshalIndent(collFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collection file: %w", err)
	}

	// Write to temporary file
	tempPath := path + ".tmp"
	tempFile, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	// Write data
	if _, err := tempFile.Write(data); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Fsync to ensure data is on disk
	if err := tempFile.Sync(); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	// Close temp file
	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, path); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// WriteDocument atomically writes a document to storage
func (e *FileStorageEngine) WriteDocument(collection string, docID core.DocumentID, doc core.Document) error {
	// Acquire write lock
	e.mu.Lock()
	defer e.mu.Unlock()

	// Acquire file lock
	lockFile, err := e.acquireFileLock(collection)
	if err != nil {
		return err
	}
	defer e.releaseFileLock(lockFile)

	// Read current collection
	collFile, err := e.readCollectionFile(collection)
	if err != nil {
		return err
	}

	// Add/update document
	collFile.Documents[string(docID)] = doc

	// Write atomically
	return e.writeCollectionFileAtomic(collection, collFile)
}

// ReadDocument retrieves a document by ID
func (e *FileStorageEngine) ReadDocument(collection string, docID core.DocumentID) (core.Document, error) {
	// Acquire read lock
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Read collection file
	collFile, err := e.readCollectionFile(collection)
	if err != nil {
		return nil, err
	}

	// Find document
	doc, exists := collFile.Documents[string(docID)]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", docID)
	}

	return doc, nil
}

// DeleteDocument removes a document from storage
func (e *FileStorageEngine) DeleteDocument(collection string, docID core.DocumentID) error {
	// Acquire write lock
	e.mu.Lock()
	defer e.mu.Unlock()

	// Acquire file lock
	lockFile, err := e.acquireFileLock(collection)
	if err != nil {
		return err
	}
	defer e.releaseFileLock(lockFile)

	// Read current collection
	collFile, err := e.readCollectionFile(collection)
	if err != nil {
		return err
	}

	// Delete document
	delete(collFile.Documents, string(docID))

	// Write atomically
	return e.writeCollectionFileAtomic(collection, collFile)
}

// ScanCollection iterates over all documents in a collection
func (e *FileStorageEngine) ScanCollection(collection string, fn func(core.DocumentID, core.Document) bool) error {
	// Acquire read lock
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Read collection file
	collFile, err := e.readCollectionFile(collection)
	if err != nil {
		return err
	}

	// Iterate over documents
	for docID, doc := range collFile.Documents {
		if !fn(core.DocumentID(docID), doc) {
			break
		}
	}

	return nil
}

// CreateCollection initializes a new collection
func (e *FileStorageEngine) CreateCollection(name string) error {
	// Acquire write lock
	e.mu.Lock()
	defer e.mu.Unlock()

	// Acquire file lock
	lockFile, err := e.acquireFileLock(name)
	if err != nil {
		return err
	}
	defer e.releaseFileLock(lockFile)

	// Check if collection already exists
	path := e.getCollectionPath(name)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("collection already exists: %s", name)
	}

	// Create empty collection
	collFile := &CollectionFile{
		Metadata: CollectionMetadata{
			Collection:    name,
			Version:       1,
			CreatedAt:     time.Now(),
			DocumentCount: 0,
		},
		Documents: make(map[string]core.Document),
	}

	// Write to disk
	return e.writeCollectionFileAtomic(name, collFile)
}

// ListCollections returns all collection names
func (e *FileStorageEngine) ListCollections() ([]string, error) {
	// Acquire read lock
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Read directory
	entries, err := os.ReadDir(e.dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read data directory: %w", err)
	}

	// Filter for .json files
	var collections []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			name := entry.Name()[:len(entry.Name())-5] // Remove .json extension
			collections = append(collections, name)
		}
	}

	return collections, nil
}

// Close flushes pending writes and releases locks
func (e *FileStorageEngine) Close() error {
	e.locksMu.Lock()
	defer e.locksMu.Unlock()

	// Release all file locks
	for collection, lockFile := range e.locks {
		if err := e.releaseFileLock(lockFile); err != nil {
			return fmt.Errorf("failed to release lock for collection %s: %w", collection, err)
		}
		if err := lockFile.Close(); err != nil {
			return fmt.Errorf("failed to close lock file for collection %s: %w", collection, err)
		}
	}

	// Clear locks map
	e.locks = make(map[string]*os.File)

	return nil
}
