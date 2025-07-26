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

| Method | Endpoint             | Description                     | Body                              | Query Params                                       |
| ------ | -------------------- | ------------------------------- | --------------------------------- | -------------------------------------------------- |
| GET    | `/health`            | Health check with instance info | -                                 | -                                                  |
| POST   | `/api/v1/tasks`      | Create new task                 | `title*`, `description`, `status` | -                                                  |
| GET    | `/api/v1/tasks`      | Get all tasks with pagination   | -                                 | `page`, `limit`, `status`, `sort_by`, `sort_order` |
| GET    | `/api/v1/tasks/{id}` | Get specific task               | -                                 | -                                                  |
| PUT    | `/api/v1/tasks/{id}` | Update existing task            | `title`, `description`, `status`  | -                                                  |
| DELETE | `/api/v1/tasks/{id}` | Delete task                     | -                                 | -                                                  |

_Fields marked with `_` are required\*

**Valid Task Statuses**: `pending` (default), `in_progress`, `completed`, `closed`  
**Pagination**: Default `page=1, limit=10`, max `limit=100`  
**Sorting**: By `id`, `title`, `status`, `created_at`, `updated_at` (asc/desc, default: `created_at desc`)

## Quick Start

### Option 1: Build from Source

**Prerequisites**:

- Docker and Docker Compose installed and running
- Make command available:
  - **Linux**: Usually pre-installed
  - **Mac**: Install with `xcode-select --install`
  - **Windows**: Use Git Bash (comes with Git for Windows) or WSL

```bash
git clone https://github.com/AashishRichhariya/task-management-api.git
cd task-management-api
```

- **With Make** (Recommended)

```
# Start with default 3 instances
make dev-scale


# Test load balancing
make test-load
```

- **Without Make** (Windows CMD/PowerShell)

```
# Start with 3 instances
docker-compose up --build --scale app=3


# Test load balancing (Windows PowerShell)
for ($i=1; $i -le 10; $i++) { curl -s http://localhost/health; echo "" }


# Test load balancing (Windows CMD)
for /l %i in (1,1,10) do (curl -s http://localhost/health & echo.)
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

## Architecture Deep Dive

### Clean Architecture Implementation

The system follows Clean Architecture principles with clear layer separation, ensuring maintainability, testability, and scalability:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │    │  Service Layer  │    │Repository Layer │    │   Database      │
│   (Handlers)    │───▶│ (Business Logic)│───▶│  (Data Access)  │───▶│  (PostgreSQL)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘    └─────────────────┘
        ▲                        ▲                        ▲
        │                        │                        │
TaskHandlerInterface    TaskServiceInterface     TaskRepository
```

_The arrows represent data flow through interfaces - each layer calls the next via interface contracts, not direct dependencies._

### Layer Responsibilities

**HTTP Layer (`internal/handlers/`)**:

- Receives HTTP requests and validates routes
- Delegates to middleware for validation
- Calls service layer for business operations
- Returns structured JSON responses
- **No business logic** - purely HTTP concerns

**Service Layer (`internal/service/`)**:

- Contains all business logic and rules
- Orchestrates data operations via repository
- Handles domain-specific validations (title trimming/length, valid status transitions)
- Manages transaction boundaries (could coordinate multiple repository operations)
- **Technology agnostic** - no HTTP or database concerns

Domain-specific validations example:

```go
func (s *TaskService) CreateTask(title, description, status string) (*models.Task, error) {
    task := &models.Task{
        Title:       strings.TrimSpace(title),     // Business rule: trim whitespace
        Description: strings.TrimSpace(description),
        Status:      models.TaskStatus(status),    // Business rule: valid status only
    }
    return s.taskRepo.CreateTask(task)
}
```

**Repository Layer (`internal/repository/`)**:

- Handles all database interactions
- Implements data persistence logic
- Manages SQL queries and transactions
- **Single responsibility** - only data access

### Interface-Driven Design

Every layer communicates through interfaces, creating clear contracts and enabling flexibility:

**TaskHandlerInterface** - HTTP-level operations:

```go
type TaskHandlerInterface interface {
    CreateTask(c *gin.Context)    // Handles HTTP request/response, JSON, status codes
    GetTask(c *gin.Context)       // Deals with Gin context and HTTP concerns
    GetAllTasks(c *gin.Context)
    UpdateTask(c *gin.Context)
    DeleteTask(c *gin.Context)
}

func NewTaskHandler(taskService service.TaskServiceInterface) TaskHandlerInterface {
    return &TaskHandler{taskService: taskService}  // Depends on interface, not concrete type
}
```

**TaskServiceInterface** - Business logic operations:

```go
type TaskServiceInterface interface {
    CreateTask(title, description, status string) (*models.Task, error)    // Pure business logic
    GetTaskByID(id int) (*models.Task, error)                             // No HTTP concerns
    GetAllTasks(page, limit int, status, sortBy, sortOrder string) (*models.PaginatedTasksResponse, error)
    UpdateTask(id int, title, description, status string) (*models.Task, error)
    DeleteTask(id int) error
}

func NewTaskService(taskRepo repository.TaskRepository) TaskServiceInterface {
    return &TaskService{taskRepo: taskRepo}  // Depends on repository interface
}
```

