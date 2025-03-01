# Comparison of Big Interface and Small Interface in Mock Testing

This document provides comparison examples of mock testing using Big Interface and Small Interface approaches.

## Impact of Adding New Features

For example, let's assume that we are adding a new feature "user activity logging" to the system. Let's see how this affects each approach.

### Case of Big Interface

To add a new feature, we need to add a new method to the `DataStore` interface:

```go
// DataStore adds a new method
type DataStore interface {
    // Existing methods
    GetUser(ctx context.Context, id string) (*domain.User, error)
    // ...other methods

    // New method
    LogUserActivity(ctx context.Context, userID string, activity string) error
}
```

This change affects all tests that use the `DataStore` interface. We need to add the new method to all existing mocks:

```go
// Need to update mock configuration in all test files
mockStore := mocks.NewMockDataStore(ctrl)
mockStore.EXPECT().GetUser(gomock.Any(), "user1").Return(mockUser, nil)
// Need to set expectations for LogUserActivity (even if not actually used)
mockStore.EXPECT().LogUserActivity(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
```

To prevent existing unit tests from breaking, we need to **update the mock configuration in all test files**. This is true even if the test doesn't actually use the new feature.

### Case of Small Interface

For a new feature, we can create a new interface or extend an existing one:

```go
// New interface
type ActivityLogger interface {
    LogUserActivity(ctx context.Context, userID string, activity string) error
}

// Existing interface remains unchanged
type UserStore interface {
    GetUser(ctx context.Context, id string) (*domain.User, error)
    // ...other existing methods
}
```

This change only affects services that **actually use the new feature**. Other services and tests are not affected:

```go
// Only services using the new feature need to inject the new interface
type UserActivityService struct {
    userStore     smallinterface.UserStore
    activityLogger smallinterface.ActivityLogger
}

// Existing tests remain unchanged!
func TestUserService_GetUser(t *testing.T) {
    mockUserStore := mocks.NewMockUserStore(ctrl)
    mockUserStore.EXPECT().GetUser(gomock.Any(), "user1").Return(mockUser, nil)
    // No need for ActivityLogger related configuration
}
```

## Comparison of Mock Management Complexity

### Scenario: Testing when a user marks a Todo as complete

#### Big Interface Approach

```go
func TestMarkTodoCompleteAndLogActivity(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Single mock object
    mockStore := mocks.NewMockDataStore(ctrl)

    // Set all expectations on the same mock
    mockStore.EXPECT().GetTodo(gomock.Any(), "todo1").Return(&domain.Todo{ID: "todo1", UserID: "user1"}, nil)
    mockStore.EXPECT().GetUser(gomock.Any(), "user1").Return(&domain.User{ID: "user1"}, nil)
    mockStore.EXPECT().MarkTodoComplete(gomock.Any(), "todo1").Return(nil)
    mockStore.EXPECT().LogUserActivity(gomock.Any(), "user1", "completed todo: todo1").Return(nil)

    service := NewAdvancedTodoService(mockStore)
    err := service.MarkTodoCompleteAndLogActivity(context.Background(), "todo1")
    assert.NoError(t, err)
}
```

#### Small Interface Approach

```go
func TestMarkTodoCompleteAndLogActivity(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Multiple mock objects
    mockTodoStore := mocks.NewMockTodoStore(ctrl)
    mockUserStore := mocks.NewMockUserStore(ctrl)
    mockActivityLogger := mocks.NewMockActivityLogger(ctrl)

    // Set only relevant expectations on each mock
    mockTodoStore.EXPECT().GetTodo(gomock.Any(), "todo1").Return(&domain.Todo{ID: "todo1", UserID: "user1"}, nil)
    mockUserStore.EXPECT().GetUser(gomock.Any(), "user1").Return(&domain.User{ID: "user1"}, nil)
    mockTodoStore.EXPECT().MarkTodoComplete(gomock.Any(), "todo1").Return(nil)
    mockActivityLogger.EXPECT().LogUserActivity(gomock.Any(), "user1", "completed todo: todo1").Return(nil)

    service := NewAdvancedTodoService(mockTodoStore, mockUserStore, mockActivityLogger)
    err := service.MarkTodoCompleteAndLogActivity(context.Background(), "todo1")
    assert.NoError(t, err)
}
```

## Comparison of Interface Change Impact

### Impact of Changes to Big Interface

When you change a method signature in the `DataStore` interface, it affects all code and tests that use that interface:

```go
// Change method arguments
GetUser(ctx context.Context, id string, includeDeleted bool) (*domain.User, error)
```

This change breaks all services that use `DataStore` and all their tests.

### Impact of Changes to Small Interface

When you change a method signature in the `UserStore` interface, it only affects code and tests that use that specific interface:

```go
// Change method arguments
GetUser(ctx context.Context, id string, includeDeleted bool) (*domain.User, error)
```

This change only affects services and tests that use the `UserStore` interface, and does not impact code that only uses `TodoStore`.

## Conclusion

- **Big Interface**:

  - Easy initial setup with fewer mock objects
  - However, changes have a wide impact range, making test maintenance difficult
  - As the interface grows, mock management becomes more complex

- **Small Interface**:
  - Requires slightly more code for initial setup
  - However, changes have a localized impact, making test maintenance easier
  - Each test clearly shows what it's testing, and dependencies are visible
