# String Analysis API

A REST API built with Go, Gin, and GORM for analyzing strings with palindrome detection, character frequency analysis, and SHA-256 hashing.

## Prerequisites

- Go 1.16 or higher
- Git (optional)

## Installation

```bash
# Create project directory
mkdir string-analysis-api && cd string-analysis-api

# Initialize Go module
go mod init string-analysis-api

# Install dependencies
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite

# Create main.go and add your code

# Run the application
go run main.go

# Or build for production
go build -o string-api
./string-api
```

Server runs on `http://localhost:8080`

## API Endpoints

### 1. Create String
```bash
POST /strings
Content-Type: application/json

{"value": "racecar"}
```

### 2. Get String
```bash
GET /strings/{string_value}
```

### 3. Get All Strings (with filters)
```bash
GET /strings?is_palindrome=true&min_length=5&max_length=20&word_count=1&contains_character=a
```

### 4. Natural Language Filter
```bash
GET /strings/filter-by-natural-language?query=all%20single%20word%20palindromic%20strings
```

### 5. Delete String
```bash
DELETE /strings/{string_value}
```

## Example Usage

```bash
# Create a string
curl -X POST http://localhost:8080/strings \
  -H "Content-Type: application/json" \
  -d '{"value": "racecar"}'

# Get palindromes
curl "http://localhost:8080/strings?is_palindrome=true"

# Natural language query
curl "http://localhost:8080/strings/filter-by-natural-language?query=strings%20longer%20than%2010%20characters"
```

## Response Format

```json
{
  "id": "sha256_hash",
  "value": "racecar",
  "properties": {
    "length": 7,
    "is_palindrome": true,
    "unique_characters": 4,
    "word_count": 1,
    "sha256_hash": "abc123...",
    "character_frequency_map": {"r": 2, "a": 2, "c": 2, "e": 1}
  },
  "created_at": "2025-10-21T10:00:00Z"
}
```

## Configuration

Change port in `main.go`:
```go
r.Run(":8080")  // Change to ":3000" etc
```

Use PostgreSQL/MySQL:
```go
import "gorm.io/driver/postgres"
dsn := "host=localhost user=user password=pass dbname=strings"
db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

## Error Codes

- `201` - Created
- `200` - OK
- `204` - No Content (deleted)
- `400` - Bad Request
- `404` - Not Found
- `409` - Conflict (duplicate)
- `422` - Unprocessable Entity
