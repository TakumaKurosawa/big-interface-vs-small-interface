package biginterface

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/mocks"
)

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 大きなインターフェースをモック化
	mockStore := mocks.NewMockDataStore(ctrl)

	// モック化するメソッドは、インターフェースで定義されている全てのメソッドではなく、
	// このテストで使用する1つのメソッドのみですが、大きなインターフェース全体をモック化する必要があります
	mockUser := &domain.User{
		ID:        "user1",
		Name:      "テスト ユーザー",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockStore.EXPECT().
		GetUser(gomock.Any(), "user1").
		Return(mockUser, nil)

	service := NewUserService(mockStore)

	ctx := context.Background()
	user, err := service.GetUser(ctx, "user1")

	require.NoError(t, err)
	assert.Equal(t, mockUser, user)
}

func TestTodoService_GetUserTodos(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 大きなインターフェースをモック化
	mockStore := mocks.NewMockDataStore(ctrl)

	mockUser := &domain.User{
		ID:        "user1",
		Name:      "テスト ユーザー",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockTodos := []*domain.Todo{
		{
			ID:          "todo1",
			UserID:      "user1",
			Title:       "テストTodo",
			Description: "これはテスト用のTodoです",
			Completed:   false,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// このテストでは2つのメソッドのみを使用しますが、
	// 大きなインターフェースの場合、モックする必要がある全ての関数を考慮する必要があります
	mockStore.EXPECT().
		GetUser(gomock.Any(), "user1").
		Return(mockUser, nil)

	mockStore.EXPECT().
		ListUserTodos(gomock.Any(), "user1").
		Return(mockTodos, nil)

	service := NewTodoService(mockStore)

	ctx := context.Background()
	todos, err := service.GetUserTodos(ctx, "user1")

	require.NoError(t, err)
	assert.Equal(t, mockTodos, todos)
}

func TestTodoService_GetUserTodos_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 大きなインターフェースをモック化
	mockStore := mocks.NewMockDataStore(ctrl)

	// ユーザーが見つからない場合のテスト
	mockStore.EXPECT().
		GetUser(gomock.Any(), "nonexistent").
		Return(nil, errors.New("user not found"))

	// このメソッドは呼ばれないはずです
	// mockStore.EXPECT().
	// 	ListUserTodos(gomock.Any(), "nonexistent").
	// 	Times(0)

	service := NewTodoService(mockStore)

	ctx := context.Background()
	todos, err := service.GetUserTodos(ctx, "nonexistent")

	require.Error(t, err)
	assert.Nil(t, todos)
	assert.Contains(t, err.Error(), "user not found")
}
