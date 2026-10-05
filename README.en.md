[Ler em Português](README.md)
# Scholarly

Academic management system built with Go and PostgreSQL.

---

## Overview

Scholarly is a REST API for academic management that handles students, teachers, courses, tests, enrollments, and grades. The project was developed as part of a Database course with the goal of demonstrating, in practice, the integration between an application and advanced PostgreSQL features.

The system addresses the need for a centralized academic management platform that goes beyond traditional CRUD operations, leveraging native database resources to process and consolidate information, such as weighted average calculation, consolidated student reports, and atomic enrollment with capacity validation.

---

## Project highlights

Beyond the minimum requirements, the project was built following practices commonly adopted in production environments:

- **Hexagonal architecture (ports and adapters)** with strict separation between domain, use cases, contracts, and infrastructure adapters.
- **Pure domain layer**, independent of frameworks, database drivers, or HTTP libraries, allowing business rules to be tested without external dependencies.
- **Password hashing with Argon2id**, using parameters above the OWASP minimum recommendation.
- **Layer-specific DTOs**, isolating the HTTP contract from domain entities and preventing exposure of sensitive fields such as `password_hash`.
- **PostgreSQL error translation**, mapping SQLSTATE codes (`23505`, `23503`, and custom application codes such as `P0S01`, `P0C01`) into domain errors.
- **Interactive documentation with Swagger UI**, allowing all endpoints to be explored and tested without additional tooling.
- **Docker Compose** with healthcheck for reliable database orchestration.
- **Continuous integration with GitHub Actions**, running `go vet`, `staticcheck`, and automated tests on every push.
- **Hot reload in development** using Air.

---

## Tech stack

| Layer | Technology |
|---|---|
| Language | Go 1.27 |
| Database | PostgreSQL 17 |
| HTTP router | chi |
| PostgreSQL driver | pgx/v5 |
| Password hashing | Argon2id |
| API documentation | Swagger (swaggo) |
| Containerization | Docker and Docker Compose |
| Testing | Testify and Testcontainers |
| CI | GitHub Actions |
| Development tooling | Air |
| Database Migrations | golang-migrate |

---

## Architecture

The project follows the principles of **hexagonal architecture**, with dependencies flowing exclusively inward:

```
┌─────────────────────────────────────────┐
│  HTTP (input adapter)                   │
│  ────────────────────────────────────   │
│  Service (use cases)                    │
│  ────────────────────────────────────   │
│  Ports (interfaces)                     │
│  ────────────────────────────────────   │
│  Domain (entities and rules)            │
└─────────────────────────────────────────┘
         ▲
         │  Output adapters (Postgres, Argon2)
```

The domain layer has no knowledge of HTTP, databases, or any infrastructure library. Services depend only on interfaces (`ports`), and adapters implement those interfaces — allowing PostgreSQL to be replaced by another database, or chi by another router, without changes to the core.

### Project structure

```
Scholarly/
├── cmd/
│   └── api/                 # Application entrypoint
├── internal/
│   ├── domain/              # Entities and business rules
│   ├── ports/               # Interfaces (contracts)
│   ├── service/             # Use cases
    ├── env/                 # Environment variable (.env) support
    ├── config/              # Basic configurations
│   └── adapters/
│       ├── http/            # HTTP handlers, DTOs, middlewares
│       ├── postgres/        # Repositories
│       └── argon2/          # Password hashing implementation
├── database/
│   ├── tables/              # Table creation scripts
│   ├── views/               # View creation scripts
│   ├── functions/           # Function creation scripts
│   ├── procedures/          # Procedure creation scripts
│   └── inserts/             # Seed data
├── docs/                    # Generated Swagger documentation
├── docker-compose.yml
├── Makefile
└── README.md
```

---

## Database

### DBMS

**PostgreSQL 17**, running through Docker.

### Main tables

| Table | Description |
|---|---|
| `users` | Central identity for students, teachers, and administrators |
| `students` | Student-specific data (table inheritance from `users`) |
| `teachers` | Teacher-specific data (table inheritance from `users`) |
| `courses` | Courses offered per semester |
| `enrollments` | Student enrollments in courses |
| `tests` | Tests per course, with weight |
| `grades` | Each student's grade per test |

The model uses **table inheritance** to separate identity (`users`) from specializations (`students`, `teachers`), avoiding duplication of authentication data and allowing a single user to hold multiple roles.

### View

**`vw_student_report_card`**

Consolidates, in a single row per student and course, the academic report data: student name, course, teacher, weighted average, and status (approved, recovery, or failed). The View encapsulates multiple joins across `users`, `students`, `enrollments`, `courses`, and `teachers`, and integrates the `fn_weighted_average` function for the average calculation.

**Usage in the application:** `GET /students/{student_id}/report`

### Function

**`fn_weighted_average(p_student_id BIGINT, p_course_id BIGINT)`**

Calculates the weighted average of a student in a course, taking into account the weight of each test. Returns a `NUMERIC` value — `0` when the student has no grades, or raises an exception (`P0S01`, `P0C01`) when the student or course does not exist.

**Usage in the application:**
- Directly through the endpoint `GET /students/{student_id}/courses/{course_id}/average`
- Internally through the View `vw_student_report_card`

### Procedure

**`sp_enroll_student(p_student_id BIGINT, p_course_id BIGINT, OUT p_enrollment_id BIGINT)`**

Enrolls a student in a course with atomic validations:

1. Student exists
2. Course exists
3. Student is not already enrolled
4. Course has available seats

