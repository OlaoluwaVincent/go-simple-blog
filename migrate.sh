#!/bin/bash

# migrate.sh - Helper script for golang-migrate
#
# Requires: DATABASE_URL environment variable (loaded from .envrc or shell)

# Load environment variables from .envrc if it exists
if [ -f .envrc ]; then
  set -a  # Export all variables defined in .envrc
  source .envrc || {
    echo "Error: Failed to load .envrc" >&2
    exit 1
  }
  set +a
fi

# Strictly require DB_URL — fail fast with clear message if missing
if [ -z "$DB_URL" ]; then
  echo -e "\033[0;31mError: DATABASE_URL is not set.\033[0m" >&2
  echo "Please define DATABASE_URL in your .envrc file or export it in your shell." >&2
  echo "" >&2
  echo "Example for .envrc:" >&2
  echo 'export DATABASE_URL="postgres://user:password@localhost:5432/dbname?sslmode=disable"' >&2
  exit 1
fi

MIGRATIONS_DIR="./cmd/migrate/migrations"

# Colors for pretty output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

case "$1" in
  create)
    if [ -z "$2" ]; then
      read -p "Enter migration name (e.g., create_users_table): " name
    else
      name="$2"
    fi

    if [ -z "$name" ]; then
      echo -e "${RED}Error: Migration name cannot be empty.${NC}"
      exit 1
    fi

    echo -e "${GREEN}Creating migration: $name${NC}"
    migrate create -seq -ext sql -dir "$MIGRATIONS_DIR" "$name"
    echo -e "${GREEN}Migration files created in $MIGRATIONS_DIR${NC}"
    ;;

  up)
    echo -e "${YELLOW}Running migrations UP...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
    ;;

  down)
    if [ -z "$2" ]; then
      read -p "How many migrations to rollback? (default: 1, or 'all'): " count
      count=${count:-1}
    else
      count="$2"
    fi

    echo -e "${YELLOW}Rolling back $count migration(s)...${NC}"

    if [ "$count" = "all" ]; then
      current_version=$(migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version 2>/dev/null | grep -o '^[0-9]*' || echo 0)
      if [ -n "$current_version" ] && [ "$current_version" -ge 0 ]; then
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$current_version" 2>/dev/null || true
      fi
      migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down -all
    else
      migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down "$count"
    fi
    ;;

  version)
    echo -e "${GREEN}Current migration version:${NC}"
    if migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version > /dev/null 2>&1; then
      migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version
    else
      echo "No schema_migrations table (database is clean/empty)"
    fi
    ;;

  force)
    if [ -z "$2" ]; then
      read -p "Enter version to force (e.g., 1): " version
    else
      version="$2"
    fi

    if [ -z "$version" ]; then
      echo -e "${RED}Error: Version cannot be empty.${NC}"
      exit 1
    fi

    echo -e "${YELLOW}Forcing migration version to $version...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$version"
    echo -e "${GREEN}Done! Database version forced to $version${NC}"
    ;;

  fix|reset)
    echo -e "${YELLOW}Fixing dirty database — forcing version 0 (clean state)...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force 0
    echo -e "${GREEN}Database is now clean at version 0${NC}"
    ;;

  drop)
    echo -e "${RED}WARNING: This will DROP the schema_migrations table and remove all migration history!${NC}"
    read -p "Are you sure you want to continue? (type 'yes' to confirm): " confirm
    if [ "$confirm" = "yes" ]; then
      echo -e "${YELLOW}Dropping schema_migrations table...${NC}"
      migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" drop
      echo -e "${GREEN}schema_migrations table dropped.${NC}"
    else
      echo "Operation cancelled."
    fi
    ;;

  *)
    echo "Usage: $0 {create|up|down|version|force|fix|drop}"
    echo ""
    echo "Commands:"
    echo "  create [name]     Create new migration"
    echo "  up                Apply all pending migrations"
    echo "  down [n|'all']    Rollback migrations"
    echo "  version           Show current version"
    echo "  force [version]   Force version"
    echo "  fix               Force version 0"
    echo "  drop              Drop migration history"
    exit 1
    ;;
esac