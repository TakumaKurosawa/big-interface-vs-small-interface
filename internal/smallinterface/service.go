package smallinterface

import (
	"context"
	"errors"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

// UserService はユーザー関連の操作を提供するサービスです
type UserService struct {
	userStore UserStore // ユーザー関連の小さなインターフェースのみを使用
}

// NewUserService は新しいUserServiceを作成します
func NewUserService(userStore UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

// GetUser はユーザーを取得します
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.userStore.GetUser(ctx, id)
}

// CreateUser は新しいユーザーを作成します
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userStore.CreateUser(ctx, user)
}

// TodoService はTodo関連の操作を提供するサービスです
type TodoService struct {
	todoStore TodoStore // Todo関連の小さなインターフェースを使用
	userStore UserStore // 必要な場合はユーザー関連の小さなインターフェースも使用
}

// NewTodoService は新しいTodoServiceを作成します
func NewTodoService(todoStore TodoStore, userStore UserStore) *TodoService {
	return &TodoService{
		todoStore: todoStore,
		userStore: userStore,
	}
}

// GetUserTodos はユーザーのTodoリストを取得します
func (s *TodoService) GetUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error) {
	// ユーザーの存在を確認する必要がある場合、
	// 特定のインターフェースを通じてユーザー情報にアクセスする
	_, err := s.userStore.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.todoStore.ListUserTodos(ctx, userID)
}

// CreateTodo は新しいTodoを作成します
func (s *TodoService) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	// ユーザー存在確認
	_, err := s.userStore.GetUser(ctx, todo.UserID)
	if err != nil {
		return errors.New("cannot create todo for non-existent user")
	}

	return s.todoStore.CreateTodo(ctx, todo)
}

// CompleteTodo はTodoを完了状態にします
func (s *TodoService) CompleteTodo(ctx context.Context, id string) error {
	return s.todoStore.MarkTodoComplete(ctx, id)
}
