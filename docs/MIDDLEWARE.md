# Middleware Documentation

This document covers all middleware implemented in this project.

## Table of Contents
- [Overview](#overview)
- [Authentication Middleware](#authentication-middleware)
- [Logging Middleware](#logging-middleware)
- [Validation Middleware](#validation-middleware)
- [Tracing Middleware](#tracing-middleware)
- [Middleware Stack](#middleware-stack)
- [Usage Examples](#usage-examples)

## Overview

Middleware in Echo processes requests before they reach your handlers. This project includes four essential middleware:

1. **Authentication** - JWT-based authentication
2. **Logging** - Structured request/response logging
3. **Validation** - Input validation utilities
4. **Tracing** - Distributed request tracing with OpenTelemetry + Jaeger

## Authentication Middleware

### Location
`internal/handlers/middleware/auth.go`

### Configuration
```go
jwtConfig := middleware.JWTConfig(secret)
```

### Features
- **JWT Token Validation**: Validates incoming JWT tokens
- **Custom Claims**: Supports custom user claims (user_id, email, roles)
- **Public Routes**: Skips authentication for public endpoints
- **Error Handling**: Returns 401 for invalid tokens

### Public Routes (No Auth Required)
```
/
/healthcheck
/swagger/index
/auth/login
/auth/register
```

### Usage in Handler

```go
// Register protected routes under /api
authGroup := e.Group("/api")
authGroup.Use(echojwt.WithConfig(jwtConfig))

authGroup.GET("/users", GetUsersHandler)
authGroup.POST("/users", CreateUserHandler)
```

### Accessing User Info in Controller

```go
func (c *UserController) GetProfile(ctx echo.Context) error {
    userID := middleware.GetUserID(ctx)
    email := middleware.GetEmail(ctx)
    roles := middleware.GetRoles(ctx)
    
    // Use user info...
    return ctx.JSON(200, map[string]interface{}{
        "user_id": userID,
        "email": email,
        "roles": roles,
    })
}
```

### Generating a JWT Token

```go
func generateToken(userID int, email string, roles []string) (string, error) {
    claims := middleware.CustomClaims{
        UserID: userID,
        Email:  email,
        Roles:  roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

## Logging Middleware

### Location
`internal/handlers/middleware/logging.go`

### Features
- **Structured Logging**: JSON-formatted logs for easy parsing
- **Request Context**: Logs method, path, IP, user agent
- **Response Info**: Logs status code, response size, duration
- **Error Tracking**: Logs errors with context
- **Built-in slog**: Uses Go's standard structured logging library

### What Gets Logged

**Incoming Request:**
```json
{
  "time": "2026-02-21T10:00:00Z",
  "level": "INFO",
  "msg": "incoming request",
  "method": "POST",
  "path": "/api/users",
  "ip": "127.0.0.1",
  "user_agent": "curl/7.68.0"
}
```

**Response:**
```json
{
  "time": "2026-02-21T10:00:00.250Z",
  "level": "INFO",
  "msg": "request completed",
  "method": "POST",
  "path": "/api/users",
  "status": 201,
  "duration": "250ms",
  "bytes_sent": 256
}
```

### Configuration

The middleware uses Go's standard `log/slog` library:

```go
// Setup in main.go (optional, default is text output)
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})
slog.SetDefault(slog.New(handler))
```

### Using Logger in Your Code

```go
import "log/slog"

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    slog.InfoContext(ctx, "creating user",
        slog.String("email", user.Email),
    )
    
    if err := s.repo.Create(user); err != nil {
        slog.ErrorContext(ctx, "failed to create user",
            slog.String("email", user.Email),
            slog.String("error", err.Error()),
        )
        return err
    }
    
    return nil
}
```

## Validation Middleware

### Location
`internal/handlers/middleware/validation.go`

### Features
- **Struct Validation**: Validates struct fields against tags
- **Custom Rules**: Support for custom validation rules
- **Error Formatting**: Formats validation errors for API responses
- **Built-in validator**: Uses `go-playground/validator/v10`

### Usage in DTO

```go
type CreateUserDTO struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,max=50"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Age      int    `json:"age" validate:"omitempty,min=18,max=120"`
}
```

### Validation in Controller

```go
func (c *UserController) Create(ctx echo.Context) error {
    var dto CreateUserDTO
    
    // Bind request to DTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.JSON(400, map[string]interface{}{
            "error": "invalid request",
        })
    }
    
    // Validate DTO
    if err := middleware.ValidateStruct(dto); err != nil {
        return ctx.JSON(422, map[string]interface{}{
            "errors": middleware.FormatValidationErrors(err),
        })
    }
    
    // Process valid data...
    return c.service.CreateUser(dto)
}
```

### Validation Tags

Common validation tags:

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field must not be empty | `validate:"required"` |
| `email` | Valid email format | `validate:"email"` |
| `min` | Minimum length/value | `validate:"min=8"` |
| `max` | Maximum length/value | `validate:"max=100"` |
| `len` | Exact length | `validate:"len=10"` |
| `numeric` | Must be numeric | `validate:"numeric"` |
| `url` | Valid URL | `validate:"url"` |
| `omitempty` | Field is optional | `validate:"omitempty,min=18"` |
| `oneof` | One of specific values | `validate:"oneof=admin user"` |

## Tracing Middleware

### Location
`internal/handlers/middleware/tracing.go`

### Features
- **Distributed Tracing**: Track requests across services
- **OpenTelemetry**: Industry standard tracing library
- **Jaeger Integration**: Send traces to Jaeger backend
- **Span Recording**: Automatic span creation for each request
- **Error Tracking**: Records errors in spans

### What Gets Traced

Each HTTP request creates a span with:
- **Operation name**: Request path (e.g., "/api/users")
- **Attributes**: Method, URL, status code, duration
- **Errors**: Any errors are recorded
- **Timestamps**: Start and end times

### Example Trace

```
Trace: POST /api/users (completed in 250ms)
├─ POST /api/users (250ms)
│  ├─ database.create_user (45ms)
│  └─ email.send (150ms)
└─ Response sent (50ms)
```

### Creating Child Spans

```go
func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // Create child span for database operation
    ctx, span := middleware.CreateChildSpan(ctx, "database.create_user")
    defer span.End()
    
    // Database operation...
    err := s.repo.Create(ctx, user)
    
    if err != nil {
        span.RecordError(err)
    }
    
    return err
}
```

### Viewing Traces

1. **Start Jaeger**:
   ```bash
   docker-compose up -d jaeger
   ```

2. **Open Jaeger UI**:
   ```
   http://localhost:16686
   ```

3. **View Your Traces**:
   - Select "ddd-golang-template" from service dropdown
   - Click on traces to view details
   - Analyze duration and dependencies

## Middleware Stack

The middleware are registered in this order in `internal/handlers/middleware.go`:

```
1. Recovery Middleware
   ↓
