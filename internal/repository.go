package internal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*Todo, error) {
	var todo Todo
	err := r.db.GetContext(ctx, &todo, "SELECT * FROM todos WHERE id = ?", id)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	return &todo, err
}

func (r *Repository) List(ctx context.Context, page, limit int) ([]Todo, int64, error) {
	offset := (page - 1) * limit

	var todos []Todo
	err := r.db.SelectContext(ctx, &todos,
		"SELECT * FROM todos ORDER BY created_at DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if todos == nil {
		todos = []Todo{}
	}

	var total int64
	err = r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM todos")
	return todos, total, err
}

func (r *Repository) BulkCreate(ctx context.Context, todos []*Todo) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rerr := tx.Rollback(); rerr != nil {
			log.Printf("Rollback failed: %v", rerr)
		}
	}()

	query := `INSERT INTO todos (title, description, due_date, completed, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	for _, todo := range todos {
		result, err := tx.ExecContext(ctx, query,
			todo.Title, todo.Description, todo.DueDate, todo.Completed, todo.CreatedAt, todo.UpdatedAt)
		if err != nil {
			if isDuplicateError(err) {
				return ErrDuplicateTitle
			}
			return err
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get last insert ID: %w", err)
		}

		todo.ID = id
	}

	return tx.Commit()
}

func (r *Repository) BulkUpdate(ctx context.Context, todos []*Todo) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if rerr := tx.Rollback(); rerr != nil {
			log.Printf("Rollback failed: %v", rerr)
		}
	}()

	query := `UPDATE todos SET title=?, description=?, due_date=?, completed=?, updated_at=? WHERE id=?`

	for _, todo := range todos {
		result, err := tx.ExecContext(ctx, query,
			todo.Title, todo.Description, todo.DueDate, todo.Completed, todo.UpdatedAt, todo.ID)
		if err != nil {
			if isDuplicateError(err) {
				return ErrDuplicateTitle
			}
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rows == 0 {
			return ErrNotFound
		}
	}

	return tx.Commit()
}

func isDuplicateError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}
