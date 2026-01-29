#!/bin/bash
#
# CrossForge Database Migration Script
# Applies SQL migrations to the SQLite database
#

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Default database path
DB_PATH="${DATABASE_PATH:-./data/crossforge.db}"
MIGRATIONS_DIR="./internal/storage/migrations"

echo "========================================"
echo "CrossForge Database Migration"
echo "========================================"
echo ""

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo -e "${RED}Error: sqlite3 is not installed${NC}"
    echo "Install it with:"
    echo "  - macOS: brew install sqlite3"
    echo "  - Ubuntu/Debian: sudo apt-get install sqlite3"
    echo "  - Fedora: sudo dnf install sqlite"
    exit 1
fi

# Create data directory if it doesn't exist
mkdir -p "$(dirname "$DB_PATH")"

# Check for migration files
if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo -e "${RED}Error: Migrations directory not found: $MIGRATIONS_DIR${NC}"
    exit 1
fi

MIGRATION_FILES=$(ls -1 "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort)

if [ -z "$MIGRATION_FILES" ]; then
    echo -e "${YELLOW}No migration files found in $MIGRATIONS_DIR${NC}"
    exit 0
fi

echo "Database: $DB_PATH"
echo "Migrations: $MIGRATIONS_DIR"
echo ""

# Get current schema version
CURRENT_VERSION=0
if [ -f "$DB_PATH" ]; then
    CURRENT_VERSION=$(sqlite3 "$DB_PATH" "SELECT COALESCE(MAX(version), 0) FROM schema_migrations;" 2>/dev/null || echo "0")
fi

echo "Current schema version: $CURRENT_VERSION"
echo ""

# Apply migrations
APPLIED=0
for migration in $MIGRATION_FILES; do
    # Extract version number from filename (e.g., 0001 from 0001_initial.sql)
    VERSION=$(basename "$migration" | sed 's/^0*//' | cut -d'_' -f1)

    if [ "$VERSION" -gt "$CURRENT_VERSION" ]; then
        echo -e "${YELLOW}Applying migration: $(basename "$migration")${NC}"

        if sqlite3 "$DB_PATH" < "$migration"; then
            echo -e "${GREEN}  ✓ Applied successfully${NC}"
            APPLIED=$((APPLIED + 1))
        else
            echo -e "${RED}  ✗ Failed to apply migration${NC}"
            exit 1
        fi
    else
        echo "  Skipping $(basename "$migration") (already applied)"
    fi
done

echo ""
if [ "$APPLIED" -gt 0 ]; then
    echo -e "${GREEN}Applied $APPLIED migration(s)${NC}"
else
    echo "No new migrations to apply"
fi

# Show current schema version
NEW_VERSION=$(sqlite3 "$DB_PATH" "SELECT MAX(version) FROM schema_migrations;")
echo "Schema version: $NEW_VERSION"
echo ""
echo "Done!"
