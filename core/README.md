# Core Package

This package contains the core types and interfaces for the Monster Backend Database.

## Types

- **Document**: A JSON document stored in the database (map[string]interface{})
- **DocumentID**: Unique identifier for a document (string)
- **Collection**: Logical grouping of documents
- **Query**: Database query with filters and options
- **Filter**: Query filter condition with operator and value
- **Transaction**: ACID transaction with buffered operations
- **Operation**: Single database operation (insert/update/delete)

## Interfaces

- **StorageEngine**: Defines storage operations (write, read, delete, scan)
- **IndexManager**: Defines index operations (create, lookup, update, persist)

## Enums

- **FilterOperator**: Comparison operators (Equal, GreaterThan, LessThan, etc.)
- **OperationType**: Operation types (Insert, Update, Delete)
