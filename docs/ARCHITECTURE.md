# Architecture Documentation

## Table of Contents
- [Overview](#overview)
- [Domain-Driven Design (DDD)](#domain-driven-design-ddd)
- [Hexagonal Architecture](#hexagonal-architecture)
- [Layer Breakdown](#layer-breakdown)
- [Data Flow](#data-flow)
- [Design Principles](#design-principles)
- [Best Practices](#best-practices)
- [Example Implementation](#example-implementation)

## Overview

This project implements a **Domain-Driven Design (DDD)** approach with **Hexagonal Architecture** (also known as Ports and Adapters). This architectural pattern ensures:

- **Separation of Concerns**: Clear boundaries between business logic and infrastructure
- **Testability**: Business logic isolated from external dependencies
- **Maintainability**: Changes in one layer don't cascade to others
- **Flexibility**: Easy to swap implementations (e.g., change databases, frameworks)

## Domain-Driven Design (DDD)

Domain-Driven Design is a software development approach that focuses on the business domain and its logic. The key concepts include:

### Core Principles

1. **Ubiquitous Language**: Use the same terminology as business stakeholders
2. **Bounded Contexts**: Clear boundaries around different parts of the domain
3. **Domain Models**: Rich models that capture business logic and rules
4. **Layers**: Separation of domain logic from infrastructure concerns

### DDD Building Blocks

- **Entities**: Objects with unique identity (e.g., User, Order)
- **Value Objects**: Immutable objects without identity (e.g., Address, Money)
- **Aggregates**: Cluster of entities treated as a single unit
- **Repositories**: Abstraction for data persistence
- **Services**: Operations that don't naturally fit in entities
- **Domain Events**: Things that happened in the domain

## Hexagonal Architecture

Hexagonal Architecture (Ports and Adapters) organizes code to separate business logic from external concerns.

### Concept

```
┌─────────────────────────────────────────────────────────┐
│                    External World                       │
│  (HTTP, CLI, Queue, Database, External APIs, etc.)     │
└─────────────────────────────────────────────────────────┘
                            │
                    ┌───────┴───────┐
                    │    Adapters    │
                    │  (Controllers, │
                    │  Repositories) │
                    └───────┬───────┘
                            │
                    ┌───────┴───────┐
                    │     Ports      │
                    │  (Interfaces)  │
                    └───────┬───────┘
                            │
                    ┌───────┴───────┐
                    │   Core/Domain  │
                    │ (Business Logic)│
                    └────────────────┘
```

### Key Concepts

- **Core/Domain**: Business logic, independent of external concerns
- **Ports**: Interfaces that define how the domain communicates
- **Adapters**: Implementations that connect the domain to the outside world
- **Primary/Driving Adapters**: Drive the application (HTTP controllers, CLI)
- **Secondary/Driven Adapters**: Driven by the application (Databases, APIs)

## Layer Breakdown

### 1. Application Layer (`internal/application/`)

**Purpose**: Handle input/output, orchestrate use cases

**Components**:
- **DTOs** (`dto/`): Data Transfer Objects for input/output
- **Controllers** (`rest/controllers/`): HTTP request handlers
- **Error Handlers** (`rest/errors/`): HTTP error formatting

**Responsibilities**:
- Validate input
- Transform DTOs to domain objects
- Call domain services
- Format responses
- Handle HTTP concerns

**Example**:
```go
// internal/application/rest/controllers/user_controller.go
func (c *UserController) CreateUser(ctx echo.Context) error {
    var dto dto.CreateUserDto
    if err := ctx.Bind(&dto); err != nil {
        return ctx.JSON(400, err)
    }
    
    // Call domain service
    user, err := c.service.CreateUser(dto)
    if err != nil {
        return ctx.JSON(500, err)
    }
    
    return ctx.JSON(201, user)
}
```

### 2. Domain Layer (`internal/domains/`)

**Purpose**: Contains core business logic and rules

**Components**:
- **Services**: Business operations and workflows
- **Domain Models**: Business entities and value objects (in `models/`)

**Responsibilities**:
- Implement business rules
- Validate business constraints
- Coordinate between repositories
- Execute business workflows
- Remain independent of infrastructure

**Example**:
```go
// internal/domains/user_service.go
func (s *UserService) CreateUser(dto dto.CreateUserDto) (*models.User, error) {
    // Business logic: check if email already exists
    if s.repo.EmailExists(dto.Email) {
        return nil, errors.New("email already registered")
    }
    
    // Business rule: username must be lowercase
    user := &models.User{
        Name:  dto.Name,
        Email: strings.ToLower(dto.Email),
    }
    
    return s.repo.Create(user)
}
```

### 3. Infrastructure Layer (`internal/infrastructure/`)

**Purpose**: Handle technical concerns and external integrations

**Components**:
- **Repositories** (`repositories/`): Database access implementations

**Responsibilities**:
- Implement data persistence
- Connect to external services
- Handle technical details (SQL, HTTP clients, etc.)
- Implement interfaces defined by the domain

**Example**:
```go
// internal/infrastructure/repositories/user_repository.go
func (r *UserRepository) Create(user *models.User) (*models.User, error) {
    result := r.db.Create(user)
    if result.Error != nil {
        return nil, result.Error
    }
    return user, nil
}

func (r *UserRepository) EmailExists(email string) bool {
    var count int64
    r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
    return count > 0
}
```

### 4. Models Layer (`internal/models/`)

**Purpose**: Define domain entities and data structures

**Components**:
- Domain entities
- Value objects
- Database models

**Example**:
```go
// internal/models/user.go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 5. Configuration Layer (`internal/config/`)

**Purpose**: Handle application and infrastructure configuration

**Components**:
- **App Config** (`app/`): Application-level settings
- **DB Config** (`db/`): Database drivers and connections
- **Migrations** (`db/migrations/`): Database schema migrations

### 6. Handlers Layer (`internal/handlers/`)

**Purpose**: Bootstrap and wire up the application

**Components**:
- Main handler initialization
- Database setup
- REST route registration
- Middleware configuration

## Data Flow

### Request Flow (Primary/Inbound)

```
1. HTTP Request
   ↓
2. Router (handlers/rest.go)
   ↓
3. Controller (application/rest/controllers/)
   - Validate input
   - Transform to DTO
   ↓
4. Domain Service (domains/)
   - Execute business logic
   - Apply business rules
   ↓
5. Repository (infrastructure/repositories/)
   - Persist/retrieve data
   ↓
6. Database
```

### Response Flow (Outbound)

```
1. Database
   ↓
2. Repository returns entity
   ↓
3. Domain Service processes result
   ↓
4. Controller transforms to response DTO
   ↓
5. HTTP Response (JSON)
```

### Complete Example Flow

```
POST /users
{
  "name": "John Doe",
  "email": "john@example.com"
}
       ↓
[Router] → /users → UserController
       ↓
[Controller] 
  • Bind JSON to CreateUserDto
  • Validate DTO
       ↓
[Domain Service]
  • Check business rules (email unique?)
  • Create User entity
  • Apply domain logic
       ↓
[Repository]
  • Save to database
  • Return User entity
       ↓
[Controller]
  • Transform entity to response
  • Return JSON
       ↓
Response 201 Created
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2026-02-21T10:00:00Z"
}
```

## Design Principles

### 1. Dependency Inversion Principle (DIP)

High-level modules (domain) should not depend on low-level modules (infrastructure). Both should depend on abstractions.

```go
// Domain service depends on interface, not implementation
type UserService struct {
    repo UserRepository  // Interface, not concrete implementation
}

// Infrastructure implements the interface
type GormUserRepository struct {
    db *gorm.DB
}

func (r *GormUserRepository) Create(user *User) error {
    // Implementation details
}
```

### 2. Separation of Concerns

Each layer has a single, well-defined responsibility:
- **Application**: Handle I/O
- **Domain**: Business logic
- **Infrastructure**: Technical details

### 3. Single Responsibility Principle (SRP)

Each component does one thing well:
- Controllers handle HTTP
- Services contain business logic
- Repositories handle persistence

### 4. Interface Segregation

Define minimal interfaces for what you need:

```go
// Good: Specific interface
type UserReader interface {
    FindByID(id uint) (*User, error)
    FindByEmail(email string) (*User, error)
}

// Avoid: Bloated interface with everything
type UserRepository interface {
    Create(*User) error
    Update(*User) error
    Delete(uint) error
    FindByID(uint) (*User, error)
    FindAll() ([]User, error)
    // ... 20 more methods
}
```

## Best Practices

### 1. Keep Domain Pure

✅ **DO**: Keep domain logic independent
```go
func (s *UserService) PromoteToAdmin(userID uint) error {
    user := s.repo.FindByID(userID)
    if user.RegistrationDate.Before(time.Now().AddDate(-1, 0, 0)) {
        user.Role = "admin"
        return s.repo.Update(user)
    }
    return errors.New("user must be registered for 1 year")
}
```

❌ **DON'T**: Mix domain logic with infrastructure
```go
func (s *UserService) PromoteToAdmin(ctx echo.Context) error {
    // Bad: HTTP concerns in domain
    userID := ctx.Param("id")
    // Bad: Direct DB access in domain
    db.Model(&User{}).Where("id = ?", userID).Update("role", "admin")
}
```

### 2. Use DTOs for Boundaries

✅ **DO**: Bind HTTP requests to DTOs (Recommended!)
```go
// DTO for input validation
type CreateUserDto struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email" validate:"required,email"`
}

// Controller binds to DTO - this is GOOD practice
func (c *UserController) Create(ctx echo.Context) error {
    var dto CreateUserDto
    if err := ctx.Bind(&dto); err != nil {  // ✅ Binding to DTO is correct
        return ctx.JSON(400, err)
    }
    
    // Transform DTO to domain model in controller
    user := models.User{
        Name:  dto.Name,
        Email: dto.Email,
    }
    
    return c.service.Create(user)
}
```

❌ **DON'T**: Bind directly to domain models
```go
// Bad: Binding directly to domain model
func (c *UserController) Create(ctx echo.Context) error {
    var user models.User
    ctx.Bind(&user)  // ❌ Exposes domain model to HTTP layer
    
    // Problems:
    // 1. HTTP tags pollute domain model
    // 2. Can't validate input separately
    // 3. Client can set fields they shouldn't (ID, CreatedAt, etc.)
    // 4. Tight coupling between HTTP and domain
}
```

**Why DTOs are important:**
- **Validation**: Apply input-specific validation rules
- **Security**: Prevent mass assignment vulnerabilities
- **Decoupling**: HTTP layer doesn't dictate domain structure
- **Flexibility**: Input format can differ from domain model
- **Clarity**: Clear contract for what clients can send

### 3. Repository Returns Domain Objects

```go
// Good: Repository works with domain models
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
    var user models.User
    err := r.db.First(&user, id).Error
    return &user, err
}
```

### 4. Keep Services Focused

Each service should handle a specific aggregate or bounded context:
- `UserService` - User operations
- `OrderService` - Order operations
- `PaymentService` - Payment operations

### 5. Use Dependency Injection

Wire dependencies from the outside:

```go
// handlers/rest.go
func NewRestHandler(h *Handler) {
    // Create dependencies
    repo := repositories.NewUserRepository(h.db)
    service := domains.NewUserService(repo)
    controller := controllers.NewUserController(service)
    
    // Register routes
    h.app.POST("/users", controller.Create)
}
```

## Example Implementation

### Complete Feature: User Registration

#### 1. Model (`internal/models/user.go`)
```go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Password  string    `json:"-" gorm:"not null"`  // Never expose in JSON
    CreatedAt time.Time `json:"created_at"`
}
```

#### 2. DTO (`internal/application/dto/user_dto.go`)
```go
type RegisterUserDto struct {
    Name     string `json:"name" validate:"required"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}
```

#### 3. Repository (`internal/infrastructure/repositories/user_repository.go`)
```go
type UserRepository interface {
    Create(user *models.User) error
    FindByEmail(email string) (*models.User, error)
}

type GormUserRepository struct {
    db *gorm.DB
}

func (r *GormUserRepository) Create(user *models.User) error {
    return r.db.Create(user).Error
}

func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
    var user models.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

#### 4. Service (`internal/domains/user_service.go`)
```go
type UserService struct {
    repo repositories.UserRepository
}

func (s *UserService) RegisterUser(dto dto.RegisterUserDto) (*models.User, error) {
    // Business rule: Check if email exists
    existing, _ := s.repo.FindByEmail(dto.Email)
    if existing != nil {
        return nil, errors.New("email already registered")
    }
    
    // Business rule: Hash password
    hashedPassword, _ := bcrypt.GenerateFromPassword(
        []byte(dto.Password), bcrypt.DefaultCost,
    )
    
    user := &models.User{
        Name:     dto.Name,
        Email:    strings.ToLower(dto.Email),
        Password: string(hashedPassword),
    }
    
    err := s.repo.Create(user)
    return user, err
}
```

#### 5. Controller (`internal/application/rest/controllers/user_controller.go`)
```go
type UserController struct {
    service *domains.UserService
}

// @Summary Register user
// @Description Register a new user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body dto.RegisterUserDto true "User registration data"
// @Success 201 {object} models.User
// @Router /users/register [post]
func (c *UserController) Register(ctx echo.Context) error {
    var dto dto.RegisterUserDto
    
    if err := ctx.Bind(&dto); err != nil {
        return ctx.JSON(400, map[string]string{"error": "invalid input"})
    }
    
    user, err := c.service.RegisterUser(dto)
    if err != nil {
        return ctx.JSON(400, map[string]string{"error": err.Error()})
    }
    
    return ctx.JSON(201, user)
}
```

#### 6. Wire Up (`internal/handlers/rest.go`)
```go
func NewRestHandler(h *Handler) {
    // Initialize dependencies
    userRepo := repositories.NewGormUserRepository(h.db)
    userService := domains.NewUserService(userRepo)
    userController := controllers.NewUserController(userService)
    
    // Register routes
    h.app.POST("/users/register", userController.Register)
}
```

## Benefits of This Architecture

### 1. **Testability**
- Test domain logic without database
- Mock repositories easily
- Test controllers without HTTP

### 2. **Maintainability**
- Clear separation of concerns
- Easy to locate code
- Changes isolated to specific layers

### 3. **Flexibility**
- Switch databases without touching business logic
- Add new input methods (GraphQL, gRPC) easily
- Replace frameworks without rewriting domain

### 4. **Scalability**
- Add features without breaking existing code
- Team can work on different layers independently
- Clear contracts between layers

### 5. **Business Focus**
- Domain layer clearly expresses business rules
- Technical details don't clutter business logic
- Easy for domain experts to review code

## Further Reading

- [Domain-Driven Design by Eric Evans](https://www.domainlanguage.com/ddd/)
- [Hexagonal Architecture by Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Implementing Domain-Driven Design by Vaughn Vernon](https://vaughnvernon.com/)

---

**Note**: This architecture may seem like overkill for simple CRUD apps, but it shines in complex domains with rich business logic. Adapt the level of abstraction to your project's needs.
