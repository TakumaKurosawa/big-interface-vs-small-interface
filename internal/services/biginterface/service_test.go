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

// -------------------------------------------------------------------------
// Big Interface Approach Testing Characteristics
// -------------------------------------------------------------------------
// 1. Single Mock Object
//    - Only one mock of the large DataStore interface is needed
//    - The same mock object can be reused for different use cases
//
// 2. Simplicity of Mock Setup
//    - Less code is needed as expectations only need to be set on a single mock object
//
// 3. Disadvantage: Dependency on Unnecessary Methods
//    - Depends on a large interface that includes methods not actually used
//    - As the interface grows, the mocks may become more complex
//
// 4. Disadvantage: Reduced Test Specificity
//    - It's often unclear which interface functionality the test depends on
// -------------------------------------------------------------------------

func setupUserExists(mock *mocks.MockDataStore) {
	mock.EXPECT().
		GetUser(gomock.Any(), "user1").
		Return(nil, nil) // Actual values specified in the test case
}

func setupUserNotFound(mock *mocks.MockDataStore) {
	mock.EXPECT().
		GetUser(gomock.Any(), "nonexistent").
		Return(nil, errors.New("user not found"))
}

func TestUserService_GetUser(t *testing.T) {
	mockUser := &domain.User{
		ID:        "user1",
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := map[string]struct {
		userID          string
		setupFunc       func(mock *mocks.MockDataStore)
		expectReturnVal *domain.User
		expectErr       error
	}{
		"Success: User exists": {
			userID:          "user1",
			setupFunc:       setupUserExists,
			expectReturnVal: mockUser,
			expectErr:       nil,
		},
		"Error: User not found": {
			userID:          "nonexistent",
			setupFunc:       setupUserNotFound,
			expectReturnVal: nil,
			expectErr:       errors.New("user not found"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockDataStore(ctrl)
			tt.setupFunc(mockStore)

			// For success case, set actual return values in the mock
			if tt.expectReturnVal != nil {
				mockStore.EXPECT().
					GetUser(gomock.Any(), tt.userID).
					Return(tt.expectReturnVal, nil)
			}

			service := NewUserService(mockStore)

			ctx := context.Background()
			user, err := service.GetUser(ctx, tt.userID)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectReturnVal, user)
			}
		})
	}
}

// In TodoService tests, both user-related and todo-related operations
// are configured through the same DataStore mock instance
// This is a characteristic of the big interface approach - one mock covers multiple related operations
func setupUserAndTodosExist(mock *mocks.MockDataStore) {
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
	mock.EXPECT().
		ListUserTodos(gomock.Any(), "user1").
		Return(nil, nil) // Actual values specified in the test case
}

func setupUserNotFoundForTodos(mock *mocks.MockDataStore) {
	mock.EXPECT().
		GetUser(gomock.Any(), "nonexistent").
		Return(nil, errors.New("user not found"))
}

func TestTodoService_GetUserTodos(t *testing.T) {
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

	tests := map[string]struct {
		userID          string
		setupFunc       func(mock *mocks.MockDataStore)
		expectReturnVal []*domain.Todo
		expectErr       error
	}{
		"Success: User and todos exist": {
			userID:          "user1",
			setupFunc:       setupUserAndTodosExist,
			expectReturnVal: mockTodos,
			expectErr:       nil,
		},
		"Error: User not found": {
			userID:          "nonexistent",
			setupFunc:       setupUserNotFoundForTodos,
			expectReturnVal: nil,
			expectErr:       errors.New("user not found"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockDataStore(ctrl)
			tt.setupFunc(mockStore)

			// For success case, set actual return values in the mock
			if tt.expectReturnVal != nil {
				mockStore.EXPECT().
					ListUserTodos(gomock.Any(), tt.userID).
					Return(tt.expectReturnVal, nil)
			}

			service := NewTodoService(mockStore)

			ctx := context.Background()
			todos, err := service.GetUserTodos(ctx, tt.userID)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
				assert.Nil(t, todos)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectReturnVal, todos)
			}
		})
	}
}

func setupTodoMarkComplete(mock *mocks.MockDataStore) {
	mock.EXPECT().
		MarkTodoComplete(gomock.Any(), "todo1").
		Return(nil)
}

func setupTodoNotFound(mock *mocks.MockDataStore) {
	mock.EXPECT().
		MarkTodoComplete(gomock.Any(), "nonexistent").
		Return(errors.New("todo not found"))
}

func TestTodoService_CompleteTodo(t *testing.T) {
	tests := map[string]struct {
		todoID    string
		setupFunc func(mock *mocks.MockDataStore)
		expectErr error
	}{
		"Success: Todo marked as complete": {
			todoID:    "todo1",
			setupFunc: setupTodoMarkComplete,
			expectErr: nil,
		},
		"Error: Todo not found": {
			todoID:    "nonexistent",
			setupFunc: setupTodoNotFound,
			expectErr: errors.New("todo not found"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockDataStore(ctrl)
			tt.setupFunc(mockStore)

			service := NewTodoService(mockStore)

			ctx := context.Background()
			err := service.CompleteTodo(ctx, tt.todoID)

			if tt.expectErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
