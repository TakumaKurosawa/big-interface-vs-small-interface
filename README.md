# Big Interface vs Small Interface

This repository is a sample project to compare the approaches of Big Interface and Small Interface.

## Overview

There are mainly two approaches in interface design in Go:

1. **Big Interface Approach**: Define a single interface with many methods
2. **Small Interface Approach**: Define multiple interfaces with only a few related methods each

This project implements the same functionality with both approaches and demonstrates their advantages and disadvantages.

## Interface Comparison

### Big Interface

```go
// DataStore defines all data operations as a single large interface
type DataStore interface {
    // User-related operations
    GetUser(ctx context.Context, id string) (*domain.User, error)
    ListUsers(ctx context.Context) ([]*domain.User, error)
    CreateUser(ctx context.Context, user *domain.User) error
    UpdateUser(ctx context.Context, user *domain.User) error
    DeleteUser(ctx context.Context, id string) error

    // Todo-related operations
    GetTodo(ctx context.Context, id string) (*domain.Todo, error)
    ListTodos(ctx context.Context) ([]*domain.Todo, error)
    ListUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error)
    CreateTodo(ctx context.Context, todo *domain.Todo) error
    UpdateTodo(ctx context.Context, todo *domain.Todo) error
    DeleteTodo(ctx context.Context, id string) error
    MarkTodoComplete(ctx context.Context, id string) error
}
```

### Small Interface

```go
// UserStore is a small interface that defines only user-related operations
type UserStore interface {
    GetUser(ctx context.Context, id string) (*domain.User, error)
    ListUsers(ctx context.Context) ([]*domain.User, error)
    CreateUser(ctx context.Context, user *domain.User) error
    UpdateUser(ctx context.Context, user *domain.User) error
    DeleteUser(ctx context.Context, id string) error
}

// TodoStore is a small interface that defines only Todo-related operations
type TodoStore interface {
    GetTodo(ctx context.Context, id string) (*domain.Todo, error)
    ListTodos(ctx context.Context) ([]*domain.Todo, error)
    ListUserTodos(ctx context.Context, userID string) ([]*domain.Todo, error)
    CreateTodo(ctx context.Context, todo *domain.Todo) error
    UpdateTodo(ctx context.Context, todo *domain.Todo) error
    DeleteTodo(ctx context.Context, id string) error
    MarkTodoComplete(ctx context.Context, id string) error
}
```

## Service Implementation Comparison

### Big Interface Approach Services

```go
// Big Interface Approach
type UserService struct {
    store biginterface.DataStore
}

type TodoService struct {
    store biginterface.DataStore
}
```

### Small Interface Approach Services

```go
// Small Interface Approach
type UserService struct {
    userStore smallinterface.UserStore
}

type TodoService struct {
    todoStore smallinterface.TodoStore
    userStore smallinterface.UserStore
}
```

## Differences in Testing

There are significant differences between the two approaches, especially in unit testing with mocks:

### Big Interface Testing Characteristics

1. **Single Mock Object**

   - Only need to create and manage a single mock for the large `DataStore` interface
   - Reuse the same mock object for different use cases

2. **Simplicity of Mock Configuration**

   - Less code tends to be required as only one mock object needs to be configured

3. **Disadvantage: Dependency on Unnecessary Methods**

   - Depends on a large interface that includes methods that aren't actually used
   - As the interface grows, mock management becomes more complex

4. **Disadvantage: Reduced Test Specificity**
   - Often unclear which interface functionality the test actually depends on

### Small Interface Testing Characteristics

1. **Multiple Specialized Mock Objects**

   - Use small interface mocks specialized for each domain (User, Todo)
   - TodoService tests require both UserStore and TodoStore mocks

2. **Clear Mock Responsibilities**

   - Each interface focuses on specific responsibilities
   - It's clear what components the test depends on

3. **Benefit: Compliance with Interface Segregation Principle**

   - Clients only depend on methods they need
   - Avoid accidental dependencies on unnecessary functionality

4. **Benefit: Test Robustness**

   - Tests are less affected when unrelated functionality changes

5. **Disadvantage: Managing Multiple Mock Objects**
   - Setup code increases slightly due to handling multiple interfaces

For more detailed comparisons with specific examples, see the [Comparative Testing Examples](internal/services/comparative_testing_example.md) document, which illustrates how both approaches handle real-world scenarios like adding new features, maintaining mocks, and dealing with interface changes.

## Test Code Examples

### Big Interface Test Example (Partial)

```go
// BigInterface approach
func TestTodoService_GetUserTodos(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Use only one mock object
    mockStore := mocks.NewMockDataStore(ctrl)

    // Configure both user-related and todo-related operations through the same mock
    mockUser := &domain.User{ID: "user1", Name: "Test User"}
    mockStore.EXPECT().GetUser(gomock.Any(), "user1").Return(mockUser, nil)
    mockStore.EXPECT().ListUserTodos(gomock.Any(), "user1").Return([]*domain.Todo{}, nil)

    service := NewTodoService(mockStore)
    // Execute test...
}
```

### Small Interface Test Example (Partial)

```go
// SmallInterface approach
func TestTodoService_GetUserTodos(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Multiple mock objects required
    mockUserStore := mocks.NewMockUserStore(ctrl)
    mockTodoStore := mocks.NewMockTodoStore(ctrl)

    // Set expectations explicitly for each store
    mockUser := &domain.User{ID: "user1", Name: "Test User"}
    mockUserStore.EXPECT().GetUser(gomock.Any(), "user1").Return(mockUser, nil)
    mockTodoStore.EXPECT().ListUserTodos(gomock.Any(), "user1").Return([]*domain.Todo{}, nil)

    service := NewTodoService(mockTodoStore, mockUserStore)
    // Execute test...
}
```

## Project Structure

```
.
├── cmd/
│   └── main.go              # Main application
├── internal/
│   ├── domain/              # Domain models
│   │   └── models.go
│   ├── biginterface/        # Big interface approach
│   │   ├── datastore.go     # Large single interface
│   │   ├── mocks/           # Interface mocks
│   │   │   └── mock_datastore.go
│   ├── smallinterface/      # Small interface approach
│   │   ├── userstore.go     # User-related small interface
│   │   ├── todostore.go     # Todo-related small interface
│   │   ├── mocks/           # Interface mocks
│   │   │   ├── mock_userstore.go
│   │   │   └── mock_todostore.go
│   ├── services/            # Service implementations
│   │   ├── biginterface/    # Services using big interface
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   └── smallinterface/  # Services using small interface
│   │       ├── service.go
│   │       └── service_test.go
│   │   └── comparative_testing_example.md  # Detailed comparison document
│   └── infra/               # Infrastructure implementations
│       └── inmemory/        # In-memory implementation
│           └── store.go     # Implements both interfaces
```

## How to Run

```bash
go run cmd/main.go
```

## Running Tests

```bash
go test ./...
```

## General Recommendations for Interface Design

1. **Prefer Small Interfaces**

   - Follow the Interface Segregation Principle (ISP)
   - Clients do not depend on unnecessary methods
   - Combine multiple interfaces as needed

2. **Design Interfaces Based on Roles**

   - Design interfaces based on client needs
   - Design from the user's perspective, not for implementation convenience

3. **Evolve Interfaces**
   - Don't try to design the perfect interface from the beginning
   - Extend gradually as needed
