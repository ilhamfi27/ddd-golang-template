# DDD Golang Template

A production-ready Golang template implementing Domain-Driven Design (DDD) principles with clean architecture. This template provides a solid foundation for building scalable REST APIs with multiple database support and comprehensive tooling.

## 🚀 Features

- **Domain-Driven Design (DDD)**: Properly structured layers following DDD principles
- **Clean Architecture**: Separation of concerns with clear boundaries between layers
- **Multiple Database Support**: MySQL, PostgreSQL, and SQLite drivers included
- **Database Migrations**: Built-in migration support using golang-migrate
- **RESTful API**: Echo web framework for high-performance REST APIs
- **API Documentation**: Auto-generated Swagger documentation
- **Hot Reload**: Air for development hot-reloading
- **Dependency Injection**: Clean dependency management across layers
- **Error Handling**: Centralized error handling mechanism
- **Environment Configuration**: Easy configuration via environment variables
- **Automated Versioning**: Semantic-release for automated version management and changelog generation

## 📁 Project Structure

```
.
├── cmd/                          # Application entry points
│   └── main.go                   # Main application entry point
├── internal/                     # Private application code
│   ├── application/              # Application layer
│   │   ├── dto/                  # Data Transfer Objects
│   │   └── rest/                 # REST API layer
│   │       ├── controllers/      # HTTP controllers
│   │       └── errors/           # REST error handling
│   ├── config/                   # Configuration
│   │   ├── app/                  # Application config (env, constants)
│   │   └── db/                   # Database config, drivers, and migrations
│   │       └── migrations/       # SQL migration files
│   ├── domains/                  # Domain layer (business logic)
│   ├── infrastructure/           # Infrastructure layer
│   │   └── repositories/         # Repository implementations
│   ├── models/                   # Data models/entities
│   ├── handlers/                 # Application handlers
│   └── utils/                    # Utility functions
├── swagger/                      # Swagger documentation (auto-generated)
├── dbs/                          # Local database files (SQLite)
├── tmp/                          # Temporary files (Air)
├── go.mod                        # Go module definition
├── Makefile                      # Build automation
└── .env.example                  # Environment variables example
```

## 🏗️ Architecture

This template follows a layered architecture based on Domain-Driven Design and Hexagonal Architecture principles.

**📖 For detailed architecture documentation, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)**

### Quick Overview

1. **Application Layer** (`internal/application/`)
   - Handles HTTP requests and responses
   - Contains DTOs for data transfer
   - REST controllers for API endpoints
   - Input validation and error formatting

2. **Domain Layer** (`internal/domains/`)
   - Core business logic
   - Domain services
   - Business rules and validations
   - Independent of external concerns

3. **Infrastructure Layer** (`internal/infrastructure/repositories/`)
   - Database repositories
   - External service integrations
   - Implementation of interfaces defined in domain layer

4. **Models** (`internal/models/`)
   - Domain entities
   - Database models
   - Shared data structures

### Data Flow

```
HTTP Request → Controller → Service (Domain) → Repository (Infrastructure) → Database
HTTP Response ← Controller ← Service (Domain) ← Repository (Infrastructure) ← Database
```

## 🛠️ Prerequisites

- **Go**: 1.24.0 or higher
- **Make**: For using Makefile commands
- **Air**: For hot-reloading in development (optional)
- **Docker**: For running databases (optional)
- **golang-migrate**: For running migrations manually (optional)

## 📦 Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/ilhamfi27/ddd-golang-template.git
   cd ddd-golang-template
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Install Air for hot-reloading** (optional):
   ```bash
   go install github.com/air-verse/air@latest
   ```

4. **Install Swag for API documentation** (optional):
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

## ⚙️ Configuration

1. **Create environment file**:
   ```bash
   cp .env.example .env
   ```

