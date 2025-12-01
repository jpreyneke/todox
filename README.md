# TodoX API

A RESTful Todo API built with Go for managing todo items with bulk operations support.

## Tech Stack

- Go 1.25.4 with Gin framework
- MySQL 8.0
- Uber-Fx for dependency injection
- Cobra for CLI
- Docker & Docker Compose

## Quick Start

### Using Docker (Recommended)

```bash
# Start everything (API + MySQL + Prometheus)
make docker-up

# API available at: http://localhost:8080
# Prometheus at: http://localhost:9090
```

The `make docker-up` command automatically runs migrations before starting the API.

### Local Development

```bash
# 1. Setup environment
cp .env.example .env

# 2. Start MySQL (update .env with your connection details)

# 3. Run migrations
make migrate

# 4. Start API
make run
```

## API Usage

### Create Todos
```bash
curl -X POST http://localhost:8080/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "todos": [
      {
        "title": "Buy groceries",
        "description": "Milk, eggs, bread",
        "due_date": "2025-12-15T10:00:00Z"
      }
    ]
  }'
```

### Update Todos
```bash
curl -X PATCH http://localhost:8080/v1/todos \
  -H "Content-Type: application/json" \
  -d '{
    "todos": [
      {"id": 1, "completed": true}
    ]
  }'
```

### List Todos (Paginated)
```bash
curl "http://localhost:8080/v1/todos?page=1&limit=10"
```

### With Authentication (Optional)
If `API_KEY` is set in environment:
```bash
curl -H "X-API-Key: your-key" http://localhost:8080/v1/todos
```

## CLI Commands

```bash
# Run API server
go run ./cmd/api api

# Run migrations
go run ./cmd/migrate

# Rollback last migration
go run ./cmd/migrate down
```

## Configuration

### Environment Variables

Configuration can be set in two ways:

**For Local Development:**
```bash
cp .env.example .env
# Edit .env with your settings
```

**For Docker:**
Environment variables are set in `docker-compose.yml` (`.env` file is not used)

## Testing

```bash
# Run all tests
make test

# With coverage
go test -v -race -coverprofile=coverage.out ./...
```

## Assumptions & Validation Rules

### Title
- **Required**: Cannot be empty
- **Unique**: Must be unique across all todos (enforced by database)
- **Max Length**: 255 characters
- Whitespace is trimmed automatically

### Description
- **Optional**: Can be omitted
- **Max Length**: No limit (stored as TEXT, ~64KB in MySQL)

### Due Date
- **Optional**: Can be omitted
- **Format**: RFC3339 timestamp (e.g., `2025-12-15T10:00:00Z`)
- **No Validation**: Past dates are allowed

### Bulk Operations
- **Minimum**: 1 todo required (empty arrays rejected)
- **Maximum**: No hard limit
- **Duplicates**: Duplicate titles or IDs within same request are rejected
- **Transactions**: All items succeed or all fail together

### Pagination
- **Default**: page=1, limit=10
- **Maximum Limit**: 100 items per page
- **Invalid Values**: Automatically corrected to defaults

### General
- All operations return todos ordered by creation date (newest first)
- Timestamps stored and returned in UTC
- Empty todo list returns empty array (not null)

## Known Limitations

- No delete endpoint (not in specification)
- No filtering on GET /todos (e.g., by completed status, due date range)
- No user/tenant isolation (single shared todo list)
- Title uniqueness is global, not per-user

## Monitoring

- **Health**: http://localhost:8080/health
- **Metrics**: http://localhost:8080/metrics (Prometheus format)
- **Prometheus UI**: http://localhost:9090 (when using Docker)

## Troubleshooting

Contact: jpreyneke1@gmail.com