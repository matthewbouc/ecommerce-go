# ecommerce-go

A learning-focused, production-style ecommerce backend built in Go. The goal is to explore backend architecture and distributed systems by building a realistic, scalable system from the ground up — comparing REST, GraphQL, and gRPC in real-world use cases, and eventually adding a frontend.

## Architecture

The system is designed as a set of independent microservices. Each service owns its own data and communicates via:

- **gRPC** — internal synchronous service-to-service calls
- **Kafka** — asynchronous event-driven communication

External APIs are exposed via:

- **REST** — public-facing endpoints
- **GraphQL** — aggregation layer for frontend clients (planned)

```
┌─────────────────┐     ┌─────────────────┐      ┌─────────────────┐
│   User Service  │     │Catalog Service  │      │  Order Service  │
│  (auth, seller  │     │(products, search│      │(order lifecycle)│
│   onboarding)   │     │ Elasticsearch)  │      │                 │
└────────┬────────┘     └────────┬────────┘      └────────┬────────┘
         │                       │                        │
         └───────────────────────┴───────────── Kafka ────┘
                                                    │
                              ┌─────────────────────┴──────────────────┐
                              │         Notification Service           │
                              │          (Email, SMS via Twilio)       │
                              └────────────────────────────────────────┘
```

## Tech Stack

| Concern        | Technology                     |
|----------------|--------------------------------|
| Language       | Go 1.25                        |
| HTTP Framework | Fiber v3                       |
| Database       | PostgreSQL (via GORM)          |
| Messaging      | Kafka *(planned)*              |
| Search         | Elasticsearch *(planned)*      |
| Payments       | Stripe *(planned)*             |
| SMS            | Twilio                         |
| Email          | SendGrid / AWS SES *(planned)* |
| Auth           | JWT (HS256)                    |

## Services

| Service      | Status      | Description                                                         |
|--------------|-------------|---------------------------------------------------------------------|
| User         | In progress | Auth, profiles, phone verification, seller onboarding, cart, orders |
| Catalog      | Planned     | Product CRUD, search indexing                                       |
| Order        | Planned     | Order lifecycle, checkout                                           |
| Transaction  | Planned     | Stripe payments, transaction history                                |
| Notification | Planned     | Email and SMS notifications                                         |

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- [nodemon](https://nodemon.io/) *(optional, for hot-reload)*
- A [Twilio](https://www.twilio.com/) account for SMS

### 1. Start the database

```bash
docker-compose up -d
```

This starts PostgreSQL on port `5432` with:
- User: `root`
- Password: `root`
- Database: `ecommerce-go`

### 2. Configure environment

Copy the example env file and fill in your values:

```bash
cp .env.example .env.dev
```

| Variable               | Description                          |
|------------------------|--------------------------------------|
| `HTTP_PORT`            | Server address e.g. `localhost:3000` |
| `DATABASE_CONFIG`      | PostgreSQL DSN                       |
| `AUTH_SECRET`          | JWT signing secret                   |
| `TWILIO_ACCOUNT_SID`   | Twilio account SID                   |
| `TWILIO_AUTH_TOKEN`    | Twilio auth token                    |
| `TWILIO_PHONE_NUMBER`  | Twilio sender phone number           |

### 3. Run the server

```bash
# With hot-reload (requires nodemon)
make server

# Single run
make run
```

Database migrations run automatically on startup via GORM `AutoMigrate`.

## API Reference

All private endpoints require a `Bearer` token in the `Authorization` header. Register or login to receive one. A Postman collection is available at `.dev/postman.txt`.

### Auth

| Method | Endpoint         | Auth | Description                          |
|--------|------------------|------|--------------------------------------|
| POST   | `/user/register` |      | Register a new user, returns JWT     |
| POST   | `/user/login`    |      | Login, returns JWT                   |

### User

| Method | Endpoint        | Auth | Description                                 |
|--------|-----------------|------|---------------------------------------------|
| DELETE | `/user`         | Yes  | Soft-delete the current user                |
| GET    | `/user/profile` | Yes  | Get current user's profile                  |
| POST   | `/user/profile` | Yes  | Update profile (firstName, lastName, phone) |

### Verification

| Method | Endpoint        | Auth | Description                                  |
|--------|-----------------|------|----------------------------------------------|
| GET    | `/user/verify`  | Yes  | Generate and SMS a 6-digit verification code |
| POST   | `/user/verify`  | Yes  | Submit the code to verify the account        |

### Seller

| Method | Endpoint               | Auth | Description                                       |
|--------|------------------------|------|---------------------------------------------------|
| POST   | `/user/become-seller`  | Yes  | Upgrade account to seller and add bank account    |

### Cart

| Method | Endpoint     | Auth | Description                                                  |
|--------|--------------|------|--------------------------------------------------------------|
| GET    | `/user/cart` | Yes  | Get all cart items for the current user                      |
| POST   | `/user/cart` | Yes  | Add an item to cart (increments quantity if already present) |

### Orders

| Method | Endpoint          | Auth | Description                          |
|--------|-------------------|------|--------------------------------------|
| GET    | `/user/order`     | Yes  | List all orders for the current user |
| GET    | `/user/order/:id` | Yes  | Get a single order by UUID           |

> **Note:** Order creation (checkout) is not yet implemented. It depends on the Catalog Service (price validation) and Transaction Service (Stripe).

## Project Structure

```
.
├── main.go
├── config/                         # Env loading and AppConfig struct
├── docker-compose.yml
├── Makefile
├── internal/
│   ├── api/
│   │   ├── server.go               # DB connect, migrations, Fiber init, route wiring
│   │   └── rest/
│   │       ├── httpHandler.go      # RestHandler — shared dependencies
│   │       └── handlers/
│   │           ├── userHandler.go
│   │           ├── catalogHandler.go    # stub
│   │           └── transactionHandler.go # stub
│   ├── domain/                     # GORM model structs
│   │   ├── User.go
│   │   ├── BankAccount.go
│   │   ├── CartItem.go
│   │   └── Order.go                # Order + OrderItem
│   ├── dto/                        # Request/response structs
│   ├── service/
│   │   └── userService.go          # Business logic
│   ├── repository/                 # GORM-backed data access
│   │   ├── userRepository.go
│   │   ├── bankAccountRepository.go
│   │   ├── cartRepository.go
│   │   └── orderRepository.go
│   └── helper/
│       ├── auth.go                 # JWT, bcrypt, middleware
│       └── randomNumber.go         # Verification code generator
└── pkg/
    └── notification/
        └── sms/
            ├── SmsClient.go        # Interface
            └── provider/
                └── twilioClient.go
```

## Authentication Flow

1. Register or login → receive a JWT (24hr expiry)
2. Include the token as `Authorization: Bearer <token>` on private routes
3. The `Authorize` middleware validates the token and injects the current user into request context
4. Calling `POST /user/become-seller` issues a new token with the `seller` role

## Development

```bash
# Build
go build ./...

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/service/... -run TestFunctionName
```

## Roadmap

- [ ] Catalog Service — product CRUD + Elasticsearch search indexing
- [ ] Order Service — checkout flow with cart-to-order conversion
- [ ] Transaction Service — Stripe payment processing
- [ ] Notification Service — Kafka consumer for email/SMS events
- [ ] gRPC contracts between services
- [ ] GraphQL gateway for frontend aggregation
- [ ] Infrastructure as code (Terraform)
- [ ] Frontend
