# SpotSync

SpotSync is a production-ready parking reservation REST API built with **Go**, **Echo**, **PostgreSQL**, and **GORM**. It follows **Clean Architecture** principles with strict separation between handlers, services, and repositories.

---

## Features

- JWT authentication with role-based authorization (`driver`, `admin`)
- Parking zone management with search, filtering, sorting, and pagination
- Reservation booking with **transactional row locking** (`FOR UPDATE`) to prevent overbooking
- Centralized error handling with standardized JSON responses
- Request validation using `validator.v10`
- Swagger API documentation
- Health check endpoint
- Graceful shutdown
- Connection pool configuration
- Soft delete support
- Database indexes for performance

---

## Tech Stack

| Component        | Technology              |
|------------------|-------------------------|
| Language         | Go 1.22+                |
| Web Framework    | Echo v4                 |
| Database         | PostgreSQL              |
| ORM              | GORM                    |
| Authentication   | JWT (golang-jwt/jwt/v5)   |
| Password Hashing | bcrypt                  |
| Validation       | go-playground/validator |
| Config           | godotenv                |
| API Docs         | swaggo/swag             |

---

## Project Structure

```
spotsync/
├── cmd/server/main.go          # Application entry point & dependency injection
├── config/                     # Environment configuration
├── database/                   # Database connection & auto-migration
├── dto/                        # Request/response data transfer objects
├── handler/                    # HTTP handlers (presentation layer)
├── middleware/                 # JWT, Admin, CORS, Logger, Recover
├── models/                     # GORM database models
├── repository/                 # Data access layer (GORM operations)
├── service/                    # Business logic layer
├── routes/                     # Route registration
├── utils/                      # JWT, password, validator utilities
├── pkg/
│   ├── response/               # Standardized JSON response helpers
│   └── errors/                 # Application error types
├── migrations/                 # SQL migration scripts
├── docs/                       # Swagger generated documentation
├── .env.example                # Environment variable template
├── go.mod
└── README.md
```

---

## Architecture

SpotSync follows **Clean Architecture** with unidirectional dependency flow:

```
HTTP Request
    ↓
Handler (Presentation)     ← Parses request, validates input, returns response
    ↓
Service (Business Logic)   ← Business rules, orchestration
    ↓
Repository (Data Access)   ← All GORM/database operations
    ↓
PostgreSQL
```

**Key principles:**
- Handlers never access the database directly
- Services contain all business logic
- Repositories handle all GORM operations
- Models represent database tables only
- DTOs decouple API contracts from internal models
- Dependencies are wired manually in `main.go`

---

## Prerequisites

- Go 1.22 or higher
- PostgreSQL 14+
- [swag](https://github.com/swaggo/swag) CLI (for regenerating Swagger docs)

---

## Installation

1. **Clone the repository**

```bash
git clone <repository-url>
cd spotsync
```

2. **Install Go dependencies**

```bash
go mod download
```

3. **Install Swagger CLI (optional, for doc regeneration)**

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

4. **Create PostgreSQL database**

```sql
CREATE DATABASE spotsync;
```

5. **Configure environment variables**

```bash
cp .env.example .env
```

Edit `.env` with your local settings:

```env
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=spotsync
DB_USER=postgres
DB_PASSWORD=postgres
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY_HOURS=24
```

---

## Run Locally

```bash
go run cmd/server/main.go
```

The server starts at `http://localhost:8080`.

- **Health check:** `GET http://localhost:8080/health`
- **Swagger UI:** `http://localhost:8080/swagger/index.html`

---

## Environment Variables

| Variable           | Description                          | Default     |
|--------------------|--------------------------------------|-------------|
| `PORT`             | HTTP server port                     | `8080`      |
| `DB_HOST`          | PostgreSQL host                      | `localhost` |
| `DB_PORT`          | PostgreSQL port                      | `5432`      |
| `DB_NAME`          | Database name                        | `spotsync`  |
| `DB_USER`          | Database user                        | `postgres`  |
| `DB_PASSWORD`      | Database password                    | `postgres`  |
| `JWT_SECRET`       | JWT signing secret (required)        | —           |
| `JWT_EXPIRY_HOURS` | Token expiry in hours                | `24`        |

---

## API Documentation

### Response Format

**Success:**
```json
{
  "success": true,
  "message": "operation successful",
  "data": {}
}
```

**Error:**
```json
{
  "success": false,
  "message": "error description",
  "errors": {}
}
```

---

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secret123",
  "role": "driver"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secret123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "role": "driver",
      "created_at": "2026-06-27T10:00:00Z"
    }
  }
}
```

---

### Parking Zones

#### Create Zone (Admin)
```http
POST /api/v1/zones
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Downtown Lot A",
  "location": "123 Main Street",
  "description": "Covered parking near city center",
  "capacity": 50,
  "hourly_rate": 5.00
}
```

#### List Zones
```http
GET /api/v1/zones?page=1&limit=10&search=downtown&sort_by=name&sort_dir=ASC
```

#### Get Zone by ID
```http
GET /api/v1/zones/1
```

---

### Reservations

#### Create Reservation (Driver)
```http
POST /api/v1/reservations
Authorization: Bearer <driver_token>
Content-Type: application/json

{
  "parking_zone_id": 1,
  "vehicle_number": "ABC-1234",
  "start_time": "2026-06-28T09:00:00Z",
  "end_time": "2026-06-28T17:00:00Z"
}
```

#### My Reservations
```http
GET /api/v1/reservations/my-reservations?page=1&limit=10
Authorization: Bearer <driver_token>
```

#### Cancel Reservation
```http
DELETE /api/v1/reservations/1
Authorization: Bearer <driver_token>
```

#### List All Reservations (Admin)
```http
GET /api/v1/reservations?page=1&limit=10&status=active
Authorization: Bearer <admin_token>
```

---

## Reservation Logic

Reservations use **PostgreSQL row-level locking** to prevent race conditions and overbooking:

1. Begin database transaction
2. Lock parking zone row with `FOR UPDATE`
3. Count overlapping active reservations for the time slot
4. If count ≥ capacity → return `409 Conflict`
5. Otherwise create reservation and commit

This guarantees no two concurrent requests can overbook a zone.

---

## Deployment

### Build Binary

```bash
go build -o spotsync cmd/server/main.go
```

### Run with Docker (PostgreSQL)

```bash
docker run -d \
  --name spotsync-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=spotsync \
  -p 5432:5432 \
  postgres:16-alpine
```

### Production Checklist

- [ ] Set a strong `JWT_SECRET` (32+ random characters)
- [ ] Use SSL for PostgreSQL (`sslmode=require` in DSN)
- [ ] Configure a reverse proxy (nginx/Caddy) with HTTPS
- [ ] Set appropriate connection pool limits for your workload
- [ ] Enable PostgreSQL backups
- [ ] Restrict CORS origins in production
- [ ] Run behind a process manager (systemd, Docker, Kubernetes)

### Regenerate Swagger Docs

```bash
swag init -g cmd/server/main.go -o docs
```

---

## License

MIT
