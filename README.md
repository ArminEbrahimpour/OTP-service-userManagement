# OTP Service - Golang Backend

A robust OTP-base authentication service built with Go , featuring user management and ratelimiting and api documentation .

## Features

- OTP Authentication : Phone number-based login and registration
- Rate Limiting : Max 3 OTP requests per phone number within 10 minutes
- JWT Tokens : Secure authentication with JWT
- User Management : RESTFUL endpoints with pagination and search
- PostgreSQL Database : Reliable Data persistence
- Docker Support : Containerized deployment
- Clean Architecture : Well-structured and maintainable structure

## Requirements

- Go 1.23 or later
- PostgreSQL (or Docker)
- Docker & Docker compose

## Architecture

The application follows this clean architecture:

```
otp-service/
├── cmd/                    # Application entry points
├── config/                 # Configuration management
├── docs/                   # Swagger documentation
├── internal/
│   ├── api/               # HTTP handlers and routes
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   └── service/           # Business logic
├── pkg/
│   ├── database/          # Database connection
│   ├── jwt/              # JWT utilities
│   └── middleware/        # HTTP middleware
├── docker-compose.yml     # Docker services
├── Dockerfile            # Container definition
└── Makefile              # Build commands
```

## Installation and Setup

Docker(recommanded)

1. clone the repository :

```
git clone github.com/ArminEbrahimpour/OTP-service-userManagement

cd OTP-service-userManagement
```

2. Run with docker compose

```
docker compose up --build
```

or run in detach mode:

```
docker compose up -d
```

## 🔧 Configuration

Configure via environment variables in docker-compose.yml file

## Database Choise : PostgreSQL

- ACID Compliance : Ensure data integrity for user accounts
- Performance : Excellent performance for read-heavy user management operations
- JSON Support: Native JSON support for future extensibility
- Scalability: Proven scalability for production applications
- Docker Integration: Easy containerization and deployment
- Community: Large ecosystem and extensive documentation

The database schema is automatically managed via GORM migrations, ensuring consistency across environments.

## API Documentation

Authentication Endpoints

### Send OTP

```
POST /api/v1/auth/send-otp
Content-Type: application/json

{
    "phone_number": "+1234567890"
}
```

Response (200):

```
json{
    "message": "OTP sent successfully"
}
```

Rate Limiting (429):

```
json{
    "error": "rate limit exceeded. Try again after 2025-01-15T10:30:00Z"
}
```

Verify OTP

```
POST /api/v1/auth/verify-otp
Content-Type: application/json

{
    "phone_number": "+1234567890",
    "code": "123456"
}
```

Response (200):

```
json{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "phone_number": "+1234567890",
        "registration_date": "2025-01-15T10:00:00Z",
        "created_at": "2025-01-15T10:00:00Z",
        "updated_at": "2025-01-15T10:00:00Z"
    }
}
```

### User Management Endpoints

Note: All user endpoints require JWT authentication. Include the token in the Authorization header:
` Authorization: Bearer <your-jwt-token>`
Get User by ID

```
GET /api/v1/users/{id}
Authorization: Bearer <jwt-token>
```

Response (200):

```
json{
    "id": 1,
    "phone_number": "+1234567890",
    "registration_date": "2025-01-15T10:00:00Z",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z"
}
```

Get Users List

```
GET /api/v1/users?page=1&page_size=10&search=+123
Authorization: Bearer <jwt-token>
```

Response (200):

```
json{
    "users": [
        {
            "id": 1,
            "phone_number": "+1234567890",
            "registration_date": "2024-01-15T10:00:00Z",
            "created_at": "2024-01-15T10:00:00Z",
            "updated_at": "2024-01-15T10:00:00Z"
        }
    ],
    "total": 1,
    "page": 1,
    "page_size": 10,
    "total_pages": 1
}
```

Health Check

```
GET /health
```

Response (200):

```
json{
    "status": "healthy"
}
```

## Security features

- Rate Limiting: Prevents OTP abuse with configurable limits
- JWT Authentication: Secure token-based authentication
- Input Validation: Request validation using Gin binding
- CORS Support: Cross-origin resource sharing configuration
- Secure Headers: Security middleware for HTTP headers

## Testing the Application

1- send OTP :

```
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}'
```

2-Check console output for OTP (it will be printed like: OTP for +1234567890: 123456)

3- Verify OTP:

```
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890", "code": "123456"}'
```

4- Use the returned JWT token:

```
# Get user details
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"

# Get users list
curl -X GET "http://localhost:8080/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

### Using Swagger UI

Visit http://localhost:8080/swagger/index.html for interactive API documentation where you can:

- Test all endpoints directly from the browser
- View request/response schemas
- Authenticate using JWT tokens
- See example requests and responses

## Example Usage Flow

1- Start the applicatoin :

```
docker-compose up --build
```

2- Register / Login a user:

```
# Send OTP
curl -X POST http://localhost:8080/api/v1/auth/send-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890"}'

# Check console for OTP, then verify
curl -X POST http://localhost:8080/api/v1/auth/verify-otp \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "+1234567890", "code": "PRINTED_OTP"}'
```

3- Use the JWT token for authenticated requests:

```
# Export the token for convenience
export JWT_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Get user profile
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $JWT_TOKEN"
```
