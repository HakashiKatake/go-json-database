# Go JSON Database

A lightweight, file-based JSON database written in Go that provides simple CRUD operations with thread-safe access and data integrity.

## 🚀 Features

- **File-based Storage**: Data is stored as JSON files in the filesystem
- **Thread-Safe Operations**: Concurrent access protection using mutexes
- **Simple API**: Easy-to-use CRUD operations
- **Collection-based Organization**: Data organized in collections (similar to tables)
- **Atomic Writes**: Temporary file writes ensure data integrity
- **Logging Support**: Built-in logging with configurable levels
  

## 📁 Project Structure

```
Go-Json-Database/
├── main.go           # Example usage and demonstration
├── db/
│   └── db.go         # Core database implementation
├── types/
│   └── user.go       # Data type definitions
├── users/            # Example data storage directory
│   ├── Alice.json
│   ├── Bob.json
│   └── ...
├── go.mod
├── go.sum
└── README.md
```

## 🛠 Installation

1. Clone the repository:
```bash
git clone https://github.com/HakashiKatake/Go-Json-Database.git
cd Go-Json-Database
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the example:
```bash
go run main.go
```

## 📖 API Reference

### Database Initialization

```go
import "github.com/HakashiKatake/Go-Json-Database/db"

// Create a new database instance
database, err := db.New("./data", nil)
if err != nil {
    log.Fatal("Error creating database:", err)
}
```

### Core Operations

#### Write Operation
Store data in a collection with a unique resource identifier.

```go
err := database.Write("users", "john_doe", user.User{
    Name:    "John Doe",
    Age:     "30",
    Contact: "1234567890",
    Company: "Tech Corp",
    Address: user.Address{
        City:    "New York",
        State:   "NY",
        Country: "USA",
        Pincode: "10001",
    },
})
```

#### Read Operation
Retrieve a specific record from a collection.

```go
var userData user.User
err := database.Read("users", "john_doe", &userData)
if err != nil {
    log.Println("Error reading user:", err)
}
```

#### ReadAll Operation
Retrieve all records from a collection.

```go
records, err := database.ReadAll("users")
if err != nil {
    log.Println("Error reading all users:", err)
}

// Parse the JSON records
var allUsers []user.User
for _, record := range records {
    var user user.User
    if err := json.Unmarshal([]byte(record), &user); err == nil {
        allUsers = append(allUsers, user)
    }
}
```

#### Delete Operation
Remove a record from a collection.

```go
err := database.Delete("users", "john_doe")
if err != nil {
    log.Println("Error deleting user:", err)
}
```

## 🏗 Data Types

### User Structure
```go
type User struct {
    Name    string      `json:"name"`
    Age     json.Number `json:"age"`
    Contact string      `json:"contact"`
    Company string      `json:"company"`
    Address Address     `json:"address"`
}
```

### Address Structure
```go
type Address struct {
    City    string      `json:"city"`
    State   string      `json:"state"`
    Country string      `json:"country"`
    Pincode json.Number `json:"pincode"`
}
```

## 🔧 Configuration Options

### Logger Configuration
```go
import "github.com/jcelliott/lumber"

options := &db.Options{
    Logger: lumber.NewConsoleLogger(lumber.DEBUG),
}

database, err := db.New("./data", options)
```

### Available Log Levels
- `lumber.TRACE`
- `lumber.DEBUG`
- `lumber.INFO`
- `lumber.WARN`
- `lumber.ERROR`
- `lumber.FATAL`

## 💡 Usage Examples

### Basic CRUD Operations
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    "github.com/HakashiKatake/Go-Json-Database/db"
    user "github.com/HakashiKatake/Go-Json-Database/types"
)

func main() {
    // Initialize database
    database, err := db.New("./", nil)
    if err != nil {
        log.Fatal("Error creating database:", err)
    }

    // Create a user
    newUser := user.User{
        Name:    "Alice Smith",
        Age:     "28",
        Contact: "9876543210",
        Company: "Innovation Labs",
        Address: user.Address{
            City:    "San Francisco",
            State:   "CA",
            Country: "USA",
            Pincode: "94102",
        },
    }

    // Write user to database
    err = database.Write("users", "alice_smith", newUser)
    if err != nil {
        log.Println("Error writing user:", err)
    }

    // Read user from database
    var retrievedUser user.User
    err = database.Read("users", "alice_smith", &retrievedUser)
    if err != nil {
        log.Println("Error reading user:", err)
    } else {
        fmt.Printf("Retrieved user: %+v\n", retrievedUser)
    }

    // Read all users
    records, err := database.ReadAll("users")
    if err != nil {
        log.Println("Error reading all users:", err)
    } else {
        fmt.Printf("Total records: %d\n", len(records))
    }

    // Delete user
    err = database.Delete("users", "alice_smith")
    if err != nil {
        log.Println("Error deleting user:", err)
    } else {
        fmt.Println("User deleted successfully")
    }
}
```

## 🔒 Thread Safety

The database implements thread-safe operations using:
- **Collection-level mutexes**: Each collection has its own mutex
- **Atomic writes**: Data is written to temporary files first, then renamed
- **Concurrent read support**: Multiple goroutines can read simultaneously

## 📂 File Storage Format

Data is stored as formatted JSON files:

```json
{
    "Name": "John Doe",
    "Age": 30,
    "Contact": "1234567890",
    "Company": "Tech Corp",
    "Address": {
        "City": "New York",
        "State": "NY",
        "Country": "USA",
        "Pincode": 10001
    }
}
```

## ⚠️ Error Handling

The database returns descriptive errors for common scenarios:
- Missing collection or resource names
- File system permission issues
- JSON marshaling/unmarshaling errors
- Non-existent records or collections

