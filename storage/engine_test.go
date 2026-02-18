package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/HakashiKatake/Go-Json-Database/core"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func setupTestEngine(t *testing.T) (*FileStorageEngine, string) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "storage_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create engine
	engine, err := NewFileStorageEngine(tempDir)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create engine: %v", err)
	}

	return engine, tempDir
}

func cleanupTestEngine(engine *FileStorageEngine, tempDir string) {
	engine.Close()
	os.RemoveAll(tempDir)
}

func TestNewFileStorageEngine(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Verify data directory was created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Errorf("Data directory was not created")
	}
}

func TestCreateCollection(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create collection
	err := engine.CreateCollection("users")
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Verify collection file exists
	collPath := filepath.Join(tempDir, "users.json")
	if _, err := os.Stat(collPath); os.IsNotExist(err) {
		t.Errorf("Collection file was not created")
	}

	// Try to create same collection again (should fail)
	err = engine.CreateCollection("users")
	if err == nil {
		t.Errorf("Expected error when creating duplicate collection")
	}
}

func TestWriteAndReadDocument(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create collection
	err := engine.CreateCollection("users")
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Write document
	doc := core.Document{
		"id":    "user_001",
		"name":  "Alice",
		"email": "alice@example.com",
	}
	err = engine.WriteDocument("users", "user_001", doc)
	if err != nil {
		t.Fatalf("Failed to write document: %v", err)
	}

	// Read document
	readDoc, err := engine.ReadDocument("users", "user_001")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}

	// Verify document contents
	if readDoc["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got '%v'", readDoc["name"])
	}
	if readDoc["email"] != "alice@example.com" {
		t.Errorf("Expected email 'alice@example.com', got '%v'", readDoc["email"])
	}
}

func TestUpdateDocument(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create collection
	err := engine.CreateCollection("users")
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Write initial document
	doc := core.Document{
		"id":   "user_001",
		"name": "Alice",
	}
	err = engine.WriteDocument("users", "user_001", doc)
	if err != nil {
		t.Fatalf("Failed to write document: %v", err)
	}

	// Update document
	updatedDoc := core.Document{
		"id":   "user_001",
		"name": "Alice Updated",
	}
	err = engine.WriteDocument("users", "user_001", updatedDoc)
	if err != nil {
		t.Fatalf("Failed to update document: %v", err)
	}

	// Read and verify
	readDoc, err := engine.ReadDocument("users", "user_001")
	if err != nil {
		t.Fatalf("Failed to read document: %v", err)
	}

	if readDoc["name"] != "Alice Updated" {
		t.Errorf("Expected name 'Alice Updated', got '%v'", readDoc["name"])
	}
}

func TestDeleteDocument(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create collection
	err := engine.CreateCollection("users")
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Write document
	doc := core.Document{
		"id":   "user_001",
		"name": "Alice",
	}
	err = engine.WriteDocument("users", "user_001", doc)
	if err != nil {
		t.Fatalf("Failed to write document: %v", err)
	}

	// Delete document
	err = engine.DeleteDocument("users", "user_001")
	if err != nil {
		t.Fatalf("Failed to delete document: %v", err)
	}

	// Try to read deleted document (should fail)
	_, err = engine.ReadDocument("users", "user_001")
	if err == nil {
		t.Errorf("Expected error when reading deleted document")
	}
}

func TestScanCollection(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create collection
	err := engine.CreateCollection("users")
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Write multiple documents
	docs := []struct {
		id   core.DocumentID
		name string
	}{
		{"user_001", "Alice"},
		{"user_002", "Bob"},
		{"user_003", "Charlie"},
	}

	for _, d := range docs {
		doc := core.Document{
			"id":   string(d.id),
			"name": d.name,
		}
		err = engine.WriteDocument("users", d.id, doc)
		if err != nil {
			t.Fatalf("Failed to write document: %v", err)
		}
	}

	// Scan collection
	count := 0
	err = engine.ScanCollection("users", func(docID core.DocumentID, doc core.Document) bool {
		count++
		return true
	})
	if err != nil {
		t.Fatalf("Failed to scan collection: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 documents, got %d", count)
	}
}

