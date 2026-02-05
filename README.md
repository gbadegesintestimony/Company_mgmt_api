# Company Management API

A production-grade REST API built with Go for managing company employees, authentication, and administrative operations. This project demonstrates enterprise-level software architecture patterns, security best practices, and comprehensive testing strategies.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup & Installation](#setup--installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Development](#development)
- [Testing](#testing)
- [Database Migrations](#database-migrations)
- [Security Considerations](#security-considerations)
- [Deployment](#deployment)
- [Contributing](#contributing)

## Overview

The Company Management API provides a robust backend infrastructure for managing organizational data with features including:

- **User Authentication & Authorization**: JWT-based authentication with role-based access control (RBAC)
- **Employee Management**: CRUD operations for employee records with audit trails
- **Security Features**: Password management, OTP verification, email verification, and session management
- **Rate Limiting**: Built-in protection against abuse
- **Comprehensive Logging**: Structured logging with request tracking
- **Database Migrations**: Version-controlled schema evolution

**Tech Stack:**

- **Language**: Go 1.24.4
- **HTTP Framework**: Chi (chi/v5)
- **Database**: PostgreSQL 15
- **Authentication**: JWT (golang-jwt)
- **Email Service**: Resend API
- **Migration Tool**: golang-migrate
- **Testing**: Testify, sqlmock
- **Logging**: Zerolog

## Architecture

### Project Structure

```
cmd/                          # Application entry point
├── main.go

config/                       # Configuration management
├── config.go

database/                     # Database connectivity and migrations
├── postgres.go
└── migrations/               # Versioned SQL migrations

handlers/                     # HTTP request handlers (controller layer)
├── auth.go
├── admin_employee_handler.go
├── password_handler.go
├── verification_handler.go
├── me_handler.go
├── health.go
└── *_integration_test.go     # Integration tests

middlewares/                  # HTTP middleware
├── auth.go                   # JWT validation
├── cors.go                   # CORS configuration
├── rbac.go                   # Role-based access control
└── request_id.go             # Request tracking

services/                     # Business logic layer
├── auth_service.go
├── employee_service.go
├── email_service.go
├── otp_service.go
└── rate_limiter.go

repositories/                 # Data access layer (repository pattern)
├── user_repositories.go
├── employee_repository.go
├── session_repository.go
├── otp_repository.go
├── audit_repository.go
└── *_test.go                 # Unit tests

models/                       # Data models
├── user.go
├── employee.go
├── profile.go
└── otp.go

routes/                       # Route definitions
├── routes.go                 # Main router setup
├── auth.go
├── employee_routes.go
├── password_routes.go
└── verification_routes.go

utils/                        # Utility functions
├── jwt.go                    # JWT token generation/validation
├── password.go               # Password hashing and verification
├── otp.go                    # OTP generation logic
├── token.go                  # Token utilities
├── errors.go                 # Error handling
└── *_test.go                 # Unit tests

logger/                       # Logging configuration
└── logger.go                 # Zerolog setup

internal/                     # Internal packages not exported
└── testutils/                # Test helpers
    ├── db.go
    ├── email.go
    └── http.go
```

### Architectural Patterns

**Layered Architecture:**

- **Handlers**: Parse requests and invoke services
- **Services**: Encapsulate business logic and orchestrate repositories
- **Repositories**: Abstract data access with a consistent interface
- **Models**: Define domain entities
- **Middlewares**: Cross-cutting concerns (auth, logging, CORS)

**Key Design Principles:**

- Separation of concerns
- Dependency injection via constructor parameters
- Repository pattern for data access abstraction
- Middleware chain for request preprocessing
- Error handling with custom error types

## Features

### Authentication & Authorization

- JWT-based token authentication (access + refresh tokens)
- Role-based access control (RBAC) with middleware enforcement
- Session management with database persistence
- Token refresh mechanism

### User Management

- User registration and login
- Email verification workflow
- Password management (change, reset)
- OTP-based verification
- User profile management

### Employee Management

- Full CRUD operations for employee records
- Admin operations with elevated privileges
- Comprehensive audit trail logging
- Request ID tracking for debugging

### Security

- Bcrypt password hashing
- JWT secret-based authentication
- CORS middleware
- Request ID middleware for tracing
- Rate limiting per endpoint
- OTP generation and validation

### Observability

- Structured logging with request context
- Request ID propagation
- Audit trails for sensitive operations
- Health check endpoint

## Prerequisites

### System Requirements

- **Go**: 1.24.4 or higher
- **PostgreSQL**: 15.x or higher
- **Docker & Docker Compose** (recommended)
- **Git**

### Environment Variables

Create a `.env` file in the project root:

```env
# Application
APP_ENV=development
APP_PORT=8080

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=company_mgmt
DB_SSLMODE=disable

# JWT
JWT_ACCESS_SECRET=your_secret_key_for_access_tokens_min_32_chars
JWT_REFRESH_SECRET=your_secret_key_for_refresh_tokens_min_32_chars

# Email Service
EMAIL_FROM=noreply@company.com
RESEND_API_KEY=your_resend_api_key
```

## Setup & Installation

### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone <repository-url>
cd Company_mgmt_api

# Build and start services
docker-compose up -d

# Run database migrations (if needed)
docker-compose exec api ./api migrate
```

The application will be available at `http://localhost:8080`

### Option 2: Local Development

```bash
# Install dependencies
go mod download

# Ensure PostgreSQL is running on localhost:5432

# Run database migrations
go run ./cmd/main.go migrate

# Start the application
go run ./cmd/main.go
```

The application will be available at `http://localhost:8080`

## Configuration

Configuration is loaded from environment variables via `.env` file. The `config.LoadConfig()` function in [config/config.go](config/config.go) manages this initialization.

### Environment Variables

| Variable             | Description                             | Required | Default |
| -------------------- | --------------------------------------- | -------- | ------- |
| `APP_ENV`            | Environment (development/production)    | Yes      | -       |
| `APP_PORT`           | Server port                             | Yes      | -       |
| `DB_HOST`            | PostgreSQL host                         | Yes      | -       |
| `DB_PORT`            | PostgreSQL port                         | Yes      | -       |
| `DB_USER`            | PostgreSQL user                         | Yes      | -       |
| `DB_PASSWORD`        | PostgreSQL password                     | Yes      | -       |
| `DB_NAME`            | PostgreSQL database name                | Yes      | -       |
| `DB_SSLMODE`         | PostgreSQL SSL mode                     | No       | disable |
| `JWT_ACCESS_SECRET`  | JWT access token secret (min 32 chars)  | Yes      | -       |
| `JWT_REFRESH_SECRET` | JWT refresh token secret (min 32 chars) | Yes      | -       |
| `EMAIL_FROM`         | Sender email address                    | Yes      | -       |
| `RESEND_API_KEY`     | Resend email service API key            | Yes      | -       |

## Running the Application

### Start with Docker Compose

```bash
docker-compose up
```

### Start Locally

```bash
go run ./cmd/main.go
```

### Verify Health

```bash
curl http://localhost:8080/v1/health
```

Expected response:

```json
{
  "status": "ok"
}
```

## API Documentation

### Base URL

```
http://localhost:8080/v1
```

### Authentication

Most endpoints require a valid JWT token in the `Authorization` header:

```
Authorization: Bearer <access_token>
```

### Core Endpoints

#### Health Check

```
GET /v1/health
```

Public endpoint for health monitoring.

#### Authentication

```
POST /v1/auth/register          # Register new user
POST /v1/auth/login             # User login
POST /v1/auth/refresh           # Refresh access token
```

#### Verification

```
POST /v1/verify/email           # Send email verification
POST /v1/verify/email/confirm   # Confirm email verification
```

#### Password Management

```
POST /v1/password/change        # Change password (authenticated)
POST /v1/password/reset         # Request password reset
POST /v1/password/reset/confirm # Confirm password reset
```

#### Employee Management (Authenticated - Admin only)

```
GET    /v1/employees            # List all employees
GET    /v1/employees/{id}       # Get employee details
POST   /v1/employees            # Create new employee
PUT    /v1/employees/{id}       # Update employee
DELETE /v1/employees/{id}       # Delete employee
```

#### User Profile (Authenticated)

```
GET    /v1/me                   # Get current user profile
PUT    /v1/me                   # Update current user profile
```

### Response Format

All responses follow a standard JSON format:

```json
{
  "success": true,
  "data": {},
  "error": null,
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## Development

### Code Style

This project follows Go conventions:

- Use `gofmt` for code formatting
- Use `go vet` for code analysis
- Use meaningful variable and function names
- Write tests for new functionality

### Running Code Quality Tools

```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Adding New Features

1. Create a new handler in `handlers/`
2. Implement service logic in `services/`
3. Create or update repository in `repositories/`
4. Add routes to `routes/`
5. Write comprehensive tests
6. Update this README with new endpoints

## Testing

### Test Coverage

- **Unit Tests**: Repository and utility functions
- **Integration Tests**: Handler and service interactions
- **Test Utilities**: Helper functions in `internal/testutils/`

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./handlers

# Run with coverage report
go test -cover ./...

# Generate coverage HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Files

Tests follow Go conventions:

- Unit tests: `*_test.go` files in the same package
- Integration tests: `*_integration_test.go` files
- Test utilities: `internal/testutils/`

### Writing Tests

Use the test utilities provided:

```go
// Create test database
db := testutils.SetupTestDB(t)
defer testutils.TeardownTestDB(t, db)

// Use mock HTTP requests
recorder := testutils.NewTestRequest(t, http.MethodGet, "/endpoint", nil)

// Assert responses
assert.Equal(t, http.StatusOK, recorder.Code)
```

## Database Migrations

Migrations are versioned SQL files in `database/migrations/`. Each migration has an `.up.sql` (apply) and `.down.sql` (rollback) file.

### Current Migrations

- **001_init**: Core schema (users, employees, profiles)
- **002_verification**: Email verification system
- **003_sessions**: Session management
- **004_audit**: Audit logging tables
- **005_otp_codes**: OTP storage
- **006_indexes**: Performance indexes

### Running Migrations

Migrations run automatically on application startup. To manually manage:

```bash
# Run migrations programmatically from Go code
import "github.com/golang-migrate/migrate/v4"
// See database/postgres.go for implementation
```

### Creating New Migrations

```bash
# Create new migration files
touch database/migrations/007_your_feature.up.sql
touch database/migrations/007_your_feature.down.sql
```

Follow the naming convention: `NNN_feature_name.{up|down}.sql`

## Security Considerations

### Authentication & Authorization

- All sensitive endpoints require JWT authentication
- JWT secrets must be minimum 32 characters in production
- Tokens should be short-lived (implement refresh token rotation)
- Implement role-based access control for admin endpoints

### Password Security

- Passwords are hashed with bcrypt (cost factor: 12)
- Never log or expose passwords
- Implement password complexity requirements
- Consider rate limiting on login attempts

### Data Protection

- Use HTTPS in production (enable SSL)
- Set `DB_SSLMODE=require` for production database connections
- Implement audit logging for sensitive operations
- Use secure session cookies with HttpOnly flag

### API Security

- Implement CORS properly for your domain
- Use rate limiting to prevent abuse
- Validate all input data
- Sanitize database queries against SQL injection
- Implement request ID tracking for security audits

### Environment Variables

- Never commit `.env` files to version control
- Rotate secrets regularly in production
- Use a secrets management system (AWS Secrets Manager, HashiCorp Vault)
- Store credentials securely in CI/CD pipelines

## Deployment

### Docker Deployment

#### Build Image

```bash
docker build -t company-mgmt-api:latest .
```

#### Run Container

```bash
docker run -d \
  --name company_mgmt_api \
  -p 8080:8080 \
  --env-file .env \
  company-mgmt-api:latest
```

### Docker Compose Deployment

```bash
docker-compose up -d
```

### Production Checklist

- [ ] Set `APP_ENV=production`
- [ ] Use strong, unique secrets for JWT keys
- [ ] Enable database SSL: `DB_SSLMODE=require`
- [ ] Configure CORS for your domain
- [ ] Set up monitoring and logging
- [ ] Enable rate limiting
- [ ] Use a reverse proxy (Nginx) with SSL/TLS
- [ ] Implement health checks and monitoring
- [ ] Set up database backups
- [ ] Configure auto-scaling if using Kubernetes
- [ ] Implement proper error handling and logging
- [ ] Set up CI/CD pipeline

### Database Backup

```bash
# Backup PostgreSQL database
pg_dump -U $DB_USER -h $DB_HOST $DB_NAME > backup.sql

# Restore from backup
psql -U $DB_USER -h $DB_HOST $DB_NAME < backup.sql
```

## Contributing

### Development Workflow

1. **Create Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Implement Feature**
   - Write clean, documented code
   - Follow Go conventions
   - Add comprehensive tests
   - Update relevant documentation

3. **Test Thoroughly**

   ```bash
   go test ./...
   go fmt ./...
   go vet ./...
   ```

4. **Submit Pull Request**
   - Include clear description of changes
   - Reference related issues
   - Ensure all tests pass
   - Request code review from team members

### Code Review Guidelines

- All code changes require review
- Tests must accompany code changes
- Documentation must be updated
- No hardcoded secrets or sensitive data
- Follow architectural patterns established in the project

### Reporting Issues

When reporting bugs, include:

- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Request ID from logs if applicable

## License

[Add your license information here]

## Support

For issues, questions, or contributions, please [add support contact information here].

---

**Last Updated**: February 2026
**Go Version**: 1.24.4
**Status**: Production Ready
