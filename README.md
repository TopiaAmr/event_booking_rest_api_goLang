# Event Booking REST API - Go

A simple REST API for event booking built with Go and Gin framework.

## Features

- CRUD operations for events
- SQLite database storage
- RESTful API endpoints
- Unit tests for all functionality

## API Endpoints

- `GET /events` - Get all events
- `GET /events/:id` - Get a specific event by ID
- `POST /event` - Create a new event
- `PUT /events/:id` - Update an existing event
- `DELETE /events/:id` - Delete an event

## Running the Application

1. Install dependencies:
```bash
go mod download
```

2. Run the application:
```bash
go run main.go
```

The server will start on port 8080.

## Running Tests

Run all tests:
```bash
go test ./... -v
```

Run tests for a specific package:
```bash
go test ./models -v
go test ./routes -v
go test ./db -v
```

## Test Coverage

The project includes comprehensive unit tests covering:

- **Models**: Event CRUD operations (Save, GetAllEvents, GetEventById, Update, Delete, GetEventsByUserId)
- **Routes**: HTTP handlers for all API endpoints
- **Database**: Database initialization and connection handling

### Test Structure

- `models/event_test.go` - Tests for Event model methods
- `routes/events_test.go` - Tests for HTTP handlers
- `db/db_test.go` - Tests for database operations
- `testutils/testutils.go` - Common testing utilities

## Database Schema

The `events` table has the following structure:

```sql
CREATE TABLE events (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    location TEXT NOT NULL,
    datetime DATETIME NOT NULL,
    user_id TEXT
);
```

## Dependencies

- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/google/uuid` - UUID generation

## Project Structure

```
├── main.go              # Application entry point
├── go.mod               # Go module file
├── go.sum               # Go module checksums
├── db/
│   ├── db.go           # Database initialization
│   └── db_test.go      # Database tests
├── models/
│   ├── event.go        # Event model and methods
│   └── event_test.go   # Event model tests
├── routes/
│   ├── routes.go       # Route registration
│   ├── events.go       # Event handlers
│   └── events_test.go  # Route handler tests
└── testutils/
    └── testutils.go    # Testing utilities
```
