# Framework

یک فریمورک جامع و قدرتمند برای ساخت اپلیکیشن‌های Go با معماری Clean Architecture و Domain-Driven Design (DDD).

## ویژگی‌ها

### 🏗️ معماری

- **Clean Architecture**: جداسازی لایه‌ها و وابستگی‌ها
- **Domain-Driven Design**: پشتیبانی از DDD patterns
- **CQRS**: جداسازی Command و Query
- **Event Sourcing**: مدیریت رویدادها و Domain Events
- **Unit of Work**: مدیریت تراکنش‌ها
- **Repository Pattern**: الگوی دسترسی به داده

### 🔧 زیرساخت‌ها

- **Database**: پشتیبانی از PostgreSQL و SQLite با GORM
- **Redis**: کش و مدیریت session
- **Elasticsearch**: جستجو و تحلیل داده
- **Kafka**: پیام‌رسانی و Event Streaming
- **WebSocket**: ارتباط Real-time با Socket.IO
- **Tracing**: ردیابی با OpenTelemetry و Jaeger
- **Logging**: لاگ‌گیری با Zerolog

### 🛠️ ابزارها

- **JWT**: احراز هویت و مدیریت Token
- **Validation**: اعتبارسنجی با go-playground/validator
- **Error Handling**: مدیریت خطاهای ساختاریافته
- **Middleware**: Logger, Request ID, Tracing
- **Specification Pattern**: پیاده‌سازی Business Rules
- **Message Bus**: مدیریت Command و Event Handlers
- **Outbox Pattern**: تضمین تحویل پیام‌ها

## نصب

```bash
go get github.com/ali-mahdavi-dev/framework
```

## استفاده

### مثال: اتصال به دیتابیس

```go
import (
    "github.com/ali-mahdavi-dev/framework/infrastructure/databases"
)

db, err := databases.New(databases.Config{
    DBType:       "postgres",
    DSN:          "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable",
    MaxOpenConns: 100,
    MaxIdleConns: 10,
    MaxLifetime:  3600,
    MaxIdleTime:  600,
    Debug:        false,
})
```

### مثال: استفاده از Logger

```go
import (
    "github.com/ali-mahdavi-dev/framework/infrastructure/logging"
)

// استفاده از Logger با Builder Pattern
logging.Info("User logged in").
    WithString("user_id", "123").
    WithString("ip", "192.168.1.1").
    Log()

// یا با Format
logging.Infof("Processing request: %s", requestID)
```

### مثال: Repository Pattern

```go
import (
    "github.com/ali-mahdavi-dev/framework/adapter"
)

type UserRepository interface {
    adapter.BaseRepository[User]
    // متدهای اضافی...
}

type User struct {
    adapter.BaseEntity
    ID    uint64
    Name  string
    Email string
}
```

### مثال: Command Handler

```go
import (
    "context"
    "github.com/ali-mahdavi-dev/framework/service_layer/command_event_handler"
)

type CreateUserCommand struct {
    Name  string
    Email string
}

handler := commandeventhandler.NewCommandHandler[CreateUserCommand](
    func(ctx context.Context, cmd *CreateUserCommand) error {
        // منطق ایجاد کاربر
        return nil
    },
)
```

### مثال: JWT

```go
import (
    "time"
    "github.com/ali-mahdavi-dev/framework/api/jwt"
)

token, err := jwt.GenerateToken(
    24*time.Hour, // مدت اعتبار
    "your-secret-key",
    123, // user_id
)
```

### مثال: Error Handling

```go
import (
    apperrors "github.com/ali-mahdavi-dev/framework/errors"
)

err := apperrors.NewValidationError(
    "INVALID_EMAIL",
    "Email format is invalid",
    "The provided email does not match the required format",
)
```

### مثال: Specification Pattern

```go
import (
    "github.com/ali-mahdavi-dev/framework/specification"
)

type ActiveUserSpec struct{}

func (s *ActiveUserSpec) IsSatisfiedBy(user User) bool {
    return user.IsActive
}

spec := specification.NewBuilder(&ActiveUserSpec{})
if spec.IsSatisfiedBy(user) {
    // کاربر فعال است
}
```

## ساختار پروژه

```
framework/
├── adapter/              # رابط‌های Repository و Unit of Work
├── api/                  # ابزارهای HTTP، JWT و Middleware
├── errors/               # مدیریت خطاهای ساختاریافته
├── helpers/              # توابع کمکی
├── infrastructure/       # اتصالات به سرویس‌های خارجی
│   ├── databases/       # PostgreSQL, SQLite
│   ├── elasticsearch/   # Elasticsearch
│   ├── kafak/           # Kafka
│   ├── logging/         # Zerolog
│   ├── redisx/          # Redis
│   ├── tracing/         # Jaeger
│   └── websocket/       # Socket.IO
├── service_layer/       # لایه سرویس
│   ├── cache/           # Redis Cache
│   ├── command_event_handler/  # CQRS Handlers
│   ├── messagebus/      # Message Bus
│   └── outbox/          # Outbox Pattern
└── specification/       # Specification Pattern
```

## وابستگی‌ها

- **Go**: 1.25.1+
- **GORM**: ORM برای دیتابیس
- **Fiber**: Web Framework
- **Zerolog**: Logger
- **Redis**: Client
- **Elasticsearch**: Client
- **Kafka**: Sarama
- **OpenTelemetry**: Tracing
- **Socket.IO**: WebSocket

## مجوز

این پروژه تحت مجوز MIT منتشر شده است.

## مشارکت

مشارکت‌ها، Issues و Pull Request‌ها خوش‌آمد هستند!

## نویسنده

[ali-mahdavi-dev](https://github.com/ali-mahdavi-dev)