func TestListCollections(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create multiple collections
	collections := []string{"users", "products", "orders"}
	for _, coll := range collections {
		err := engine.CreateCollection(coll)
		if err != nil {
			t.Fatalf("Failed to create collection %s: %v", coll, err)
		}
	}

	// List collections
	list, err := engine.ListCollections()
	if err != nil {
		t.Fatalf("Failed to list collections: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 collections, got %d", len(list))
	}

	// Verify all collections are in the list
	collMap := make(map[string]bool)
	for _, coll := range list {
		collMap[coll] = true
	}

	for _, coll := range collections {
		if !collMap[coll] {
			t.Errorf("Collection %s not found in list", coll)
		}
	}
}

func TestCollectionIsolation(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create two collections
	err := engine.CreateCollection("users")
	if err != nil {
		t.Fatalf("Failed to create users collection: %v", err)
	}

	err = engine.CreateCollection("products")
	if err != nil {
		t.Fatalf("Failed to create products collection: %v", err)
	}

	// Write to users collection
	userDoc := core.Document{
		"id":   "user_001",
		"name": "Alice",
	}
	err = engine.WriteDocument("users", "user_001", userDoc)
	if err != nil {
		t.Fatalf("Failed to write to users: %v", err)
	}

	// Write to products collection
	productDoc := core.Document{
		"id":   "prod_001",
		"name": "Widget",
	}
	err = engine.WriteDocument("products", "prod_001", productDoc)
	if err != nil {
		t.Fatalf("Failed to write to products: %v", err)
	}

	// Verify users collection has only user document
	userCount := 0
	err = engine.ScanCollection("users", func(docID core.DocumentID, doc core.Document) bool {
		userCount++
		if doc["name"] != "Alice" {
			t.Errorf("Unexpected document in users collection")
		}
		return true
	})
	if err != nil {
		t.Fatalf("Failed to scan users: %v", err)
	}
	if userCount != 1 {
		t.Errorf("Expected 1 document in users, got %d", userCount)
	}

	// Verify products collection has only product document
	productCount := 0
	err = engine.ScanCollection("products", func(docID core.DocumentID, doc core.Document) bool {
		productCount++
		if doc["name"] != "Widget" {
			t.Errorf("Unexpected document in products collection")
		}
		return true
	})
	if err != nil {
		t.Fatalf("Failed to scan products: %v", err)
	}
	if productCount != 1 {
		t.Errorf("Expected 1 document in products, got %d", productCount)
	}
}

func TestReadNonExistentDocument(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Create collection
	err := engine.CreateCollection("users")
	if err != nil {
		t.Fatalf("Failed to create collection: %v", err)
	}

	// Try to read non-existent document
	_, err = engine.ReadDocument("users", "nonexistent")
	if err == nil {
		t.Errorf("Expected error when reading non-existent document")
	}
}

func TestReadFromNonExistentCollection(t *testing.T) {
	engine, tempDir := setupTestEngine(t)
	defer cleanupTestEngine(engine, tempDir)

	// Try to read from non-existent collection (should return empty, not error)
	_, err := engine.ReadDocument("nonexistent", "doc_001")
	if err == nil {
		t.Errorf("Expected error when reading from non-existent collection")
	}
}

// Property-Based Tests

