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

	// 共通のデータストアを作成
	store := inmemory.NewStore()

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
	bigUserService := bigservice.NewUserService(store)
	bigTodoService := bigservice.NewTodoService(store)

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
	smallUserService := smallservice.NewUserService(store)
	smallTodoService := smallservice.NewTodoService(store, store)

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
