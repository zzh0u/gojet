# Common Patterns

## API Response Format

```go
// Standard API response structure for Go web applications
type ApiResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Meta    *MetaInfo   `json:"meta,omitempty"`
}

type MetaInfo struct {
    Total int `json:"total,omitempty"`
    Page  int `json:"page,omitempty"`
    Limit int `json:"limit,omitempty"`
}

// Usage example in Gin handler
func GetUsers(c *gin.Context) {
    users, total, err := service.GetUsers()
    if err != nil {
        response.Error(c, err)
        return
    }
    response.SuccessWithMeta(c, users, MetaInfo{Total: total})
}
```

## Middleware Pattern

```go
// Middleware for authentication, logging, etc.
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }

        // Validate token
        userID, err := jwt.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        c.Set("userID", userID)
        c.Next()
    }
}

// Usage in router
router := gin.Default()
router.Use(AuthMiddleware())
```

## Repository Pattern

```go
// Repository interface for data access
type UserRepository interface {
    FindAll(filters Filters) ([]User, error)
    FindByID(id uint) (*User, error)
    Create(user User) error
    Update(user User) error
    Delete(id uint) error
}

// GORM implementation
type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) FindAll(filters Filters) ([]User, error) {
    var users []User
    query := r.db.Model(&User{})
    // Apply filters
    if err := query.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
```

## Skeleton Projects

When implementing new functionality:
1. Search for battle-tested skeleton projects
2. Use parallel agents to evaluate options:
   - Security assessment
   - Extensibility analysis
   - Relevance scoring
   - Implementation planning
3. Clone best match as foundation
4. Iterate within proven structure