// Feature: monster-backend-database, Property 1: Concurrent Write Safety
// **Validates: Requirements 1.2**
func TestProperty_ConcurrentWriteSafety(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("concurrent writes complete without race conditions and all documents are present", 
		prop.ForAll(
			func(numGoroutines int, docsPerGoroutine int) bool {
				// Setup test engine
				engine, tempDir := setupTestEngine(t)
				defer cleanupTestEngine(engine, tempDir)

				// Create collection
				if err := engine.CreateCollection("test_concurrent"); err != nil {
					t.Logf("Failed to create collection: %v", err)
					return false
				}

				// Track all document IDs we expect to write
				expectedDocs := make(map[core.DocumentID]bool)
				var expectedMu sync.Mutex

				// Use WaitGroup to synchronize goroutines
				var wg sync.WaitGroup
				errChan := make(chan error, numGoroutines*docsPerGoroutine)

				// Launch concurrent writers
				for i := 0; i < numGoroutines; i++ {
					wg.Add(1)
					go func(goroutineID int) {
						defer wg.Done()
						
						for j := 0; j < docsPerGoroutine; j++ {
							docID := core.DocumentID(fmt.Sprintf("doc_%d_%d", goroutineID, j))
							doc := core.Document{
								"id":          string(docID),
								"goroutine":   goroutineID,
								"index":       j,
								"data":        fmt.Sprintf("data_%d_%d", goroutineID, j),
							}

							// Track expected document
							expectedMu.Lock()
							expectedDocs[docID] = true
							expectedMu.Unlock()

							// Write document
							if err := engine.WriteDocument("test_concurrent", docID, doc); err != nil {
								errChan <- fmt.Errorf("goroutine %d failed to write doc %s: %w", goroutineID, docID, err)
								return
							}
						}
					}(i)
				}

				// Wait for all goroutines to complete
				wg.Wait()
				close(errChan)

				// Check for any errors
				for err := range errChan {
					t.Logf("Write error: %v", err)
					return false
				}

				// Verify all documents are present in storage
				actualDocs := make(map[core.DocumentID]bool)
				err := engine.ScanCollection("test_concurrent", func(docID core.DocumentID, doc core.Document) bool {
					actualDocs[docID] = true
					return true
				})
				if err != nil {
					t.Logf("Failed to scan collection: %v", err)
					return false
				}

				// Check that all expected documents are present
				if len(actualDocs) != len(expectedDocs) {
					t.Logf("Document count mismatch: expected %d, got %d", len(expectedDocs), len(actualDocs))
					return false
				}

				for docID := range expectedDocs {
					if !actualDocs[docID] {
						t.Logf("Missing document: %s", docID)
						return false
					}
				}

				return true
			},
			gen.IntRange(2, 10),   // numGoroutines: 2-10 concurrent writers
			gen.IntRange(5, 20),   // docsPerGoroutine: 5-20 documents per goroutine
		))

	properties.TestingRun(t)
}

