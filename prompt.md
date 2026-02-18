# MASTER BUILD PROMPT — GO JSON DATABASE → MONSTER BACKEND PROJECT

You are a senior backend systems engineer and database architect.

Your task is to evolve the existing **Go JSON Database** project into a production-grade, high-performance lightweight database engine that demonstrates advanced backend and systems engineering.

This is NOT a beginner CRUD project.  
This must become a technically impressive, recruiter-level backend system.

All design decisions must prioritize:
- performance
- concurrency safety
- scalability
- clean architecture
- production realism
- engineering depth

Avoid shortcuts. Build like this will be used in real systems.

---

# CORE OBJECTIVE
Transform the current file-based JSON DB into:

> A lightweight concurrent embedded database engine with indexing, query engine, transactions, benchmarking, and REST interface built in Go.

The codebase must reflect strong backend/system design skills.

---

# REQUIRED SYSTEM ARCHITECTURE

Use clean modular architecture:

/core
/storage
/index
/query
/transaction
/wal
/api
/benchmark
/tests


Each module must be decoupled and testable.

Follow idiomatic Go patterns and clean code practices.

---

# PHASE 1 — STORAGE ENGINE IMPROVEMENTS

## 1. Thread-safe storage layer
Implement:
- file locking
- mutex protected operations
- safe concurrent read/write
- atomic file writes

Support:
- collections
- document-based storage
- configurable storage path

Ensure zero data corruption under concurrent writes.

---

# PHASE 2 — INDEXING SYSTEM (CRITICAL)

Implement indexing to avoid full file scans.

### Required:
- primary key index (hash map)
- optional secondary index
- auto update index on insert/update/delete
- persistent index storage

Target:
> O(1) key lookup performance

Add functions:
- CreateIndex(collection, field)
- DropIndex()
- RebuildIndex()

---

# PHASE 3 — QUERY ENGINE

Add advanced query support.

Example:
Find("users", {
age: { $gt: 18 },
city: "Mumbai"
})


Support:
- equality filters
- > < >= <=
- AND conditions
- limit/offset
- basic sorting

Must use indexes where available.

Avoid full scan when indexed.

---

# PHASE 4 — TRANSACTION SYSTEM

Implement basic ACID-like transactions.

### Required:
- BeginTransaction()
- Commit()
- Rollback()

Use:
- in-memory write buffer
- apply on commit
- discard on rollback

Ensure isolation during concurrent transactions.

---

# PHASE 5 — WRITE AHEAD LOG (WAL)

Implement crash recovery.

Before write:
1. write operation to WAL
2. apply change to DB
3. clear log on success

On restart:
- replay WAL
- restore consistent state

Ensure durability.

---

# PHASE 6 — REST API LAYER

Expose database as HTTP service.

Endpoints:
- POST /insert
- GET /find
- PUT /update
- DELETE /delete
- GET /collections

Add:
- API key auth
- rate limiting
- JSON responses
- error handling

Use:
Go net/http or Gin (prefer standard lib if possible)

---

# PHASE 7 — PERFORMANCE & BENCHMARKING (VERY IMPORTANT)

Create benchmark suite.

Test:
- 1k inserts
- 10k inserts
- concurrent writes (100 goroutines)
- indexed vs non-indexed queries
- memory usage
- latency

Compare with:
- SQLite (JSON mode)
- plain file system

Output results in:
/benchmarks/results.md


Include charts if possible.

---

# PHASE 8 — LOAD TESTING

Simulate:
- 100 concurrent users
- 1000 writes/min
- heavy read/write mixed workload

Ensure:
- no race conditions
- no corruption
- stable latency

---

# PHASE 9 — DOCKERIZATION

Provide:
- Dockerfile
- docker-compose
- run as standalone DB service

Command:
docker run jsondb


---

# PHASE 10 — DOCUMENTATION (CRUCIAL)

Create elite README with:

## Sections:
1. What this project is
2. Why it exists
3. Architecture diagram
4. Storage engine design
5. Indexing design
6. Transaction model
7. Benchmarks
8. Tradeoffs vs Mongo/Postgres
9. Use cases
10. Future roadmap

Write like system engineer, not student.

---

# ENGINEERING STANDARDS

## Code quality
- idiomatic Go
- comments explaining WHY not WHAT
- modular design
- no messy monolith files

## Performance mindset
- avoid unnecessary allocations
- use goroutines properly
- profile when needed
- optimize critical paths

## Git discipline
- meaningful commits
- feature branches
- clean history
- version tags (v1.0 later)

---

# FINAL GOAL

This project must look like it was built by:
> a serious backend engineer obsessed with systems

When someone opens GitHub, reaction should be:
> “Why is a student building database engines?”

Do NOT build quickly.  
Build correctly.

Focus on engineering depth, performance, and architecture clarity.