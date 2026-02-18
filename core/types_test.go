package core

import (
	"testing"
)

// TestDocumentCreation verifies Document type can be created and used
func TestDocumentCreation(t *testing.T) {
	doc := Document{
		"id":   "test_001",
		"name": "Test Document",
		"age":  25,
	}

	if doc["id"] != "test_001" {
		t.Errorf("Expected id to be 'test_001', got %v", doc["id"])
	}

	if doc["name"] != "Test Document" {
		t.Errorf("Expected name to be 'Test Document', got %v", doc["name"])
	}

	if doc["age"] != 25 {
		t.Errorf("Expected age to be 25, got %v", doc["age"])
	}
}

// TestDocumentID verifies DocumentID type
func TestDocumentID(t *testing.T) {
	var id DocumentID = "doc_123"
	if string(id) != "doc_123" {
		t.Errorf("Expected DocumentID to be 'doc_123', got %s", id)
	}
}

// TestFilterOperators verifies FilterOperator constants
func TestFilterOperators(t *testing.T) {
	tests := []struct {
		op       FilterOperator
		expected int
	}{
		{OpEqual, 0},
		{OpGreaterThan, 1},
		{OpLessThan, 2},
		{OpGreaterThanOrEqual, 3},
		{OpLessThanOrEqual, 4},
	}

	for _, tt := range tests {
		if int(tt.op) != tt.expected {
			t.Errorf("Expected operator value %d, got %d", tt.expected, int(tt.op))
		}
	}
}

// TestOperationTypes verifies OperationType constants
func TestOperationTypes(t *testing.T) {
	tests := []struct {
		op       OperationType
		expected int
	}{
		{OpInsert, 0},
		{OpUpdate, 1},
		{OpDelete, 2},
	}

	for _, tt := range tests {
		if int(tt.op) != tt.expected {
			t.Errorf("Expected operation type value %d, got %d", tt.expected, int(tt.op))
		}
	}
}

// TestQueryCreation verifies Query type can be created
func TestQueryCreation(t *testing.T) {
	query := Query{
		Collection: "users",
		Filters: []Filter{
			{Field: "age", Operator: OpGreaterThan, Value: 18},
		},
		Sort:   &SortOption{Field: "name", Descending: false},
		Limit:  10,
		Offset: 0,
	}

	if query.Collection != "users" {
		t.Errorf("Expected collection to be 'users', got %s", query.Collection)
	}

	if len(query.Filters) != 1 {
		t.Errorf("Expected 1 filter, got %d", len(query.Filters))
	}

	if query.Limit != 10 {
		t.Errorf("Expected limit to be 10, got %d", query.Limit)
	}
}

// TestTransactionCreation verifies Transaction type can be created
func TestTransactionCreation(t *testing.T) {
	txn := Transaction{
		ID: "txn_001",
		Operations: []Operation{
			{
				Type:       OpInsert,
				Collection: "users",
				DocID:      "user_001",
				Document:   Document{"name": "Alice"},
			},
		},
		Committed: false,
	}

	if txn.ID != "txn_001" {
		t.Errorf("Expected transaction ID to be 'txn_001', got %s", txn.ID)
	}

	if len(txn.Operations) != 1 {
		t.Errorf("Expected 1 operation, got %d", len(txn.Operations))
	}

	if txn.Committed {
		t.Error("Expected transaction to not be committed")
	}
}
