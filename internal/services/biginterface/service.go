package biginterface

import (
	"context"
	"errors"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/biginterface"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

// UserService is a service that provides user-related operations
type UserService struct {
	store biginterface.DataStore // Using the big interface
}

// NewUserService creates a new UserService
func NewUserService(store biginterface.DataStore) *UserService {
	return &UserService{
		store: store,
	}
}

// GetUser retrieves a user
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.store.GetUser(ctx, id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.store.CreateUser(ctx, user)
}

// TodoService is a service that provides Todo-related operations
type TodoService struct {
	store biginterface.DataStore // Using the same big interface
}

// NewTodoService creates a new TodoService
func NewTodoService(store biginterface.DataStore) *TodoService {
	return &TodoService{
		store: store,
	}
}

// GetUserTodos retrieves a user's Todo list
func (s *TodoService) GetUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error) {
	// When we need to check if a user exists,
	// we can access user information through the big interface here as well
	_, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.store.ListUserTodos(ctx, userID)
}

// CreateTodo creates a new Todo
func (s *TodoService) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	// Check if user exists
	_, err := s.store.GetUser(ctx, todo.UserID)
	if err != nil {
		return errors.New("cannot create todo for non-existent user")
	}

	return s.store.CreateTodo(ctx, todo)
}

// CompleteTodo marks a Todo as complete
func (s *TodoService) CompleteTodo(ctx context.Context, id string) error {
	return s.store.MarkTodoComplete(ctx, id)
}
