#!/bin/bash

# =============================================================================
# Lumera App - Deployment Script
# =============================================================================
# Automatiza el proceso completo de deployment:
# 1. Actualiza código desde git
# 2. Construye y levanta contenedores
# 3. Ejecuta migraciones de base de datos
# 4. Verifica estado del deployment
#
# Uso:
#   ./scripts/deploy.sh [--skip-git] [--skip-build] [--skip-migrations]
#
# Opciones:
#   --skip-git         No hace git pull
#   --skip-build       No reconstruye las imágenes Docker
#   --skip-migrations  No ejecuta migraciones (no recomendado)
#   --env=<file>       Usa archivo de entorno específico (default: .env)
# =============================================================================

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default options
SKIP_GIT=false
SKIP_BUILD=false
SKIP_MIGRATIONS=false
ENV_FILE=".env"

# Parse arguments
for arg in "$@"; do
    case $arg in
        --skip-git)
            SKIP_GIT=true
            shift
            ;;
        --skip-build)
            SKIP_BUILD=true
            shift
            ;;
        --skip-migrations)
            SKIP_MIGRATIONS=true
            shift
            ;;
        --env=*)
            ENV_FILE="${arg#*=}"
            shift
            ;;
        --help)
            grep "^#" "$0" | grep -v "^#!/" | sed 's/^# //'
            exit 0
            ;;
        *)
            echo -e "${RED}Error: Argumento desconocido: $arg${NC}"
            echo "Usa --help para ver opciones disponibles"
            exit 1
            ;;
    esac
done

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1"
}

