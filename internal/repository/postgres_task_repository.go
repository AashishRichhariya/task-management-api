package repository

import (
	"database/sql"
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

// CreateTask inserts a new task into database
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

// GetAllTasks retrieves all tasks
func (r *PostgresTaskRepository) GetAllTasks() ([]models.Task, error) {
	query := `
		SELECT id, title, description, status, created_at, updated_at
		FROM tasks 
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	
	return tasks, nil
}

// UpdateTask updates an existing task
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

// DeleteTask removes a task by ID
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