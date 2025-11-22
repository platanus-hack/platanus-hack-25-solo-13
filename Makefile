# Lumera App - Makefile de Desarrollo
# Comandos rápidos para acelerar desarrollo

.PHONY: help up down logs restart ps clean db-shell db-query db-reset db-seed migrate-create migrate-up migrate-down migrate-version backend-logs frontend-logs landing-logs postgres-logs backend-shell test-api docs

# Default target
.DEFAULT_GOAL := help

## help: Mostrar esta ayuda
help:
	@echo "Comandos disponibles:"
	@echo ""
	@echo "🚀 Servicios:"
	@echo "  make up              - Levantar todos los servicios"
	@echo "  make down            - Detener todos los servicios"
	@echo "  make restart         - Reiniciar todos los servicios"
	@echo "  make ps              - Ver estado de servicios"
	@echo "  make clean           - Limpiar todo (⚠️  elimina volúmenes)"
	@echo ""
	@echo "📝 Logs:"
	@echo "  make logs            - Ver logs de todos los servicios"
	@echo "  make backend-logs    - Ver logs del backend"
	@echo "  make frontend-logs   - Ver logs del frontend"
	@echo "  make landing-logs    - Ver logs del landing"
	@echo "  make postgres-logs   - Ver logs de PostgreSQL"
	@echo ""
	@echo "🗄️  Base de Datos:"
	@echo "  make db-shell        - Conectar a PostgreSQL (psql)"
	@echo "  make db-query SQL='' - Ejecutar query SQL"
	@echo "  make db-reset        - Resetear BD (⚠️  elimina datos)"
	@echo "  make db-seed         - Cargar datos de prueba"
	@echo "  make db-ui           - Abrir Adminer (GUI web)"
	@echo ""
	@echo "🔄 Migraciones:"
	@echo "  make migrate-create name=nombre  - Crear migración"
	@echo "  make migrate-up                  - Aplicar migraciones"
	@echo "  make migrate-down                - Revertir última migración"
	@echo "  make migrate-version             - Ver versión actual"
	@echo ""
	@echo "🔧 Desarrollo:"
	@echo "  make backend-shell   - Shell en contenedor backend"
	@echo "  make test-api        - Probar endpoint de health"
	@echo "  make docs            - Regenerar documentación Swagger"
	@echo "  make rebuild         - Rebuild completo"
	@echo ""

## up: Levantar servicios
up:
	@echo "🚀 Levantando servicios..."
	docker compose up -d
	@echo "✅ Servicios levantados"
	@echo ""
	@echo "Landing:   http://localhost:4321  (Astro)"
	@echo "Frontend:  http://localhost:5173  (Svelte)"
	@echo "Backend:   http://localhost:8080  (Go API)"
	@echo "Adminer:   http://localhost:8081  (DB GUI)"
	@echo "Database:  localhost:5432"
	@echo ""
	@echo "💡 Tip: Usa 'make db-ui' para abrir Adminer en el navegador"

## down: Detener servicios
down:
	@echo "🛑 Deteniendo servicios..."
	docker compose down

## logs: Ver todos los logs
logs:
	docker compose logs -f

## restart: Reiniciar servicios
restart:
	@echo "🔄 Reiniciando servicios..."
	docker compose restart

## ps: Ver estado de servicios
ps:
	docker compose ps

## clean: Limpiar todo (⚠️  ELIMINA VOLÚMENES)
clean:
	@echo "⚠️  ADVERTENCIA: Esto eliminará TODOS los datos de la BD"
	@read -p "¿Estás seguro? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker compose down -v; \
		echo "✅ Todo limpio"; \
	else \
		echo "❌ Cancelado"; \
	fi

## db-shell: Conectar a PostgreSQL
db-shell:
	@echo "🗄️  Conectando a PostgreSQL..."
	@echo "Comandos útiles:"
	@echo "  \\dt          - Listar tablas"
	@echo "  \\d tabla     - Describir tabla"
	@echo "  \\q           - Salir"
	@echo ""
	docker compose exec postgres psql -U admin -d hackathon

