package smallinterface

import (
	"context"
	"errors"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/smallinterface"
)

// UserService is a service that provides user-related operations
type UserService struct {
	userStore smallinterface.UserStore // Using only the small user interface
}

// NewUserService creates a new UserService
func NewUserService(userStore smallinterface.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

// GetUser retrieves a user
func (s *UserService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.userStore.GetUser(ctx, id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userStore.CreateUser(ctx, user)
}

// TodoService is a service that provides Todo-related operations
type TodoService struct {
	todoStore smallinterface.TodoStore // Using the small Todo interface
	userStore smallinterface.UserStore // Also using the small user interface when needed
}

// NewTodoService creates a new TodoService
func NewTodoService(todoStore smallinterface.TodoStore, userStore smallinterface.UserStore) *TodoService {
	return &TodoService{
		todoStore: todoStore,
		userStore: userStore,
	}
}

// GetUserTodos retrieves a user's Todo list
func (s *TodoService) GetUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error) {
	// When we need to check if a user exists,
	// we access user information through the specific interface
	_, err := s.userStore.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.todoStore.ListUserTodos(ctx, userID)
}

// CreateTodo creates a new Todo
func (s *TodoService) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	// Check if user exists
	_, err := s.userStore.GetUser(ctx, todo.UserID)
	if err != nil {
		return errors.New("cannot create todo for non-existent user")
	}

	return s.todoStore.CreateTodo(ctx, todo)
}

// CompleteTodo marks a Todo as complete
func (s *TodoService) CompleteTodo(ctx context.Context, id string) error {
	return s.todoStore.MarkTodoComplete(ctx, id)
}
