# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Lumera App is a hackathon project (Platanus Hack 25) optimized for **rapid development** (33.5 hours)

**Stack:** Go 1.23 + Svelte 5 + Astro 5 + PostgreSQL 16 + Docker Compose

## Development Commands

**IMPORTANT:** A comprehensive Makefile is available at the project root with shortcuts for all common tasks. Run `make help` to see all available commands.

### Quick Start

```bash
# See all available commands
make help

# Start everything
make up

# View all logs
make logs

# Stop everything
make down
```

### Starting the Application

```bash
# Full stack (preferred)
make up
# or: docker compose up

# Rebuild and start
docker compose up --build

# Run in background
make up
# or: docker compose up -d
```

**Services started:**
- Landing (Astro): `http://localhost:4321`
- Frontend (Svelte): `http://localhost:5173`
- Backend API (Go): `http://localhost:8080`
- PostgreSQL: `localhost:5432`
- Adminer (DB GUI): `http://localhost:8081`

### Viewing Logs

```bash
# All services
make logs

# Specific service
make backend-logs
make frontend-logs
make landing-logs
make postgres-logs

# Or directly:
docker compose logs -f backend
```

### Database Quick Access

```bash
# Execute SQL query directly
make db-query SQL="SELECT * FROM users LIMIT 5"

# Connect to interactive psql shell
make db-shell

# Load test data
make db-seed

# Reset database (⚠️ deletes all data)
make db-reset
```

**Common queries:**
```bash
# Count records
make db-query SQL="SELECT COUNT(*) FROM users"

# View recent records
make db-query SQL="SELECT * FROM users ORDER BY created_at DESC LIMIT 10"

# Search
make db-query SQL="SELECT * FROM users WHERE email LIKE '%@example.com'"

# List all tables
make db-query SQL="\dt"
```

See `scripts/db-queries.md` for more query examples.

### Automated Deployment

**IMPORTANT:** A deployment script automates the entire deployment process including migrations.

```bash
# Development deployment (dev servers with hot reload)
./scripts/deploy.sh

# Production deployment (optimized builds with nginx)
./scripts/deploy.sh --production

# Development mode (skip git pull)
./scripts/deploy.sh --skip-git

# Fast deployment (skip Docker rebuild)
./scripts/deploy.sh --skip-build

# View all options
./scripts/deploy.sh --help
```

**Key differences:**
- **Development mode** (`docker-compose.yml`): Uses dev servers (Air, Vite, Astro) with hot reload and volumes
- **Production mode** (`docker-compose.prod.yml`): Uses optimized builds with nginx for landing/frontend, compiled Go binary for backend

**What the script does automatically:**
1. ✅ Updates code from git (optional with `--skip-git`)
2. ✅ Stops old containers
3. ✅ Builds and starts new containers
4. ✅ **Runs database migrations automatically** (includes seed data)
5. ✅ Verifies services are healthy
6. ✅ Shows service status and URLs

**Script features:**
- Color-coded output (info, success, warnings, errors)
- Pre-deployment checks (Docker running, .env exists, etc.)
- Detects uncommitted git changes (asks for confirmation)
- Loads DB credentials from `.env` file
- Shows migration version before and after
- Checks backend health endpoint
- Interactive prompt to view logs after deployment

**Production deployment on VPS/Cloud:**
```bash
# SSH to server
ssh user@your-server.com

# Clone or pull latest
git clone [repo] && cd lumera_app
# or: git pull

# Run production deployment
./scripts/deploy.sh --production

# Script will:
# - Pull latest code
# - Build optimized Docker images (nginx + compiled binaries)
# - Rebuild containers
# - Run migrations (migrations 17-20 include seed data)
# - Verify everything is working
```

**⚠️ IMPORTANT:** Always use `--production` flag when deploying to servers. Using dev mode in production causes:
- Constant page reloading and "?" appearing in URLs
- Astro dev toolbar visible at the bottom
- Slower performance and larger memory usage
- Hot reload attempting to connect to dev servers

**Migrations included in seed data:**
- Migration 17: Cursos y Materias (1 curso, 2 materias)
- Migration 18: OAs de Lengua (23 objetivos)
- Migration 19: Bloom Objectives (134 objectives)
- Migration 20: Questions (672 preguntas)

All migrations use `ON CONFLICT DO UPDATE` for idempotency, so they can be safely re-run.

**Manual migration commands** (if needed):
```bash
# Apply migrations
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' up

# Check version
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' version

# Rollback last migration
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' down 1
```

**Troubleshooting deployment:**
- Check logs: `docker compose logs -f backend`
- Verify containers: `docker compose ps`
- Check health: `curl http://localhost:8080/api/health`
- Reset everything: `docker compose down -v && ./scripts/deploy.sh`

