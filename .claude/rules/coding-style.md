# Coding Style

## Immutability (CRITICAL)

ALWAYS create new objects, NEVER mutate existing ones:

```go
// WRONG: Mutation
func updateUser(user *User, name string) {
    user.Name = name  // MUTATION!
}

// CORRECT: Immutability - return a new instance
func updateUser(user User, name string) User {
    return User{
        ID:   user.ID,
        Name: name,
        // copy other fields
    }
}
```

## File Organization

MANY SMALL FILES > FEW LARGE FILES:
- High cohesion, low coupling
- 200-400 lines typical, 1000 max for Go files
- Extract utilities and helpers into separate files
- Organize by feature/domain, not by type
- Follow Go standard project layout conventions

## Error Handling

ALWAYS handle errors comprehensively:

```go
// WRONG: Ignoring errors
result, _ := riskyOperation()

// CORRECT: Proper error handling with context
result, err := riskyOperation()
if err != nil {
    return fmt.Errorf("riskyOperation failed: %w", err)
}

// For functions that return errors
func ProcessData(data []byte) error {
    if err := validate(data); err != nil {
        return fmt.Errorf("validate: %w", err)
    }
    // processing logic
    return nil
}
```

## Input Validation

ALWAYS validate user input:

```go
// Using struct tags for validation
type UserRequest struct {
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"required,min=0,max=150"`
}

// Manual validation in service layer
func validateUser(req UserRequest) error {
    if req.Age < 0 || req.Age > 150 {
        return fmt.Errorf("age must be between 0 and 150")
    }
    if !isValidEmail(req.Email) {
        return fmt.Errorf("invalid email format")
    }
    return nil
}
```

## Code Quality Checklist

Before marking work complete:
- [ ] Code is readable and well-named (idiomatic Go names)
- [ ] Functions are small and focused (<50 lines)
- [ ] Files are focused and not too large (<1000 lines)
- [ ] No deep nesting (>4 levels)
- [ ] Proper error handling (errors are checked and wrapped)
- [ ] No fmt.Println/fmt.Printf for debugging in production code
- [ ] No hardcoded values (use constants or configuration)
- [ ] No mutation of function parameters (prefer return values)
- [ ] Follow Go conventions (go fmt, go vet, go lint)
- [ ] Documentation for exported functions and types
