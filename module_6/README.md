# Notes API

A REST API for managing personal notes with categorization. Built with Go, Fiber v3, and PostgreSQL.

## Features

- User authentication with JWT
- Personal notes management
- Category organization
- Full CRUD operations
- Input validation
- Structured logging
- Comprehensive testing
- Docker support

## Tech Stack

- **Language**: Go 1.23+
- **Framework**: Fiber v3
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Containerization**: Docker + Docker Compose
- **Testing**: Go testing + testify
- **Validation**: go-playground/validator

## API Endpoints

### Authentication
- `POST /api/auth/signup` - User registration
- `POST /api/auth/signin` - User login

### Categories (requires authentication)
- `GET /api/categories` - Get user categories
- `POST /api/categories` - Create category
- `GET /api/categories/:id` - Get category by ID
- `PUT /api/categories/:id` - Update category
- `DELETE /api/categories/:id` - Delete category

### Notes (requires authentication)
- `GET /api/notes` - Get user notes (with filters)
- `POST /api/notes` - Create note
- `GET /api/notes/:id` - Get note by ID
- `PUT /api/notes/:id` - Update note
- `DELETE /api/notes/:id` - Delete note

### Health Check
- `GET /health` - Application health status

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository
2. Copy environment file:
   ```bash
   cp .env.example .env
   ```
3. Start services:
   ```bash
   docker-compose up -d
   ```
4. Run migrations:
   ```bash
   bash scripts/migrate.sh up
   ```

### Manual Setup

1. **Prerequisites**:
   - Go 1.23+
   - PostgreSQL
   - golang-migrate CLI

2. **Environment Setup**:
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Database Setup**:
   ```bash
   # Start PostgreSQL
   docker-compose up -d postgres
   
   # Run migrations
   bash scripts/migrate.sh up
   ```

4. **Start Application**:
   ```bash
   go run cmd/server/main.go
   ```

The API will be available at `http://localhost:8888`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8888 |
| `ENV` | Environment (development/production) | development |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | postgres |
| `DB_NAME` | Database name | notes_db |
| `JWT_SECRET` | JWT signing secret | *required* |
| `JWT_EXPIRE_HOURS` | JWT expiration hours | 24 |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | info |
| `LOG_FORMAT` | Log format (json/text) | json |
| `LOG_OUTPUT` | Log output (stdout/stderr/file path) | stdout |

## Database Schema

```sql
-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, name)
);

-- Notes table
CREATE TABLE notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API Usage Examples

### Registration
```bash
curl -X POST http://localhost:8888/api/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "name": "username",
    "password": "password123"
  }'
```

### Create Note
```bash
curl -X POST http://localhost:8888/api/notes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "My Note",
    "content": "Note content here",
    "category_id": 1
  }'
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Migration Commands

```bash
# Run all migrations
bash scripts/migrate.sh up

# Rollback one migration
bash scripts/migrate.sh down-one

# Check migration version
bash scripts/migrate.sh version

# Create new migration
bash scripts/migrate.sh create migration_name
```

## Project Structure

```
├── cmd/server/          # Application entry point
│   ├── app/            # App initialization
│   ├── handlers/       # HTTP handlers
│   ├── middleware/     # Custom middleware
│   └── routes/         # Route definitions
├── internal/           # Private application code
│   ├── config/         # Configuration
│   ├── database/       # Database layer
│   ├── logger/         # Logging interface
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   └── validator/      # Input validation
├── scripts/            # Build and deployment scripts
└── docs/              # API documentation
```

## Development

### Adding New Features

1. Create model in `internal/models/`
2. Add repository interface and implementation
3. Implement service layer in `internal/services/`
4. Create handlers in `cmd/server/handlers/`
5. Add routes in `cmd/server/routes/`
6. Write tests for all layers

### Code Quality

- Follow Go conventions and best practices
- Write unit tests for all new features
- Use structured logging
- Validate all inputs
- Handle errors properly

## License

MIT License - see LICENSE file for details

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request