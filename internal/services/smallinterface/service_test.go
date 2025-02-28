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
	// テストケースの定義
	tests := map[string]struct {
		userID      string
		mockSetup   func(mock *mocks.MockUserStore) *domain.User
		expectErr   bool
		expectedErr string
	}{
		"正常系：ユーザーが取得できる場合": {
			userID: "user1",
			mockSetup: func(mock *mocks.MockUserStore) *domain.User {
				mockUser := &domain.User{
					ID:        "user1",
					Name:      "テスト ユーザー",
					Email:     "test@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mock.EXPECT().
					GetUser(gomock.Any(), "user1").
					Return(mockUser, nil)
				return mockUser
			},
			expectErr: false,
		},
		"異常系：ユーザーが存在しない場合": {
			userID: "nonexistent",
			mockSetup: func(mock *mocks.MockUserStore) *domain.User {
				mock.EXPECT().
					GetUser(gomock.Any(), "nonexistent").
					Return(nil, errors.New("user not found"))
				return nil
			},
			expectErr:   true,
			expectedErr: "user not found",
		},
	}

	// 各テストケースを実行
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// モックのセットアップ
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockUserStore(ctrl)
			expectedUser := tt.mockSetup(mockStore)

			// サービスの作成
			service := NewUserService(mockStore)

			// テスト対象メソッドの実行
			ctx := context.Background()
			user, err := service.GetUser(ctx, tt.userID)

			// 結果の検証
			if tt.expectErr {
				require.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expectedUser, user)
			}
		})
	}
}

func TestTodoService_GetUserTodos(t *testing.T) {
	// テストケースの定義
	tests := map[string]struct {
		userID        string
		mockUserStore func(mock *mocks.MockUserStore)
		mockTodoStore func(mock *mocks.MockTodoStore) []*domain.Todo
		expectErr     bool
		expectedErr   string
	}{
		"正常系：ユーザーとTodoが存在する場合": {
			userID: "user1",
			mockUserStore: func(mock *mocks.MockUserStore) {
				mockUser := &domain.User{
					ID:        "user1",
					Name:      "テスト ユーザー",
					Email:     "test@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mock.EXPECT().
					GetUser(gomock.Any(), "user1").
					Return(mockUser, nil)
			},
			mockTodoStore: func(mock *mocks.MockTodoStore) []*domain.Todo {
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
				mock.EXPECT().
					ListUserTodos(gomock.Any(), "user1").
					Return(mockTodos, nil)
				return mockTodos
			},
			expectErr: false,
		},
		"異常系：ユーザーが存在しない場合": {
			userID: "nonexistent",
			mockUserStore: func(mock *mocks.MockUserStore) {
				mock.EXPECT().
					GetUser(gomock.Any(), "nonexistent").
					Return(nil, errors.New("user not found"))
			},
			mockTodoStore: func(mock *mocks.MockTodoStore) []*domain.Todo {
				// この場合、TodoStoreのメソッドは呼ばれないので、モックセットアップは不要
				return nil
			},
			expectErr:   true,
			expectedErr: "user not found",
		},
	}

	// 各テストケースを実行
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// モックのセットアップ
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserStore := mocks.NewMockUserStore(ctrl)
			mockTodoStore := mocks.NewMockTodoStore(ctrl)

			tt.mockUserStore(mockUserStore)
			expectedTodos := tt.mockTodoStore(mockTodoStore)

			// サービスの作成
			service := NewTodoService(mockTodoStore, mockUserStore)

			// テスト対象メソッドの実行
			ctx := context.Background()
			todos, err := service.GetUserTodos(ctx, tt.userID)

			// 結果の検証
			if tt.expectErr {
				require.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
				assert.Nil(t, todos)
			} else {
				require.NoError(t, err)
				assert.Equal(t, expectedTodos, todos)
			}
		})
	}
}

func TestTodoService_CompleteTodo(t *testing.T) {
	// テストケースの定義
	tests := map[string]struct {
		todoID      string
		mockSetup   func(mock *mocks.MockTodoStore)
		expectErr   bool
		expectedErr string
	}{
		"正常系：Todoを完了状態にできる場合": {
			todoID: "todo1",
			mockSetup: func(mock *mocks.MockTodoStore) {
				mock.EXPECT().
					MarkTodoComplete(gomock.Any(), "todo1").
					Return(nil)
			},
			expectErr: false,
		},
		"異常系：Todoが存在しない場合": {
			todoID: "nonexistent",
			mockSetup: func(mock *mocks.MockTodoStore) {
				mock.EXPECT().
					MarkTodoComplete(gomock.Any(), "nonexistent").
					Return(errors.New("todo not found"))
			},
			expectErr:   true,
			expectedErr: "todo not found",
		},
	}

	// 各テストケースを実行
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// モックのセットアップ
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockTodoStore := mocks.NewMockTodoStore(ctrl)
			mockUserStore := mocks.NewMockUserStore(ctrl)
			tt.mockSetup(mockTodoStore)

			// サービスの作成
			service := NewTodoService(mockTodoStore, mockUserStore)

			// テスト対象メソッドの実行
			ctx := context.Background()
			err := service.CompleteTodo(ctx, tt.todoID)

			// 結果の検証
			if tt.expectErr {
				require.Error(t, err)
				if tt.expectedErr != "" {
					assert.Contains(t, err.Error(), tt.expectedErr)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
