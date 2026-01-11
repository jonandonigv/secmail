# AGENTS.md - Secmail Project Guidelines

This file contains essential information for agentic coding assistants working on the Secmail project. Follow these guidelines to maintain consistency and quality.

## Project Overview

Secmail is a secure email system built in Go that provides end-to-end encryption for email messages. It uses asymmetric encryption (RSA) for key exchange and symmetric encryption (age) for message bodies.

## Build and Development Commands

### Basic Commands
- **Install dependencies**: `go mod tidy`
- **Run the server**: `go run main.go`
- **Build executable**: `go build`
- **Format code**: `gofmt -w .`
- **Lint code**: `go vet ./...`
- **Run tests**: `go test ./...` (when tests are added)

### Database Setup
- **Database**: PostgreSQL
- **Default DSN**: `host=localhost user=postgres password=postgres dbname=secmail port=5432 sslmode=disable`
- **Create database**: `createdb secmail` (requires PostgreSQL running)

### Testing
Currently no tests exist. When adding tests:
- **Run all tests**: `go test ./...`
- **Run single test**: `go test -run TestFunctionName ./package/path`
- **Run with coverage**: `go test -cover ./...`
- **Run with verbose output**: `go test -v ./...`

## Code Style Guidelines

### General Go Conventions
- Follow standard Go formatting with `gofmt`
- Use `go vet` for static analysis
- Maximum line length: 100 characters
- Use meaningful variable names
- Prefer early returns for error handling

### Package Organization
```
secmail/
├── main.go                 # Application entry point
├── internal/
│   ├── auth/              # Authentication logic
│   ├── database/          # Database connection and setup
│   ├── handlers/          # HTTP request handlers
│   ├── email/             # Email send/receive logic
│   ├── crypto/            # Encryption/decryption utilities
│   └── models/            # Database models
└── go.mod
```

### Naming Conventions
- **Packages**: lowercase, single word (e.g., `auth`, `email`)
- **Exported functions/types**: PascalCase (e.g., `SendMessage`, `User`)
- **Unexported functions/types**: camelCase (e.g., `sendMessage`, `user`)
- **Variables**: camelCase (e.g., `userID`, `encryptedBody`)
- **Constants**: PascalCase (e.g., `JWTExpiration`)

### Import Organization
```go
import (
    // Standard library imports (alphabetical)
    "encoding/json"
    "net/http"
    "time"

    // Third-party imports (alphabetical)
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "gorm.io/gorm"

    // Local imports
    "secmail/internal/models"
)
```

### Struct Tags
- **JSON**: `json:"field_name"`
- **GORM**: `gorm:"column_name;constraints"`
- **Validation**: `binding:"required,email"`

### Error Handling
- Return errors from functions rather than panicking
- Handle errors immediately in calling code
- Use descriptive error messages
- Log errors appropriately

```go
// Good
func processData(data []byte) error {
    if len(data) == 0 {
        return errors.New("data cannot be empty")
    }
    // ... processing logic
    return nil
}

// In handler
if err := processData(input); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

### Function Signatures
- Keep functions focused and single-purpose
- Use pointer receivers for structs when mutation is needed
- Return errors as the last return value
- Use context.Context for cancellable operations

### Database Operations
- Use GORM for database interactions
- Handle GORM errors appropriately
- Use transactions for multi-step operations
- Validate data before database operations

### Security Best Practices
- **NEVER** log sensitive information (passwords, private keys, tokens)
- Use environment variables for secrets (JWT secret, database credentials)
- Validate and sanitize all user inputs
- Use prepared statements (handled by GORM)
- Implement proper authentication/authorization

### Comments and Documentation
- Add comments for exported functions and types
- Use complete sentences starting with the name
- Keep comments concise but informative

```go
// SendMessage sends an encrypted email from sender to recipients.
func SendMessage(senderID uint, recipients []uint, subject, body string, db *gorm.DB) error {
    // Implementation...
}
```

### Constants and Configuration
- Define constants for magic numbers and strings
- Use environment variables for runtime configuration
- Group related constants together

### Testing Guidelines (Future)
- Write unit tests for all exported functions
- Use table-driven tests for multiple test cases
- Mock external dependencies
- Test error conditions
- Use descriptive test names: `TestFunctionName_Scenario_Result`

### API Design
- Use RESTful conventions
- Return appropriate HTTP status codes
- Use JSON for request/response bodies
- Include error messages in responses
- Use middleware for cross-cutting concerns (auth, logging)

### Encryption Guidelines
- Use established crypto libraries (golang.org/x/crypto, filippo.io/age)
- Generate strong keys (RSA 2048+ bits)
- Never store private keys in plaintext
- Use hybrid encryption (RSA + AES)
- Validate cryptographic operations

### Logging
- Use the standard `log` package for simple logging
- Log errors and important events
- Don't log sensitive information
- Consider structured logging for production

### Dependencies
- Keep dependencies minimal and well-maintained
- Review licenses and security vulnerabilities
- Use Go modules for dependency management
- Update dependencies regularly

### Git Workflow
- Use descriptive commit messages
- Commit related changes together
- Use branches for features/bugs
- Write clear pull request descriptions

## Common Patterns

### Handler Structure
```go
func HandlerName(c *gin.Context, db *gorm.DB) {
    // Extract and validate user ID from context
    userIDVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    userID := userIDVal.(uint)

    // Parse request
    var req RequestType
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Business logic
    result, err := businessFunction(userID, req, db)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Success response
    c.JSON(http.StatusOK, gin.H{"data": result})
}
```

### Model Structure
```go
type ModelName struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    FieldName string         `gorm:"not null" json:"field_name" binding:"required"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

## Quality Checks

Before committing code:
1. Run `gofmt -w .` to format code
2. Run `go vet ./...` to check for issues
3. Run `go build` to ensure compilation
4. Test manually if applicable
5. Review code for security issues

## Future Improvements
- Add comprehensive test suite
- Implement proper logging framework
- Add API documentation
- Set up CI/CD pipeline
- Add security audit tools (gosec)
- Implement proper configuration management</content>
<parameter name="filePath">AGENTS.md