2. Request ID Middleware
   ↓
3. Structured Logging Middleware
   ↓
4. CORS Middleware
   ↓
5. Tracing Middleware
   ↓
6. JWT Middleware (for /api routes only)
   ↓
7. Your Handler
```

### Middleware Order Matters

- **Recovery** first: Catches panics before anything else
- **Logging** early: Logs all requests
- **Tracing** before handlers: Traces all operations
- **JWT** on protected routes: Validates auth on protected endpoints

## Usage Examples

### Example 1: Protected Endpoint with Validation

```go
// DTO with validation
type CreateUserDTO struct {
    Email string `json:"email" validate:"required,email"`
    Name  string `json:"name" validate:"required,min=2"`
}

// Handler
func (c *UserController) Create(ctx echo.Context) error {
    // Bind and validate
    var dto CreateUserDTO
    if err := ctx.Bind(&dto); err != nil {
        return ctx.JSON(400, map[string]interface{}{
            "error": "invalid request",
        })
    }
    
    if err := middleware.ValidateStruct(dto); err != nil {
        return ctx.JSON(422, map[string]interface{}{
            "errors": middleware.FormatValidationErrors(err),
        })
    }
    
    // Use authenticated user
    userID := middleware.GetUserID(ctx)
    
    // Call service
    user, err := c.service.CreateUser(dto)
    if err != nil {
        return ctx.JSON(500, map[string]interface{}{
            "error": err.Error(),
        })
    }
    
    return ctx.JSON(201, user)
}

// Register route (protected)
authGroup := e.Group("/api")
authGroup.POST("/users", userController.Create)
```

### Example 2: Using Tracing in Service

```go
func (s *UserService) CreateUser(ctx context.Context, dto dto.CreateUserDTO) (*models.User, error) {
    // Log operation
    slog.InfoContext(ctx, "creating user", slog.String("email", dto.Email))
    
    // Check if email exists (with child span)
    ctx, span := middleware.CreateChildSpan(ctx, "database.check_email")
    exists, err := s.repo.EmailExists(ctx, dto.Email)
    span.End()
    
    if exists {
        return nil, errors.New("email already registered")
    }
    
    // Create user (with child span)
    ctx, span = middleware.CreateChildSpan(ctx, "database.create_user")
    defer span.End()
    
    user := &models.User{
        Email: dto.Email,
        Name:  dto.Name,
    }
    
    err = s.repo.Create(ctx, user)
    if err != nil {
        span.RecordError(err)
        slog.ErrorContext(ctx, "failed to create user", slog.String("error", err.Error()))
        return nil, err
    }
    
    slog.InfoContext(ctx, "user created successfully", slog.Int("user_id", int(user.ID)))
    return user, nil
}
```

### Example 3: Public Endpoint

```go
// Public endpoints don't need JWT

func (c *BaseController) Health(ctx echo.Context) error {
    return ctx.JSON(200, map[string]string{
        "status": "ok",
    })
}

// No JWT protection needed
e.GET("/healthcheck", baseController.Health)
```

## Configuration

Set these in `.env`:

```env
# JWT Secret (change in production!)
JWT_SECRET=your-super-secret-key

# Jaeger endpoint
JAEGER_ENDPOINT=http://localhost:4318

# Log level (optional)
APP_ENV=development
```

## Troubleshooting

### Traces not appearing in Jaeger?
1. Ensure Jaeger is running: `docker-compose up -d jaeger`
2. Check endpoint: `JAEGER_ENDPOINT=http://localhost:4318`
3. Verify connectivity: `curl http://localhost:16686`

### JWT not validating?
1. Check JWT secret in `.env`
2. Verify token hasn't expired
3. Ensure route is under `/api` group

### Validation not working?
1. Initialize validator: `middleware.InitValidator()`
2. Check validation tags syntax
3. Ensure DTOs are passed to `ValidateStruct()`

## Further Reading

- [Echo Middleware Documentation](https://echo.labstack.com/docs/middleware)
- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger Tracing](https://www.jaegertracing.io/docs/)
- [Validator Package](https://github.com/go-playground/validator)
- [Go slog](https://pkg.go.dev/log/slog)