### Backend Development

```bash
cd backend

# After modifying go.mod dependencies
docker compose exec backend go mod tidy
docker compose restart backend

# Run locally without Docker (requires local PostgreSQL)
# Set DB_HOST=localhost in .env
air  # or: go run cmd/main.go
```

**Hot reload:** Air automatically recompiles on `.go` file changes (<2 seconds)

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Run locally without Docker
npm run dev

# Production build
npm run build
```

**Hot reload:** Vite HMR provides instant updates (<1 second)

### Database Access

```bash
# Connect to PostgreSQL
docker compose exec postgres psql -U admin -d hackathon

# From local machine (if psql installed)
psql -h localhost -U admin -d hackathon
```

**Default credentials:** `admin` / `(see .env file)` / `hackathon` (user/password/database)

## Architecture

### Backend Structure (Go + Chi + GORM)

```
backend/
├── cmd/main.go                    # Entry point, router setup, middleware
├── internal/
│   ├── db/db.go                   # PostgreSQL connection with GORM
│   ├── handlers/                  # HTTP handlers (currently: health.go)
│   └── models/                    # GORM models (create here for new entities)
└── .air.toml                      # Hot reload configuration
```

**Key patterns:**
- Chi router with standard middleware (RequestID, Logger, Recoverer, Timeout)
- CORS configured for localhost:5173 (Vite dev server)
- Database connection retry logic (10 attempts, 2s interval)
- All routes prefixed with `/api`

**Adding new endpoints:**
1. Create handler in `internal/handlers/`
2. Register route in `cmd/main.go` router
3. Hot reload will pick up changes automatically

### Frontend Structure (Svelte 5 + Vite + Tailwind)

```
frontend/
├── src/
│   ├── App.svelte                 # Main component (health check UI)
│   ├── main.js                    # Entry point
│   ├── app.css                    # Tailwind imports
│   ├── components/                # Create reusable components here
│   └── lib/                       # Utilities and helpers
├── vite.config.js                 # Proxy: /api → backend:8080
└── tailwind.config.js
```

**Key patterns:**
- Svelte 5 uses runes: `$state`, `$derived`, `$effect`
- API calls use `/api/*` (proxied to backend via Vite)
- GSAP for animations (imported in components with `onMount`)
- Tailwind for styling

**API integration example:**
```javascript
// Frontend calls /api/health
const response = await fetch('/api/health');
// Vite proxy → http://backend:8080/api/health
```

### Landing Structure (Astro 5 + Tailwind)

```
landing/
├── src/
│   ├── pages/                     # Astro pages (file-based routing)
│   ├── components/                # Reusable components
│   └── layouts/                   # Page layouts
├── public/                        # Static assets
├── astro.config.mjs              # Astro configuration
└── Dockerfile                    # Multi-stage build
```

**Key patterns:**
- Astro 5 for static site generation
- Tailwind 4.x for styling
- File-based routing in `src/pages/`
- Zero JavaScript by default (only when needed)
- Optimized for marketing/landing pages

**Development:**
```bash
cd landing
npm install
npm run dev          # Runs on http://localhost:4321
```

**Hot reload:** Astro dev server provides instant updates

### Docker Architecture

**Development mode:**
- Backend: Volume-mounted `/app` with Air hot reload
- Frontend: Volume-mounted `/app` with Vite HMR
- Landing: Volume-mounted `/app` with Astro dev server
- PostgreSQL: Persistent volume `postgres_data`

**Important:** Backend depends on PostgreSQL health check. If backend fails to start, check postgres is healthy:
```bash
docker compose ps
```

## Database (PostgreSQL + GORM)

**Connection configuration:**
- In Docker: `DB_HOST=postgres` (service name)
- Local development: `DB_HOST=localhost`

**Creating models:**
1. Define struct in `backend/internal/models/`
2. Import in `main.go` or migration file
3. Use `db.DB.AutoMigrate(&Model{})` for schema
4. GORM handles SQL generation

**Database connection:**
- Accessed globally via `db.DB` (exported from `internal/db`)
- Connection pool: 10 idle, 100 max open connections

## Database Migrations

**CRITICAL:** This project uses **golang-migrate** for database schema versioning. **NEVER use GORM AutoMigrate in production.**

### When to Create a Migration

Create a migration whenever you need to:
- Create a new table
- Add/remove/modify columns
- Add/remove indexes
- Change constraints
- Insert seed data

**Do NOT modify existing migration files that have been applied.** Always create a new migration for changes.

### Creating a New Migration

```bash
# From project root
make -C backend migrate-create name=add_products_table

# This creates two files in backend/migrations/:
# - XXXXXX_add_products_table.up.sql   (apply changes)
# - XXXXXX_add_products_table.down.sql (rollback changes)
```

**Migration naming conventions:**
- `create_users_table` - Creating new table
- `add_status_to_orders` - Adding column
- `remove_old_field_from_products` - Removing column
- `add_index_to_users_email` - Adding index

### Writing Migration Files

**Example: Creating a table**

`000002_add_products_table.up.sql`:
```sql
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_name ON products(name);
```

`000002_add_products_table.down.sql`:
```sql
DROP INDEX IF EXISTS idx_products_name;
DROP TABLE IF EXISTS products;
```

**Example: Adding a column**

`000003_add_category_to_products.up.sql`:
```sql
ALTER TABLE products
ADD COLUMN category VARCHAR(100);

CREATE INDEX idx_products_category ON products(category);
```

`000003_add_category_to_products.down.sql`:
```sql
DROP INDEX IF EXISTS idx_products_category;
ALTER TABLE products DROP COLUMN IF EXISTS category;
```

### Applying Migrations

```bash
# Apply all pending migrations
make -C backend migrate-up

# Check current version
make -C backend migrate-version

# View migration status in database
docker compose exec postgres psql -U admin -d hackathon -c 'SELECT * FROM schema_migrations;'
```

**After applying migrations:**
- Verify tables exist: `docker compose exec postgres psql -U admin -d hackathon -c '\dt'`
- Test rollback in development: `make -C backend migrate-down` then `make -C backend migrate-up`

### Rollback (Development Only)

```bash
# Revert the last applied migration
make -C backend migrate-down

# Reapply it
make -C backend migrate-up
```

**WARNING:** Only use rollback in development. In production, create a new forward migration to fix issues.

### Integration with GORM Models

After creating tables via migrations, create corresponding GORM models:

1. **Create the migration** (as shown above)
2. **Apply migration:** `make -C backend migrate-up`
3. **Create GORM model in `backend/internal/models/`:**

```go
// backend/internal/models/product.go
package models

import "time"

type Product struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    Description string    `json:"description"`
    Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
    Stock       int       `json:"stock" gorm:"default:0"`
    Category    string    `json:"category"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**Important:** The GORM model should match the migration exactly. Use migrations for schema changes, GORM for queries.

### Production Deployment

**Automated (recommended):**
```bash
# Script handles migrations automatically
./scripts/deploy.sh
```

**Manual on VPS:**
```bash
git pull
docker compose up -d --build
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' up
```

### Migration Troubleshooting

**"Dirty database version" error:**
```bash
# Check version
make -C backend migrate-version

# Force to last known good version
make -C backend migrate-force version=N
```

**Migration fails midway:**
1. Check backend logs: `docker compose logs backend`
2. Connect to DB and inspect: `docker compose exec postgres psql -U admin -d hackathon`
3. Fix SQL in a NEW migration (don't edit the broken one if applied)
4. Force clean the version if needed: `make -C backend migrate-force version=N`
5. Apply new fix migration

**Testing migrations before production:**
```bash
# Create test environment
docker compose down
docker volume rm lumera_app_postgres_data
docker compose up -d
make -C backend migrate-up

# Verify everything works
# Then test rollback
make -C backend migrate-down
make -C backend migrate-up
```

### Migration Best Practices

1. **Always write both up and down migrations**
   - Even if you never plan to rollback, write the down migration

2. **Test rollback locally before deploying**
   ```bash
   make -C backend migrate-up
   make -C backend migrate-down
   make -C backend migrate-up
   ```

3. **One logical change per migration**
   - ✅ `add_products_table` - creates products table
   - ✅ `add_category_to_products` - adds one field
   - ❌ `update_schema` - too vague

4. **Use transactions for safety**
   ```sql
   BEGIN;
   CREATE TABLE products (...);
   CREATE INDEX idx_products_name ON products(name);
   COMMIT;
   ```

5. **Never modify applied migrations**
   - If migration is in production, create a new migration to fix it

6. **Use IF EXISTS / IF NOT EXISTS**
   ```sql
   CREATE TABLE IF NOT EXISTS products (...);
   DROP TABLE IF EXISTS old_table;
   ```

### Full Migration Workflow Example

User asks: "Add a products table with name, price, and stock"

**Step 1: Create migration**
```bash
make -C backend migrate-create name=create_products_table
```

**Step 2: Write up migration** (`backend/migrations/XXXXXX_create_products_table.up.sql`)
```sql
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_name ON products(name);
```

**Step 3: Write down migration** (`backend/migrations/XXXXXX_create_products_table.down.sql`)
```sql
DROP INDEX IF EXISTS idx_products_name;
DROP TABLE IF EXISTS products;
```

**Step 4: Apply migration**
```bash
make -C backend migrate-up
```

**Step 5: Create GORM model** (`backend/internal/models/product.go`)
```go
package models

import "time"

type Product struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Price     float64   `json:"price" gorm:"type:decimal(10,2)"`
    Stock     int       `json:"stock" gorm:"default:0"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Step 6: Create handler** (`backend/internal/handlers/products.go`)
```go
package handlers

import (
    "net/http"
    "encoding/json"
    "github.com/platanus-hack-25/lumera_app/internal/db"
    "github.com/platanus-hack-25/lumera_app/internal/models"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
    var products []models.Product
    db.DB.Find(&products)
    json.NewEncoder(w).Encode(products)
}
```

**Step 7: Register route** (in `cmd/main.go`)
```go
r.Get("/api/products", handlers.GetProducts)
```

**Complete documentation:** See `backend/migrations/README.md`

## Rapid CRUD Generation

For hackathon speed, use the CRUD scaffold script to generate complete REST endpoints in seconds.

### Generate Complete CRUD

```bash
./scripts/scaffold-crud.sh EntityName field1:type field2:type ...
```

**Example: Create a Product entity**
```bash
./scripts/scaffold-crud.sh Product name:string description:text price:decimal stock:int
```

This automatically generates:
1. ✅ Migration files (up and down SQL)
2. ✅ GORM model with proper types
3. ✅ Complete CRUD handlers (GET all, GET one, POST, PUT, DELETE)
4. ✅ Provides route registration code

**Supported field types:**
- `string` → VARCHAR(255)
- `text` → TEXT
- `int` → INTEGER
- `decimal` → DECIMAL(10,2)
- `bool` → BOOLEAN
- `timestamp` → TIMESTAMP

**After running scaffold:**
1. Apply migration: `make migrate-up`
2. Add routes to `cmd/main.go` (script provides exact code)
3. Restart backend: `docker compose restart backend`
4. Test endpoints immediately

**Full workflow example:**

```bash
# 1. Generate CRUD for products
./scripts/scaffold-crud.sh Product name:string price:decimal stock:int

# 2. Apply migration
make migrate-up

# 3. Add routes (copy from script output)
# Edit backend/cmd/main.go and add:
# r.Route("/api/products", func(r chi.Router) {
#     r.Get("/", handlers.GetProducts)
#     r.Post("/", handlers.CreateProduct)
#     r.Get("/{id}", handlers.GetProduct)
#     r.Put("/{id}", handlers.UpdateProduct)
#     r.Delete("/{id}", handlers.DeleteProduct)
# })

# 4. Restart
docker compose restart backend

# 5. Test
curl http://localhost:8080/api/products
```

This saves 10-15 minutes per entity during a hackathon.

## Environment Variables

**Root `.env` (Docker Compose):**
```
DB_USER=admin
DB_PASSWORD=your_secure_password
DB_NAME=hackathon
PORT=8080
```

**Backend `.env` (local development):**
```
DB_HOST=localhost  # or 'postgres' in Docker
DB_PORT=5432
PORT=8080
```

**Frontend:**
- Use `VITE_` prefix for build-time variables
- Runtime API URL configured in `vite.config.js` proxy

## Hackathon Priorities

Per STACK.md philosophy:

1. **Velocity over quality:** Fast iteration beats perfect code
2. **Visual demos:** GSAP animations, polished UI for presentations
3. **Avoid over-engineering:** No microservices, complex auth, or heavy frameworks
4. **Control over convenience:** Own infrastructure (no Supabase/Firebase)

**Antipatterns to avoid:**
- Heavy frameworks (Next.js, NestJS)
- Complex authentication (stick to simple JWT/magic links if needed)
- External services that could fail during demo
- Premature optimization

## Common Issues

**Backend won't compile after dependency changes:**
```bash
docker compose exec backend go mod tidy
docker compose restart backend
```

**Frontend can't reach backend:**
- Check proxy in `frontend/vite.config.js` points to `http://backend:8080`
- Verify CORS origins in `backend/cmd/main.go` include `http://localhost:5173`

**PostgreSQL connection errors:**
- Ensure DB_HOST=postgres in Docker environment
- Check postgres container is healthy: `docker compose ps`
- View backend logs: `docker compose logs backend`

**Port conflicts (5432, 8080, 5173):**
- Edit port mappings in `docker-compose.yml`
- Update frontend proxy and CORS accordingly

## Testing Health

```bash
# Backend API
curl http://localhost:8080/api/health

# Expected response:
# {"status":"ok","database":"connected","timestamp":"..."}

# Frontend
open http://localhost:5173
# Should display UI with API status and database connection
```