2. **Configure environment variables** in `.env`:
   ```env
   # Application
   PORT=1321
   APP_ENV=development  # development or production
   
   # Database
   DATABASE_DRIVER=sqlite  # mysql, postgres, or sqlite
   DATABASE_HOST=localhost
   DATABASE_PORT=3306      # 3306 for MySQL, 5432 for PostgreSQL
   DATABASE_USER=root
   DATABASE_PASS=password
   DATABASE_NAME=ddd_golang
   
   # Migrations
   MIGRATION_DIR=./internal/config/db/migrations
   ```

### Database Drivers

#### MySQL
```env
DATABASE_DRIVER=mysql
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=root
DATABASE_PASS=password
DATABASE_NAME=ddd_golang
```

#### PostgreSQL
```env
DATABASE_DRIVER=postgres
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASS=password
DATABASE_NAME=ddd_golang
```

#### SQLite
```env
DATABASE_DRIVER=sqlite
DATABASE_NAME=ddd_golang
```

## 🚀 Usage

### Development

#### Start the application with hot-reload:
```bash
make start
```
Or directly with Air:
```bash
air
```

#### Run the application without hot-reload:
```bash
go run cmd/main.go
```

### Production

#### Build the application:
```bash
go build -o main cmd/main.go
```

#### Run the built binary:
```bash
./main
```

### Database Migrations

#### Run migrations:
```bash
make run-migration-up
```
Or:
```bash
go run cmd/main.go --migrate
```

#### Create a new migration:
```bash
make new-migration NAME=create_users_table
```

This creates two files in `internal/config/db/migrations/`:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

### Generate Swagger Documentation

```bash
make swag
```

After generation, access the documentation at:
```
http://localhost:1321/swagger/index.html
```

## 📚 API Documentation

Once the application is running, you can access the Swagger UI documentation at:

```
http://localhost:1321/swagger/index.html
```

### Example Endpoints

#### GET /example
Get example data.

**Query Parameters:**
- `example` (string, required): Example parameter

**Response:**
```json
{
  "name": "example"
}
```

## 🔧 Makefile Commands

| Command | Description |
|---------|-------------|
| `make init` | Initialize project (clean) |
| `make tidy` | Tidy Go modules |
| `make vendor` | Create vendor directory |
| `make clean` | Clean build artifacts and databases |
| `make clean-dbs` | Remove local database files |
| `make clean-build` | Remove built binary |
| `make start` | Start development server with hot-reload |
| `make run-migration-up` | Run database migrations |
| `make new-migration NAME=name` | Create new migration files |
| `make swag` | Generate Swagger documentation |

## 🏗️ Development Workflow

### Adding a New Feature

1. **Create the model** in `internal/models/`:
   ```go
   // internal/models/user.go
   package models
   
   type User struct {
       ID    uint   `json:"id" gorm:"primaryKey"`
       Name  string `json:"name"`
       Email string `json:"email"`
   }
   ```

2. **Create the DTO** in `internal/application/dto/`:
   ```go
   // internal/application/dto/user_dto.go
   package dto
   
   type CreateUserDto struct {
       Name  string `json:"name" validate:"required"`
       Email string `json:"email" validate:"required,email"`
   }
   ```

3. **Create the repository** in `internal/infrastructure/repositories/`:
   ```go
   // internal/infrastructure/repositories/user_repository.go
   package repositories
   
   import (
       "github.com/ilhamfi27/ddd-golang-template/internal/models"
       "gorm.io/gorm"
   )
   
   type UserRepository struct {
       db *gorm.DB
   }
   
   func NewUserRepository(db *gorm.DB) *UserRepository {
       return &UserRepository{db: db}
   }
   
   func (r *UserRepository) Create(user *models.User) error {
       return r.db.Create(user).Error
   }
   ```

4. **Create the service** in `internal/domains/`:
   ```go
   // internal/domains/user_service.go
   package domains
   
   import (
       "github.com/ilhamfi27/ddd-golang-template/internal/application/dto"
       "github.com/ilhamfi27/ddd-golang-template/internal/infrastructure/repositories"
       "github.com/ilhamfi27/ddd-golang-template/internal/models"
   )
   
   type UserService struct {
       repo *repositories.UserRepository
   }
   
   func NewUserService(repo *repositories.UserRepository) *UserService {
       return &UserService{repo: repo}
   }
   
   func (s *UserService) CreateUser(data dto.CreateUserDto) error {
       user := &models.User{
           Name:  data.Name,
           Email: data.Email,
       }
       return s.repo.Create(user)
   }
   ```

