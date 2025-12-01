package internal

import (
	"strings"
	"time"
)

type Todo struct {
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description,omitempty" db:"description"`
	ID          int64      `json:"id" db:"id"`
	Completed   bool       `json:"completed" db:"completed"`
}

type CreateTodoInput struct {
	DueDate     *time.Time `json:"due_date"`
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
}

func (c *CreateTodoInput) Validate() error {
	c.Title = strings.TrimSpace(c.Title)
	if c.Title == "" {
		return ErrTitleRequired
	}
	if len(c.Title) > 255 {
		return ErrTitleMaxLength
	}
	return nil
}

type UpdateTodoInput struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	Completed   *bool      `json:"completed"`
	ID          int64      `json:"id" binding:"required"`
}

func (u *UpdateTodoInput) Validate() error {
	if u.ID <= 0 {
		return ErrInvalidID
	}
	if u.Title != nil {
		*u.Title = strings.TrimSpace(*u.Title)
		if *u.Title == "" {
			return ErrTitleEmpty
		}
		if len(*u.Title) > 255 {
			return ErrTitleMaxLength
		}
	}
	return nil
}
