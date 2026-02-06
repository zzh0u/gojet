---
name: security-review
description: Use this skill when adding authentication, handling user input, working with secrets, creating API endpoints, or implementing payment/sensitive features. Provides comprehensive security checklist and patterns.
---

# Security Review Skill (Go Edition)

This skill ensures Go code follows security best practices and identifies potential vulnerabilities specific to Go applications.

## When to Activate

- Implementing authentication or authorization in Go web applications
- Handling user input or file uploads in Gin/Go HTTP handlers
- Creating new API endpoints in Go
- Working with secrets or credentials in Go configuration
- Implementing payment or sensitive features in Go
- Storing or transmitting sensitive data
- Integrating third-party APIs with Go clients

## Security Checklist

### 1. Secrets Management

#### ❌ NEVER Do This
```go
// Hardcoded secrets in source code
const apiKey = "sk-proj-xxxxx"
const dbPassword = "password123"

// In configuration structs
type Config struct {
    JWTSecret string `json:"jwt_secret"`  // Will be hardcoded in JSON
}
```

#### ✅ ALWAYS Do This
```go
// Use environment variables or dedicated secret management
import (
    "os"
    "fmt"
)

// Load from environment
jwtSecret := os.Getenv("JWT_SECRET")
if jwtSecret == "" {
    return fmt.Errorf("JWT_SECRET environment variable not set")
}

// Or use configuration struct with validation
type Config struct {
    JWTSecret   string `env:"JWT_SECRET,required"`
    DatabaseURL string `env:"DATABASE_URL,required"`
    APIKey      string `env:"API_KEY"`
}

// Use packages like github.com/caarlos0/env for structured env loading
import "github.com/caarlos0/env/v6"

var cfg Config
if err := env.Parse(&cfg); err != nil {
    log.Fatal("Failed to parse config:", err)
}
```

#### Verification Steps
- [ ] No hardcoded API keys, tokens, or passwords in source code
- [ ] All secrets loaded from environment variables or secret managers
- [ ] Configuration files with secrets excluded from git (`.env`, `config/local.yaml`)
- [ ] No secrets in git history (check with `git log -p -S "password"`)
- [ ] Production secrets managed by platform (Kubernetes Secrets, AWS Secrets Manager, etc.)
- [ ] Secrets validated at application startup

### 2. Input Validation

#### Always Validate User Input
```go
// Use struct tags for validation (Gin framework example)
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Name     string `json:"name" binding:"required,min=1,max=100"`
    Age      int    `json:"age" binding:"required,min=0,max=150"`
    Password string `json:"password" binding:"required,min=8"`
}

// In Gin handler
func CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid input", "details": err.Error()})
        return
    }

    // Proceed with validated request
    user, err := service.CreateUser(req)
    if err != nil {
        c.JSON(500, gin.H{"error": "Internal server error"})
        return
    }

    c.JSON(200, gin.H{"data": user})
}

// Custom validation with go-playground/validator
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func validateUser(req CreateUserRequest) error {
    if err := validate.Struct(req); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }

    // Custom business logic validation
    if strings.Contains(strings.ToLower(req.Name), "admin") {
        return fmt.Errorf("name cannot contain 'admin'")
    }

    return nil
}
```

#### File Upload Validation
```go
import (
    "mime/multipart"
    "path/filepath"
    "strings"
)

func validateFileUpload(fileHeader *multipart.FileHeader) error {
    // Size check (5MB max)
    const maxSize = 5 * 1024 * 1024 // 5MB
    if fileHeader.Size > maxSize {
        return fmt.Errorf("file too large (max 5MB)")
    }

    // Content type check
    contentType := fileHeader.Header.Get("Content-Type")
    allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
    validType := false
    for _, t := range allowedTypes {
        if contentType == t {
            validType = true
            break
        }
    }
    if !validType {
        return fmt.Errorf("invalid file type: %s", contentType)
    }

    // Extension check (additional safety)
    ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
    allowedExts := []string{".jpg", ".jpeg", ".png", ".gif"}
    validExt := false
    for _, e := range allowedExts {
        if ext == e {
            validExt = true
            break
        }
    }
    if !validExt {
        return fmt.Errorf("invalid file extension: %s", ext)
    }

    // Check filename for path traversal
    if strings.Contains(fileHeader.Filename, "..") || strings.Contains(fileHeader.Filename, "/") {
        return fmt.Errorf("invalid filename")
    }

    return nil
}

// In Gin handler for file upload
func UploadFile(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }

    if err := validateFileUpload(file); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Save the file securely
    safeFilename := generateSafeFilename(file.Filename)
    dst := filepath.Join("uploads", safeFilename)
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(500, gin.H{"error": "Failed to save file"})
        return
    }

    c.JSON(200, gin.H{"message": "File uploaded successfully"})
}
```

#### Verification Steps
- [ ] All user inputs validated with struct tags or validation libraries
- [ ] File uploads restricted by size, content type, and extension
- [ ] No direct concatenation of user input in queries or commands
- [ ] Whitelist validation (allow known good values) instead of blacklist
- [ ] Error messages generic, no sensitive information exposed
- [ ] Path traversal prevention for file uploads
- [ ] Input length limits enforced to prevent DoS attacks

### 3. SQL Injection Prevention

#### ❌ NEVER Concatenate SQL
```go
// DANGEROUS - SQL Injection vulnerability
func getUserByEmailUnsafe(email string) (*User, error) {
    query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
    var user User
    err := db.Raw(query).Scan(&user).Error
    return &user, err
}

// Also dangerous: using Sprintf with query conditions
func searchUsersUnsafe(search string) ([]User, error) {
    condition := ""
    if search != "" {
        condition = fmt.Sprintf("WHERE name LIKE '%%%s%%'", search)
    }
    query := fmt.Sprintf("SELECT * FROM users %s", condition)
    var users []User
    err := db.Raw(query).Scan(&users).Error
    return users, err
}
```

