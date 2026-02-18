package core

// Document represents a JSON document stored in the database
type Document map[string]interface{}

// DocumentID is the unique identifier for a document
type DocumentID string

// Collection represents a logical grouping of documents
type Collection struct {
	Name    string
	Storage StorageEngine
	Indexes IndexManager
}

// Query represents a database query with filters and options
type Query struct {
	Collection string
	Filters    []Filter
	Sort       *SortOption
	Limit      int
	Offset     int
}

// Filter represents a query filter condition
type Filter struct {
	Field    string
	Operator FilterOperator
	Value    interface{}
}

// FilterOperator defines comparison operators
type FilterOperator int

const (
	OpEqual FilterOperator = iota
	OpGreaterThan
	OpLessThan
	OpGreaterThanOrEqual
	OpLessThanOrEqual
)

// SortOption defines sorting configuration
type SortOption struct {
	Field      string
	Descending bool
}

// Transaction represents an ACID transaction
type Transaction struct {
	ID         string
	Operations []Operation
	Committed  bool
}

// Operation represents a single database operation
type Operation struct {
	Type       OperationType
	Collection string
	DocID      DocumentID
	Document   Document
}

// OperationType defines operation types
type OperationType int

const (
	OpInsert OperationType = iota
	OpUpdate
	OpDelete
)

// StorageEngine interface defines storage operations
type StorageEngine interface {
	// WriteDocument atomically writes a document to storage
	WriteDocument(collection string, docID DocumentID, doc Document) error

	// ReadDocument retrieves a document by ID
	ReadDocument(collection string, docID DocumentID) (Document, error)

	// DeleteDocument removes a document from storage
	DeleteDocument(collection string, docID DocumentID) error

	// ScanCollection iterates over all documents in a collection
	ScanCollection(collection string, fn func(DocumentID, Document) bool) error

	// CreateCollection initializes a new collection
	CreateCollection(name string) error

	// ListCollections returns all collection names
	ListCollections() ([]string, error)

	// Close flushes pending writes and releases locks
	Close() error
}

// IndexManager interface defines index operations
type IndexManager interface {
	// CreatePrimaryIndex builds the primary key index
	CreatePrimaryIndex(collection string) error

	// CreateSecondaryIndex builds an index on a field
	CreateSecondaryIndex(collection string, field string) error

	// LookupPrimary performs O(1) lookup by primary key
	LookupPrimary(collection string, docID DocumentID) (Document, error)

	// LookupSecondary finds documents matching a field value
	LookupSecondary(collection string, field string, value interface{}) ([]Document, error)

	// UpdateIndexes updates all indexes after a write operation
	UpdateIndexes(collection string, docID DocumentID, doc Document, op OperationType) error

	// PersistIndexes writes indexes to disk
	PersistIndexes(collection string) error

	// LoadIndexes reads indexes from disk
	LoadIndexes(collection string) error

	// RebuildIndexes reconstructs indexes from storage
	RebuildIndexes(collection string) error
}
