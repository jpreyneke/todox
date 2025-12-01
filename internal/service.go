package internal

import (
	"context"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) BulkCreate(ctx context.Context, inputs []CreateTodoInput) ([]*Todo, error) {
	if len(inputs) == 0 {
		return nil, ErrEmptyList
	}

	seen := make(map[string]bool)
	now := time.Now()
	todos := make([]*Todo, 0, len(inputs))

	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, err
		}
		if seen[input.Title] {
			return nil, ErrDuplicateInRequest
		}
		seen[input.Title] = true

		todos = append(todos, &Todo{
			Title:       input.Title,
			Description: input.Description,
			DueDate:     input.DueDate,
			Completed:   false,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}

	if err := s.repo.BulkCreate(ctx, todos); err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *Service) List(ctx context.Context, page, limit int) ([]Todo, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		return nil, 0, ErrLimitExceeded
	}
	return s.repo.List(ctx, page, limit)
}

func (s *Service) BulkUpdate(ctx context.Context, inputs []UpdateTodoInput) ([]*Todo, error) {
	if len(inputs) == 0 {
		return nil, ErrEmptyList
	}

	seenIDs := make(map[int64]bool)
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, err
		}
		if seenIDs[input.ID] {
			return nil, ErrDuplicateInRequest
		}
		seenIDs[input.ID] = true
	}

	todos := make([]*Todo, 0, len(inputs))
	now := time.Now()

	for _, input := range inputs {
		todo, err := s.repo.GetByID(ctx, input.ID)
		if err != nil {
			return nil, err
		}

		if input.Title != nil {
			todo.Title = *input.Title
		}
		if input.Description != nil {
			todo.Description = *input.Description
		}
		if input.DueDate != nil {
			todo.DueDate = input.DueDate
		}
		if input.Completed != nil {
			todo.Completed = *input.Completed
		}
		todo.UpdatedAt = now

		todos = append(todos, todo)
	}

	if err := s.repo.BulkUpdate(ctx, todos); err != nil {
		return nil, err
	}
	return todos, nil
}