#### ✅ ALWAYS Use Parameterized Queries
```go
// Safe - GORM parameterized queries
func getUserByEmailSafe(email string) (*User, error) {
    var user User
    err := db.Where("email = ?", email).First(&user).Error
    return &user, err
}

// Safe - GORM with map conditions
func searchUsersSafe(search string) ([]User, error) {
    dbQuery := db.Model(&User{})
    if search != "" {
        dbQuery = dbQuery.Where("name LIKE ?", "%"+search+"%")
    }
    var users []User
    err := dbQuery.Find(&users).Error
    return users, err
}

// Safe - Raw SQL with parameterized queries
func getUserRawSafe(email string) (*User, error) {
    var user User
    err := db.Raw("SELECT * FROM users WHERE email = ?", email).Scan(&user).Error
    return &user, err
}

// Safe - PostgreSQL style numbered parameters
func getUsersByIDs(ids []uint) ([]User, error) {
    var users []User
    query := "SELECT * FROM users WHERE id IN (?)"
    err := db.Raw(query, ids).Scan(&users).Error
    return users, err
}

// Safe - Using GORM's Exec with parameters
func updateUserEmail(userID uint, newEmail string) error {
    result := db.Exec("UPDATE users SET email = ? WHERE id = ?", newEmail, userID)
    return result.Error
}
```

#### Verification Steps
- [ ] All database queries use parameterized queries (?, $1, @param placeholders)
- [ ] No string concatenation or fmt.Sprintf for SQL query construction
- [ ] ORM (GORM) used correctly with Where(), Find(), Raw() with parameters
- [ ] User input never directly interpolated into SQL strings
- [ ] Raw SQL queries always use parameter binding
- [ ] SQL query builders used for complex dynamic queries
- [ ] Database driver's parameterized query support utilized

### 4. Authentication & Authorization

#### JWT Token Handling and Storage
```go
// ❌ WRONG: Storing tokens in localStorage (vulnerable to XSS in SPA)
// Frontend JavaScript: localStorage.setItem('token', token)
// Go backend returning token in JSON response (insecure for web)
func LoginHandler(c *gin.Context) {
    // ... authenticate user ...
    token, err := generateJWT(user)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate token"})
        return
    }
    // Insecure: token in JSON response (client stores in localStorage)
    c.JSON(200, gin.H{"token": token})  // ❌ VULNERABLE TO XSS
}

// ✅ CORRECT: httpOnly, Secure, SameSite cookies
func LoginHandlerSecure(c *gin.Context) {
    // ... authenticate user ...
    token, err := generateJWT(user)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate token"})
        return
    }

    // Set secure httpOnly cookie
    c.SetCookie(
        "token",           // name
        token,             // value
        3600,              // max age in seconds
        "/",               // path
        ".example.com",    // domain (set appropriately)
        true,              // secure (HTTPS only)
        true,              // httpOnly (inaccessible to JavaScript)
    )

    c.JSON(200, gin.H{"message": "Login successful"})
}

// JWT generation and validation using golang-jwt/jwt
import "github.com/golang-jwt/jwt/v4"

type Claims struct {
    UserID uint   `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func generateJWT(user User) (string, error) {
    claims := &Claims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Gin middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get token from cookie or Authorization header
        tokenString, err := c.Cookie("token")
        if err != nil {
            // Fallback to Authorization header
            authHeader := c.GetHeader("Authorization")
            if authHeader == "" {
                c.JSON(401, gin.H{"error": "Authorization required"})
                c.Abort()
                return
            }
            // Extract Bearer token
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                c.JSON(401, gin.H{"error": "Invalid authorization header"})
                c.Abort()
                return
            }
            tokenString = parts[1]
        }

        // Parse and validate token
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.JSON(401, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(*Claims)
        if !ok {
            c.JSON(401, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        // Store user info in context for downstream handlers
        c.Set("userID", claims.UserID)
        c.Set("userRole", claims.Role)
        c.Next()
    }
}
```

#### Authorization Checks
```go
// Role-based authorization middleware
func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("userRole")
        if !exists {
            c.JSON(500, gin.H{"error": "User role not found in context"})
            c.Abort()
            return
        }

        if userRole != requiredRole {
            c.JSON(403, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// Resource-level authorization in service layer
func DeleteUser(userID uint, requesterID uint) error {
    // Get requester from database
    var requester User
    if err := db.First(&requester, requesterID).Error; err != nil {
        return fmt.Errorf("requester not found: %w", err)
    }

    // Authorization check: only admins can delete users
    if requester.Role != "admin" {
        return fmt.Errorf("unauthorized: admin role required")
    }

    // Additional check: users cannot delete themselves (business rule)
    if requesterID == userID {
        return fmt.Errorf("cannot delete your own account")
    }

    // Proceed with deletion
    if err := db.Delete(&User{}, userID).Error; err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    return nil
}

// Usage in Gin handler with combined authentication and authorization
func DeleteUserHandler(c *gin.Context) {
    userIDStr := c.Param("id")
    userID, err := strconv.ParseUint(userIDStr, 10, 32)
    if err != nil {
        c.JSON(400, gin.H{"error": "Invalid user ID"})
        return
    }

    requesterID, exists := c.Get("userID")
    if !exists {
        c.JSON(500, gin.H{"error": "User ID not found in context"})
        return
    }

    if err := DeleteUser(uint(userID), requesterID.(uint)); err != nil {
        // Check error type to provide appropriate response
        if strings.Contains(err.Error(), "unauthorized") {
            c.JSON(403, gin.H{"error": err.Error()})
        } else {
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
        return
    }

    c.JSON(200, gin.H{"message": "User deleted successfully"})
}
```

#### Row Level Security (Application-Level)
```go
// GORM scopes for row-level security
func UserScope(userID uint) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("user_id = ?", userID)
    }
}

// Usage: users can only access their own data
func GetUserOrders(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(500, gin.H{"error": "User ID not found in context"})
        return
    }

    var orders []Order
    // Apply scope to restrict to user's own orders
    err := db.Scopes(UserScope(userID.(uint))).Find(&orders).Error
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch orders"})
        return
    }

    c.JSON(200, gin.H{"data": orders})
}

// Service layer row-level security
type UserRepository struct {
    db *gorm.DB
    currentUserID uint
}

func (r *UserRepository) FindAll() ([]User, error) {
    var users []User
    // Admins see all users, regular users only see themselves
    if r.isAdmin() {
        err := r.db.Find(&users).Error
        return users, err
    } else {
        err := r.db.Where("id = ?", r.currentUserID).Find(&users).Error
        return users, err
    }
}

func (r *UserRepository) isAdmin() bool {
    // Check if current user is admin (implementation depends on your auth system)
    // This is a simplified example
    return false
}
```

#### Verification Steps
- [ ] JWT tokens stored in httpOnly, Secure, SameSite cookies (not localStorage or response JSON)
- [ ] Authorization checks performed before all sensitive operations
- [ ] Row-level security implemented at application level (scopes, repository patterns)
- [ ] Role-based access control (RBAC) implemented with clear role definitions
- [ ] Session management secure (appropriate token expiration, refresh mechanisms)
- [ ] User context properly propagated through request chain (Gin context)
- [ ] Principle of least privilege: users have minimum necessary permissions
- [ ] Regular review of authorization rules and access controls

### 5. XSS Prevention

#### HTML Template Auto-Escaping
```go
// Go's html/template package automatically escapes HTML by default
import (
    "html/template"
    "net/http"
)