5. **Create the controller** in `internal/application/rest/controllers/`:
   ```go
   // internal/application/rest/controllers/user_controller.go
   package controllers
   
   import (
       "net/http"
       
       "github.com/ilhamfi27/ddd-golang-template/internal/application/dto"
       "github.com/ilhamfi27/ddd-golang-template/internal/domains"
       "github.com/labstack/echo/v4"
   )
   
   type UserController struct {
       service *domains.UserService
   }
   
   func NewUserController(s *domains.UserService) *UserController {
       return &UserController{service: s}
   }
   
   // @Summary Create user
   // @Description Create a new user
   // @Tags Users
   // @Accept json
   // @Produce json
   // @Param user body dto.CreateUserDto true "User data"
   // @Success 201 {object} models.User
   // @Router /users [post]
   func (c *UserController) CreateUser(ctx echo.Context) error {
       var data dto.CreateUserDto
       if err := ctx.Bind(&data); err != nil {
           return ctx.JSON(http.StatusBadRequest, err)
       }
       
       if err := c.service.CreateUser(data); err != nil {
           return ctx.JSON(http.StatusInternalServerError, err)
       }
       
       return ctx.JSON(http.StatusCreated, map[string]string{"message": "User created"})
   }
   ```

6. **Register the routes** in `internal/handlers/rest.go`

7. **Generate Swagger docs**:
   ```bash
   make swag
   ```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## 📝 Code Style

This project follows standard Go conventions:
- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines
- Write meaningful commit messages using [Conventional Commits](https://www.conventionalcommits.org/)
- Document public functions and types

## 🔄 Versioning & Releases

This project uses [semantic-release](https://github.com/semantic-release/semantic-release) for automated versioning and changelog generation.

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>
```

**Types:**
- `feat`: New feature (triggers minor release)
- `fix`: Bug fix (triggers patch release)
- `docs`: Documentation changes (triggers patch release)
- `perf`: Performance improvement (triggers patch release)
- `refactor`: Code refactoring (triggers patch release)
- `chore`: Maintenance tasks (no release)
- `test`: Adding tests (no release)

**Examples:**
```bash
feat(auth): add JWT authentication
fix(db): resolve connection timeout
docs: update installation guide
```

For detailed information, see [SEMANTIC_RELEASE.md](./docs/SEMANTIC_RELEASE.md).

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes using conventional commits (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

## 👤 Author

**Ilham Fadhilah**

- Email: r.ilhamfadhilah@gmail.com
- GitHub: [@ilhamfi27](https://github.com/ilhamfi27)

## 🙏 Acknowledgments

- [Echo Framework](https://echo.labstack.com/) - High performance, extensible, minimalist Go web framework
- [GORM](https://gorm.io/) - The fantastic ORM library for Golang
- [Swaggo](https://github.com/swaggo/swag) - Automatically generate RESTful API documentation with Swagger 2.0
- [Air](https://github.com/air-verse/air) - Live reload for Go apps

## 📚 Additional Resources

- [Architecture Documentation](docs/ARCHITECTURE.md) - Detailed DDD and Hexagonal Architecture guide
- [Semantic Release Guide](docs/SEMANTIC_RELEASE.md) - Versioning and release workflow
- [Domain-Driven Design](https://martinfowler.com/tags/domain%20driven%20design.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

## 🐛 Known Issues

- None at the moment

## 🗺️ Roadmap

- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Add Docker support for the application
- [ ] Add CI/CD pipeline examples
- [ ] Add authentication middleware
- [ ] Add request validation middleware
- [ ] Add logging middleware
- [ ] Add rate limiting
- [ ] Add caching layer
- [ ] Add GraphQL support

---

**Made with ❤️ using Go**
