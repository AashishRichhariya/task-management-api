package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/AashishRichhariya/task-management-api/internal/models"
	_ "github.com/lib/pq"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

// Constructor - creates new repository instance
func NewPostgresTaskRepository(db *sql.DB) TaskRepository {
	return &PostgresTaskRepository{db: db}
}

// Inserts a new task into database
func (r *PostgresTaskRepository) CreateTask(task *models.Task) error {
	query := `
		INSERT INTO tasks (title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	
	err := r.db.QueryRow(query, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).Scan(&task.ID)
	if err != nil {
		return err
	}
	
	return nil
}

// GetTaskByID retrieves a single task by ID
func (r *PostgresTaskRepository) GetTaskByID(id int) (*models.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks 
		WHERE id = $1`
	
	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Task not found
		}
		return nil, err
	}
	
	return task, nil
}

// Retrieves all tasks
func (r *PostgresTaskRepository) GetAllTasks(limit, page int, status, sortBy, sortOrder string) ([]models.Task, int, error) {
	// Convert page to offset for database
	offset := (page - 1) * limit

	// Build the WHERE clause for filtering
	whereClause := ""
	args := []any{limit, offset}  // $1 = limit, $2 = offset
	argIndex := 3

	if status != "" {
		whereClause = "WHERE status = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Build the complete query with COUNT
	query := fmt.Sprintf(`
		SELECT 
			id, title, description, status, created_at, updated_at,
			COUNT(*) OVER() as total_count
		FROM tasks 
		%s
		ORDER BY %s %s
		LIMIT $1 OFFSET $2
	`, whereClause, sortBy, sortOrder)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	tasks := []models.Task{}
	var totalCount int

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
			&totalCount,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	// Handle case where no rows returned (high page number)
	if len(tasks) == 0 {
			// Do a separate count query to get total
			countQuery := "SELECT COUNT(*) FROM tasks"
			countArgs := []any{}
			
			if status != "" {
					countQuery += " WHERE status = $1"
					countArgs = append(countArgs, status)
			}
			
			err = r.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
			if err != nil {
					return nil, 0, fmt.Errorf("failed to get count: %w", err)
			}
	}

	return tasks, totalCount, nil
}

// Updates an existing task
func (r *PostgresTaskRepository) UpdateTask(task *models.Task) error {
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5`
	
	task.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(query, task.Title, task.Description, task.Status, task.UpdatedAt, task.ID)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows // Task not found
	}
	
	return nil
}

// Removes a task by ID
func (r *PostgresTaskRepository) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows // Task not found
	}
	
	return nil
}