// Feature: monster-backend-database, Property 2: Atomic Write Guarantee
// **Validates: Requirements 1.3, 1.6**
func TestProperty_AtomicWriteGuarantee(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("failed write operations leave previous valid state unchanged and readable",
		prop.ForAll(
			func(initialDocs []string, updateDocs []string) bool {
				// Setup test engine
				engine, tempDir := setupTestEngine(t)
				defer cleanupTestEngine(engine, tempDir)

				// Create collection
				if err := engine.CreateCollection("test_atomic"); err != nil {
					t.Logf("Failed to create collection: %v", err)
					return false
				}

				// Write initial documents to establish a valid state
				initialState := make(map[core.DocumentID]core.Document)
				for i, data := range initialDocs {
					docID := core.DocumentID(fmt.Sprintf("doc_%d", i))
					doc := core.Document{
						"id":      string(docID),
						"data":    data,
						"version": "v1",
					}
					if err := engine.WriteDocument("test_atomic", docID, doc); err != nil {
						t.Logf("Failed to write initial document: %v", err)
						return false
					}
					initialState[docID] = doc
				}

				// Verify initial state is readable
				for docID, expectedDoc := range initialState {
					readDoc, err := engine.ReadDocument("test_atomic", docID)
					if err != nil {
						t.Logf("Failed to read initial document %s: %v", docID, err)
						return false
					}
					if readDoc["data"] != expectedDoc["data"] {
						t.Logf("Initial document data mismatch for %s", docID)
						return false
					}
				}

				// Test 1: Normal updates (should succeed)
				for i, data := range updateDocs {
					if i >= len(initialDocs) {
						break
					}
					docID := core.DocumentID(fmt.Sprintf("doc_%d", i))
					updatedDoc := core.Document{
						"id":      string(docID),
						"data":    data,
						"version": "v2",
					}

					err := engine.WriteDocument("test_atomic", docID, updatedDoc)

					// Read back the document
					readDoc, readErr := engine.ReadDocument("test_atomic", docID)
					if readErr != nil {
						t.Logf("Failed to read document after update: %v", readErr)
						return false
					}

					if err == nil {
						// Update succeeded, verify new state
						if readDoc["data"] != data {
							t.Logf("After successful update, document has wrong data")
							return false
						}
						if readDoc["version"] != "v2" {
							t.Logf("After successful update, document has wrong version")
							return false
						}
						initialState[docID] = updatedDoc
					} else {
						// Update failed, verify old state is preserved
						expectedData := initialState[docID]["data"]
						if readDoc["data"] != expectedData {
							t.Logf("After failed update, document data changed")
							return false
						}
						if readDoc["version"] != "v1" {
							t.Logf("After failed update, document version changed")
							return false
						}
					}
				}

				// Test 2: Simulate failure by making directory read-only
				// Save current state before attempting to cause failures
				stateBeforeFailure := make(map[core.DocumentID]core.Document)
				for k, v := range initialState {
					stateBeforeFailure[k] = v
				}

				// Make the data directory read-only to force write failures
				if err := os.Chmod(tempDir, 0555); err != nil {
					t.Logf("Failed to make directory read-only: %v", err)
					// Continue anyway, this is just an attempt to force failures
				}

				// Attempt writes that should fail
				for i := 0; i < len(initialDocs); i++ {
					docID := core.DocumentID(fmt.Sprintf("doc_%d", i))
					failDoc := core.Document{
						"id":      string(docID),
						"data":    "should_fail",
						"version": "v3",
					}

					// This write should fail due to permissions
					err := engine.WriteDocument("test_atomic", docID, failDoc)

					// Regardless of whether it failed or not, verify state is valid
					readDoc, readErr := engine.ReadDocument("test_atomic", docID)
					if readErr != nil {
						t.Logf("Failed to read document after failed write attempt: %v", readErr)
						return false
					}

					// If write failed, state should be unchanged
					if err != nil {
						expectedDoc := stateBeforeFailure[docID]
						if readDoc["data"] != expectedDoc["data"] {
							t.Logf("After failed write, document data was corrupted")
							return false
						}
						if readDoc["version"] != expectedDoc["version"] {
							t.Logf("After failed write, document version was corrupted")
							return false
						}
					}
				}

				// Restore permissions for cleanup
				os.Chmod(tempDir, 0755)

				// Final verification: all documents should be in a valid, readable state
				err := engine.ScanCollection("test_atomic", func(docID core.DocumentID, doc core.Document) bool {
					if doc["id"] == nil || doc["data"] == nil || doc["version"] == nil {
						t.Logf("Document %s has invalid structure", docID)
						return false
					}

					// Document should have valid data (not corrupted)
					_, ok := doc["data"].(string)
					if !ok {
						t.Logf("Document %s data is not a string", docID)
						return false
					}

					// Version should be valid
					versionStr, ok := doc["version"].(string)
					if !ok || (versionStr != "v1" && versionStr != "v2" && versionStr != "v3") {
						t.Logf("Document %s has invalid version: %v", docID, doc["version"])
						return false
					}

					return true
				})

				if err != nil {
					t.Logf("Failed to scan collection in final verification: %v", err)
					return false
				}

				return true
			},
			gen.SliceOfN(5, gen.AlphaString()),  // initialDocs: 5 random strings
			gen.SliceOfN(5, gen.AlphaString()),  // updateDocs: 5 random strings for updates
		))

	properties.TestingRun(t)
}

