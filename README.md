# Interface Size and Cohesion Comparison

This project is a sample to demonstrate the difference in code cohesion based on interface size.

## Overview

This repository compares two approaches:

1. **Big Interface Approach**:

   - Uses a large single interface (`datastore`)
   - Combines all data operations into one interface
   - Low cohesion, requiring definition of unnecessary function behaviors

2. **Small Interface Approach**:
   - Uses small specialized interfaces (`userstore`, `todostore`)
   - Defines only operations necessary for each domain logic
   - High cohesion, allowing use of only what's needed

## Installation

```bash
# Clone the repository
git clone https://github.com/TakumaKurosawa/big-interface-vs-small-interface.git

# Move to the directory
cd big-interface-vs-small-interface

# Install dependencies
go mod download
```

## Running Tests

```bash
# Run all tests
go test ./...
```

## Project Structure

```
.
├── cmd/                        # Main application entry point
├── internal/                   # Packages that should not be imported externally
│   ├── domain/                 # Domain models
│   ├── biginterface/           # Big Interface approach definition
│   │   ├── datastore.go        # Large interface definition
│   │   └── mocks/              # Big Interface mocks
│   │       └── mock_datastore.go # DataStore mock
│   ├── smallinterface/         # Small Interface approach definition
│   │   ├── userstore.go        # User-related small interface
│   │   ├── todostore.go        # Todo-related small interface
│   │   └── mocks/              # Small Interface mocks
│   │       ├── mock_userstore.go # UserStore mock
│   │       └── mock_todostore.go # TodoStore mock
│   ├── infra/                  # Infrastructure layer implementation
│   │   └── inmemory/           # In-memory implementation
│   │       └── store.go        # In-memory store implementation
│   └── services/               # Service implementations
│       ├── biginterface/       # Big Interface approach service implementation
│       │   ├── service.go      # Service implementation
│       │   └── service_test.go # Service tests
│       └── smallinterface/     # Small Interface approach service implementation
│           ├── service.go      # Service implementation
│           └── service_test.go # Service tests
├── pkg/                        # Packages that can be imported externally
│   └── greeting/               # Common modules
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
└── README.md                   # This file
```

## License

MIT