func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
    user := getUserFromRequest(r)

    // Define template (usually loaded from files)
    tmpl := `
    <!DOCTYPE html>
    <html>
    <head><title>User Profile</title></head>
    <body>
        <h1>User Profile</h1>
        <p>Name: {{.Name}}</p>  <!-- Auto-escaped -->
        <p>Bio: {{.Bio}}</p>    <!-- Auto-escaped -->
        <p>Email: {{.Email}}</p> <!-- Auto-escaped -->
    </body>
    </html>
    `

    t, err := template.New("profile").Parse(tmpl)
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    // Auto-escaping happens automatically when executing the template
    // Even if user.Bio contains <script>alert('xss')</script>, it will be rendered as text
    err = t.Execute(w, user)
    if err != nil {
        http.Error(w, "Template execution error", http.StatusInternalServerError)
    }
}

// Safe HTML rendering when you intentionally want to render HTML
func RenderSafeHTML(w http.ResponseWriter, r *http.Request) {
    userContent := getUserContent() // Might contain HTML from user

    tmpl := `
    <!DOCTYPE html>
    <html>
    <body>
        <div class="user-content">
            <!-- Use template.HTML type to mark content as safe -->
            {{.SafeContent}}
        </div>
    </body>
    </html>
    `

    t, err := template.New("content").Parse(tmpl)
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }

    // Only use template.HTML if content is sanitized!
    data := struct {
        SafeContent template.HTML
    }{
        SafeContent: template.HTML(sanitizeHTML(userContent)),
    }

    err = t.Execute(w, data)
    if err != nil {
        http.Error(w, "Template execution error", http.StatusInternalServerError)
    }
}
```

#### HTML Sanitization for User Content
```go
// Use bluemonday for HTML sanitization
import (
    "github.com/microcosm-cc/bluemonday"
)

func sanitizeHTML(input string) string {
    // Create a strict policy (only allows simple formatting)
    p := bluemonday.StrictPolicy()

    // Or create a more permissive policy with controlled tags
    p = bluemonday.UGCPolicy() // User Generated Content policy

    // Customize allowed tags and attributes
    p.AllowAttrs("class").OnElements("span", "div")
    p.AllowElements("br", "hr")

    // Remove disallowed tags and attributes
    return p.Sanitize(input)
}

// Example with Gin
func SubmitCommentHandler(c *gin.Context) {
    var req struct {
        Content string `json:"content" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid input"})
        return
    }

    // Sanitize HTML before storage
    sanitizedContent := sanitizeHTML(req.Content)

    // Store sanitized content
    comment := Comment{
        Content: sanitizedContent,
        UserID:  getUserID(c),
    }

    if err := db.Create(&comment).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to save comment"})
        return
    }

    c.JSON(200, gin.H{"message": "Comment submitted"})
}
```

#### Content Security Policy Headers
```go
// Gin middleware for security headers including CSP
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Content Security Policy
        c.Header("Content-Security-Policy",
            "default-src 'self'; " +
            "script-src 'self' 'unsafe-inline' https://cdn.example.com; " +
            "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; " +
            "img-src 'self' data: https:; " +
            "font-src 'self' https://fonts.gstatic.com; " +
            "connect-src 'self' https://api.example.com; " +
            "frame-ancestors 'none'; " +
            "form-action 'self';")

        // Other security headers
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

        c.Next()
    }
}

// Usage in Gin router setup
func main() {
    r := gin.Default()

    // Apply security headers middleware
    r.Use(SecurityHeaders())

    // ... routes ...

    r.Run(":8080")
}
```

#### Verification Steps
- [ ] Go's `html/template` package used for all HTML rendering (auto-escapes by default)
- [ ] User-provided HTML sanitized with bluemonday or similar library before storage/display
- [ ] Content Security Policy headers configured via middleware
- [ ] No use of `text/template` for HTML rendering (does not auto-escape)
- [ ] `template.HTML` type used only for pre-sanitized content
- [ ] Additional security headers set (X-Content-Type-Options, X-Frame-Options, etc.)
- [ ] JSON responses use appropriate Content-Type headers to prevent script execution
- [ ] User input never directly concatenated into HTML/JavaScript strings

### 6. CSRF Protection

#### CSRF Token Implementation
```go
// Simple CSRF token generation and validation
import (
    "crypto/rand"
    "encoding/base64"
    "github.com/gin-gonic/gin"
    "time"
)

// Generate CSRF token
func generateCSRFToken() (string, error) {
    bytes := make([]byte, 32)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

// Gin middleware for CSRF protection
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip CSRF for GET, HEAD, OPTIONS (safe methods)
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
            c.Next()
            return
        }

        // Get token from header or form
        token := c.GetHeader("X-CSRF-Token")
        if token == "" {
            token = c.PostForm("csrf_token")
        }

        // Get expected token from session/cookie
        expectedToken, err := c.Cookie("csrf_token")
        if err != nil || expectedToken == "" {
            c.JSON(403, gin.H{"error": "CSRF token missing from session"})
            c.Abort()
            return
        }

        // Compare tokens
        if token != expectedToken {
            c.JSON(403, gin.H{"error": "Invalid CSRF token"})
            c.Abort()
            return
        }

        c.Next()
    }
}

// Handler to provide CSRF token to client
func GetCSRFToken(c *gin.Context) {
    token, err := generateCSRFToken()
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to generate CSRF token"})
        return
    }

    // Set token in httpOnly cookie (not accessible to JavaScript for double-submit)
    // Alternatively, return in JSON for SPAs that will include in headers
    c.SetCookie(
        "csrf_token",
        token,
        3600, // 1 hour
        "/",
        "",
        true,  // secure
        true,  // httpOnly
    )

    // Also return in JSON for JavaScript clients (for double-submit pattern)
    c.JSON(200, gin.H{"csrf_token": token})
}

// Handler requiring CSRF protection
func ChangePasswordHandler(c *gin.Context) {
    // CSRF middleware will validate token before this handler executes
    var req struct {
        NewPassword string `json:"new_password" binding:"required,min=8"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "Invalid input"})
        return
    }

    // Change password logic
    // ...

    c.JSON(200, gin.H{"message": "Password changed successfully"})
}
```

#### SameSite Cookie Configuration
```go
// Setting secure cookies with SameSite in Gin
func SetAuthCookie(c *gin.Context, token string) {
    // Set secure session cookie with SameSite=Strict
    c.SetCookie(
        "session",        // name
        token,            // value
        24*3600,          // max age (24 hours)
        "/",              // path
        ".example.com",   // domain
        true,             // secure (HTTPS only)
        true,             // httpOnly (inaccessible to JavaScript)
    )

    // Note: Gin's SetCookie doesn't have direct SameSite parameter in older versions
    // For Gin v1.7.7+, you can use c.SetSameSite(http.SameSiteStrictMode)
}