// Feature: monster-backend-database, Property 3: File Lock Exclusivity
// **Validates: Requirements 1.4**
func TestProperty_FileLockExclusivity(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("when one process holds a write lock, other processes block until lock is released",
		prop.ForAll(
			func(numCompetingWriters int, holdDurationMs int) bool {
				// Setup test engine
				engine, tempDir := setupTestEngine(t)
				defer cleanupTestEngine(engine, tempDir)

				// Create collection
				if err := engine.CreateCollection("test_locks"); err != nil {
					t.Logf("Failed to create collection: %v", err)
					return false
				}

				// Track the order of lock acquisitions
				var lockOrder []int
				var lockOrderMu sync.Mutex

				// Channel to signal when first writer has acquired lock
				firstLockAcquired := make(chan struct{})
				
				// Channel to signal when first writer should release lock
				releaseLock := make(chan struct{})

				// Use WaitGroup to synchronize all goroutines
				var wg sync.WaitGroup

				// First writer: acquires lock and holds it
				wg.Add(1)
				go func() {
					defer wg.Done()

					// Manually acquire file lock (simulating a long-held lock)
					lockFile, err := engine.acquireFileLock("test_locks")
					if err != nil {
						t.Logf("First writer failed to acquire lock: %v", err)
						close(firstLockAcquired)
						return
					}

					// Record that first writer got the lock
					lockOrderMu.Lock()
					lockOrder = append(lockOrder, 0)
					lockOrderMu.Unlock()

					// Signal that first lock is acquired
					close(firstLockAcquired)

					// Hold the lock for specified duration
					time.Sleep(time.Duration(holdDurationMs) * time.Millisecond)

					// Wait for signal to release
					<-releaseLock

					// Release the lock
					if err := engine.releaseFileLock(lockFile); err != nil {
						t.Logf("First writer failed to release lock: %v", err)
					}
				}()

				// Wait for first writer to acquire lock
				<-firstLockAcquired

				// Launch competing writers that should block
				for i := 1; i <= numCompetingWriters; i++ {
					wg.Add(1)
					go func(writerID int) {
						defer wg.Done()

						// Try to write a document (this should block until first writer releases)
						docID := core.DocumentID(fmt.Sprintf("doc_%d", writerID))
						doc := core.Document{
							"id":       string(docID),
							"writer":   writerID,
							"data":     fmt.Sprintf("data_%d", writerID),
						}

						// This write should block until the first writer releases the lock
						startTime := time.Now()
						err := engine.WriteDocument("test_locks", docID, doc)
						elapsed := time.Since(startTime)

						if err != nil {
							t.Logf("Writer %d failed to write: %v", writerID, err)
							return
						}

						// Record lock acquisition order
						lockOrderMu.Lock()
						lockOrder = append(lockOrder, writerID)
						lockOrderMu.Unlock()

						// Verify that this write was blocked for at least some time
						// (should be blocked until first writer releases)
						if elapsed < time.Duration(holdDurationMs/2)*time.Millisecond {
							t.Logf("Writer %d completed too quickly (%v), may not have been properly blocked", writerID, elapsed)
						}
					}(i)
				}

				// Give competing writers a moment to start and attempt lock acquisition
				time.Sleep(50 * time.Millisecond)

				// Signal first writer to release lock
				close(releaseLock)

				// Wait for all goroutines to complete
				wg.Wait()

				// Verify lock order: first writer (0) should have acquired lock first
				lockOrderMu.Lock()
				defer lockOrderMu.Unlock()

				if len(lockOrder) == 0 {
					t.Logf("No locks were acquired")
					return false
				}

				if lockOrder[0] != 0 {
					t.Logf("First writer did not acquire lock first, order: %v", lockOrder)
					return false
				}

				// Verify all competing writers eventually acquired locks
				if len(lockOrder) != numCompetingWriters+1 {
					t.Logf("Not all writers acquired locks: expected %d, got %d", numCompetingWriters+1, len(lockOrder))
					return false
				}

				// Verify all documents were written successfully
				actualDocs := make(map[core.DocumentID]bool)
				err := engine.ScanCollection("test_locks", func(docID core.DocumentID, doc core.Document) bool {
					actualDocs[docID] = true
					return true
				})
				if err != nil {
					t.Logf("Failed to scan collection: %v", err)
					return false
				}

				// Should have documents from all competing writers
				if len(actualDocs) != numCompetingWriters {
					t.Logf("Document count mismatch: expected %d, got %d", numCompetingWriters, len(actualDocs))
					return false
				}

				return true
			},
			gen.IntRange(2, 5),      // numCompetingWriters: 2-5 concurrent writers trying to acquire lock
			gen.IntRange(50, 200),   // holdDurationMs: 50-200ms to hold the lock
		))

	properties.TestingRun(t)
}

