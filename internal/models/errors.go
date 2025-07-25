package models

import "fmt"

type TaskNotFoundError struct {
	ID int
}

func (e TaskNotFoundError) Error() string {
	return fmt.Sprintf("task with id %d not found", e.ID)
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type BusinessError struct {
	Message string
}

func (e BusinessError) Error() string {
	return e.Message
}