**TaskRepository** - Data access operations:

```go
type TaskRepository interface {
    CreateTask(task *models.Task) error                                                    // Pure SQL operations
    GetTaskByID(id int) (*models.Task, error)                                            // Database interactions only
    GetAllTasks(limit, page int, status, sortBy, sortOrder string) ([]models.Task, int, error)
    UpdateTask(task *models.Task) error
    DeleteTask(id int) error
}

func NewPostgresTaskRepository(db *sql.DB) TaskRepository {
    return &PostgresTaskRepository{db: db}  // Concrete implementation
}
```

**Benefits Achieved**:

**Dependency Injection**: Each layer receives interfaces, not concrete implementations

```go
// In main.go - easy to swap implementations
taskRepo := repository.NewPostgresTaskRepository(db)     // Could be NewMongoTaskRepository
taskService := service.NewTaskService(taskRepo)
taskHandler := handlers.NewTaskHandler(taskService)
```

**Testing Isolation**: Each layer can be tested with mocks

```go
// Service tests - no database needed
mockRepo := newMockTaskRepository()
service := NewTaskService(mockRepo)

// Handler tests - no business logic or database needed
mockService := new(MockTaskService)
handler := NewTaskHandler(mockService)
```

**Technology Flexibility**:

- Switch PostgreSQL → MySQL: Implement `TaskRepository` interface
- Add GraphQL: Create GraphQL handlers using same `TaskServiceInterface`
- Add gRPC: Create gRPC handlers calling existing service layer

### Middleware Architecture

Centralized middleware handles cross-cutting concerns:

**Validation Middleware**: Validates and transforms requests before reaching handlers

```go
// Clean handler - only business operations
func (h *TaskHandler) CreateTask(c *gin.Context) {
    req := middleware.GetCreateTaskRequest(c)  // Pre-validated
    task, err := h.taskService.CreateTask(req.Title, req.Description, req.Status)
    // ...
}
```

**Error Middleware**: Converts typed errors to appropriate HTTP responses

```go
// Typed business errors
return models.TaskNotFoundError{ID: id}  // → 404 JSON response
return models.ValidationError{...}       // → 400 JSON response
```

### Microservices Enablement

Current architecture demonstrates core microservice principles:

- **Self-Contained Service**: Task service with dedicated database, isolated from other services
- **Independent Scaling**: Can scale task service independently based on demand
- **Technology Isolation**: Each microservice can use different tech stack (Go + PostgreSQL for tasks, Node.js + Redis for notifications)
- **Fault Isolation**: Task service failures don't affect user service or notification service
- **Independent Deployment**: Can deploy task service updates without affecting other services

### Database Schema & Design

**Task Table Structure**:

```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,                    -- Auto-incrementing unique identifier
    title VARCHAR(255) NOT NULL,             -- Required task title (max 255 chars)
    description TEXT,                         -- Optional detailed description
    status VARCHAR(50) NOT NULL DEFAULT 'pending',  -- Task status with default
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Business rule enforcement at database level
    CONSTRAINT valid_status CHECK (status IN ('pending', 'in_progress', 'completed', 'closed'))
);

-- Performance optimization indexes
CREATE INDEX idx_tasks_status ON tasks(status);           -- Fast filtering by status
CREATE INDEX idx_tasks_created_at ON tasks(created_at DESC);  -- Fast sorting by date
```

**Design Decisions**:

- **Database Constraints**: Enforce valid status values at DB level (defense in depth)
- **Timestamps**: Automatic audit trail for created/updated tracking
- **Indexes**: Optimized for common query patterns (status filtering, date sorting)
- **Text vs VARCHAR**: TEXT for descriptions (unlimited), VARCHAR for constrained fields

## Inter-Service Communication Strategy

**Assignment Question**: How would two microservices (e.g., Task Service + User Service) communicate?

**Answer**: Using the same interface pattern with different implementation strategies:

```go
// Define service interface (same across all communication methods)
type UserService interface {
    GetUser(id int) (*User, error)
    ValidateUser(id int) (bool, error)
}

// Implementation options:

// 1. REST Communication
userService := rest.NewUserServiceClient("http://user-service:8080")

// 2. gRPC Communication
userService := grpc.NewUserServiceClient(grpcConn)

// 3. Message Queue Communication
userService := queue.NewUserServiceClient(rabbitMQConn)
```

**Communication Method Comparison**:

| Method             | Pros                                      | Cons                             | Use Case                             |
| ------------------ | ----------------------------------------- | -------------------------------- | ------------------------------------ |
| **REST**           | Simple, HTTP-based, wide adoption         | Higher latency, verbose JSON     | CRUD operations, external APIs       |
| **gRPC**           | Fast, type-safe, bi-directional streaming | Learning curve, HTTP/2 required  | Internal service communication       |
| **Message Queues** | Async, decoupled, fault-tolerant          | Complexity, eventual consistency | Event-driven, high-volume processing |

