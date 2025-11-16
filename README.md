<div align="center">

# 🚀 Framework

**A comprehensive and powerful framework for building Go applications with Clean Architecture and Domain-Driven Design**

[![Go Version](https://img.shields.io/badge/Go-1.25.1+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ali-mahdavi-dev/shikposh-framework?style=for-the-badge)](https://goreportcard.com/report/github.com/ali-mahdavi-dev/shikposh-framework)

[📖 Documentation](#-documentation) • [🚀 Quick Start](#-quick-start) • [💡 Examples](#-examples) • [🏗️ Architecture](#️-architecture) • [🤝 Contributing](#-contributing)

---

</div>

## 📋 Table of Contents

- [✨ Features](#-features)
- [🏗️ Architecture](#️-architecture)
- [🔧 Infrastructure](#-infrastructure)
- [🛠️ Tools](#️-tools)
- [🚀 Quick Start](#-quick-start)
- [💡 Examples](#-examples)
- [📁 Project Structure](#-project-structure)
- [📦 Dependencies](#-dependencies)
- [📖 Documentation](#-documentation)
- [🤝 Contributing](#-contributing)

---

## ✨ Features

This framework is designed to provide a complete and professional solution for building Go applications. By leveraging the best patterns and architectures available, developers can build scalable and maintainable applications.

### 🎯 Key Benefits

- ✅ **Clean Architecture**: Complete separation of layers and dependencies
- ✅ **Scalability**: Designed for large and complex applications
- ✅ **Testability**: Proper structure for Unit Testing and Integration Testing
- ✅ **High Performance**: Optimized for efficiency and speed
- ✅ **Comprehensive Documentation**: Extensive examples and documentation
- ✅ **Flexibility**: Customizable and extensible

---

## 🏗️ Architecture

This framework supports modern architectures and advanced design patterns:

### 🏛️ Clean Architecture

Complete layer separation with unidirectional dependencies toward the Domain Layer

```
┌─────────────────────────────────────┐
│      Presentation Layer (API)       │
├─────────────────────────────────────┤
│      Application Layer (Service)    │
├─────────────────────────────────────┤
│      Domain Layer (Business Logic)  │
├─────────────────────────────────────┤
│   Infrastructure Layer (Database)  │
└─────────────────────────────────────┘
```

### 📐 Domain-Driven Design (DDD)

- **Entities**: Domain entities with unique identifiers
- **Value Objects**: Objects without identity
- **Aggregates**: Grouping of Entities and Value Objects
- **Domain Events**: Domain events for communication between Bounded Contexts
- **Repositories**: Data access pattern

### ⚡ CQRS (Command Query Responsibility Segregation)

Separation of read and write operations for performance optimization:

- **Commands**: State-changing operations
- **Queries**: Data reading operations
- **Event Handlers**: Event processing

### 🔄 Event Sourcing

- Event and Domain Event management
- State reconstruction from events
- Complete audit trail

### 🔒 Unit of Work Pattern

Transaction management and consistency guarantee

### 📚 Repository Pattern

Data access pattern with separation of business logic from data access

---

## 🔧 Infrastructure

The framework supports a wide range of infrastructure technologies:

| Technology           | Description                   | Usage                      |
| -------------------- | ----------------------------- | -------------------------- |
| 🗄️ **PostgreSQL**    | Powerful relational database  | Primary data storage       |
| 💾 **SQLite**        | Lightweight and fast database | Development and testing    |
| ⚡ **Redis**         | Cache and Session Management  | Performance improvement    |
| 🔍 **Elasticsearch** | Search engine and analytics   | Advanced search            |
| 📨 **Kafka**         | Messaging and Event Streaming | Asynchronous communication |
| 🔌 **WebSocket**     | Real-time communication       | With Socket.IO             |
| 📊 **Jaeger**        | Distributed Tracing           | Tracking and Debugging     |
| 📝 **Zerolog**       | Fast and structured logger    | Logging                    |

---

## 🛠️ Tools

### 🔐 Authentication & Security

- **JWT**: Token generation and validation
- **Middleware**: Authentication and Authorization
- **Validation**: Validation with `go-playground/validator`

### 🎯 Error Management

- Structured errors
- Standard error codes
- Multilingual error messages

### 📋 Middleware

- **Logger**: Request logging
- **Request ID**: Request tracking
- **Tracing**: Tracking with OpenTelemetry

### 🧩 Design Patterns

- **Specification Pattern**: Business Rules implementation
- **Message Bus**: Command and Event Handler management
- **Outbox Pattern**: Message delivery guarantee

---

## 🚀 Quick Start

### 📦 Installation

```bash
go get github.com/ali-mahdavi-dev/shikposh-framework
```

### 🔧 Prerequisites

- Go 1.25.1 or higher
- PostgreSQL (for Production)
- Redis (optional - for Cache)
- Elasticsearch (optional - for Search)

---

## 💡 Examples

### 📊 Database Connection

```go
package main

import (
    "log"
    "github.com/ali-mahdavi-dev/shikposh-framework/infrastructure/databases"
)

func main() {
    db, err := databases.New(databases.Config{
        DBType:       "postgres",
        DSN:          "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable",
        MaxOpenConns: 100,
        MaxIdleConns: 10,
        MaxLifetime:  3600, // seconds
        MaxIdleTime:  600,  // seconds
        Debug:        false,
    })

    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Use db...
}
```

### 📝 Using Logger

```go
package main

import (
    "github.com/ali-mahdavi-dev/shikposh-framework/infrastructure/logging"
)

func main() {
    // Using Logger with Builder Pattern
    logging.Info("User logged in").
        WithString("user_id", "123").
        WithString("ip", "192.168.1.1").
        WithInt("status_code", 200).
        Log()

    // Or with Format
    logging.Infof("Processing request: %s", requestID)

    // Error logging
    logging.Error("Failed to process request").
        WithString("error", err.Error()).
        Log()
}
```

### 🗄️ Repository Pattern

```go
package main

import (
    "gorm.io/gorm"
    "github.com/ali-mahdavi-dev/shikposh-framework/adapter"
)

// Define Entity
type User struct {
    adapter.BaseEntity
    ID    uint64 `gorm:"primaryKey"`
    Name  string
    Email string
}

// Define Repository Interface
type UserRepository interface {
    adapter.BaseRepository[User]
    FindByEmail(email string) (*User, error)
    FindActiveUsers() ([]User, error)
}

// Implement Repository
type userRepository struct {
    adapter.GormRepository[User]
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{
        GormRepository: adapter.NewGormRepository[User](db),
    }
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
    var user User
    err := r.DB.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

### ⚡ Command Handler (CQRS)

```go
package main

import (
    "context"
    "github.com/ali-mahdavi-dev/shikposh-framework/service_layer/command_event_handler"
)

// Define Command
type CreateUserCommand struct {
    Name  string `validate:"required,min=3"`
    Email string `validate:"required,email"`
}

// Create Handler
func setupUserHandlers() {
    handler := commandeventhandler.NewCommandHandler[CreateUserCommand](
        func(ctx context.Context, cmd *CreateUserCommand) error {
            // User creation logic
            user := &User{
                Name:  cmd.Name,
                Email: cmd.Email,
            }

            // Save to database
            return userRepo.Create(ctx, user)
        },
    )

    // Register Handler
    messageBus.RegisterCommandHandler(handler)
}
```

### 🔐 JWT Authentication

```go
package main

import (
    "time"
    "github.com/ali-mahdavi-dev/shikposh-framework/api/jwt"
)

func generateToken(userID uint64) (string, error) {
    token, err := jwt.GenerateToken(
        24*time.Hour,           // Expiration time
        "your-secret-key",      // Secret Key
        userID,                 // User ID
    )

    if err != nil {
        return "", err
    }

    return token, nil
}

func validateToken(tokenString string) (uint64, error) {
    userID, err := jwt.ValidateToken(tokenString, "your-secret-key")
    return userID, err
}
```

### ⚠️ Error Handling

```go
package main

import (
    "log"
    apperrors "github.com/ali-mahdavi-dev/shikposh-framework/errors"
)

func validateEmail(email string) error {
    if !isValidEmailFormat(email) {
        return apperrors.NewValidationError(
            "INVALID_EMAIL",
            "Email format is invalid",
            "The provided email does not match the required format",
        )
    }
    return nil
}

func handleError(err error) {
    switch e := err.(type) {
    case *apperrors.ValidationError:
        // Handle validation error
        log.Printf("Validation error: %s", e.Message)
    case *apperrors.NotFoundError:
        // Handle not found error
        log.Printf("Not found: %s", e.Message)
    default:
        // Handle general error
        log.Printf("Error: %v", err)
    }
}
```

### 🎯 Specification Pattern

```go
package main

import (
    "time"
    "github.com/ali-mahdavi-dev/shikposh-framework/specification"
)

// Define Specification
type ActiveUserSpec struct{}

func (s *ActiveUserSpec) IsSatisfiedBy(user User) bool {
    return user.IsActive && !user.IsDeleted
}

type PremiumUserSpec struct{}

func (s *PremiumUserSpec) IsSatisfiedBy(user User) bool {
    return user.IsPremium && user.SubscriptionExpiresAt.After(time.Now())
}

// Using Specification
func filterUsers(users []User) []User {
    activeSpec := specification.NewBuilder(&ActiveUserSpec{})
    premiumSpec := specification.NewBuilder(&PremiumUserSpec{})

    var result []User
    for _, user := range users {
        if activeSpec.IsSatisfiedBy(user) && premiumSpec.IsSatisfiedBy(user) {
            result = append(result, user)
        }
    }

    return result
}
```

---

## 📁 Project Structure

```
framework/
├── 📂 adapter/                    # Repository and Unit of Work interfaces
│   ├── interface_entity.go
│   ├── interface_gorm_repository.go
│   ├── interface_repository.go
│   └── unit_of_work.go
│
├── 📂 api/                        # HTTP utilities, JWT, and Middleware
│   ├── http/                      # HTTP Utilities
│   │   ├── errors.go
│   │   ├── schema.go
│   │   └── utils.go
│   ├── jwt/                       # JWT Authentication
│   │   └── jwt.go
│   └── middleware/                # HTTP Middleware
│       ├── logger.go
│       ├── request_id.go
│       └── tracing.go
│
├── 📂 errors/                     # Structured error management
│   ├── constructors.go
│   ├── phrases/                   # Error messages
│   │   ├── error_code.go
│   │   └── message.go
│   ├── types.go
│   └── utils.go
│
├── 📂 helpers/                    # Helper functions
│   ├── jsonhelper/
│   ├── kind/
│   └── utils.go
│
├── 📂 infrastructure/             # External service connections
│   ├── databases/                 # PostgreSQL, SQLite
│   │   └── postgres_connection.go
│   ├── elasticsearch/             # Elasticsearch Client
│   │   └── connection.go
│   ├── kafak/                     # Kafka Producer/Consumer
│   │   ├── kafka.go
│   │   └── topic.go
│   ├── logging/                   # Zerolog Logger
│   │   ├── logger.go
│   │   ├── types.go
│   │   └── zerolog_adapter.go
│   ├── redisx/                    # Redis Client
│   │   └── connection.go
│   ├── tracing/                   # OpenTelemetry & Jaeger
│   │   └── jaeger.go
│   └── websocket/                 # Socket.IO
│       └── socketio.go
│
├── 📂 service_layer/              # Service layer
│   ├── cache/                     # Redis Cache
│   │   └── redis_cache.go
│   ├── command_event_handler/     # CQRS Handlers
│   │   ├── command_middleware/
│   │   ├── command.go
│   │   └── event.go
│   ├── messagebus/                # Message Bus
│   │   └── messagebus.go
│   ├── outbox/                    # Outbox Pattern
│   │   ├── consumer.go
│   │   ├── entity.go
│   │   ├── processor.go
│   │   └── repository.go
│   └── unit_of_work/
│
└── 📂 specification/              # Specification Pattern
    ├── composite.go
    └── specification.go
```

---

## 📦 Dependencies

### Core

| Package            | Version | Description              |
| ------------------ | ------- | ------------------------ |
| **Go**             | 1.25.1+ | Programming language     |
| **GORM**           | latest  | ORM for database         |
| **Fiber**          | v3      | Web Framework            |
| **Zerolog**        | latest  | Structured logger        |
| **Redis**          | v9      | Client for Redis         |
| **Elasticsearch**  | v8      | Client for Elasticsearch |
| **Kafka (Sarama)** | latest  | Client for Kafka         |
| **OpenTelemetry**  | latest  | Distributed Tracing      |
| **Socket.IO**      | latest  | WebSocket Communication  |
| **JWT**            | v5      | Authentication           |

### Details

For the complete list of dependencies, check the `go.mod` file.

---

## 📖 Documentation

### 🎓 Learning Resources

- [Clean Architecture in Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)

### 📚 Additional Resources

- Complete documentation is being prepared
- More examples in the `examples/` folder (coming soon)

---

## 🤝 Contributing

Contributions, issues, and pull requests are welcome! To contribute:

1. ⭐ Star this project
2. 🍴 Fork it
3. 🌿 Create a new branch (`git checkout -b feature/AmazingFeature`)
4. 💾 Commit your changes (`git commit -m 'Add some AmazingFeature'`)
5. 📤 Push to the branch (`git push origin feature/AmazingFeature`)
6. 🔄 Open a Pull Request

### 📝 Contribution Guidelines

- Format code with `gofmt`
- Add tests
- Update documentation
- Use clear commit messages

---

## 📄 License

This project is licensed under the **MIT License**. For more details, see the `LICENSE` file.

---

## 👤 Author

**Ali Mahdavi**

- 🌐 GitHub: [@ali-mahdavi-dev](https://github.com/ali-mahdavi-dev)
- 📧 Email: [Coming soon]

---

<div align="center">

### ⭐ If this project was helpful to you, give it a star!

**Made with ❤️ in Iran**

[⬆ Back to top](#-framework)

</div>