## db-query: Ejecutar query SQL
db-query:
	@if [ -z "$(SQL)" ]; then \
		echo "❌ Error: Debes especificar SQL=''"; \
		echo "Ejemplo: make db-query SQL='SELECT * FROM users LIMIT 5'"; \
		exit 1; \
	fi
	@docker compose exec postgres psql -U admin -d hackathon -c "$(SQL)"

## db-reset: Resetear base de datos
db-reset:
	@echo "⚠️  ADVERTENCIA: Esto eliminará TODOS los datos"
	@read -p "¿Estás seguro? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		echo "🗑️  Eliminando volumen..."; \
		docker compose down; \
		docker volume rm lumera_app_postgres_data || true; \
		echo "🚀 Levantando servicios..."; \
		docker compose up -d; \
		echo "⏳ Esperando PostgreSQL..."; \
		sleep 5; \
		echo "🔄 Aplicando migraciones..."; \
		$(MAKE) migrate-up; \
		echo "✅ BD reseteada y migraciones aplicadas"; \
	else \
		echo "❌ Cancelado"; \
	fi

## db-seed: Cargar datos de prueba
db-seed:
	@echo "🌱 Cargando datos de prueba..."
	@if [ -f scripts/seed.sql ]; then \
		docker compose exec -T postgres psql -U admin -d hackathon < scripts/seed.sql; \
		echo "✅ Datos cargados"; \
	else \
		echo "⚠️  No existe scripts/seed.sql"; \
		echo "Crear archivo con datos de ejemplo"; \
	fi

## db-ui: Abrir Adminer en el navegador
db-ui:
	@echo "🌐 Abriendo Adminer..."
	@echo ""
	@echo "URL: http://localhost:8081"
	@echo ""
	@echo "Credenciales:"
	@echo "  System:   PostgreSQL"
	@echo "  Server:   postgres"
	@echo "  Username: admin"
	@echo "  Password: (ver archivo .env)"
	@echo "  Database: hackathon"
	@echo ""
	@open http://localhost:8081 2>/dev/null || xdg-open http://localhost:8081 2>/dev/null || echo "Abre http://localhost:8081 en tu navegador"

## backend-logs: Ver logs del backend
backend-logs:
	docker compose logs -f backend

## frontend-logs: Ver logs del frontend
frontend-logs:
	docker compose logs -f frontend

## landing-logs: Ver logs del landing
landing-logs:
	docker compose logs -f landing

## postgres-logs: Ver logs de PostgreSQL
postgres-logs:
	docker compose logs -f postgres

## backend-shell: Shell en contenedor backend
backend-shell:
	docker compose exec backend sh

## test-api: Probar endpoint de health
test-api:
	@echo "🧪 Probando API..."
	@curl -s http://localhost:8080/api/health | python3 -m json.tool || echo "❌ API no responde"

## rebuild: Rebuild completo
rebuild:
	@echo "🔨 Rebuilding..."
	docker compose down
	docker compose build --no-cache
	docker compose up -d
	@echo "✅ Rebuild completo"

## migrate-create: Crear nueva migración
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Debes especificar name=nombre"; \
		echo "Ejemplo: make migrate-create name=add_products_table"; \
		exit 1; \
	fi
	$(MAKE) -C backend migrate-create name=$(name)

## migrate-up: Aplicar migraciones
migrate-up:
	$(MAKE) -C backend migrate-up

## migrate-down: Revertir última migración
migrate-down:
	$(MAKE) -C backend migrate-down

## migrate-version: Ver versión de migración
migrate-version:
	$(MAKE) -C backend migrate-version

## docs: Regenerar documentación Swagger
docs:
	@echo "📚 Regenerando documentación Swagger..."
	docker compose exec backend /go/bin/swag init -g cmd/main.go -o ./docs
	@echo "✅ Documentación generada en backend/docs/"
	@echo ""
	@echo "Archivos generados:"
	@echo "  - backend/docs/swagger.json"
	@echo "  - backend/docs/swagger.yaml"
	@echo ""
	@echo "💡 Para visualizar: Importa swagger.json en https://editor.swagger.io"
