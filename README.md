# Go Marketplace

Backend marketplace-приложения на Go с PostgreSQL, Redis, JWT-аутентификацией, ролями, корзиной, заказами и загрузкой изображений.

Проект создан как практический pet-project для углубления навыков backend-разработки на Go и работы с PostgreSQL, транзакциями, HTTP API и архитектурой приложения.

## Features

* Регистрация и аутентификация пользователей
* JWT access / refresh tokens
* Refresh token rotation
* Redis для хранения refresh tokens
* Role-based authorization
* Пользователи и роли
* Категории товаров
* Создание, изменение и удаление товаров
* Поиск и пагинация товаров
* Сортировка товаров
* Загрузка изображений товаров
* Корзина пользователя
* Изменение количества товаров в корзине
* Создание заказов
* Транзакции через Unit of Work
* PostgreSQL
* Redis
* Rate limiting
* Structured logging через `slog`
* Request ID
* HTTP middleware
* Graceful shutdown
* Swagger
* Unit и integration tests

## Tech Stack

### Backend

* Go
* Chi
* PostgreSQL
* SQLX
* Redis
* JWT
* Decimal arithmetic
* Swagger

### Infrastructure

* Docker
* Docker Compose

## Architecture

Проект построен вокруг разделения ответственности между HTTP-слоем, бизнес-логикой и инфраструктурой.

```text
HTTP Handler
     │
     ▼
  Service
     │
     ▼
Repository
     │
     ▼
 PostgreSQL
```

Дополнительные инфраструктурные компоненты:

```text
                    ┌─────────────┐
                    │    HTTP     │
                    │   Handler   │
                    └──────┬──────┘
                           │
                           ▼
                    ┌─────────────┐
                    │   Service   │
                    └──────┬──────┘
                           │
             ┌─────────────┴─────────────┐
             │                           │
             ▼                           ▼
      ┌─────────────┐             ┌─────────────┐
      │ Repository  │             │    Redis     │
      └──────┬──────┘             └─────────────┘
             │
             ▼
      ┌─────────────┐
      │ PostgreSQL  │
      └─────────────┘
```

Для операций, которые изменяют несколько связанных сущностей, используется Unit of Work.

## Project Structure

```text
backend/
├── config/
├── constants/
├── internal/
│       ├── config/
│       ├── domain/
│       ├── infrastructure/
│       ├── logger/
│       ├── repository/
│       ├── domain/
│       ├── service/
│       ├── transport/
│       │   ├── dto/
│       │   ├── handler/
│       │   ├── middleware/
│       │   └── validator/
│       └── utils/
├── docs/
├── migrations/
├── tests/
├── uploads/
├── Dockerfile
├── main.go
└── go.mod
```

## Domain

Основные сущности приложения:

* User
* Role
* Category
* Product
* Cart
* CartItem
* Order
* OrderItem
* RefreshToken
* AccessToken

Основные пользовательские сценарии:

```text
User
 │
 ├── Authentication
 │
 ├── Browse products
 │
 ├── Add product to cart
 │
 ├── Change cart quantity
 │
 └── Create order
```

Для продавца:

```text
Seller
 │
 ├── Create product
 ├── Update product
 ├── Delete product
 └── Manage own products
```

## Transactions

Для операций, затрагивающих несколько репозиториев, используется Unit of Work.

Пример создания заказа:

```text
BEGIN
  │
  ├── Get cart
  ├── Lock required data
  ├── Get cart items
  ├── Calculate order total
  ├── Create order
  ├── Create order items
  └── Clear cart
  │
COMMIT
```

При ошибке:

```text
BEGIN
  │
  ├── ...
  ├── error
  │
ROLLBACK
```

## Money

Для денежных значений используется decimal arithmetic вместо `float64`.

```go
decimal.Decimal
```

Это позволяет избежать типичных проблем бинарной арифметики floating-point чисел при работе с ценами и суммами.
Цена товара фиксируется в `OrderItem` в момент создания заказа.
Таким образом, изменение цены товара после оформления заказа не изменяет стоимость уже созданного заказа.

## Orders

Заказ содержит snapshot необходимых данных товара:

* product price
* product name
* seller
* quantity
* total

Это позволяет сохранить информацию, необходимую для отображения исторического заказа, даже если исходный товар впоследствии изменён или удалён.

## Authentication

Используется JWT-аутентификация.

Основные endpoints:

```text
POST /api/v1/login
POST /api/v1/refresh
POST /api/v1/logout
```

Используются два типа токенов:

```text
Access Token
    │
    └── используется для авторизации API-запросов

Refresh Token
    │
    └── используется для получения нового Access Token
```

Refresh tokens хранятся в Redis.
Присутствует чёрный список access tokens.

## Authorization

Для ограничения доступа используется RBAC authorization.

Примеры:

```text
User
 ├── browse products
 ├── manage own cart
 └── create orders

Seller
 ├── create products
 ├── update own products
 └── delete own products

Admin
 └── administrative operations
```

Проверка прав выполняется как на HTTP-уровне, так и в бизнес-логике там, где это необходимо.

## Rate Limiting

Для отдельных чувствительных или потенциально дорогих операций используется rate limiting.

Например, ограничивается частота создания товаров.

Redis используется как хранилище состояния rate limiter.

## File Uploads

Изображения товаров сохраняются отдельно от данных PostgreSQL.

Общий flow:

```text
HTTP multipart request
        │
        ▼
 Validate file
        │
        ▼
 Detect content type
        │
        ▼
 Save file
        │
        ▼
 Store file reference
        │
        ▼
 PostgreSQL
```

При ошибках бизнес-операции выполняется удаление уже сохранённого файла.

Файлы не хранятся непосредственно в PostgreSQL.

## Database

Основная база данных — PostgreSQL.

Используется SQLX для работы с SQL напрямую.

Вместо ORM основной упор сделан на:

* SQL queries
* transactions
* constraints
* foreign keys
* indexes
* atomic updates
* row-level locking

Пример ограничения:

```sql
UNIQUE(cart_id, product_id)
```

гарантирует, что один и тот же товар не может находиться в корзине пользователя более одного раза.

## Database Integrity

Часть бизнес-инвариантов обеспечивается непосредственно базой данных.

Например:

```sql
CHECK(quantity > 0)
```

и:

```sql
UNIQUE(cart_id, product_id)
```

Также используются foreign keys и соответствующие `ON DELETE` стратегии.

## Product Listing

Получение товаров поддерживает:

* offset pagination
* search
* category filtering
* sorting

Для сортировки используется whitelist допустимых SQL-полей, чтобы пользовательский input не мог напрямую попасть в SQL query.

Пример:

```text
GET /api/v1/products?page=1
GET /api/v1/products?search=phone
GET /api/v1/products?category_id=1
GET /api/v1/products?sort_by=price
```

## API

Основные группы endpoints:

```text
/api/v1
│
├── /login
├── /refresh
├── /logout
│
├── /users
├── /roles
├── /categories
├── /products
├── /seller/products
├── /cart
├── /orders
│
├── /uploads
└── /swagger
```

Health check:

```text
GET /health
```

## Testing

В проекте используются:

* unit tests
* repository integration tests

Для unit-тестов зависимости сервисов представлены интерфейсами.

Это позволяет тестировать бизнес-логику без необходимости поднимать PostgreSQL или Redis для каждого теста.

Repository tests проверяют взаимодействие с реальной базой данных.

## Error Handling
Используются Sentinel ошибки.
Ошибки разделяются между слоями приложения.

```text
Repository
    │
    ▼
Domain / Service error
    │
    ▼
HTTP Handler
    │
    ▼
HTTP status + response
```

## Logging

Для structured logging используется стандартный пакет:

```go
log/slog
```

В запросах используется request ID, позволяющий связать несколько log entries с одним HTTP request.

Основные middleware:

* request logging
* request ID
* authentication
* CORS
* body size limit
* recoverer
* timeout

## Graceful Shutdown

Приложение корректно обрабатывает завершение процесса и освобождает используемые ресурсы.

Основные ресурсы:

* HTTP server
* PostgreSQL connection pool
* Redis connection
* другие инфраструктурные зависимости

## Running Locally

### Requirements

Перед запуском необходимо установить:

* Go
* Docker
* Docker Compose

### Clone

```bash
git clone https://github.com/CosmosDevX/Go-Marketplace.git
cd Go-Marketplace
```

### Environment

Создайте `.env` файл и укажите необходимые параметры подключения к PostgreSQL, Redis и остальные настройки приложения.

Пример структуры в .env.example

Названия переменных должны соответствовать конфигурации проекта.

### Docker Compose

Запустить инфраструктуру:

```bash
docker compose up -d
```

### Run Backend

```bash
cd backend
go run .
```

### Tests

Запуск тестов:

```bash
go test ./...
```

## Swagger

Swagger используется для документирования HTTP API.

После запуска приложения документация доступна через:

```text
/api/v1/swagger
```

## Design Goals

Основные цели проекта:

* понимать границы ответственности между слоями
* работать с PostgreSQL без ORM
* корректно использовать SQL transactions
* использовать Redis не только как cache
* проектировать HTTP API
* писать unit и integration tests
* применять database constraints как часть бизнес-инвариантов
* работать с authentication и authorization
* практиковать dependency inversion через interfaces

## Known Limitations

Проект находится в стадии разработки и не является production-ready marketplace.

Некоторые вещи намеренно упрощены.

В частности:

* полноценная inventory/stock система пока не реализована;
* отсутствует полноценная payment integration;
* файловое хранилище пока реализовано через локальную файловую систему;
* pagination основана на offset pagination;
* поиск товаров использует SQL `LIKE`;
* нет полноценной системы доставки;
* нет распределённой очереди событий;
* нет полноценного order state machine;

## What I Practiced

В рамках проекта я практиковал:

### Go

* interfaces
* dependency injection
* context propagation
* error handling
* HTTP servers
* middleware
* testing

### PostgreSQL

* transactions
* foreign keys
* constraints
* indexes
* atomic updates
* pagination
* SQL joins

### Backend Architecture

* handler / service / repository separation
* dependency inversion
* Unit of Work
* domain validation
* error propagation

### Infrastructure

* Docker
* Redis
* JWT
* rate limiting
* structured logging
* graceful shutdown

## Future Improvements

Планируемые направления развития:

* inventory management
* payment integration
* explicit checkout idempotency
* object storage вместо локальной файловой системы
* cursor-based pagination
* более эффективный поиск
* order state machine
* background jobs
* event-driven interactions
* расширенное observability
* metrics
* tracing

## License

This project is intended for educational and portfolio purposes.