// Using http package directly for full control
import "net/http"

func SetSecureCookie(w http.ResponseWriter, name, value string) {
    cookie := &http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        Domain:   ".example.com",
        MaxAge:   3600,
        Secure:   true,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
    }
    http.SetCookie(w, cookie)
}

// Gin v1.7.7+ with SameSite support
func SetSecureCookieGin(c *gin.Context, name, value string) {
    cookie := &http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        Domain:   ".example.com",
        MaxAge:   3600,
        Secure:   true,
        HttpOnly: true,
        SameSite: http.SameSiteStrictMode,
    }
    http.SetCookie(c.Writer, cookie)
}
```

#### Double-Submit Cookie Pattern
```go
// Alternative: Double-submit cookie pattern for SPAs
func SetupCSRFDoubleSubmit() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Skip safe methods
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
            // Generate and set CSRF token in cookie
            token, _ := generateCSRFToken()
            c.SetCookie(
                "csrf_token",
                token,
                3600,
                "/",
                "",
                true,
                false, // httpOnly=false so JavaScript can read it
            )
            c.Next()
            return
        }

        // For state-changing methods, check double-submit
        headerToken := c.GetHeader("X-CSRF-Token")
        cookieToken, _ := c.Cookie("csrf_token")

        if headerToken == "" || cookieToken == "" || headerToken != cookieToken {
            c.JSON(403, gin.H{"error": "Invalid CSRF token"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

#### Verification Steps
- [ ] CSRF tokens required for all state-changing operations (POST, PUT, DELETE, PATCH)
- [ ] SameSite=Strict or SameSite=Lax attribute set on all cookies
- [ ] Double-submit cookie pattern implemented for Single Page Applications (SPAs)
- [ ] CSRF tokens cryptographically random and unpredictable
- [ ] CSRF tokens expire appropriately (session-based or time-based)
- [ ] CSRF protection exempts safe methods (GET, HEAD, OPTIONS)
- [ ] Token validation compares against session-stored value (not just presence check)
- [ ] CSRF tokens regenerated after login to prevent session fixation

### 7. Rate Limiting

#### Basic Rate Limiting Middleware
```go
// Simple in-memory rate limiter for Go
import (
    "github.com/gin-gonic/gin"
    "sync"
    "time"
)

type RateLimiter struct {
    requests map[string][]time.Time
    mu       sync.RWMutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

func (rl *RateLimiter) Allow(identifier string) bool {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    now := time.Now()
    cutoff := now.Add(-rl.window)

    // Clean old requests
    requests := rl.requests[identifier]
    var validRequests []time.Time
    for _, t := range requests {
        if t.After(cutoff) {
            validRequests = append(validRequests, t)
        }
    }

    // Check if limit exceeded
    if len(validRequests) >= rl.limit {
        return false
    }

    // Add current request
    validRequests = append(validRequests, now)
    rl.requests[identifier] = validRequests

    return true
}

// Gin middleware using the rate limiter
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Use IP address as identifier (add X-Forwarded-For support for proxies)
        identifier := c.ClientIP()

        if !limiter.Allow(identifier) {
            c.Header("Retry-After", limiter.window.String())
            c.JSON(429, gin.H{
                "error": "Too many requests",
                "retry_after": limiter.window.Seconds(),
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// Usage in Gin router
func main() {
    r := gin.Default()

    // Global rate limiter: 100 requests per 15 minutes per IP
    globalLimiter := NewRateLimiter(100, 15*time.Minute)
    r.Use(RateLimitMiddleware(globalLimiter))

    // Stricter rate limiter for login endpoints
    loginLimiter := NewRateLimiter(10, 5*time.Minute)
    r.POST("/login", RateLimitMiddleware(loginLimiter), LoginHandler)

    r.Run(":8080")
}
```

#### Using ulule/limiter Package (Recommended)
```go
// Production-ready rate limiting with ulule/limiter
import (
    "github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/store/memory"
    limitergin "github.com/ulule/limiter/v3/drivers/middleware/gin"
    "github.com/gin-gonic/gin"
    "time"
)

func setupRateLimiting() {
    // Define rate limits
    globalRate := limiter.Rate{
        Period: 15 * time.Minute,
        Limit:  100,
    }

    loginRate := limiter.Rate{
        Period: 5 * time.Minute,
        Limit:  10,
    }

    searchRate := limiter.Rate{
        Period: 1 * time.Minute,
        Limit:  30,
    }

    // Create stores (memory for example, use Redis in production)
    globalStore := memory.NewStore()
    loginStore := memory.NewStore()
    searchStore := memory.NewStore()

    // Create limiters
    globalLimiter := limiter.New(globalStore, globalRate)
    loginLimiter := limiter.New(loginStore, loginRate)
    searchLimiter := limiter.New(searchStore, searchRate)

    // Create Gin middleware
    globalMiddleware := limitergin.NewMiddleware(globalLimiter)
    loginMiddleware := limitergin.NewMiddleware(loginLimiter)
    searchMiddleware := limitergin.NewMiddleware(searchLimiter)

    // Apply middleware
    r := gin.Default()

    // Global middleware (all routes)
    r.Use(globalMiddleware)

    // Specific route with stricter limits
    r.POST("/login", loginMiddleware, LoginHandler)
    r.POST("/register", loginMiddleware, RegisterHandler)

    // Search endpoint with its own limits
    r.GET("/search", searchMiddleware, SearchHandler)
}
```

#### User-Based Rate Limiting
```go
// Rate limiting based on user ID (authenticated users)
func UserRateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        var identifier string

        // Try to get user ID from context (if authenticated)
        userID, exists := c.Get("userID")
        if exists {
            // Use user ID for authenticated users
            identifier = fmt.Sprintf("user:%v", userID)
        } else {
            // Fall back to IP for anonymous users
            identifier = fmt.Sprintf("ip:%s", c.ClientIP())
        }

        if !limiter.Allow(identifier) {
            c.Header("Retry-After", limiter.window.String())
            c.JSON(429, gin.H{
                "error": "Too many requests",
                "retry_after": limiter.window.Seconds(),
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

// Different limits for different user roles
func RoleBasedRateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userRole, exists := c.Get("userRole")

        var limiter *RateLimiter
        if exists && userRole == "admin" {
            // Admins have higher limits
            limiter = NewRateLimiter(1000, time.Minute)
        } else if exists {
            // Authenticated users
            limiter = NewRateLimiter(100, time.Minute)
        } else {
            // Anonymous users
            limiter = NewRateLimiter(10, time.Minute)
        }

        // Apply the appropriate limiter
        var identifier string
        if userID, exists := c.Get("userID"); exists {
            identifier = fmt.Sprintf("user:%v", userID)
        } else {
            identifier = c.ClientIP()
        }

        if !limiter.Allow(identifier) {
            c.Header("Retry-After", limiter.window.String())
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

#### Verification Steps
- [ ] Rate limiting implemented on all API endpoints
- [ ] Stricter limits applied to expensive operations (search, file processing, etc.)
- [ ] IP-based rate limiting for anonymous access
- [ ] User-based rate limiting for authenticated users (prevents single user abuse)
- [ ] Role-based rate limiting (different limits for users, admins, etc.)
- [ ] Appropriate Retry-After headers in 429 responses
- [ ] Rate limit information in response headers (X-RateLimit-Limit, X-RateLimit-Remaining, X-RateLimit-Reset)
- [ ] Rate limiters use appropriate storage (memory for single instance, Redis for distributed)
- [ ] Login endpoints have strict rate limits to prevent brute force attacks
- [ ] Rate limit configuration adjustable per environment (higher limits in development)

### 8. Sensitive Data Exposure

#### Logging
```go
// ❌ WRONG: Logging sensitive data in Go
log.Printf("User login: email=%s, password=%s", email, password)
slog.Info("Payment processing", "card_number", cardNumber, "cvv", cvv)

// ✅ CORRECT: Redact sensitive data in Go logs
log.Printf("User login attempt: email=%s", email)
slog.Info("Payment processed", "user_id", userID, "last4", last4Digits)

// Using structured logging with Go's log/slog package
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("User authenticated",
    "user_id", user.ID,
    "email_hash", hashEmail(user.Email), // Hash instead of raw email
    "action", "login",
)

// Example function to hash sensitive data for logging
func hashEmail(email string) string {
    h := sha256.New()
    h.Write([]byte(email))
    return hex.EncodeToString(h.Sum(nil))[:8] // First 8 chars of hash
}
```

#### Error Messages
```go
// ❌ WRONG: Exposing internal details in Go
func getUserHandler(c *gin.Context) {
    userID := c.Param("id")
    user, err := service.GetUser(userID)
    if err != nil {
        // Exposing database errors and internal details
        c.JSON(500, gin.H{
            "error": fmt.Sprintf("Database error: %v", err),
            "query": "SELECT * FROM users WHERE id = ?",
        })
        return
    }
    c.JSON(200, user)
}

// ✅ CORRECT: Generic error messages in Go
func getUserHandlerSecure(c *gin.Context) {
    userID := c.Param("id")
    user, err := service.GetUser(userID)
    if err != nil {
        // Log detailed error internally
        slog.Error("Failed to get user",
            "user_id", userID,
            "error", err,
            "trace_id", getTraceID(c),
        )

        // Return generic error to client
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(404, gin.H{"error": "用户不存在"})
        } else {
            c.JSON(500, gin.H{"error": "内部服务器错误"})
        }
        return
    }
    c.JSON(200, user)
}

// Custom error types for better control
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Err     error  `json:"-"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewAppError(code int, message string, err error) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Err:     err,
    }
}

// Usage in handler
func processPayment(c *gin.Context) {
    if err := service.ProcessPayment(); err != nil {
        var appErr *AppError
        if errors.As(err, &appErr) {
            // Log internally with full details
            slog.Error("Payment processing failed",
                "error", appErr.Err,
                "user_id", getUserID(c),
                "payment_id", getPaymentID(c),
            )
            // Return sanitized error to client
            c.JSON(appErr.Code, gin.H{"error": appErr.Message})
        } else {
            // Generic fallback
            slog.Error("Unknown payment error", "error", err)
            c.JSON(500, gin.H{"error": "支付处理失败"})
        }
        return
    }
    c.JSON(200, gin.H{"message": "支付成功"})
}
```

#### Verification Steps
- [ ] No passwords, tokens, API keys, or secrets in Go logs (use `log/slog` with redaction)
- [ ] Error messages generic for users (Chinese messages for this project: "内部服务器错误", "用户不存在" etc.)
- [ ] Detailed errors logged internally with structured logging (`slog.Error` with context)
- [ ] No stack traces exposed to users in HTTP responses
- [ ] Sensitive data hashed or masked in logs (emails, IDs, etc.)
- [ ] Use custom error types (`AppError`) to separate internal errors from user messages
- [ ] Database query details not exposed in error responses
- [ ] HTTP status codes appropriate (404 for not found, 500 for internal errors)

### 9. Cryptography and Key Management

#### Secure Key Storage and Usage
```go
// ❌ WRONG: Hardcoded cryptographic keys in source code
const encryptionKey = "my-secret-key-12345"

// ✅ CORRECT: Load keys from secure sources
import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "io"
)

// Load encryption key from environment or secure key management service
func getEncryptionKey() ([]byte, error) {
    keyBase64 := os.Getenv("ENCRYPTION_KEY")
    if keyBase64 == "" {
        return nil, fmt.Errorf("ENCRYPTION_KEY environment variable not set")
    }
    key, err := base64.StdEncoding.DecodeString(keyBase64)
    if err != nil {
        return nil, fmt.Errorf("failed to decode encryption key: %w", err)
    }
    // Ensure key is correct length for AES-256
    if len(key) != 32 {
        return nil, fmt.Errorf("encryption key must be 32 bytes for AES-256")
    }
    return key, nil
}

// Secure encryption/decryption using AES-GCM
func encryptData(plaintext []byte) ([]byte, error) {
    key, err := getEncryptionKey()
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

func decryptData(ciphertext []byte) ([]byte, error) {
    key, err := getEncryptionKey()
    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}
```

#### Password Hashing with Argon2
```go
// ❌ WRONG: Weak password hashing (MD5, SHA-1, SHA-256 without salt)
import "crypto/sha256"
func hashPasswordWeak(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}

// ✅ CORRECT: Use Argon2 for password hashing
import "golang.org/x/crypto/argon2"
import "golang.org/x/crypto/bcrypt"

// Using bcrypt (simpler, built-in)
func hashPasswordBcrypt(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hashedPassword), nil
}

func verifyPasswordBcrypt(hashedPassword, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil
}

// Using Argon2 (more configurable, recommended for new applications)
type Argon2Params struct {
    Memory      uint32
    Iterations  uint32
    Parallelism uint8
    SaltLength  uint32
    KeyLength   uint32
}

var DefaultArgon2Params = &Argon2Params{
    Memory:      64 * 1024, // 64 MB
    Iterations:  3,
    Parallelism: 2,
    SaltLength:  16,
    KeyLength:   32,
}

func hashPasswordArgon2(password string, p *Argon2Params) (string, error) {
    salt := make([]byte, p.SaltLength)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

    // Encode hash and salt together
    encodedHash := base64.RawStdEncoding.EncodeToString(hash)
    encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

    return fmt.Sprintf("%s:%s", encodedHash, encodedSalt), nil
}

func verifyPasswordArgon2(password, encodedHash string, p *Argon2Params) (bool, error) {
    parts := strings.Split(encodedHash, ":")
    if len(parts) != 2 {
        return false, fmt.Errorf("invalid encoded hash format")
    }

    hash, err := base64.RawStdEncoding.DecodeString(parts[0])
    if err != nil {
        return false, err
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[1])
    if err != nil {
        return false, err
    }

    computedHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

    // Constant-time comparison to prevent timing attacks
    return subtle.ConstantTimeCompare(hash, computedHash) == 1, nil
}
```

#### Digital Signatures and Verification
```go
// Using Ed25519 for digital signatures
import (
    "crypto/ed25519"
    "crypto/rand"
)

func generateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
    publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        return nil, nil, err
    }
    return publicKey, privateKey, nil
}

func signMessage(privateKey ed25519.PrivateKey, message []byte) []byte {
    signature := ed25519.Sign(privateKey, message)
    return signature
}

func verifySignature(publicKey ed25519.PublicKey, message, signature []byte) bool {
    return ed25519.Verify(publicKey, message, signature)
}

// Example: API request signing
func signAPIRequest(privateKey ed25519.PrivateKey, method, path, body string, timestamp int64) (string, error) {
    message := fmt.Sprintf("%s\n%s\n%s\n%d", method, path, body, timestamp)
    signature := signMessage(privateKey, []byte(message))
    return base64.StdEncoding.EncodeToString(signature), nil
}

func verifyAPIRequest(publicKey ed25519.PublicKey, method, path, body string, timestamp int64, signatureBase64 string) bool {
    signature, err := base64.StdEncoding.DecodeString(signatureBase64)
    if err != nil {
        return false
    }
    message := fmt.Sprintf("%s\n%s\n%s\n%d", method, path, body, timestamp)
    return verifySignature(publicKey, []byte(message), signature)
}
```

#### Verification Steps
- [ ] Cryptographic keys stored securely (environment variables, secret managers)
- [ ] Passwords hashed with strong algorithms (Argon2, bcrypt with appropriate cost)
- [ ] Encryption uses authenticated modes (AES-GCM, not ECB or CBC without HMAC)
- [ ] Random values generated with cryptographically secure RNG (`crypto/rand`)
- [ ] Digital signatures used for critical operations (API requests, data integrity)
- [ ] Key rotation policies implemented for long-lived keys
- [ ] No hardcoded keys, passwords, or secrets in source code
- [ ] Encryption keys properly sized (32 bytes for AES-256, 64 bytes for HMAC-SHA512)
- [ ] Salt values unique per password and generated with secure RNG
- [ ] Timing attack protection for comparisons (`subtle.ConstantTimeCompare`)

### 10. Go Dependency Security

#### Regular Updates and Vulnerability Scanning
```bash
# Check for available updates in Go modules
go list -m -u all

# Update all dependencies to latest versions
go get -u ./...

# Update specific dependency
go get -u github.com/gin-gonic/gin

# Check for known vulnerabilities using govulncheck (Go 1.18+)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Alternative: Use gosec for static security analysis
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...

# Check for outdated dependencies
go list -m -f '{{if not .Indirect}}{{.Path}} {{.Version}}{{end}}' all

# Tidy go.mod file
go mod tidy

# Verify dependencies
go mod verify
```

#### go.mod and go.sum Files
```bash
# ALWAYS commit both go.mod and go.sum files
git add go.mod go.sum

# Use vendor directory for reproducible builds (optional)
go mod vendor

# In CI/CD, use vendored dependencies or download fresh
go mod download
go build -mod=readonly ./...

# For production builds, use vendored dependencies
go build -mod=vendor ./...

# Verify checksums match go.sum
go mod verify
```

#### Dependency Analysis Tools
```go
// Example: Programmatically check module versions
package main

import (
    "fmt"
    "golang.org/x/mod/modfile"
    "io/ioutil"
    "os"
)

func checkDependencies() error {
    // Read go.mod file
    data, err := ioutil.ReadFile("go.mod")
    if err != nil {
        return fmt.Errorf("failed to read go.mod: %w", err)
    }

    f, err := modfile.Parse("go.mod", data, nil)
    if err != nil {
        return fmt.Errorf("failed to parse go.mod: %w", err)
    }

    // Analyze dependencies
    fmt.Println("Direct dependencies:")
    for _, req := range f.Require {
        if !req.Indirect {
            fmt.Printf("  %s %s\n", req.Mod.Path, req.Mod.Version)
        }
    }

    fmt.Println("\nIndirect dependencies:")
    for _, req := range f.Require {
        if req.Indirect {
            fmt.Printf("  %s %s\n", req.Mod.Path, req.Mod.Version)
        }
    }

    return nil
}

// Using go/packages to analyze dependencies in code
import "golang.org/x/tools/go/packages"

func analyzePackageImports() error {
    cfg := &packages.Config{
        Mode: packages.NeedImports | packages.NeedDeps,
    }
    pkgs, err := packages.Load(cfg, "./...")
    if err != nil {
        return fmt.Errorf("failed to load packages: %w", err)
    }

    for _, pkg := range pkgs {
        fmt.Printf("Package: %s\n", pkg.ID)
        for _, imp := range pkg.Imports {
            fmt.Printf("  Import: %s\n", imp.ID)
        }
    }
    return nil
}
```

#### Security Scanning in CI/CD
```yaml
# Example GitHub Actions workflow for Go security scanning
name: Go Security Scan

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Run govulncheck
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

      - name: Run gosec
        run: |
          go install github.com/securego/gosec/v2/cmd/gosec@latest
          gosec ./...

      - name: Check for outdated dependencies
        run: |
          go list -m -u all | grep -E "\[.*\]"

      - name: Run go mod tidy and verify
        run: |
          go mod tidy
          go mod verify
          git diff --exit-code go.mod go.sum
```

#### Verification Steps
- [ ] Go modules up to date (`go list -m -u all` shows no updates available or updates applied)
- [ ] No known vulnerabilities (`govulncheck ./...` reports no issues)
- [ ] Static security analysis clean (`gosec ./...` passes)
- [ ] Both `go.mod` and `go.sum` files committed to version control
- [ ] `go mod tidy` run and changes committed
- [ ] `go mod verify` passes (checksums valid)
- [ ] Dependabot or Renovate configured for Go modules on GitHub
- [ ] Regular security updates applied (at least monthly)
- [ ] Indirect dependencies reviewed for security implications
- [ ] Minimum Go version specified in go.mod (`go 1.21`)
- [ ] CI/CD pipeline includes security scanning steps
- [ ] Vendor directory used for production builds if needed for reproducibility
- [ ] Dependency licenses reviewed for compliance
- [ ] No dependencies on unmaintained or suspicious packages

## Security Testing for Go Applications

### Automated Security Tests with Go Testing
```go
// Test authentication middleware
func TestAuthMiddleware(t *testing.T) {
    // Setup test router
    r := gin.New()
    r.Use(AuthMiddleware())
    r.GET("/protected", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "protected"})
    })

    // Test without authentication
    req := httptest.NewRequest("GET", "/protected", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    assert.Contains(t, w.Body.String(), "Authorization required")

    // Test with invalid token
    req = httptest.NewRequest("GET", "/protected", nil)
    req.Header.Set("Authorization", "Bearer invalid-token")
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    assert.Contains(t, w.Body.String(), "Invalid or expired token")
}

// Test authorization (role-based access)
func TestRequireAdminRole(t *testing.T) {
    r := gin.New()

    // Mock user with non-admin role
    r.GET("/admin", func(c *gin.Context) {
        c.Set("userRole", "user") // Non-admin role
        RequireRole("admin")(c)
        if c.IsAborted() {
            return
        }
        c.JSON(200, gin.H{"message": "admin access"})
    })

    req := httptest.NewRequest("GET", "/admin", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusForbidden, w.Code)
    assert.Contains(t, w.Body.String(), "Insufficient permissions")
}

// Test input validation
func TestInputValidation(t *testing.T) {
    r := gin.New()
    r.POST("/users", CreateUserHandler)

    // Test invalid email
    jsonStr := `{"email": "not-an-email", "name": "test", "age": 25, "password": "password123"}`
    req := httptest.NewRequest("POST", "/users", strings.NewReader(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    assert.Contains(t, w.Body.String(), "Invalid input")

    // Test missing required field
    jsonStr = `{"email": "test@example.com", "age": 25}`
    req = httptest.NewRequest("POST", "/users", strings.NewReader(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test SQL injection prevention
func TestSQLInjectionPrevention(t *testing.T) {
    // Setup test database
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    // Expect parameterized query
    mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\?").
        WithArgs("test@example.com").
        WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

    // Call function that should use parameterized query
    _, err = getUserByEmailSafe(db, "test@example.com")
    assert.NoError(t, err)

    // Verify all expectations were met
    assert.NoError(t, mock.ExpectationsWereMet())
}

// Test rate limiting
func TestRateLimiting(t *testing.T) {
    limiter := NewRateLimiter(5, time.Minute) // 5 requests per minute
    identifier := "test-ip"

    // First 5 requests should succeed
    for i := 0; i < 5; i++ {
        assert.True(t, limiter.Allow(identifier), fmt.Sprintf("Request %d should be allowed", i+1))
    }

    // 6th request should be blocked
    assert.False(t, limiter.Allow(identifier), "6th request should be blocked")

    // Test with different identifiers
    assert.True(t, limiter.Allow("different-ip"), "Different IP should be allowed")
}

// Test CSRF protection
func TestCSRFProtection(t *testing.T) {
    r := gin.New()
    r.Use(CSRFMiddleware())
    r.POST("/change-password", ChangePasswordHandler)

    // Test without CSRF token
    req := httptest.NewRequest("POST", "/change-password", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusForbidden, w.Code)
    assert.Contains(t, w.Body.String(), "CSRF token")

    // Test with invalid CSRF token (setup would require session/cookie mocking)
    // This is simplified - actual test would need more setup
}

// Test XSS prevention in templates
func TestXSSPrevention(t *testing.T) {
    tmpl := `{{.UserInput}}`
    tpl, err := template.New("test").Parse(tmpl)
    require.NoError(t, err)

    var buf bytes.Buffer
    data := map[string]interface{}{
        "UserInput": `<script>alert('xss')</script>`,
    }

    err = tpl.Execute(&buf, data)
    require.NoError(t, err)

    // html/template should escape the script tags
    output := buf.String()
    assert.Contains(t, output, "&lt;script&gt;")
    assert.Contains(t, output, "&lt;/script&gt;")
    assert.NotContains(t, output, "<script>")
}

// Test file upload validation
func TestFileUploadValidation(t *testing.T) {
    // Test valid file
    fileHeader := &multipart.FileHeader{
        Filename: "test.jpg",
        Size:     1024 * 1024, // 1MB
        Header:   map[string][]string{"Content-Type": {"image/jpeg"}},
    }

    err := validateFileUpload(fileHeader)
    assert.NoError(t, err)

    // Test file too large
    fileHeader.Size = 10 * 1024 * 1024 // 10MB
    err = validateFileUpload(fileHeader)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "too large")

    // Test invalid file type
    fileHeader.Size = 1024 * 1024
    fileHeader.Header["Content-Type"] = []string{"application/exe"}
    err = validateFileUpload(fileHeader)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid file type")

    // Test path traversal in filename
    fileHeader.Header["Content-Type"] = []string{"image/jpeg"}
    fileHeader.Filename = "../../etc/passwd"
    err = validateFileUpload(fileHeader)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid filename")
}

// Test password hashing
func TestPasswordHashing(t *testing.T) {
    password := "mySecurePassword123"

    // Test bcrypt
    hashed, err := hashPasswordBcrypt(password)
    require.NoError(t, err)
    assert.True(t, verifyPasswordBcrypt(hashed, password))
    assert.False(t, verifyPasswordBcrypt(hashed, "wrongpassword"))

    // Test Argon2
    params := &Argon2Params{
        Memory:      64 * 1024,
        Iterations:  3,
        Parallelism: 2,
        SaltLength:  16,
        KeyLength:   32,
    }
    hashed, err = hashPasswordArgon2(password, params)
    require.NoError(t, err)
    valid, err := verifyPasswordArgon2(password, hashed, params)
    require.NoError(t, err)
    assert.True(t, valid)
}
```

## Pre-Deployment Security Checklist for Go Applications

Before ANY production deployment of Go web applications:

### Authentication & Authorization
- [ ] **JWT Tokens**: Stored in httpOnly, Secure, SameSite cookies (not localStorage or response JSON)
- [ ] **Token Validation**: JWT signatures verified, expiration checked, issuer validated
- [ ] **Role-Based Access**: Authorization checks for all sensitive operations
- [ ] **Row-Level Security**: Application-level data access control (GORM scopes, repository patterns)

### Input & Data Security
- [ ] **Input Validation**: All user inputs validated with struct tags (`binding:"required,email"`) or custom validation
- [ ] **SQL Injection Prevention**: All database queries use parameterized queries (GORM `Where("email = ?", email)`, never string concatenation)
- [ ] **XSS Prevention**: `html/template` used for HTML rendering (auto-escapes), user HTML sanitized with bluemonday
- [ ] **File Uploads**: Validated for size, content type, extension; path traversal prevented

### API & Network Security
- [ ] **HTTPS Enforcement**: All production traffic forced to HTTPS (HSTS headers, redirect HTTP → HTTPS)
- [ ] **CORS Configuration**: Properly configured for frontend origins, no wildcard (`*`) in production
- [ ] **Rate Limiting**: Enabled on all endpoints (IP-based, user-based, stricter limits for login/search)
- [ ] **CSRF Protection**: Tokens required for state-changing operations, SameSite cookies configured
- [ ] **Security Headers**: Content-Security-Policy, X-Frame-Options, X-Content-Type-Options, X-XSS-Protection

### Application Security
- [ ] **Error Handling**: Generic error messages to users, detailed errors logged internally (no stack traces exposed)
- [ ] **Logging Security**: No passwords, tokens, or secrets in logs; sensitive data hashed/masked
- [ ] **Session Management**: Appropriate token expiration, secure session storage
- [ ] **Password Security**: Strong hashing (Argon2 or bcrypt with appropriate cost), no plaintext storage

### Dependency & Configuration Security
- [ ] **Go Modules**: `go.mod` and `go.sum` committed, `go mod verify` passes
- [ ] **Vulnerability Scanning**: `govulncheck ./...` reports no known vulnerabilities
- [ ] **Static Analysis**: `gosec ./...` passes security checks
- [ ] **Dependency Updates**: Dependencies up to date (`go list -m -u all` reviewed)
- [ ] **Secrets Management**: No hardcoded secrets; all in environment variables or secret managers
- [ ] **Configuration**: Sensitive config loaded from env vars (`caarlos0/env` or similar), config files excluded from git

### Infrastructure & Deployment Security
- [ ] **Database Security**: PostgreSQL connection uses SSL/TLS, strong passwords, limited privileges
- [ ] **Container Security**: Docker images built from minimal base images, no root user, secrets not in layers
- [ ] **Network Security**: Firewall rules limit access, database not publicly accessible
- [ ] **Monitoring & Alerting**: Log aggregation, security event monitoring, intrusion detection

### Go-Specific Security
- [ ] **Memory Safety**: No use of `unsafe` package unless absolutely necessary and reviewed
- [ ] **Concurrency Security**: Race condition prevention (mutexes, channels), deadlock avoidance
- [ ] **Buffer Security**: Input size limits to prevent memory exhaustion
- [ ] **TLS Configuration**: Modern TLS versions (1.2+), secure cipher suites
- [ ] **Cryptography**: Use standard library crypto packages, not custom implementations

### Compliance & Verification
- [ ] **Security Testing**: Unit tests for security features (auth, validation, rate limiting)
- [ ] **Code Review**: Security-focused code review completed
- [ ] **Penetration Testing**: External security assessment considered for critical applications
- [ ] **Incident Response**: Process defined for security incident handling

## Go Security Resources

### Official Documentation
- [Go Security Best Practices](https://go.dev/security/) - Official Go security guide
- [Go Cryptography Package](https://pkg.go.dev/crypto) - Standard library crypto packages
- [Go Vulnerability Database](https://vuln.go.dev/) - Go vulnerability database

### OWASP & Web Security
- [OWASP Top 10](https://owasp.org/www-project-top-ten/) - Top web application security risks
- [OWASP Go Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_Security_Cheat_Sheet.html)
- [OWASP Web Security Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)

### Go Security Tools
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) - Go vulnerability checker
- [gosec](https://github.com/securego/gosec) - Go security checker
- [staticcheck](https://staticcheck.io/) - Advanced Go linter with security checks
- [golangci-lint](https://golangci-lint.run/) - Fast linters runner with security rules

### Gin Framework Security
- [Gin Security Best Practices](https://gin-gonic.com/docs/examples/security/)
- [Gin JWT Middleware](https://github.com/appleboy/gin-jwt) - JWT middleware for Gin
- [Gin CORS Middleware](https://github.com/gin-contrib/cors) - CORS middleware for Gin

### Database Security
- [GORM Security](https://gorm.io/docs/security.html) - GORM security practices
- [PostgreSQL Security](https://www.postgresql.org/docs/current/security.html) - PostgreSQL security documentation

### Cryptography Resources
- [Go crypto/rand vs math/rand](https://go.dev/blog/cryptorand) - Cryptographic randomness in Go
- [Go Subtle Package](https://pkg.go.dev/crypto/subtle) - Constant-time comparison utilities
- [Argon2 in Go](https://pkg.go.dev/golang.org/x/crypto/argon2) - Password hashing with Argon2

### Books & Articles
- ["Secure Programming with Go"](https://leanpub.com/secureprogrammingwithgo) - Book on Go security
- [Go Security Pitfalls](https://github.com/kkHAIKE/awesome-golang-security) - Collection of Go security resources
- [Practical Go Security](https://pragprog.com/titles/d-apgs/practical-go-security/) - Book on practical Go security

---

**Remember**: Security is not optional. One vulnerability can compromise the entire platform. When in doubt, err on the side of caution.