**Implementation in Task Service**:

```go
// Task service calling User service
func (s *TaskService) AssignTask(taskID, userID int) error {
    // Validate user exists via interface (implementation abstracted)
    valid, err := s.userService.ValidateUser(userID)
    if !valid {
        return errors.New("invalid user")
    }

    // Update task assignment
    return s.taskRepo.AssignTask(taskID, userID)
}
```

**Benefits of Interface Approach**:

- **Swap communication protocols** without changing business logic
- **Mock external services** for testing
- **Gradual migration** from REST to gRPC as services mature

## Horizontal Scaling Demonstration

```bash
# Scale to multiple instances
make dev-scale INSTANCES=5

# Test load distribution
make test-load
# Output shows requests distributed across different container instances
```

### Architecture Components

**Nginx Load Balancer**:

```nginx
# Dynamic service discovery using Docker DNS
resolver 127.0.0.11 valid=10s;
set $upstream app:8080;  # Docker service name
proxy_pass http://$upstream;
```

**Benefits**:

- **Stateless Services**: Each app instance is identical and independent
- **Dynamic Scaling**: Add/remove instances without configuration changes
- **Fault Tolerance**: Failed instances automatically removed from load balancing

### Production-Grade Improvements

**Container Orchestration**:

- **Kubernetes/Docker Swarm**: Advanced container management and auto-scaling
- **Service Discovery**: Automatic load balancer configuration as services scale

**Load Balancing Optimization**:

- **Algorithm Selection**: Round Robin, Least Connections, or IP Hash depending on requirements
- **Multiple Load Balancers**: Scale load balancers themselves for high availability (HA Proxy = High Availability Proxy)

**Database Scaling**:

- **Read Replicas**: Distribute read operations across multiple database instances
- **Write Replicas**: Master-slave configuration with consistency requirements consideration
- **Connection Pooling**: Optimize database connection management per service instance

**Advanced Resilience**:

- **Health Checks**: Automatic removal of unhealthy instances
- **Circuit Breakers**: Prevent cascade failures across services

## Testing Architecture

**Layer-Specific Testing Strategy**:

```go
// Unit Tests (Service Layer) - Fast, isolated business logic testing
func TestTaskService_CreateTask(t *testing.T) {
    mockRepo := newMockTaskRepository()           // No database dependency
    service := NewTaskService(mockRepo)           // Pure business logic testing
    task, err := service.CreateTask("Test", "", "pending")
    // Verify business rules, validations, transformations
}

// Integration Tests (Repository Layer) - Real database operations
func TestPostgresTaskRepository_CreateTask(t *testing.T) {
    db := SetupTestDB(t)                         // Real PostgreSQL test database
    repo := NewPostgresTaskRepository(db)        // Actual SQL execution
    err := repo.CreateTask(task)                 // Tests real database interactions
    // Verify data persistence, constraints, indexing
}

// Handler Tests (HTTP Layer) - Protocol testing with mocks
func TestCreateTask_Success(t *testing.T) {
    mockService := new(MockTaskService)          // No business logic dependency
    handler := NewTaskHandler(mockService)       // Pure HTTP protocol testing
    // Test JSON parsing, status codes, response formatting
}
```

**Test Isolation Benefits**:

- **Service Tests**: Business logic validation without database overhead
- **Repository Tests**: SQL query validation with real PostgreSQL and ephemeral test database
- **Handler Tests**: HTTP protocol validation without business complexity
- **Load Tests**: Multi-instance distribution verification via `make test-load`

**Testing Commands & Coverage**:

```bash
# Fast tests (no database)
make test-unit        # Service + Handler tests (~10 seconds)

# Integration tests (with database)
make test-repository  # Real PostgreSQL operations (~45 seconds)

# Complete test suite
make test-all         # Unit + Integration + Handler tests

# Load balancing verification
make test-load        # Verify request distribution across instances
```

**Benefits**: Each test type validates specific concerns without interference, enabling confident refactoring and rapid development cycles.

## Assignment Scope Limitations

**Intentionally not implemented for assignment focus**:

- **Database Migrations**: Using single migration file executed at container startup via Docker's initdb.d; production would need versioned migration tools like golang-migrate for schema evolution
- **Logging**: No structured logging middleware or in-code logging implemented
- **Advanced Load Balancing**: Currently using Docker's internal DNS with 0-second cache for demonstration; production needs Kubernetes/Docker Swarm
- **Authentication/Authorization**: Assumed all requests are valid (out of scope)
- **Rate Limiting**: No DDoS protection implemented
- **CORS Configuration**: No frontend integration setup

## Summary

This Task Management API demonstrates microservices architecture principles through:

- Clean separation of concerns enabling independent scaling
- Interface-based design supporting multiple communication protocols
- Comprehensive testing strategy across all architectural layers
- Production-ready containerization with horizontal scaling capabilities

The implementation exceeds assignment requirements by providing a foundation for microservices evolution while maintaining simplicity and clarity.
