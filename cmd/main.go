package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/biginterface"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/smallinterface"
)

// InMemoryStore は両方のインターフェースを満たすための実装
type InMemoryStore struct {
	users map[string]*domain.User
	todos map[string]*domain.Todo
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		users: make(map[string]*domain.User),
		todos: make(map[string]*domain.Todo),
	}
}

// ユーザー関連の操作
func (s *InMemoryStore) GetUser(ctx context.Context, id string) (*domain.User, error) {
	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	return user, nil
}

func (s *InMemoryStore) ListUsers(ctx context.Context) ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *InMemoryStore) CreateUser(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryStore) UpdateUser(ctx context.Context, user *domain.User) error {
	if _, ok := s.users[user.ID]; !ok {
		return fmt.Errorf("user not found: %s", user.ID)
	}
	s.users[user.ID] = user
	return nil
}

func (s *InMemoryStore) DeleteUser(ctx context.Context, id string) error {
	if _, ok := s.users[id]; !ok {
		return fmt.Errorf("user not found: %s", id)
	}
	delete(s.users, id)
	return nil
}

// Todo関連の操作
func (s *InMemoryStore) GetTodo(ctx context.Context, id string) (*domain.Todo, error) {
	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	return todo, nil
}

func (s *InMemoryStore) ListTodos(ctx context.Context) ([]*domain.Todo, error) {
	todos := make([]*domain.Todo, 0, len(s.todos))
	for _, todo := range s.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (s *InMemoryStore) ListUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error) {
	todos := make([]*domain.Todo, 0)
	for _, todo := range s.todos {
		if todo.UserID == userID {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

func (s *InMemoryStore) CreateTodo(ctx context.Context, todo *domain.Todo) error {
	if todo.ID == "" {
		return fmt.Errorf("todo ID cannot be empty")
	}
	s.todos[todo.ID] = todo
	return nil
}

func (s *InMemoryStore) UpdateTodo(ctx context.Context, todo *domain.Todo) error {
	if _, ok := s.todos[todo.ID]; !ok {
		return fmt.Errorf("todo not found: %s", todo.ID)
	}
	s.todos[todo.ID] = todo
	return nil
}

func (s *InMemoryStore) DeleteTodo(ctx context.Context, id string) error {
	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("todo not found: %s", id)
	}
	delete(s.todos, id)
	return nil
}

func (s *InMemoryStore) MarkTodoComplete(ctx context.Context, id string) error {
	todo, ok := s.todos[id]
	if !ok {
		return fmt.Errorf("todo not found: %s", id)
	}
	todo.Completed = true
	todo.UpdatedAt = time.Now()
	return nil
}

func main() {
	ctx := context.Background()

	// 共通のデータストアを作成
	store := NewInMemoryStore()

	// サンプルデータを作成
	user := &domain.User{
		ID:        "user1",
		Name:      "山田 太郎",
		Email:     "yamada@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	todo := &domain.Todo{
		ID:          "todo1",
		UserID:      "user1",
		Title:       "インターフェース設計の比較",
		Description: "大きなインターフェースと小さなインターフェースの違いを検証する",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// データを保存
	if err := store.CreateUser(ctx, user); err != nil {
		log.Fatalf("ユーザー作成エラー: %v", err)
	}

	if err := store.CreateTodo(ctx, todo); err != nil {
		log.Fatalf("Todo作成エラー: %v", err)
	}

	fmt.Println("===== 大きなインターフェースアプローチ =====")
	// 大きなインターフェースアプローチ
	bigUserService := biginterface.NewUserService(store)
	bigTodoService := biginterface.NewTodoService(store)

	// ユーザー情報を取得
	fetchedUser, err := bigUserService.GetUser(ctx, "user1")
	if err != nil {
		log.Fatalf("ユーザー取得エラー: %v", err)
	}
	fmt.Printf("ユーザー: %s (%s)\n", fetchedUser.Name, fetchedUser.Email)

	// Todoリストを取得
	todos, err := bigTodoService.GetUserTodos(ctx, "user1")
	if err != nil {
		log.Fatalf("Todo取得エラー: %v", err)
	}

	fmt.Println("Todoリスト:")
	for _, t := range todos {
		status := "未完了"
		if t.Completed {
			status = "完了"
		}
		fmt.Printf("- %s: %s (%s)\n", t.Title, t.Description, status)
	}

	fmt.Println("\n===== 小さなインターフェースアプローチ =====")
	// 小さなインターフェースアプローチ
	smallUserService := smallinterface.NewUserService(store)
	smallTodoService := smallinterface.NewTodoService(store, store)

	// ユーザー情報を取得
	fetchedUser, err = smallUserService.GetUser(ctx, "user1")
	if err != nil {
		log.Fatalf("ユーザー取得エラー: %v", err)
	}
	fmt.Printf("ユーザー: %s (%s)\n", fetchedUser.Name, fetchedUser.Email)

	// Todoリストを取得
	todos, err = smallTodoService.GetUserTodos(ctx, "user1")
	if err != nil {
		log.Fatalf("Todo取得エラー: %v", err)
	}

	fmt.Println("Todoリスト:")
	for _, t := range todos {
		status := "未完了"
		if t.Completed {
			status = "完了"
		}
		fmt.Printf("- %s: %s (%s)\n", t.Title, t.Description, status)
	}

	// Todoを完了にする
	if err := smallTodoService.CompleteTodo(ctx, "todo1"); err != nil {
		log.Fatalf("Todo更新エラー: %v", err)
	}

	fmt.Println("\n===== Todo完了後 =====")
	todos, _ = smallTodoService.GetUserTodos(ctx, "user1")
	for _, t := range todos {
		status := "未完了"
		if t.Completed {
			status = "完了"
		}
		fmt.Printf("- %s: %s (%s)\n", t.Title, t.Description, status)
	}
}
