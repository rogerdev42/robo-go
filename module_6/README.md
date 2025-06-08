# Notes API

REST API for managing personal notes with categorization. Built with Go, Fiber v3, and PostgreSQL.

## Tech Stack

- **Language**: Go 1.23+
- **Framework**: Fiber v3
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Containerization**: Docker + Docker Compose

## Quick Start

### Using Docker Compose

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

1. Start PostgreSQL:
   ```bash
   docker-compose up -d postgres
   ```

2. Run migrations:
   ```bash
   bash scripts/migrate.sh up
   ```

3. Start application:
   ```bash
   go run cmd/server/main.go
   ```

The API will be available at `http://localhost:8888`

## API Endpoints

### Authentication
- `POST /api/auth/signup` - User registration
- `POST /api/auth/signin` - User login

### Categories (requires authentication)
- `GET /api/categories` - Get user categories
- `POST /api/categories` - Create category
- `PUT /api/categories/:id` - Update category
- `DELETE /api/categories/:id` - Delete category

### Notes (requires authentication)
- `GET /api/notes` - Get user notes
- `POST /api/notes` - Create note
- `GET /api/notes/:id` - Get note by ID
- `PUT /api/notes/:id` - Update note
- `DELETE /api/notes/:id` - Delete note

### Documentation
- `GET /docs/` - Swagger UI (development only)

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8888 | Server port |
| `ENV` | development | Environment |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | notes_db | Database name |
| `JWT_SECRET` | *required* | JWT signing secret |
| `JWT_EXPIRE_HOURS` | 24 | JWT expiration hours |
| `LOG_LEVEL` | info | Log level |
| `LOG_FORMAT` | json | Log format |

## Usage Examples

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

```bash
go test ./...
```

## Migration Commands

```bash
# Run migrations
bash scripts/migrate.sh up

# Rollback one migration
bash scripts/migrate.sh down-one

# Create new migration
bash scripts/migrate.sh create migration_name
```