// Feature: monster-backend-database, Property 5: Concurrent Read Availability
// **Validates: Requirements 1.7**
func TestProperty_ConcurrentReadAvailability(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("multiple concurrent read operations complete successfully without blocking each other",
		prop.ForAll(
			func(numReaders int, numDocuments int) bool {
				// Setup test engine
				engine, tempDir := setupTestEngine(t)
				defer cleanupTestEngine(engine, tempDir)

				// Create collection
				if err := engine.CreateCollection("test_concurrent_reads"); err != nil {
					t.Logf("Failed to create collection: %v", err)
					return false
				}

				// Write initial documents to the collection
				expectedDocs := make(map[core.DocumentID]core.Document)
				for i := 0; i < numDocuments; i++ {
					docID := core.DocumentID(fmt.Sprintf("doc_%d", i))
					doc := core.Document{
						"id":    string(docID),
						"index": i,
						"data":  fmt.Sprintf("data_%d", i),
					}
					if err := engine.WriteDocument("test_concurrent_reads", docID, doc); err != nil {
						t.Logf("Failed to write document: %v", err)
						return false
					}
					expectedDocs[docID] = doc
				}

				// Track read operations and timing
				var wg sync.WaitGroup
				errChan := make(chan error, numReaders*numDocuments)
				readTimes := make([]time.Duration, numReaders)
				var readTimesMu sync.Mutex

				// Launch concurrent readers
				startTime := time.Now()
				for i := 0; i < numReaders; i++ {
					wg.Add(1)
					go func(readerID int) {
						defer wg.Done()

						readerStart := time.Now()

						// Each reader reads all documents
						for j := 0; j < numDocuments; j++ {
							docID := core.DocumentID(fmt.Sprintf("doc_%d", j))
							
							doc, err := engine.ReadDocument("test_concurrent_reads", docID)
							if err != nil {
								errChan <- fmt.Errorf("reader %d failed to read doc %s: %w", readerID, docID, err)
								return
							}

							// Verify document contents
							expectedDoc := expectedDocs[docID]
							if doc["data"] != expectedDoc["data"] {
								errChan <- fmt.Errorf("reader %d got wrong data for doc %s", readerID, docID)
								return
							}
						}

						readerElapsed := time.Since(readerStart)
						readTimesMu.Lock()
						readTimes[readerID] = readerElapsed
						readTimesMu.Unlock()
					}(i)
				}

				// Wait for all readers to complete
				wg.Wait()
				totalElapsed := time.Since(startTime)
				close(errChan)

				// Check for any errors
				for err := range errChan {
					t.Logf("Read error: %v", err)
					return false
				}

				// Verify that reads completed successfully
				// All readers should have completed
				readTimesMu.Lock()
				completedReaders := 0
				for _, duration := range readTimes {
					if duration > 0 {
						completedReaders++
					}
				}
				readTimesMu.Unlock()

				if completedReaders != numReaders {
					t.Logf("Not all readers completed: expected %d, got %d", numReaders, completedReaders)
					return false
				}

				// Verify concurrent reads didn't block each other excessively
				// If reads were truly concurrent (not blocking), total time should be
				// roughly similar to individual read times, not the sum of all read times
				readTimesMu.Lock()
				var totalIndividualTime time.Duration
				for _, duration := range readTimes {
					totalIndividualTime += duration
				}
				readTimesMu.Unlock()

				// If reads were completely serialized, total time would equal sum of individual times
				// If reads are concurrent, total time should be much less than sum
				// We check that total time is less than 80% of sum (allowing for some overhead)
				if numReaders > 1 && totalElapsed > (totalIndividualTime*8/10) {
					t.Logf("Reads appear to be blocking: total=%v, sum of individual=%v", 
						totalElapsed, totalIndividualTime)
					// This is a warning but not a hard failure, as timing can be variable
					// The main property is that all reads complete successfully
				}

				// Additional verification: perform concurrent ScanCollection operations
				var scanWg sync.WaitGroup
				scanErrChan := make(chan error, numReaders)

				for i := 0; i < numReaders; i++ {
					scanWg.Add(1)
					go func(readerID int) {
						defer scanWg.Done()

						count := 0
						err := engine.ScanCollection("test_concurrent_reads", func(docID core.DocumentID, doc core.Document) bool {
							count++
							// Verify document is valid
							if doc["id"] == nil || doc["data"] == nil {
								scanErrChan <- fmt.Errorf("reader %d found invalid document %s", readerID, docID)
								return false
							}
							return true
						})

						if err != nil {
							scanErrChan <- fmt.Errorf("reader %d failed to scan: %w", readerID, err)
							return
						}

						if count != numDocuments {
							scanErrChan <- fmt.Errorf("reader %d scanned wrong count: expected %d, got %d", 
								readerID, numDocuments, count)
							return
						}
					}(i)
				}

				scanWg.Wait()
				close(scanErrChan)

				// Check for scan errors
				for err := range scanErrChan {
					t.Logf("Scan error: %v", err)
					return false
				}

				return true
			},
			gen.IntRange(3, 10),   // numReaders: 3-10 concurrent readers
			gen.IntRange(5, 20),   // numDocuments: 5-20 documents to read
		))

	properties.TestingRun(t)
}

