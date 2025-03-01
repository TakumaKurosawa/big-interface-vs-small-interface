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

// -------------------------------------------------------------------------
// Small Interface Approach Testing Characteristics
// -------------------------------------------------------------------------
// 1. Multiple Specialized Mock Objects
//    - Uses small interface mocks specialized for each domain (User and Todo)
//    - For TodoService tests, both UserStore and TodoStore mocks are needed
//
// 2. Clear Mock Responsibilities
//    - Each interface focuses on specific responsibilities, making
//      mock roles and expectation setup clearer
//    - It's clear what the test is testing and which components it depends on
//
// 3. Benefit: Compliance with Interface Segregation Principle (ISP)
//    - Clients only depend on methods they need
//    - Avoids accidental dependencies on unnecessary methods
//
// 4. Benefit: Test Robustness
//    - Tests are less affected when unrelated functionality changes
//
// 5. Disadvantage: Managing Multiple Mock Objects
//    - Setup code increases slightly due to handling multiple interfaces
// -------------------------------------------------------------------------

func setupUserExists(mock *mocks.MockUserStore) {
	mock.EXPECT().
		GetUser(gomock.Any(), "user1").
		Return(nil, nil) // Actual values specified in the test case
}

func setupUserNotFound(mock *mocks.MockUserStore) {
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
		setupFunc       func(mock *mocks.MockUserStore)
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
			mockStore := mocks.NewMockUserStore(ctrl)
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

// TodoService tests use separate UserStore and TodoStore mocks
// This is a characteristic of the small interface approach - each interface has clear responsibilities
func setupUserExistsForTodos(mock *mocks.MockUserStore) {
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
}

func setupUserNotFoundForTodos(mock *mocks.MockUserStore) {
	mock.EXPECT().
		GetUser(gomock.Any(), "nonexistent").
		Return(nil, errors.New("user not found"))
}

func setupTodosExist(mock *mocks.MockTodoStore) {
	mock.EXPECT().
		ListUserTodos(gomock.Any(), "user1").
		Return(nil, nil) // Actual values specified in the test case
}

func setupNoTodos(mock *mocks.MockTodoStore) {
	// Do nothing
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
		setupUserFunc   func(mock *mocks.MockUserStore)
		setupTodoFunc   func(mock *mocks.MockTodoStore)
		expectReturnVal []*domain.Todo
		expectErr       error
	}{
		"Success: User and todos exist": {
			userID:          "user1",
			setupUserFunc:   setupUserExistsForTodos,
			setupTodoFunc:   setupTodosExist,
			expectReturnVal: mockTodos,
			expectErr:       nil,
		},
		"Error: User not found": {
			userID:          "nonexistent",
			setupUserFunc:   setupUserNotFoundForTodos,
			setupTodoFunc:   setupNoTodos,
			expectReturnVal: nil,
			expectErr:       errors.New("user not found"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Small interface approach requires multiple mock objects
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserStore := mocks.NewMockUserStore(ctrl)
			mockTodoStore := mocks.NewMockTodoStore(ctrl)

			tt.setupUserFunc(mockUserStore)
			tt.setupTodoFunc(mockTodoStore)

			// For success case, set actual return values in the mock
			if tt.expectReturnVal != nil {
				mockTodoStore.EXPECT().
					ListUserTodos(gomock.Any(), tt.userID).
					Return(tt.expectReturnVal, nil)
			}

			// Note that TodoService depends on two different interfaces
			service := NewTodoService(mockTodoStore, mockUserStore)

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

func setupTodoMarkComplete(mock *mocks.MockTodoStore) {
	mock.EXPECT().
		MarkTodoComplete(gomock.Any(), "todo1").
		Return(nil)
}

func setupTodoNotFound(mock *mocks.MockTodoStore) {
	mock.EXPECT().
		MarkTodoComplete(gomock.Any(), "nonexistent").
		Return(errors.New("todo not found"))
}

func TestTodoService_CompleteTodo(t *testing.T) {
	tests := map[string]struct {
		todoID    string
		setupFunc func(mock *mocks.MockTodoStore)
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
			// Here too, both TodoStore and UserStore mocks are needed,
			// but UserStore is not actually used in this test.
			// This is a common pattern with the small interface approach.
			mockTodoStore := mocks.NewMockTodoStore(ctrl)
			mockUserStore := mocks.NewMockUserStore(ctrl)
			tt.setupFunc(mockTodoStore)

			service := NewTodoService(mockTodoStore, mockUserStore)

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