separator() {
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# =============================================================================
# Pre-deployment checks
# =============================================================================
separator
log_info "Lumera App - Deployment Automatizado"
separator

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    log_error "Docker no está corriendo. Por favor inicia Docker Desktop."
    exit 1
fi
log_success "Docker está corriendo"

# Check if docker-compose exists
if ! command -v docker &> /dev/null; then
    log_error "docker compose no está instalado"
    exit 1
fi
log_success "docker compose disponible"

# Check if we're in the right directory
if [ ! -f "docker-compose.yml" ]; then
    log_error "No se encuentra docker-compose.yml. Ejecuta este script desde la raíz del proyecto."
    exit 1
fi
log_success "Directorio correcto detectado"

# Check environment file
if [ ! -f "$ENV_FILE" ]; then
    log_warning "Archivo de entorno $ENV_FILE no encontrado"
    read -p "¿Deseas continuar de todas formas? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_info "Deployment cancelado"
        exit 0
    fi
else
    log_success "Archivo de entorno $ENV_FILE encontrado"
fi

# =============================================================================
# Step 1: Update code from git
# =============================================================================
if [ "$SKIP_GIT" = false ]; then
    separator
    log_info "Paso 1/4: Actualizando código desde git..."

    # Check if git repo
    if [ ! -d ".git" ]; then
        log_warning "No es un repositorio git, saltando git pull"
    else
        # Check for uncommitted changes
        if [[ -n $(git status -s) ]]; then
            log_warning "Hay cambios sin commitear en el repositorio"
            git status -s
            read -p "¿Deseas continuar de todas formas? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log_info "Deployment cancelado"
                exit 0
            fi
        fi

        # Get current branch
        CURRENT_BRANCH=$(git branch --show-current)
        log_info "Branch actual: $CURRENT_BRANCH"

        # Pull latest changes
        if git pull origin "$CURRENT_BRANCH"; then
            log_success "Código actualizado desde git"
        else
            log_error "Error al hacer git pull"
            exit 1
        fi
    fi
else
    log_warning "Saltando actualización de git (--skip-git)"
fi

# =============================================================================
# Step 2: Build and start containers
# =============================================================================
separator
log_info "Paso 2/4: Construyendo y levantando contenedores..."

# Stop existing containers
log_info "Deteniendo contenedores existentes..."
docker compose down

# Build and start
if [ "$SKIP_BUILD" = false ]; then
    log_info "Construyendo imágenes Docker..."
    if docker compose up -d --build; then
        log_success "Contenedores levantados exitosamente"
    else
        log_error "Error al levantar contenedores"
        exit 1
    fi
else
    log_warning "Saltando rebuild de imágenes (--skip-build)"
    if docker compose up -d; then
        log_success "Contenedores levantados exitosamente"
    else
        log_error "Error al levantar contenedores"
        exit 1
    fi
fi

# Wait for containers to be healthy
log_info "Esperando a que los contenedores estén listos..."
sleep 5

# Check if containers are running
if ! docker compose ps | grep -q "Up"; then
    log_error "Los contenedores no están corriendo correctamente"
    docker compose ps
    exit 1
fi
log_success "Todos los contenedores están corriendo"

# =============================================================================
# Step 3: Run database migrations
# =============================================================================
if [ "$SKIP_MIGRATIONS" = false ]; then
    separator
    log_info "Paso 3/4: Ejecutando migraciones de base de datos..."

    # Wait a bit more for PostgreSQL to be fully ready
    log_info "Esperando a que PostgreSQL esté completamente listo..."
    sleep 3

    # Load DB credentials from .env
    if [ -f "$ENV_FILE" ]; then
        export $(grep -v '^#' "$ENV_FILE" | xargs)
    fi

    # Use default values if not set
    DB_USER=${DB_USER:-admin}
    DB_PASSWORD=${DB_PASSWORD:-hackathon2025}
    DB_NAME=${DB_NAME:-hackathon}

    # Build connection string
    DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable"

    # Check current migration version
    log_info "Versión actual de migraciones:"
    docker compose exec backend migrate \
        -path=/app/migrations \
        -database="$DB_URL" \
        version || log_warning "No se pudo obtener versión (puede ser primera ejecución)"

    # Run migrations
    log_info "Aplicando migraciones pendientes..."
    if docker compose exec backend migrate \
        -path=/app/migrations \
        -database="$DB_URL" \
        up; then
        log_success "Migraciones aplicadas exitosamente"
    else
        log_error "Error al aplicar migraciones"
        log_info "Revisa los logs: docker compose logs backend"
        exit 1
    fi

    # Show final migration version
    log_info "Versión final de migraciones:"
    docker compose exec backend migrate \
        -path=/app/migrations \
        -database="$DB_URL" \
        version
else
    log_warning "Saltando migraciones de base de datos (--skip-migrations) - NO RECOMENDADO"
fi

# =============================================================================
# Step 4: Verify deployment
# =============================================================================
separator
log_info "Paso 4/4: Verificando deployment..."

# Wait for backend to be ready
log_info "Esperando a que el backend esté listo..."
sleep 2

# Check health endpoint
log_info "Verificando endpoint de salud..."
if curl -f -s http://localhost:8080/api/health > /dev/null; then
    log_success "Backend respondiendo correctamente"

    # Show health status
    HEALTH_RESPONSE=$(curl -s http://localhost:8080/api/health)
    echo "$HEALTH_RESPONSE" | jq '.' 2>/dev/null || echo "$HEALTH_RESPONSE"
else
    log_warning "Backend aún no está respondiendo (puede tomar unos segundos más)"
fi

# Show running containers
log_info "Contenedores activos:"
docker compose ps

# =============================================================================
# Deployment complete
# =============================================================================
separator
log_success "¡Deployment completado exitosamente!"
separator

echo ""
log_info "Servicios disponibles:"
echo "  • Backend API:    http://localhost:8080"
echo "  • Frontend:       http://localhost:5173"
echo "  • PostgreSQL:     localhost:5432"
echo ""
log_info "Comandos útiles:"
echo "  • Ver logs:           docker compose logs -f"
echo "  • Ver logs backend:   docker compose logs -f backend"
echo "  • Ver logs frontend:  docker compose logs -f frontend"
echo "  • Detener servicios:  docker compose down"
echo "  • Estado servicios:   docker compose ps"
echo ""

# Ask if user wants to see logs
read -p "¿Deseas ver los logs del backend? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    docker compose logs -f backend
fi
