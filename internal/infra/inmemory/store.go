package inmemory

import (
	"context"
	"fmt"
	"time"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/biginterface"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/smallinterface"
)

// Store は両方のインターフェースを満たすための実装
type Store struct {
	users map[string]*domain.User
	todos map[string]*domain.Todo
}

var _ biginterface.DataStore = (*Store)(nil)
var _ smallinterface.UserStore = (*Store)(nil)
var _ smallinterface.TodoStore = (*Store)(nil)

// NewStore は新しいインメモリストアを生成します
func NewStore() *Store {
	return &Store{
		users: make(map[string]*domain.User),
		todos: make(map[string]*domain.Todo),
	}
}

// ユーザー関連の操作
func (s *Store) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return user, nil
}

func (s *Store) ListUsers(ctx context.Context) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *Store) CreateUser(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	s.users[user.ID] = user
	return nil
}

func (s *Store) UpdateUser(ctx context.Context, user *domain.User) error {
	if _, ok := s.users[user.ID]; !ok {
		return fmt.Errorf("user not found: %s", user.ID)
	}
	s.users[user.ID] = user
	return nil
}

func (s *Store) DeleteUser(ctx context.Context, id string) error {
	if _, ok := s.users[id]; !ok {
		return fmt.Errorf("user not found: %s", id)
	}
	delete(s.users, id)
	return nil
}

// Todo関連の操作
func (s *Store) GetTodo(ctx context.Context, id string) (*domain.Todo, error) {
	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	return todo, nil
}

func (s *Store) ListTodos(ctx context.Context) ([]*domain.Todo, error) {
	todos := make([]*domain.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s *Store) ListUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error) {
	todos := make([]*domain.Todo, 0)
	for _, todo := range s.todos {
		if todo.UserID == userID {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (s *Store) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	if todo.ID == "" {
		return fmt.Errorf("todo ID cannot be empty")
	}
	s.todos[todo.ID] = todo
	return nil
}

func (s *Store) UpdateTodo(ctx context.Context, todo *domain.Todo) error {
	if _, ok := s.todos[todo.ID]; !ok {
		return fmt.Errorf("todo not found: %s", todo.ID)
	}
	s.todos[todo.ID] = todo
	return nil
}

func (s *Store) DeleteTodo(ctx context.Context, id string) error {
	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("todo not found: %s", id)
	}
	delete(s.todos, id)
	return nil
}

func (s *Store) MarkTodoComplete(ctx context.Context, id string) error {
	todo, ok := s.todos[id]
	if !ok {
		return fmt.Errorf("todo not found: %s", id)
	}
	todo.Completed = true
	todo.UpdatedAt = time.Now()
	return nil
}
