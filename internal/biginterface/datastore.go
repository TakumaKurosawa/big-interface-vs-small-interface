package biginterface

import (
	"context"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

//go:generate mockgen -destination=./mocks/mock_datastore.go -package=mocks github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/biginterface DataStore

// DataStore defines all data operations as a single large interface
// This is an example of a low cohesion approach
type DataStore interface {
	// User-related operations
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error

	// Todo-related operations
	GetTodo(ctx context.Context, id string) (*domain.Todo, error)
	ListTodos(ctx context.Context) ([]*domain.Todo, error)
	ListUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error)
	CreateTodo(ctx context.Context, todo *domain.Todo) error
	UpdateTodo(ctx context.Context, todo *domain.Todo) error
	DeleteTodo(ctx context.Context, id string) error
	MarkTodoComplete(ctx context.Context, id string) error
}
