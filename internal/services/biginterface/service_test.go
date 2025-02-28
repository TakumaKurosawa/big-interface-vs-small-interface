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
	tests := map[string]struct {
		userID      string
		mockSetup   func(mock *mocks.MockDataStore) *domain.User
		expectErr   bool
		expectedErr string
	}{
		"Success: User exists": {
			userID: "user1",
			mockSetup: func(mock *mocks.MockDataStore) *domain.User {
				mockUser := &domain.User{
					ID:        "user1",
					Name:      "Test User",
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
		"Error: User not found": {
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockDataStore(ctrl)
			expectedUser := tt.mockSetup(mockStore)

			service := NewUserService(mockStore)

			ctx := context.Background()
			user, err := service.GetUser(ctx, tt.userID)

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
	tests := map[string]struct {
		userID      string
		mockSetup   func(mock *mocks.MockDataStore) []*domain.Todo
		expectErr   bool
		expectedErr string
	}{
		"Success: User and todos exist": {
			userID: "user1",
			mockSetup: func(mock *mocks.MockDataStore) []*domain.Todo {
				mockUser := &domain.User{
					ID:        "user1",
					Name:      "Test User",
					Email:     "test@example.com",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				mockTodos := []*domain.Todo{
					{
						ID:          "todo1",
						UserID:      "user1",
						Title:       "Test Todo",
						Description: "This is a test todo",
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
		"Error: User not found": {
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockDataStore(ctrl)
			expectedTodos := tt.mockSetup(mockStore)

			service := NewTodoService(mockStore)

			ctx := context.Background()
			todos, err := service.GetUserTodos(ctx, tt.userID)

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
	tests := map[string]struct {
		todoID      string
		mockSetup   func(mock *mocks.MockDataStore)
		expectErr   bool
		expectedErr string
	}{
		"Success: Todo marked as complete": {
			todoID: "todo1",
			mockSetup: func(mock *mocks.MockDataStore) {
				mock.EXPECT().
					MarkTodoComplete(gomock.Any(), "todo1").
					Return(nil)
			},
			expectErr: false,
		},
		"Error: Todo not found": {
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

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockDataStore(ctrl)
			tt.mockSetup(mockStore)

			service := NewTodoService(mockStore)

			ctx := context.Background()
			err := service.CompleteTodo(ctx, tt.todoID)

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