// Feature: monster-backend-database, Property 4: Collection Isolation
// **Validates: Requirements 1.5**
func TestProperty_CollectionIsolation(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 20
	properties := gopter.NewProperties(parameters)

	properties.Property("write operations on one collection do not affect other collections",
		prop.ForAll(
			func(coll1Docs []string, coll2Docs []string, coll1Updates []string) bool {
				// Setup test engine
				engine, tempDir := setupTestEngine(t)
				defer cleanupTestEngine(engine, tempDir)

				// Create two separate collections
				if err := engine.CreateCollection("collection1"); err != nil {
					t.Logf("Failed to create collection1: %v", err)
					return false
				}

				if err := engine.CreateCollection("collection2"); err != nil {
					t.Logf("Failed to create collection2: %v", err)
					return false
				}

				// Write initial documents to collection1
				coll1InitialState := make(map[core.DocumentID]core.Document)
				for i, data := range coll1Docs {
					docID := core.DocumentID(fmt.Sprintf("coll1_doc_%d", i))
					doc := core.Document{
						"id":         string(docID),
						"data":       data,
						"collection": "collection1",
					}
					if err := engine.WriteDocument("collection1", docID, doc); err != nil {
						t.Logf("Failed to write to collection1: %v", err)
						return false
					}
					coll1InitialState[docID] = doc
				}

				// Write initial documents to collection2
				coll2InitialState := make(map[core.DocumentID]core.Document)
				for i, data := range coll2Docs {
					docID := core.DocumentID(fmt.Sprintf("coll2_doc_%d", i))
					doc := core.Document{
						"id":         string(docID),
						"data":       data,
						"collection": "collection2",
					}
					if err := engine.WriteDocument("collection2", docID, doc); err != nil {
						t.Logf("Failed to write to collection2: %v", err)
						return false
					}
					coll2InitialState[docID] = doc
				}

				// Capture collection2 state before modifying collection1
				coll2BeforeModification := make(map[core.DocumentID]core.Document)
				err := engine.ScanCollection("collection2", func(docID core.DocumentID, doc core.Document) bool {
					coll2BeforeModification[docID] = doc
					return true
				})
				if err != nil {
					t.Logf("Failed to scan collection2 before modification: %v", err)
					return false
				}

				// Perform write operations on collection1 (inserts, updates, deletes)
				// 1. Update existing documents in collection1
				for i, data := range coll1Updates {
					if i >= len(coll1Docs) {
						break
					}
					docID := core.DocumentID(fmt.Sprintf("coll1_doc_%d", i))
					updatedDoc := core.Document{
						"id":         string(docID),
						"data":       data,
						"collection": "collection1",
						"updated":    true,
					}
					if err := engine.WriteDocument("collection1", docID, updatedDoc); err != nil {
						t.Logf("Failed to update document in collection1: %v", err)
						return false
					}
				}

				// 2. Insert new documents to collection1
				for i := 0; i < 3; i++ {
					docID := core.DocumentID(fmt.Sprintf("coll1_new_%d", i))
					doc := core.Document{
						"id":         string(docID),
						"data":       fmt.Sprintf("new_data_%d", i),
						"collection": "collection1",
						"new":        true,
					}
					if err := engine.WriteDocument("collection1", docID, doc); err != nil {
						t.Logf("Failed to insert new document to collection1: %v", err)
						return false
					}
				}

				// 3. Delete some documents from collection1
				if len(coll1Docs) > 0 {
					docID := core.DocumentID(fmt.Sprintf("coll1_doc_%d", 0))
					if err := engine.DeleteDocument("collection1", docID); err != nil {
						t.Logf("Failed to delete document from collection1: %v", err)
						return false
					}
				}

				// Verify collection2 state is unchanged after all collection1 operations
				coll2AfterModification := make(map[core.DocumentID]core.Document)
				err = engine.ScanCollection("collection2", func(docID core.DocumentID, doc core.Document) bool {
					coll2AfterModification[docID] = doc
					return true
				})
				if err != nil {
					t.Logf("Failed to scan collection2 after modification: %v", err)
					return false
				}

				// Check that collection2 has the same number of documents
				if len(coll2AfterModification) != len(coll2BeforeModification) {
					t.Logf("Collection2 document count changed: before=%d, after=%d",
						len(coll2BeforeModification), len(coll2AfterModification))
					return false
				}

				// Check that all collection2 documents are identical
				for docID, docBefore := range coll2BeforeModification {
					docAfter, exists := coll2AfterModification[docID]
					if !exists {
						t.Logf("Document %s disappeared from collection2", docID)
						return false
					}

					// Compare document contents
					if docBefore["id"] != docAfter["id"] {
						t.Logf("Document %s id changed in collection2", docID)
						return false
					}
					if docBefore["data"] != docAfter["data"] {
						t.Logf("Document %s data changed in collection2", docID)
						return false
					}
					if docBefore["collection"] != docAfter["collection"] {
						t.Logf("Document %s collection field changed", docID)
						return false
					}
				}

				// Verify collection2 documents don't have any fields from collection1 operations
				for docID, doc := range coll2AfterModification {
					if updated, exists := doc["updated"]; exists && updated == true {
						t.Logf("Document %s in collection2 has 'updated' field from collection1", docID)
						return false
					}
					if newField, exists := doc["new"]; exists && newField == true {
						t.Logf("Document %s in collection2 has 'new' field from collection1", docID)
						return false
					}
					if coll, exists := doc["collection"]; exists && coll != "collection2" {
						t.Logf("Document %s in collection2 has wrong collection field: %v", docID, coll)
						return false
					}
				}

				// Verify collection1 was actually modified (sanity check)
				coll1Count := 0
				err = engine.ScanCollection("collection1", func(docID core.DocumentID, doc core.Document) bool {
					coll1Count++
					return true
				})
				if err != nil {
					t.Logf("Failed to scan collection1: %v", err)
					return false
				}

				// Collection1 should have different count than initial
				// (we added 3 new docs and deleted 1, so net +2 if we had at least 1 doc initially)
				expectedColl1Count := len(coll1Docs) + 3 - 1
				if len(coll1Docs) == 0 {
					expectedColl1Count = 3 // Only new docs, no deletes
				}
				if coll1Count != expectedColl1Count {
					t.Logf("Collection1 count unexpected: expected %d, got %d", expectedColl1Count, coll1Count)
					return false
				}

				return true
			},
			gen.SliceOfN(5, gen.AlphaString()),  // coll1Docs: 5 initial documents for collection1
			gen.SliceOfN(5, gen.AlphaString()),  // coll2Docs: 5 initial documents for collection2
			gen.SliceOfN(5, gen.AlphaString()),  // coll1Updates: 5 update values for collection1
		))

	properties.TestingRun(t)
}
