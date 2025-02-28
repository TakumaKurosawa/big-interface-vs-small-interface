package biginterface

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/biginterface/mocks"
	"github.com/TakumaKurosawa/big-interface-vs-small-interface/internal/domain"
)

func TestUserService_GetUser(t *testing.T) {
	// テストケースの定義
	tests := map[string]struct {
		userID      string
		mockSetup   func(mock *mocks.MockDataStore) *domain.User
		expectErr   bool
		expectedErr string
	}{
		"正常系：ユーザーが取得できる場合": {
			userID: "user1",
			mockSetup: func(mock *mocks.MockDataStore) *domain.User {
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
			mockSetup: func(mock *mocks.MockDataStore) *domain.User {
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
			mockStore := mocks.NewMockDataStore(ctrl)
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
		userID      string
		mockSetup   func(mock *mocks.MockDataStore) []*domain.Todo
		expectErr   bool
		expectedErr string
	}{
		"正常系：ユーザーとTodoが存在する場合": {
			userID: "user1",
			mockSetup: func(mock *mocks.MockDataStore) []*domain.Todo {
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

				mock.EXPECT().
					GetUser(gomock.Any(), "user1").
					Return(mockUser, nil)
				mock.EXPECT().
					ListUserTodos(gomock.Any(), "user1").
					Return(mockTodos, nil)

				return mockTodos
			},
			expectErr: false,
		},
		"異常系：ユーザーが存在しない場合": {
			userID: "nonexistent",
			mockSetup: func(mock *mocks.MockDataStore) []*domain.Todo {
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
			mockStore := mocks.NewMockDataStore(ctrl)
			expectedTodos := tt.mockSetup(mockStore)

			// サービスの作成
			service := NewTodoService(mockStore)

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
		mockSetup   func(mock *mocks.MockDataStore)
		expectErr   bool
		expectedErr string
	}{
		"正常系：Todoを完了状態にできる場合": {
			todoID: "todo1",
			mockSetup: func(mock *mocks.MockDataStore) {
				mock.EXPECT().
					MarkTodoComplete(gomock.Any(), "todo1").
					Return(nil)
			},
			expectErr: false,
		},
		"異常系：Todoが存在しない場合": {
			todoID: "nonexistent",
			mockSetup: func(mock *mocks.MockDataStore) {
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
			mockStore := mocks.NewMockDataStore(ctrl)
			tt.mockSetup(mockStore)

			// サービスの作成
			service := NewTodoService(mockStore)

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
