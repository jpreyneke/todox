package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTodoInput_Validate(t *testing.T) {
	tests := []struct {
		input   CreateTodoInput
		name    string
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   CreateTodoInput{Title: "Test Todo"},
			wantErr: false,
		},
		{
			name:    "empty title",
			input:   CreateTodoInput{Title: ""},
			wantErr: true,
		},
		{
			name:    "whitespace title",
			input:   CreateTodoInput{Title: "   "},
			wantErr: true,
		},
		{
			name:    "title too long",
			input:   CreateTodoInput{Title: string(make([]byte, 300))},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUpdateTodoInput_Validate(t *testing.T) {
	tests := []struct {
		input   UpdateTodoInput
		name    string
		wantErr bool
	}{
		{
			name:    "valid input",
			input:   UpdateTodoInput{ID: 1, Title: strPtr("Updated")},
			wantErr: false,
		},
		{
			name:    "zero id",
			input:   UpdateTodoInput{ID: 0},
			wantErr: true,
		},
		{
			name:    "negative id",
			input:   UpdateTodoInput{ID: -1},
			wantErr: true,
		},
		{
			name:    "empty title",
			input:   UpdateTodoInput{ID: 1, Title: strPtr("")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_BulkCreate_EmptyList(t *testing.T) {
	service := &Service{repo: nil}

	_, err := service.BulkCreate(context.Background(), []CreateTodoInput{})

	assert.ErrorIs(t, err, ErrEmptyList)
}

func TestService_BulkCreate_DuplicateTitlesInRequest(t *testing.T) {
	service := &Service{repo: nil}

	inputs := []CreateTodoInput{
		{Title: "Same Title"},
		{Title: "Same Title"},
	}

	_, err := service.BulkCreate(context.Background(), inputs)

	assert.ErrorIs(t, err, ErrDuplicateInRequest)
}

func TestService_List_LimitTooHigh(t *testing.T) {
	service := &Service{repo: nil}

	_, _, err := service.List(context.Background(), 1, 200)

	assert.ErrorIs(t, err, ErrLimitExceeded)
}

func strPtr(s string) *string {
	return &s
}
