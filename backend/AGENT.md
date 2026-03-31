# Backend Development Standards

This document defines the development and testing standards for the SportLink backend. All code and tests must follow these guidelines to ensure consistency, maintainability, and quality.

## Table of Contents

1. [SOLID Principles](#solid-principles)
2. [Code Against Interfaces](#code-against-interfaces)
3. [High Cohesion and Low Coupling](#high-cohesion-and-low-coupling)
4. [Mock Generation](#mock-generation)
5. [Test Naming Convention](#test-naming-convention)
6. [Test Structure](#test-structure)
7. [Matchers and Assertions](#matchers-and-assertions)
8. [Deterministic Tests](#deterministic-tests)
9. [Test Organization](#test-organization)

## SOLID Principles

All code in the backend **MUST** follow SOLID principles:

### Single Responsibility Principle (SRP)

- **Each interface and struct should have ONE reason to change**
- **Interfaces must be small and focused** - Avoid "god interfaces"
- **Each function should do ONE thing well**

✅ **Good** (Single responsibility):
```go
type QueryParser interface {
    ParseSports(sportsQuery string) ([]common.Sport, error)
    ParseCategories(categoriesQuery string) ([]common.Category, error)
    ParseStatuses(statusesQuery string) ([]matchannouncement.Status, error)
    ParseDate(dateQuery string) (time.Time, error)
    ParseLocation(country, province, locality string) *matchannouncement.Location
}
```

❌ **Bad** (Multiple responsibilities):
```go
type Parser interface {
    ParseSports(...)
    ParseCategories(...)
    SaveToDatabase(...)  // Wrong! Persistence is not parsing responsibility
    SendEmail(...)        // Wrong! Email is not parsing responsibility
}
```

### Open/Closed Principle (OCP)

- **Open for extension, closed for modification**
- Use interfaces to allow extension without modifying existing code
- Prefer composition over inheritance

### Liskov Substitution Principle (LSP)

- **Implementations must be substitutable for their interfaces**
- Any implementation of an interface should work wherever the interface is expected

### Interface Segregation Principle (ISP)

- **Clients should not depend on interfaces they don't use**
- **Split large interfaces into smaller, more specific ones**
- **Interfaces should be cohesive** - methods should be related

✅ **Good** (Segregated interfaces):
```go
type Reader interface {
    Read() ([]Entity, error)
}

type Writer interface {
    Write(entity Entity) error
}

type Repository interface {
    Reader
    Writer
}
```

❌ **Bad** (Fat interface):
```go
type Repository interface {
    Read() ([]Entity, error)
    Write(entity Entity) error
    Delete(id string) error
    Update(entity Entity) error
    SendNotification() error  // Wrong! Not repository responsibility
    GenerateReport() error    // Wrong! Not repository responsibility
}
```

### Dependency Inversion Principle (DIP)

- **Depend on abstractions (interfaces), not concretions (implementations)**
- **High-level modules should not depend on low-level modules**
- **Both should depend on abstractions**

## Code Against Interfaces

### Always Use Interfaces

- **Code against interfaces, NOT implementations**
- **Define interfaces in the domain/application layer**
- **Implementations belong in the infrastructure layer**

### Interface Definition Rules

1. **Interfaces define WHAT, not HOW**
2. **Interfaces should be in the same package as the domain entity or use case**
3. **Interfaces should be small and focused (Single Responsibility)**
4. **Use interfaces for all external dependencies** (repositories, services, clients)

### Example Structure

```
domain/
  └── matchannouncement/
      ├── entity.go
      ├── repository.go          # Interface definition
      └── status.go

infrastructure/
  └── persistence/
      └── matchannouncement/
          └── repository.go      # Implementation of interface
```

✅ **Good** (Interface in domain):
```go
// domain/matchannouncement/repository.go
package matchannouncement

type Repository interface {
    Save(entity Entity) error
    Find(query DomainQuery) ([]Entity, error)
}
```

❌ **Bad** (Interface in infrastructure):
```go
// infrastructure/persistence/matchannouncement/repository.go
package matchannouncement

type Repository interface {  // Wrong! Interface should be in domain
    Save(entity Entity) error
}
```

### Dependency Injection

- **Always inject dependencies through constructors**
- **Use interfaces as parameter types**
- **Never create concrete implementations inside functions**

✅ **Good** (Dependency injection):
```go
type DefaultController struct {
    createMatchAnnouncementUC application.UseCase[Entity, Entity]
    queryParser               parser.QueryParser  // Interface
}

func NewController(
    createMatchAnnouncementUC application.UseCase[Entity, Entity],
    queryParser parser.QueryParser,  // Interface, not implementation
) Controller {
    return &DefaultController{
        createMatchAnnouncementUC: createMatchAnnouncementUC,
        queryParser:               queryParser,
    }
}
```

❌ **Bad** (Creating concrete implementation):
```go
func (sc *DefaultController) FindMatchAnnouncements(c *gin.Context) {
    queryParser := parser.NewQueryParser()  // Wrong! Creating concrete implementation
    // ...
}
```

## High Cohesion and Low Coupling

### High Cohesion

- **Related functionality should be grouped together**
- **Each module/package should have a single, well-defined purpose**
- **Functions in a package should work together toward a common goal**

✅ **Good** (High cohesion):
```go
// parser/parser.go - All parsing-related functions together
type QueryParser interface {
    ParseSports(...)
    ParseCategories(...)
    ParseStatuses(...)
    ParseDate(...)
    ParseLocation(...)
}
```

❌ **Bad** (Low cohesion):
```go
// parser/parser.go - Mixed responsibilities
type Parser interface {
    ParseSports(...)
    SaveToDatabase(...)  // Wrong! Not parsing responsibility
    SendEmail(...)       // Wrong! Not parsing responsibility
}
```

### Low Coupling

- **Modules should depend on abstractions, not concrete implementations**
- **Minimize dependencies between modules**
- **Use dependency injection to reduce coupling**

✅ **Good** (Low coupling):
```go
// Controller depends on interface, not concrete implementation
type DefaultController struct {
    createMatchAnnouncementUC application.UseCase[Entity, Entity]  // Interface
    queryParser               parser.QueryParser                    // Interface
}
```

❌ **Bad** (High coupling):
```go
// Controller depends on concrete implementation
type DefaultController struct {
    repository *DynamoDBRepository  // Wrong! Concrete type
    parser     *DefaultQueryParser   // Wrong! Concrete type
}
```

### Best Practices

1. **Extract helper functions** - Don't let one function do too many things
2. **Use composition** - Build complex behavior from simple components
3. **Keep functions small** - Each function should have a single responsibility
4. **Use interfaces liberally** - They reduce coupling and increase testability

## Mock Generation

### Using Mockery

- **Always use `mockery` to generate mocks** - Never create mocks manually
- **Mocks location**: All mocks must be placed in the `mocks` directory at the root of the backend
- **Package structure**: Mocks must maintain the same package structure as the original code
  - Example: `api/domain/team/repository.go` → `mocks/api/domain/team/repository_mock.go`
- **Mock package naming**: Use the `mocks` package prefix (e.g., `tmocks "sportlink/mocks/api/domain/team"`)

### Generating Mocks

```bash
# Generate mock for a specific interface
mockery --name=Repository --dir=api/domain/team --output=mocks/api/domain/team --outpkg=mocks --case=underscore --filename=repository_mock.go
```

## Test Naming Convention

### Given-When-Then Format (Sentences)

All test names **MUST** follow the `given_when_then` convention as **complete sentences** (no underscores) and be written in **business terms**, not implementation details.

**Format**: `given <initial_state> when <action> then <expected_outcome>`

### Examples

✅ **Good** (Business terms, complete sentences):
```go
"given valid match announcement when creating then returns created announcement"
"given team does not exist when creating then returns error"
"given missing required fields when creating then returns validation error"
"given use case returns error when creating then returns internal server error"
```

❌ **Bad** (Implementation details, underscores):
```go
"test_create_match_announcement_success"
"test_controller_returns_201"
"test_use_case_invoke_called"
"given_valid_match_announcement_when_creating_then_returns_created_announcement" // Has underscores

```

### Guidelines

- Use **complete sentences** with spaces, not underscores
- Use **business language** that describes what the system should do from a user/business perspective
- Avoid technical terms like "controller", "use case", "repository", "HTTP status", etc.
- Focus on **what** is being tested, not **how** it's implemented
- Be specific about the scenario (e.g., "team does not exist" not just "error")
- Write as if describing the scenario to a business stakeholder

## Test Structure

### Standard Test Case Structure

All tests **MUST** use a table-driven approach with the following structure:

```go
    testCases := []struct {
        name    string
        payload RequestType  // or setup data
        on      func(t *testing.T, mock1 *Mock1, mock2 *Mock2)
        then    func(t *testing.T, responseCode int, response map[string]interface{})
    }{
        {
            name: "given <state> when <action> then <outcome>",
        payload: RequestType{...},
        on: func(t *testing.T, mock1 *Mock1, mock2 *Mock2) {
            // Configure mocks and setup scenario
        },
        then: func(t *testing.T, responseCode int, response map[string]interface{}) {
            // Assertions
        },
    },
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T) {
        // Setup
        mock1 := NewMock1(t)
        mock2 := NewMock2(t)
        controller := NewController(mock1, mock2, validator)
        
        // Given
        tc.on(t, mock1, mock2)
        
        // When
        // Execute the code under test
        
        // Then
        tc.then(t, responseCode, response)
    })
}
```

### Function Responsibilities

#### `on` Function

- **Purpose**: Configure mocks and set up the test scenario
- **Parameters**: Receives `*testing.T` and all mocks needed for the test
- **Responsibilities**:
  - Configure mock expectations using specific matchers
  - Set up any required test data
  - Define what the mocks should return

#### `then` Function

- **Purpose**: Evaluate the result and perform assertions
- **Parameters**: Receives `*testing.T`, response code, and response data
- **Responsibilities**:
  - Assert HTTP status codes
  - Assert response structure and values
  - Assert error codes and messages
  - Verify business logic outcomes

## Matchers and Assertions

### Using Matchers

- **ALWAYS use specific matchers** - `mock.Anything` is a **TERRIBLE PRACTICE** that must be avoided
- **Use `mock.MatchedBy()`** for complex validation logic
- **Be specific** about what you're matching - validate the actual values being passed

### Why mock.Anything is BAD

Using `mock.Anything` is a **code smell** that indicates you're not really testing anything meaningful:

❌ **Why it's terrible**:
- You're only testing that the mock was **called**, not **how** it was called
- You're not validating the actual data being passed
- Bugs can slip through because you're not verifying the correct values
- Tests become meaningless - they pass even when the code is wrong
- It's essentially saying "I don't care what data is passed" which defeats the purpose of testing

### Examples

✅ **Good** (Specific matchers - validates the actual data):
```go
useCaseMock.On("Invoke", mock.MatchedBy(func(entity domain.Entity) bool {
    return entity.TeamName == "Boca" &&
           entity.Sport == common.Paddle &&
           entity.Location.Country == "Argentina"
})).Return(expectedEntity, nil)
```

✅ **Good** (For DynamoDB queries - validates query structure):
```go
mockClient.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
    return input.Limit != nil && 
           *input.Limit == 100 &&
           input.FilterExpression != nil
})).Return(&dynamodb.QueryOutput{...}, nil)
```

❌ **Bad** (Using mock.Anything - doesn't validate data):
```go
useCaseMock.On("Invoke", mock.Anything).Return(expectedEntity, nil)
// This test passes even if you pass completely wrong data!
```

❌ **Bad** (Using mock.Anything for query inputs):
```go
mockClient.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{...}, nil)
// You're not testing what query parameters are being used!
```

### Exception: Context Parameters ONLY

The **ONLY** exception is for `context.Context` parameters, where `mock.Anything` is acceptable:

```go
repository.On("Save", mock.Anything, mock.MatchedBy(func(entity Entity) bool {
    return entity.Name == "Boca" &&
           entity.Category == common.L5
})).Return(nil)
```

**Important**: Even though we use `mock.Anything` for context, we **still validate** all other parameters with specific matchers.

## Deterministic Tests

### Requirements

- **Tests MUST be deterministic** - Same input always produces same output
- **No random data** - Use fixed, predictable test data
- **No time-dependent logic** - Mock time if necessary
- **No external dependencies** - All external calls must be mocked
- **Isolated tests** - Each test should be independent and not rely on other tests

### Best Practices

- Use fixed dates, IDs, and values
- Mock all external services and databases
- Use `t.Parallel()` when tests are truly independent
- Clean up after each test if needed

## Test Organization

### File Naming

- Test files must be named `*_test.go`
- Place test files in the same package as the code being tested
- Use `_test` package suffix when testing from external package perspective

### Package Structure

```
backend/
├── api/
│   └── infrastructure/
│       └── rest/
│           └── matchannouncement/
│               ├── controller.go
│               ├── controller_test.go          # Same package tests
│               └── create_match_announcement_controller_test.go
└── mocks/
    └── api/
        └── domain/
            └── matchannouncement/
                └── repository_mock.go
```

### Test Coverage

- Aim for high test coverage of business logic
- Focus on testing behavior, not implementation
- Test happy paths, error cases, and edge cases
- Use table-driven tests to cover multiple scenarios efficiently

## Example: Complete Test File

```go
package matchannouncement_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"

    "sportlink/api/application/matchannouncement/request"
    "sportlink/api/domain/matchannouncement"
    "sportlink/api/infrastructure/middleware"
    "sportlink/api/infrastructure/rest/matchannouncement"
    amocks "sportlink/mocks/api/application"
)

type UseCaseMock = amocks.UseCase[domain.Entity, domain.Entity]

func TestCreateMatchAnnouncement(t *testing.T) {
    validator := validator.New()

    testCases := []struct {
        name    string
        payload request.NewMatchAnnouncementRequest
        on      func(t *testing.T, useCaseMock *UseCaseMock)
        then    func(t *testing.T, responseCode int, response map[string]interface{})
    }{
        {
            name: "given valid match announcement when creating then returns created announcement",
            payload: request.NewMatchAnnouncementRequest{
                TeamName: "Boca",
                Sport:    "Paddle",
                // ... other fields
            },
            on: func(t *testing.T, useCaseMock *UseCaseMock) {
                useCaseMock.On("Invoke", mock.MatchedBy(func(entity domain.Entity) bool {
                    return entity.TeamName == "Boca" && entity.Sport == common.Paddle
                })).Return(expectedEntity, nil)
            },
            then: func(t *testing.T, responseCode int, response map[string]interface{}) {
                assert.Equal(t, http.StatusCreated, responseCode)
                assert.Equal(t, "Boca", response["team_name"])
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            // Setup
            useCaseMock := amocks.NewUseCase[domain.Entity, domain.Entity](t)
            controller := matchannouncement.NewController(useCaseMock, nil, validator)

            gin.SetMode(gin.TestMode)
            router := gin.Default()
            router.Use(middleware.ErrorHandler())
            router.POST("/match-announcement", controller.CreateMatchAnnouncement)

            // Given
            tc.on(t, useCaseMock)
            jsonData, _ := json.Marshal(tc.payload)
            req, _ := http.NewRequest("POST", "/match-announcement", bytes.NewBuffer(jsonData))
            req.Header.Set("Content-Type", "application/json")
            resp := httptest.NewRecorder()

            // When
            router.ServeHTTP(resp, req)

            // Then
            response := createMapResponse(resp)
            tc.then(t, resp.Code, response)
        })
    }
}

func createMapResponse(resp *httptest.ResponseRecorder) map[string]interface{} {
    var response map[string]interface{}
    json.Unmarshal(resp.Body.Bytes(), &response)
    return response
}
```

## Summary Checklist

When writing code, ensure:

- [ ] **SOLID Principles** are followed
- [ ] **Code against interfaces**, not implementations
- [ ] **Interfaces have single responsibility** (small and focused)
- [ ] **High cohesion** - related functionality grouped together
- [ ] **Low coupling** - depend on abstractions, use dependency injection
- [ ] **Interfaces defined in domain/application layer**
- [ ] **Implementations in infrastructure layer**

When writing tests, ensure:

- [ ] Mocks are generated using `mockery` and placed in `mocks/` directory
- [ ] Test names follow `given when then` convention as **complete sentences** (no underscores) in business terms
- [ ] Test structure uses table-driven approach with `on` and `then` functions
- [ ] Specific matchers are used (no `mock.Anything` except for contexts)
- [ ] Tests are deterministic (no random data, fixed values)
- [ ] Tests are isolated and independent
- [ ] All external dependencies are mocked
- [ ] Both happy paths and error cases are covered

## References

- [Mockery Documentation](https://github.com/vektra/mockery)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Go Testing Best Practices](https://golang.org/doc/effective_go#testing)

