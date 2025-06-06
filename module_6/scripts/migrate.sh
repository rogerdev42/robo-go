#!/bin/bash

# Load variables from .env
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

# Build database URL
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Migrations path
MIGRATIONS_PATH="./internal/database/migrations"

# Check if migrate command exists
if ! command -v migrate &> /dev/null; then
    echo "migrate command not found. Please install golang-migrate"
    exit 1
fi

# Execute command
case "$1" in
    up)
        echo "Running migrations up..."
        migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" up
        ;;
    down)
        echo "Running migrations down..."
        migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" down
        ;;
    down-one)
        echo "Rolling back one migration..."
        migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" down 1
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Usage: $0 create <migration_name>"
            exit 1
        fi
        echo "Creating migration: $2"
        migrate create -ext sql -dir "$MIGRATIONS_PATH" -seq "$2"
        ;;
    version)
        echo "Current migration version:"
        migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" version
        ;;
    *)
        echo "Usage: $0 {up|down|down-one|create|version}"
        exit 1
        ;;
esac