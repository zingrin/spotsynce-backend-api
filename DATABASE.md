# 🗄️ Database Schema

<!-- **[ERD LINK](https://)** -->

### User

| Field     | Description     |
| --------- | --------------- |
| id        | UUID (PK)       |
| name      | String          |
| email     | String (unique) |
| password  | String          |
| role      | admin / driver  |
| phone     | String?         |
| createdAt | DateTime        |
| updatedAt | DateTime        |
| deletedAt | DateTime?       |

### Zone

| Field          | Description                     |
| -------------- | ------------------------------- |
| id             | UUID (PK)                       |
| name           | String                          |
| type           | general / ev_charging / covered |
| total_capacity | Int                             |
| price_per_hour | float64                         |
| createdAt      | DateTime                        |
| updatedAt      | DateTime                        |
| deletedAt      | DateTime?                       |

### Reservation

| Field         | Description                   |
| ------------- | ----------------------------- |
| id            | UUID (PK)                     |
| userId        | UUID (FK)                     |
| zone_id       | UUID (FK)                     |
| license_plate | String                        |
| Status        | ACTIVE / COMPLETED / CANCELED |
| createdAt     | DateTime                      |
| updatedAt     | DateTime                      |
| deletedAt     | DateTime?                     |
