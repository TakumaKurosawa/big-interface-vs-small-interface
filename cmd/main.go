package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/infra/inmemory"
	bigservice "github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/services/biginterface"
	smallservice "github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/services/smallinterface"
)

// InMemoryStoreの実装をすべて削除

func main() {
	ctx := context.Background()

	// Create a common data store
	store := inmemory.NewStore()

	// Create sample data
	user := &domain.User{
		ID:        "user1",
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	todo := &domain.Todo{
		ID:          "todo1",
		UserID:      "user1",
		Title:       "Interface Design Comparison",
		Description: "Verify the differences between large and small interfaces",
		Completed:   false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save data
	if err := store.CreateUser(ctx, user); err != nil {
		log.Fatalf("User creation error: %v", err)
	}

	if err := store.CreateTodo(ctx, todo); err != nil {
		log.Fatalf("Todo creation error: %v", err)
	}

	fmt.Println("===== Big Interface Approach =====")
	// Big interface approach
	bigUserService := bigservice.NewUserService(store)
	bigTodoService := bigservice.NewTodoService(store)

	// Get user information
	fetchedUser, err := bigUserService.GetUser(ctx, "user1")
	if err != nil {
		log.Fatalf("User retrieval error: %v", err)
	}
	fmt.Printf("User: %s (%s)\n", fetchedUser.Name, fetchedUser.Email)

	// Get Todo list
	todos, err := bigTodoService.GetUserTodos(ctx, "user1")
	if err != nil {
		log.Fatalf("Todo retrieval error: %v", err)
	}

	fmt.Println("Todo List:")
	for _, t := range todos {
		status := "Incomplete"
		if t.Completed {
			status = "Complete"
		}
		fmt.Printf("- %s: %s (%s)\n", t.Title, t.Description, status)
	}

	fmt.Println("\n===== Small Interface Approach =====")
	// Small interface approach
	smallUserService := smallservice.NewUserService(store)
	smallTodoService := smallservice.NewTodoService(store, store)

	// Get user information
	fetchedUser, err = smallUserService.GetUser(ctx, "user1")
	if err != nil {
		log.Fatalf("User retrieval error: %v", err)
	}
	fmt.Printf("User: %s (%s)\n", fetchedUser.Name, fetchedUser.Email)

	// Get Todo list
	todos, err = smallTodoService.GetUserTodos(ctx, "user1")
	if err != nil {
		log.Fatalf("Todo retrieval error: %v", err)
	}

	fmt.Println("Todo List:")
	for _, t := range todos {
		status := "Incomplete"
		if t.Completed {
			status = "Complete"
		}
		fmt.Printf("- %s: %s (%s)\n", t.Title, t.Description, status)
	}

	// Mark Todo as complete
	if err := smallTodoService.CompleteTodo(ctx, "todo1"); err != nil {
		log.Fatalf("Todo update error: %v", err)
	}

	fmt.Println("\n===== After Todo Completion =====")
	todos, _ = smallTodoService.GetUserTodos(ctx, "user1")
	for _, t := range todos {
		status := "Incomplete"
		if t.Completed {
			status = "Complete"
		}
		fmt.Printf("- %s: %s (%s)\n", t.Title, t.Description, status)
	}
}
