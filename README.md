#SpotSync – SERVER

A high-performance backend service built with Go, Echo, GORM, and PostgreSQL for managing parking reservations and EV charging stations. The platform enables drivers to reserve parking spaces while allowing administrators to manage parking zones, pricing, and reservations through a secure role-based system.

<!-- - **[LIVE LINK](https://)** -->
<!-- - **[CLIENT REPOSITORY](https://github.com/Captain-Kanak/)** -->

---

## Overview

SpotSync is a centralized parking management system designed for environments with limited parking capacity such as airports, shopping malls, and commercial complexes.

The backend provides:

- Secure JWT Authentication
- Role-Based Access Control (Driver & Admin)
- Parking Zone Management
- EV Charging Reservation System
- Concurrency-safe Reservation Processing
- Transactional Database Operations
- Clean Architecture with Dependency Injection

The system guarantees that parking zones never exceed their configured capacity, even when multiple users attempt to reserve the last available space simultaneously.

---

## Technology Stack

- Language: GO
- Framework: Echo
- Database: PostgreSQL (Neon DB)
- ORM: GORM
- Authentication: JWT

---

# Authentication & Authorization

### Authentication is implemented using JWT.

After successful login:

- Server generates a signed JWT
- Client stores the token
- Protected routes require:

Authorization: Bearer <JWT_TOKEN>

### Role-Based Access Control (RBAC)

- ADMIN – Create, update, and delete parking zones, configure pricing, view all reservations, and manage overall parking operations.
- DRIVER – Register and log in, browse parking zones and real-time availability, reserve parking or EV charging spots, view personal reservations, and cancel their own reservations.

---

## Security Considerations

The application includes several security mechanisms:

- Passwords are hashed using bcrypt.
- JWT is used for stateless authentication.
- Request validation prevents invalid payloads.
- Protected routes require authentication.
- Role verification prevents unauthorized actions.
- Sensitive fields (passwords) are never returned in API responses.

### Concurrency-safe Reservation

Parking reservation is implemented using:

- Database Transactions
- Row-Level Locking (FOR UPDATE)

This prevents race conditions when multiple users reserve the last available parking spot simultaneously.

---

## Installation & Setup

Prerequisites:

- Go 1.22+
- PostgreSQL / (Neon DB)
- Git

Clone Repository:

```bash
git clone
cd SpotSync
```

Install Dependencies:

```bash
go mod tidy
```

Environment Variables:
Create a `.env` file in the root of your project and add the following:

```env
ENV="development"
PORT="8080"
DSN='database-source-name'
JWT_SECRET="jwt-secret"
FRONTEND_URL="http://localhost:3000"
```

---

Start Server:

```bash
go run cmd/main.go
```

---

## 👤 Author

**zingrin**

> Software Engineer

- TypeScript
- Express.js
- PostgreSQL
- Docker
- GO
- Echo

---
