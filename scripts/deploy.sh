#!/bin/bash
set -e  # Exit on error

echo "🚀 Deploying Lumera App..."

# 1. Pull latest code
echo "📥 Pulling latest code from git..."
git pull origin main

# 2. Rebuild images
echo "🔨 Building Docker images..."
docker compose build

# 3. Stop old containers
echo "🛑 Stopping old containers..."
docker compose down

# 4. Start new containers
echo "▶️  Starting new containers..."
docker compose up -d

# 5. Wait for PostgreSQL to be ready
echo "⏳ Waiting for PostgreSQL to be healthy..."
timeout 30 sh -c 'until docker compose exec postgres pg_isready -U admin; do sleep 1; done'

# 6. Run database migrations
echo "🗄️  Running database migrations..."
docker compose exec backend migrate -path=/app/migrations \
  -database="postgres://admin:${DB_PASSWORD}@postgres:5432/hackathon?sslmode=disable" up

# 7. Verify services are running
echo "✅ Verifying services..."
docker compose ps

echo ""
echo "✨ Deploy completed successfully!"
echo ""
echo "Services:"
echo "  - Frontend: http://localhost:5173"
echo "  - Backend:  http://localhost:8080"
echo "  - Database: localhost:5432"
echo ""
echo "To view logs: docker compose logs -f"
