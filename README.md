<div align="center">

# 🚀 Framework

**یک فریمورک جامع و قدرتمند برای ساخت اپلیکیشن‌های Go با معماری Clean Architecture و Domain-Driven Design**

[![Go Version](https://img.shields.io/badge/Go-1.25.1+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/ali-mahdavi-dev/framework?style=for-the-badge)](https://goreportcard.com/report/github.com/ali-mahdavi-dev/framework)

[📖 مستندات](#-مستندات) • [🚀 شروع سریع](#-شروع-سریع) • [💡 مثال‌ها](#-مثال‌ها) • [🏗️ معماری](#️-معماری) • [🤝 مشارکت](#-مشارکت)

---

</div>

## 📋 فهرست مطالب

- [✨ ویژگی‌ها](#-ویژگی‌ها)
- [🏗️ معماری](#️-معماری)
- [🔧 زیرساخت‌ها](#-زیرساخت‌ها)
- [🛠️ ابزارها](#️-ابزارها)
- [🚀 شروع سریع](#-شروع-سریع)
- [💡 مثال‌ها](#-مثال‌ها)
- [📁 ساختار پروژه](#-ساختار-پروژه)
- [📦 وابستگی‌ها](#-وابستگی‌ها)
- [📖 مستندات](#-مستندات)
- [🤝 مشارکت](#-مشارکت)

---

## ✨ ویژگی‌ها

این فریمورک با هدف ارائه یک راه‌حل کامل و حرفه‌ای برای ساخت اپلیکیشن‌های Go طراحی شده است. با استفاده از بهترین الگوها و معماری‌های موجود، توسعه‌دهندگان می‌توانند اپلیکیشن‌های مقیاس‌پذیر و قابل نگهداری بسازند.

### 🎯 مزایای کلیدی

- ✅ **معماری تمیز**: جداسازی کامل لایه‌ها و وابستگی‌ها
- ✅ **مقیاس‌پذیری**: طراحی شده برای اپلیکیشن‌های بزرگ و پیچیده
- ✅ **قابلیت تست**: ساختار مناسب برای Unit Testing و Integration Testing
- ✅ **عملکرد بالا**: بهینه‌سازی شده برای کارایی و سرعت
- ✅ **مستندات کامل**: مثال‌ها و مستندات جامع
- ✅ **انعطاف‌پذیری**: قابل سفارشی‌سازی و توسعه

---

## 🏗️ معماری

این فریمورک از معماری‌های مدرن و الگوهای طراحی پیشرفته پشتیبانی می‌کند:

### 🏛️ Clean Architecture

جداسازی کامل لایه‌ها با وابستگی‌های یک‌طرفه به سمت Domain Layer

```
┌─────────────────────────────────────┐
│      Presentation Layer (API)       │
├─────────────────────────────────────┤
│      Application Layer (Service)    │
├─────────────────────────────────────┤
│      Domain Layer (Business Logic)  │
├─────────────────────────────────────┤
│   Infrastructure Layer (Database)   │
└─────────────────────────────────────┘
```

### 📐 Domain-Driven Design (DDD)

- **Entities**: موجودیت‌های Domain با شناسه یکتا
- **Value Objects**: اشیاء بدون شناسه
- **Aggregates**: گروه‌بندی Entities و Value Objects
- **Domain Events**: رویدادهای Domain برای ارتباط بین Bounded Contexts
- **Repositories**: الگوی دسترسی به داده

### ⚡ CQRS (Command Query Responsibility Segregation)

جداسازی عملیات خواندن و نوشتن برای بهینه‌سازی عملکرد:

- **Commands**: عملیات تغییردهنده State
- **Queries**: عملیات خواندن داده
- **Event Handlers**: پردازش رویدادها

### 🔄 Event Sourcing

- مدیریت رویدادها و Domain Events
- بازسازی State از رویدادها
- Audit Trail کامل

### 🔒 Unit of Work Pattern

مدیریت تراکنش‌ها و تضمین Consistency

### 📚 Repository Pattern

الگوی دسترسی به داده با جداسازی منطق Business از دسترسی به داده

---

## 🔧 زیرساخت‌ها

فریمورک از طیف گسترده‌ای از تکنولوژی‌های زیرساختی پشتیبانی می‌کند:

| تکنولوژی             | توضیحات                      | استفاده                  |
| -------------------- | ---------------------------- | ------------------------ |
| 🗄️ **PostgreSQL**    | دیتابیس رابطه‌ای قدرتمند     | ذخیره‌سازی داده‌های اصلی |
| 💾 **SQLite**        | دیتابیس سبک و سریع           | توسعه و تست              |
| ⚡ **Redis**         | Cache و Session Management   | بهبود عملکرد             |
| 🔍 **Elasticsearch** | موتور جستجو و تحلیل          | جستجوی پیشرفته           |
| 📨 **Kafka**         | پیام‌رسانی و Event Streaming | ارتباط Asynchronous      |
| 🔌 **WebSocket**     | ارتباط Real-time             | با Socket.IO             |
| 📊 **Jaeger**        | Distributed Tracing          | ردیابی و Debugging       |
| 📝 **Zerolog**       | Logger سریع و ساختاریافته    | لاگ‌گیری                 |

---

## 🛠️ ابزارها

### 🔐 احراز هویت و امنیت

- **JWT**: تولید و اعتبارسنجی Token
- **Middleware**: Authentication و Authorization
- **Validation**: اعتبارسنجی با `go-playground/validator`

### 🎯 مدیریت خطا

- خطاهای ساختاریافته
- کدهای خطای استاندارد
- پیام‌های خطای چندزبانه

### 📋 Middleware

- **Logger**: لاگ‌گیری درخواست‌ها
- **Request ID**: ردیابی درخواست‌ها
- **Tracing**: ردیابی با OpenTelemetry

### 🧩 الگوهای طراحی

- **Specification Pattern**: پیاده‌سازی Business Rules
- **Message Bus**: مدیریت Command و Event Handlers
- **Outbox Pattern**: تضمین تحویل پیام‌ها

---

## 🚀 شروع سریع

### 📦 نصب

```bash
go get github.com/ali-mahdavi-dev/framework
```

### 🔧 پیش‌نیازها

- Go 1.25.1 یا بالاتر
- PostgreSQL (برای Production)
- Redis (اختیاری - برای Cache)
- Elasticsearch (اختیاری - برای جستجو)

---

## 💡 مثال‌ها

### 📊 اتصال به دیتابیس

```go
package main

import (
    "log"
    "github.com/ali-mahdavi-dev/framework/infrastructure/databases"
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

    // استفاده از db...
}
```

### 📝 استفاده از Logger

```go
package main

import (
    "github.com/ali-mahdavi-dev/framework/infrastructure/logging"
)

func main() {
    // استفاده از Logger با Builder Pattern
    logging.Info("User logged in").
        WithString("user_id", "123").
        WithString("ip", "192.168.1.1").
        WithInt("status_code", 200).
        Log()

    // یا با Format
    logging.Infof("Processing request: %s", requestID)

    // لاگ خطا
    logging.Error("Failed to process request").
        WithString("error", err.Error()).
        Log()
}
```

### 🗄️ Repository Pattern

```go
package main

import (
    "github.com/ali-mahdavi-dev/framework/adapter"
)

// تعریف Entity
type User struct {
    adapter.BaseEntity
    ID    uint64 `gorm:"primaryKey"`
    Name  string
    Email string
}

// تعریف Repository Interface
type UserRepository interface {
    adapter.BaseRepository[User]
    FindByEmail(email string) (*User, error)
    FindActiveUsers() ([]User, error)
}

// پیاده‌سازی Repository
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
    "github.com/ali-mahdavi-dev/framework/service_layer/command_event_handler"
)

// تعریف Command
type CreateUserCommand struct {
    Name  string `validate:"required,min=3"`
    Email string `validate:"required,email"`
}

// ایجاد Handler
func setupUserHandlers() {
    handler := commandeventhandler.NewCommandHandler[CreateUserCommand](
        func(ctx context.Context, cmd *CreateUserCommand) error {
            // منطق ایجاد کاربر
            user := &User{
                Name:  cmd.Name,
                Email: cmd.Email,
            }

            // ذخیره در دیتابیس
            return userRepo.Create(ctx, user)
        },
    )

    // ثبت Handler
    messageBus.RegisterCommandHandler(handler)
}
```

### 🔐 JWT Authentication

```go
package main

import (
    "time"
    "github.com/ali-mahdavi-dev/framework/api/jwt"
)

func generateToken(userID uint64) (string, error) {
    token, err := jwt.GenerateToken(
        24*time.Hour,           // مدت اعتبار
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
    apperrors "github.com/ali-mahdavi-dev/framework/errors"
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
        // مدیریت خطای اعتبارسنجی
        log.Printf("Validation error: %s", e.Message)
    case *apperrors.NotFoundError:
        // مدیریت خطای پیدا نشدن
        log.Printf("Not found: %s", e.Message)
    default:
        // مدیریت خطای عمومی
        log.Printf("Error: %v", err)
    }
}
```

### 🎯 Specification Pattern

```go
package main

import (
    "github.com/ali-mahdavi-dev/framework/specification"
)

// تعریف Specification
type ActiveUserSpec struct{}

func (s *ActiveUserSpec) IsSatisfiedBy(user User) bool {
    return user.IsActive && !user.IsDeleted
}

type PremiumUserSpec struct{}

func (s *PremiumUserSpec) IsSatisfiedBy(user User) bool {
    return user.IsPremium && user.SubscriptionExpiresAt.After(time.Now())
}

// استفاده از Specification
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

## 📁 ساختار پروژه

```
framework/
├── 📂 adapter/                    # رابط‌های Repository و Unit of Work
│   ├── interface_entity.go
│   ├── interface_gorm_repository.go
│   ├── interface_repository.go
│   └── unit_of_work.go
│
├── 📂 api/                        # ابزارهای HTTP، JWT و Middleware
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
├── 📂 errors/                     # مدیریت خطاهای ساختاریافته
│   ├── constructors.go
│   ├── phrases/                   # پیام‌های خطا
│   │   ├── error_code.go
│   │   └── message.go
│   ├── types.go
│   └── utils.go
│
├── 📂 helpers/                    # توابع کمکی
│   ├── jsonhelper/
│   ├── kind/
│   └── utils.go
│
├── 📂 infrastructure/             # اتصالات به سرویس‌های خارجی
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
├── 📂 service_layer/              # لایه سرویس
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

## 📦 وابستگی‌ها

### اصلی

| پکیج               | نسخه    | توضیحات                   |
| ------------------ | ------- | ------------------------- |
| **Go**             | 1.25.1+ | زبان برنامه‌نویسی         |
| **GORM**           | latest  | ORM برای دیتابیس          |
| **Fiber**          | v3      | Web Framework             |
| **Zerolog**        | latest  | Logger ساختاریافته        |
| **Redis**          | v9      | Client برای Redis         |
| **Elasticsearch**  | v8      | Client برای Elasticsearch |
| **Kafka (Sarama)** | latest  | Client برای Kafka         |
| **OpenTelemetry**  | latest  | Distributed Tracing       |
| **Socket.IO**      | latest  | WebSocket Communication   |
| **JWT**            | v5      | احراز هویت                |

### جزئیات

برای مشاهده لیست کامل وابستگی‌ها، فایل `go.mod` را بررسی کنید.

---

## 📖 مستندات

### 🎓 یادگیری

- [Clean Architecture در Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)

### 📚 منابع بیشتر

- مستندات کامل در حال آماده‌سازی است
- مثال‌های بیشتر در پوشه `examples/` (به زودی)

---

## 🤝 مشارکت

مشارکت‌ها، Issues و Pull Request‌ها خوش‌آمد هستند! برای مشارکت:

1. ⭐ این پروژه را Star کنید
2. 🍴 Fork کنید
3. 🌿 یک Branch جدید ایجاد کنید (`git checkout -b feature/AmazingFeature`)
4. 💾 تغییرات را Commit کنید (`git commit -m 'Add some AmazingFeature'`)
5. 📤 Branch را Push کنید (`git push origin feature/AmazingFeature`)
6. 🔄 یک Pull Request باز کنید

### 📝 دستورالعمل‌های مشارکت

- کد را با `gofmt` فرمت کنید
- تست‌ها را اضافه کنید
- مستندات را به‌روز کنید
- از Commit Message های واضح استفاده کنید

---

## 📄 مجوز

این پروژه تحت مجوز **MIT** منتشر شده است. برای جزئیات بیشتر، فایل `LICENSE` را مطالعه کنید.

---

## 👤 نویسنده

**Ali Mahdavi**

- 🌐 GitHub: [@ali-mahdavi-dev](https://github.com/ali-mahdavi-dev)
- 📧 Email: [در حال آماده‌سازی]

---

<div align="center">

### ⭐ اگر این پروژه برای شما مفید بود، یک Star بدهید!

**ساخته شده با ❤️ در ایران**

[⬆ به بالا برگردید](#-framework)

</div>
