package biginterface

import (
	"context"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

//go:generate mockgen -destination=../mocks/mock_datastore.go -package=mocks github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/biginterface DataStore

// DataStore は全てのデータ操作を一つの大きなインターフェースとして定義します
// これは凝集度が低いアプローチの例です
type DataStore interface {
	// ユーザー関連の操作
	GetUser(ctx context.Context, id string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id string) error

	// Todo関連の操作
	GetTodo(ctx context.Context, id string) (*domain.Todo, error)
	ListTodos(ctx context.Context) ([]*domain.Todo, error)
	ListUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error)
	CreateTodo(ctx context.Context, todo *domain.Todo) error
	UpdateTodo(ctx context.Context, todo *domain.Todo) error
	DeleteTodo(ctx context.Context, id string) error
	MarkTodoComplete(ctx context.Context, id string) error
}
