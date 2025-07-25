# Task Management API - Assignment Submission

## Assignment Context

**Problem Statement**: Build a simple Task Management System in Go demonstrating microservices architecture with CRUD operations, pagination, filtering, and horizontal scalability.

**Solution**: A production-ready RESTful API with clean architecture, comprehensive testing, Docker containerization, nginx load balancing, and horizontal scaling.

## Key Design Decisions

- **Clean Architecture**: HTTP → Service → Repository → Database layers for maintainability and testability
- **Interface-Based Design**: All layers use interfaces, enabling dependency injection, easy mocking for testing, and flexible implementation swapping (e.g., switching from PostgreSQL to MySQL, REST to GraphQL/gRPC by implementing the same interfaces)
- **Validation Middleware**: Centralized request validation using struct tags and go-playground/validator, ensuring only valid data reaches handlers with relevant defaults, resulting in clean, minimal endpoint code
- **Error Handling**: Typed errors with middleware for consistent API responses, handling both application errors and panic recovery
- **Docker-First**: Everything containerized with docker-compose for easy deployment and scaling
- **Separate-Container Architecture**: Load balancer (nginx), application servers, and database (postgres) run in separate containers, enabling horizontal scaling with configurable application server instances
- **Testing Strategy**: Unit tests (service layer with mocks), integration tests (repository layer with real DB), handler tests (HTTP layer with mocks) - all commands available in Makefile

## API Endpoints Overview

| Method | Endpoint | Description | Body | Query Params |
|--------|----------|-------------|------|--------------|
| GET | `/health` | Health check with instance info | - | - |
| POST | `/api/v1/tasks` | Create new task | `title*`, `description`, `status` | - |
| GET | `/api/v1/tasks` | Get all tasks with pagination | - | `page`, `limit`, `status`, `sort_by`, `sort_order` |
| GET | `/api/v1/tasks/{id}` | Get specific task | - | - |
| PUT | `/api/v1/tasks/{id}` | Update existing task | `title`, `description`, `status` | - |
| DELETE | `/api/v1/tasks/{id}` | Delete task | - | - |

*Fields marked with `*` are required*

**Valid Task Statuses**: `pending` (default), `in_progress`, `completed`, `closed`  
**Pagination**: Default `page=1, limit=10`, max `limit=100`  
**Sorting**: By `id`, `title`, `status`, `created_at`, `updated_at` (asc/desc, default: `created_at desc`)

## Assignment Scope Limitations

**Intentionally not implemented for assignment focus**:
- **Database Migrations**: Using single migration file executed at container startup via Docker's initdb.d; production would need versioned migration tools like golang-migrate for schema evolution
- **Logging**: No structured logging middleware or in-code logging implemented
- **Advanced Load Balancing**: Currently using Docker's internal DNS with 0-second cache for demonstration; production needs Kubernetes/Docker Swarm
- **Authentication/Authorization**: Assumed all requests are valid (out of scope)
- **Rate Limiting**: No DDoS protection implemented
- **CORS Configuration**: No frontend integration setup


## Quick Start

### Option 1: Build from Source
**Prerequisites**: Docker and Docker Compose installed and running

```bash
git clone https://github.com/AashishRichhariya/task-management-api.git
cd task-management-api

# Start with default 3 instances
make dev-scale

# Test load balancing
make test-load
```

## Horizontal Scaling Demo
```bash
# Scale to 5 instances
make dev-scale INSTANCES=5

# Verify distribution across instances
make test-load
# (Verify different requests being sent to different instances)
```

**Repository**: [https://github.com/AashishRichhariya/task-management-api](https://github.com/AashishRichhariya/task-management-api)

## API Testing Examples

### 1. Health Check & Load Balancing
```bash
# Check service health
curl http://localhost/health

# Test load balancing (see different instance IDs)
make test-load
```

### 2. Create Tasks with Different Statuses
```bash
# Create sample tasks
curl -X POST http://localhost/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Go","description":"Study Go fundamentals","status":"pending"}'

curl -X POST http://localhost/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Build API","description":"Create REST API","status":"in_progress"}'

curl -X POST http://localhost/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Write Tests","description":"Unit and integration tests","status":"completed"}'

# Minimal task (title only, status defaults to pending)
curl -X POST http://localhost/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Minimal Task"}'
```

### 3. Get All Tasks (Pagination & Filtering)
```bash
# Get all tasks (default pagination)
curl http://localhost/api/v1/tasks

# Specific page and limit
curl "http://localhost/api/v1/tasks?page=1&limit=5"

# Filter by status
curl "http://localhost/api/v1/tasks?status=pending"
curl "http://localhost/api/v1/tasks?status=completed"

# Sorting
curl "http://localhost/api/v1/tasks?sort_by=title&sort_order=asc"
curl "http://localhost/api/v1/tasks?sort_by=created_at&sort_order=desc"

# Complex filtering
curl "http://localhost/api/v1/tasks?status=completed&sort_by=created_at&sort_order=desc&limit=3"
```

### 4. Get Specific Task by ID
```bash
# Get task by ID
curl http://localhost/api/v1/tasks/1

# Test non-existent task (404 error)
curl http://localhost/api/v1/tasks/999
```

### 5. Update Task
```bash
# Full update
curl -X PUT http://localhost/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated Title","description":"Updated description","status":"completed"}'

# Partial update (status only)
curl -X PUT http://localhost/api/v1/tasks/2 \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'

# Fetch updated task
curl http://localhost/api/v1/tasks/1
```

### 6. Delete Task
```bash
# Delete task
curl -X DELETE http://localhost/api/v1/tasks/1

# Verify deletion (should return 404)
curl http://localhost/api/v1/tasks/1
```

### 7. Error Handling Examples
```bash
# Validation error (400) - empty title
curl -X POST http://localhost/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"","status":"invalid"}'

# Invalid ID parameter (400)
curl http://localhost/api/v1/tasks/abc

# Invalid query parameters (400)
curl "http://localhost/api/v1/tasks?status=invalid_status"

# Not found error (404)
curl http://localhost/api/v1/tasks/999
```
