# Tiger FastTrack Card API

A modern Go REST API built with Gin framework and PostgreSQL database.

## Project Structure

```
tiger-fasttrack-card/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies checksums
├── Makefile                # Common development tasks
├── docker-compose.yml      # PostgreSQL development environment
├── .env.example           # Environment variables template
├── .gitignore             # Git ignore file
├── README.md              # This file
└── internal/              # Private application code
    ├── config/            # Configuration management
    │   └── config.go
    ├── database/          # Database connection and setup
    │   └── database.go
    ├── handlers/          # HTTP request handlers
    │   └── handlers.go
    ├── middleware/        # HTTP middleware
    │   └── middleware.go
    ├── repository/        # Data access layer
    │   └── repository.go
    ├── routes/            # Route definitions
    │   └── routes.go
    └── service/           # Business logic layer
        └── service.go
```

## Tech Stack

- **Go 1.21+** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library
- **PostgreSQL** - Database
- **Docker** - Containerization (for development database)

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for database)
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd tiger-fasttrack-card
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Start PostgreSQL database:
   ```bash
   make db-up
   # or
   docker-compose up -d
   ```

4. Install dependencies:
   ```bash
   make deps
   # or
   go mod tidy
   ```

5. Run the application:
   ```bash
   make dev
   # or
   go run main.go
   ```

The API will be available at `http://localhost:8080`

### Environment Variables

Copy `.env.example` to `.env` and configure the following variables:

#### Server Configuration
- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Application environment (development/production)

#### Database Configuration
You can use either individual database settings or a DATABASE_URL:

**Individual settings:**
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database username (default: postgres)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: tiger_fasttrack)
- `DB_SSLMODE`: SSL mode (default: disable)

**Or use DATABASE_URL:**
- `DATABASE_URL`: Full PostgreSQL connection string
  - Example: `postgres://username:password@localhost:5432/tiger_fasttrack?sslmode=disable`

#### Authentication
- `JWT_SECRET`: Secret key for JWT tokens

## API Endpoints

### Health Check
- `GET /health` - Check API health status

### Cards
- `GET /api/v1/cards` - Get all cards
- `GET /api/v1/cards/:id` - Get card by ID
- `POST /api/v1/cards` - Create a new card
- `PUT /api/v1/cards/:id` - Update card by ID
- `DELETE /api/v1/cards/:id` - Delete card by ID

## Development

### Quick Start with Make

```bash
# Start PostgreSQL database
make db-up

# Run in development mode
make dev

# Build the application
make build

# Run tests
make test

# Stop database
make db-down

# See all available commands
make help
```

### Project Architecture

This project follows a clean architecture pattern:

- **`main.go`**: Entry point, dependency injection, server setup
- **`internal/config`**: Configuration management using environment variables
- **`internal/database`**: Database connection and configuration (PostgreSQL with GORM)
- **`internal/repository`**: Data access layer - direct database operations
- **`internal/service`**: Business logic layer - contains application logic
- **`internal/handlers`**: HTTP request handlers - handles HTTP requests/responses
- **`internal/middleware`**: HTTP middleware for CORS, authentication, logging, etc.
- **`internal/routes`**: Route definitions and groupings

### Adding New Features

1. **Define your model** (when you create model files)
2. **Add repository methods** in `internal/repository/`
3. **Add service methods** in `internal/service/`
4. **Create handlers** in `internal/handlers/`
5. **Add routes** in `internal/routes/`
6. **Update configuration** in `internal/config/` if needed

### Database Migrations

When you create models, update the `Migrate()` function in `internal/database/database.go`:

```go
func (d *Database) Migrate() error {
    return d.DB.AutoMigrate(&models.User{}, &models.Card{})
}
```

## Building for Production

```bash
go build -o tiger-fasttrack-card main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
