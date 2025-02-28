package smallinterface

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/smallinterface/mocks"
)

func TestUserService_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 小さなインターフェースをモック化
	mockStore := mocks.NewMockUserStore(ctrl)

	mockUser := &domain.User{
		ID:        "user1",
		Name:      "テスト ユーザー",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 小さなインターフェースは、ユーザー関連のメソッドのみを持つため、
	// モック化するのは必要なメソッドだけです
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

	// 小さなインターフェースをそれぞれモック化
	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockTodoStore := mocks.NewMockTodoStore(ctrl)

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

	// 小さなインターフェースごとに必要なメソッドのみをモック化
	mockUserStore.EXPECT().
		GetUser(gomock.Any(), "user1").
		Return(mockUser, nil)

	mockTodoStore.EXPECT().
		ListUserTodos(gomock.Any(), "user1").
		Return(mockTodos, nil)

	service := NewTodoService(mockTodoStore, mockUserStore)

	ctx := context.Background()
	todos, err := service.GetUserTodos(ctx, "user1")

	require.NoError(t, err)
	assert.Equal(t, mockTodos, todos)
}

func TestTodoService_GetUserTodos_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 小さなインターフェースをそれぞれモック化
	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockTodoStore := mocks.NewMockTodoStore(ctrl)

	// ユーザーが見つからない場合のテスト
	mockUserStore.EXPECT().
		GetUser(gomock.Any(), "nonexistent").
		Return(nil, errors.New("user not found"))

	// このメソッドは呼ばれないはずです
	// Todo関連のインターフェースをモック化するときに、
	// User関連のメソッドを気にする必要はありません

	service := NewTodoService(mockTodoStore, mockUserStore)

	ctx := context.Background()
	todos, err := service.GetUserTodos(ctx, "nonexistent")

	require.Error(t, err)
	assert.Nil(t, todos)
	assert.Contains(t, err.Error(), "user not found")
}
