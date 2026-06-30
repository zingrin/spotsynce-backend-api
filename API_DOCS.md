# 🔗 API Endpoints

### Authentication

| Method | Endpoint              | Access | Description            |
| ------ | --------------------- | ------ | ---------------------- |
| POST   | /api/v1/auth/register | PUBLIC | register new user      |
| POST   | /api/v1/auth/login    | PUBLIC | log in registered user |

---

### Zone

| Method | Endpoint          | Access | Description                     |
| ------ | ----------------- | ------ | ------------------------------- |
| POST   | /api/v1/zones     | ADMIN  | create a new zone               |
| GET    | /api/v1/zones     | PUBLIC | get all zones                   |
| GET    | /api/v1/zones/:id | PUBLIC | get zone by id                  |
| PATCH  | /api/v1/zones/:id | ADMIN  | update zone by id               |
| DELETE | /api/v1/zones/:id | ADMIN  | delete zone by id (soft delete) |

---

### Reservation

| Method | Endpoint                             | Access         | Description               |
| ------ | ------------------------------------ | -------------- | ------------------------- |
| POST   | /api/v1/reservations                 | ADMIN / DRIVER | create as new reservation |
| GET    | /api/v1/reservations/my-reservations | ADMIN / DRIVER | get own reservations      |
| GET    | /api/v1/reservations                 | ADMIN          | get all reservations      |
| DELETE | /api/v1/reservations/:id             | ADMIN / DRIVER | cancel reservation by id  |

---
