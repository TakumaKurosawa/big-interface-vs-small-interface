package biginterface

import (
	"context"
	"errors"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/biginterface"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

// UserService はユーザー関連の操作を提供するサービスです
type UserService struct {
	store biginterface.DataStore // 大きなインターフェースを使用
}

// NewUserService は新しいUserServiceを作成します
func NewUserService(store biginterface.DataStore) *UserService {
	return &UserService{
		store: store,
	}
}

// GetUser はユーザーを取得します
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.store.GetUser(ctx, id)
}

// CreateUser は新しいユーザーを作成します
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.store.CreateUser(ctx, user)
}

// TodoService はTodo関連の操作を提供するサービスです
type TodoService struct {
	store biginterface.DataStore // 同じ大きなインターフェースを使用
}

// NewTodoService は新しいTodoServiceを作成します
func NewTodoService(store biginterface.DataStore) *TodoService {
	return &TodoService{
		store: store,
	}
}

// GetUserTodos はユーザーのTodoリストを取得します
func (s *TodoService) GetUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error) {
	// ユーザーの存在を確認する必要がある場合、
	// ここでも大きなインターフェースを通じてユーザー情報にアクセスできる
	_, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.store.ListUserTodos(ctx, userID)
}

// CreateTodo は新しいTodoを作成します
func (s *TodoService) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	// ユーザー存在確認
	_, err := s.store.GetUser(ctx, todo.UserID)
	if err != nil {
		return errors.New("cannot create todo for non-existent user")
	}

	return s.store.CreateTodo(ctx, todo)
}

// CompleteTodo はTodoを完了状態にします
func (s *TodoService) CompleteTodo(ctx context.Context, id string) error {
	return s.store.MarkTodoComplete(ctx, id)
}