Returns the created enrollment `id` through an `OUT` parameter.

**Usage in the application:** `POST /enrollments`

---

## API endpoints

### Infrastructure

| Method | Route | Description |
|---|---|---|
| GET | `/health` | Health check |
| GET | `/swagger/*` | Swagger UI |

### Users

| Method | Route | Description |
|---|---|---|
| POST | `/users` | Create user |
| GET | `/users` | List users with filters |
| GET | `/users/{id}` | Get user by ID |
| PUT | `/users/{id}` | Update user |
| DELETE | `/users/{id}` | Delete user |

### Students

| Method | Route | Description |
|---|---|---|
| POST | `/students` | Create student |
| GET | `/students/enrollment/{enrollment}` | Get student by enrollment number |
| GET | `/students/{id}` | Get student by ID |
| PUT | `/students/{id}` | Update student |
| DELETE | `/students/{id}` | Delete student |

### Teachers

| Method | Route | Description |
|---|---|---|
| POST | `/teachers` | Create teacher |
| GET | `/teachers` | List teachers |
| GET | `/teachers/{id}` | Get teacher by ID |
| PUT | `/teachers/{id}` | Update teacher |
| DELETE | `/teachers/{id}` | Delete teacher |

### Courses

| Method | Route | Description |
|---|---|---|
| POST | `/courses` | Create course |
| GET | `/courses` | List courses |
| GET | `/courses/{id}` | Get course by ID |
| PUT | `/courses/{id}` | Update course |
| DELETE | `/courses/{id}` | Delete course |
| POST | `/courses/enroll/`| Enrolls student in a course (utilizing a procedure)| 

### Enrollments

| Method | Route | Description |
|---|---|---|

| GET | `/enrollments` | List enrollments |
| GET | `/enrollments/{id}` | Get enrollment by ID |
| PUT | `/enrollments/{id}` | Update enrollment |
| DELETE | `/enrollments/{id}` | Delete enrollment |

### Tests

| Method | Route | Description |
|---|---|---|
| POST | `/tests` | Create test |
| GET | `/tests` | List tests |
| GET | `/tests/{id}` | Get test by ID |
| PUT | `/tests/{id}` | Update test |
| DELETE | `/tests/{id}` | Delete test |

### Grades

| Method | Route | Description |
|---|---|---|
| POST | `/grades` | Create grade |
| GET | `/grades` | List grades |
| GET | `/grades/{id}` | Get grade by ID |
| PUT | `/grades/{id}` | Update grade |
| DELETE | `/grades/{id}` | Delete grade |

### Reports

| Method | Route | Database resource |
|---|---|---|
| GET | `/students/{student_id}/report` | View `vw_student_report_card` |
| GET | `/students/{student_id}/courses/{course_id}/average` | Function `fn_weighted_average` |

Full documentation with payloads and response codes is available through Swagger UI at `http://localhost:8080/swagger/index.html`.

---

## Getting started

### Requirements

- Go 1.27 or later
- Docker and Docker Compose
- Make
- golang-migrate (https://github.com/golang-migrate/migrate)

### Steps

**1. Clone the repository**

```bash
git clone https://github.com/alvarolucio2007/Scholarly.git
cd Scholarly
```

**2. Start the database**

```bash
docker compose up -d
```

The container includes a healthcheck that validates database availability before accepting connections.

**3. Apply the schema**

```bash
make migrateup
```

This executes the table, View, Function, and Procedure creation scripts in the correct order.

**5. Run the application**

```bash
make run
```

For development with hot reload:

```bash
air
```

**6. Access Swagger UI**

```
http://localhost:8080/swagger/index.html
```

---

## Design decisions

### Hexagonal architecture

Hexagonal architecture was chosen to isolate the domain from external dependencies. This enables business rules to be tested without a database, adapters to be replaced without impacting the core, and the project to evolve without accumulating coupling.

### Input and output DTOs

The application uses HTTP-specific DTOs, kept separate from domain entities. This avoids exposing sensitive fields such as `password_hash`, decouples the API contract from internal entities, and simplifies Swagger documentation generation.

### Domain errors

Errors are defined in the `domain` package (`ErrNotFound`, `ErrEmailAlreadyExists`, `ErrCourseFull`, among others) and translated in the PostgreSQL adapter from SQLSTATE codes. The HTTP handler maps each domain error to the appropriate HTTP status, keeping the core independent from infrastructure details.

### Use of PostgreSQL features

The View, Function, and Procedure were chosen to solve problems that make sense at the database level:

- **View:** consolidated queries with multiple joins and aggregations that would otherwise be duplicated across the application.
- **Function:** weighted average calculation, executed both directly and inside the View.
- **Procedure:** student enrollment, involving multiple validations and an insert that must be atomic.

---

## Evolution from the previous semester

This project is a rewrite of an academic system developed in the previous semester, now focused on architectural quality and advanced database usage.

| Aspect | Previous version | Current version |
|---|---|---|
| Architecture | Layered Architecture | Hexagonal (ports and adapters) |
| Persistence | No foreign keys, no indexes | Foreign keys, indexes, constraints |
| SQL features | No View, Function, or Procedure | View, Function, and Procedure integrated |
| Password Hashing | None | Argon2id |
| Documentation | None | Swagger |
| Testing | Unit Tests | Unit tests |
| CI/CD | None | GitHub Actions |
| Containerization | Docker and Docker Compose | Docker and Docker Compose |

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