## 🚦 System Architecture

- DB's System Structure and Methods

<img width="816" alt="Screenshot 2025-06-10 at 1 15 55 PM" src="https://github.com/user-attachments/assets/12acd6ce-f193-48a5-934b-1037c36b3c4e" />

- Methods 

<img width="355" alt="Screenshot 2025-06-10 at 1 16 27 PM" src="https://github.com/user-attachments/assets/c7add492-0017-4647-8ff3-ad09a30806e8" />

- Data integrity 

<img width="394" alt="Screenshot 2025-06-10 at 1 16 44 PM" src="https://github.com/user-attachments/assets/af7e0b77-443d-49b6-bfc7-ad619de0b34c" />

# Full design of the system.

![db drawio](https://github.com/user-attachments/assets/1d1787e9-e8ba-4fde-8d91-c5b7d61cb65f)


### Core Components

1. **Driver**: Main database engine that handles all operations
2. **Collections**: Logical groupings of related data (like database tables)
3. **Resources**: Individual records within collections
4. **Mutexes**: Thread-safety mechanisms for concurrent access
5. **Logger**: Configurable logging system for debugging and monitoring

### Database Operations Flow

```
Write Operation:
User Data → JSON Marshal → Temporary File → Atomic Rename → Final File

Read Operation:
File Path → File Read → JSON Unmarshal → User Data

Delete Operation:
Resource Path → File/Directory Removal
```

### Thread Safety Model

- Each collection gets its own mutex
- Write operations are fully locked
- Read operations can happen concurrently
- Atomic file operations prevent data corruption

## 🔍 Database Methods

### Core Methods

| Method | Description | Parameters | Returns |
|--------|-------------|------------|---------|
| `New()` | Initialize database | `dir string, options *Options` | `*Driver, error` |
| `Write()` | Store record | `collection, resource string, v interface{}` | `error` |
| `Read()` | Retrieve record | `collection, resource string, v interface{}` | `error` |
| `ReadAll()` | Get all records | `collection string` | `[]string, error` |
| `Delete()` | Remove record | `collection, resource string` | `error` |

### Internal Methods

| Method | Description | Purpose |
|--------|-------------|---------|
| `getOrCreateMutex()` | Manage collection mutexes | Thread safety |
| `stat()` | File existence check | File operations |

## 📊 Performance Characteristics

### Strengths
- **Simple API**: Easy to learn and use
- **No Dependencies**: Minimal external requirements
- **File-based**: Human-readable storage format
- **Thread-safe**: Concurrent access support
- **Atomic Operations**: Data integrity guaranteed

### Limitations
- **File I/O Bound**: Performance limited by disk speed
- **Memory Usage**: Entire records loaded into memory
- **No Indexing**: Linear search for record retrieval
- **No Transactions**: No multi-operation atomicity
- **No Query Language**: Basic key-value access only

### Recommended Use Cases
- **Prototyping**: Quick database setup for development
- **Small Applications**: < 10,000 records per collection
- **Configuration Storage**: Application settings and metadata
- **Logging**: Structured log data storage
- **Testing**: Mock database for unit tests

## 🛡️ Best Practices

### Data Management
1. **Use descriptive collection names** (e.g., "users", "products", "orders")
2. **Choose meaningful resource identifiers** (e.g., user IDs, product SKUs)
3. **Validate data before writing** to prevent corruption
4. **Use json.Number for numeric fields** to maintain precision
5. **Implement proper error handling** for all operations

### Performance Optimization
1. **Batch operations** when possible to reduce I/O
2. **Use appropriate data structures** for your use case
3. **Monitor disk space** usage regularly
4. **Consider data archiving** for old records
5. **Implement caching** for frequently accessed data

### Security Considerations
1. **Validate input data** to prevent injection attacks
2. **Set appropriate file permissions** on data directories
3. **Backup data regularly** to prevent data loss
4. **Monitor access patterns** for unusual activity
5. **Consider encryption** for sensitive data

## 🔧 Troubleshooting

### Common Issues

**Permission Denied Errors**
```bash
# Fix file permissions
chmod 755 ./data
chmod 644 ./data/collection/*.json
```

**JSON Marshal/Unmarshal Errors**
```go
// Ensure struct fields are exported (capitalized)
type User struct {
    Name string `json:"name"` // ✅ Correct
    age  string `json:"age"`  // ❌ Won't marshal
}
```

**Concurrent Access Issues**
```go
// The database handles this automatically
// No additional locking required in user code
```

### Debug Mode
```go
import "github.com/jcelliott/lumber"

options := &db.Options{
    Logger: lumber.NewConsoleLogger(lumber.DEBUG),
}
database, err := db.New("./data", options)
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup
```bash
git clone https://github.com/HakashiKatake/Go-Json-Database.git
cd Go-Json-Database
go mod tidy
go test ./...
```

## 📝 License

This project is open source and available under the [MIT License](LICENSE).

## 🔗 Dependencies

- [lumber](https://github.com/jcelliott/lumber) - Logging library for Go

## 📈 Roadmap

### Planned Features
- [ ] Query language support
- [ ] Indexing for faster lookups
- [ ] Transaction support
- [ ] Data compression
- [ ] Backup/restore utilities
- [ ] Web interface for data management
- [ ] Replication support
- [ ] Schema validation

### Version History
- **v1.0.1**: Current stable release
- **v1.0.0**: Initial release with basic CRUD operations

---

**Version**: 1.0.1  
**Go Version**: 1.24.2  
**Maintained by**: [HakashiKatake](https://github.com/HakashiKatake)
