package smallinterface

import (
	"context"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

//go:generate mockgen -destination=../mocks/mock_todostore.go -package=mocks github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/smallinterface TodoStore

// TodoStore はTodo関連の操作のみを定義した小さなインターフェースです
// これは凝集度が高いアプローチの例です
type TodoStore interface {
	GetTodo(ctx context.Context, id string) (*domain.Todo, error)
	ListTodos(ctx context.Context) ([]*domain.Todo, error)
	ListUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error)
	CreateTodo(ctx context.Context, todo *domain.Todo) error
	UpdateTodo(ctx context.Context, todo *domain.Todo) error
	DeleteTodo(ctx context.Context, id string) error
	MarkTodoComplete(ctx context.Context, id string) error
}
