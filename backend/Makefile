.PHONY: migrate-create migrate-up migrate-down migrate-force migrate-version

# Crear nueva migración
# Uso: make migrate-create name=create_users_table
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: Debes especificar un nombre. Uso: make migrate-create name=nombre_migracion"; \
		exit 1; \
	fi
	docker compose exec backend sh -c "migrate create -ext sql -dir /app/migrations -seq $(name)"

# Aplicar todas las migraciones pendientes
migrate-up:
	docker compose exec backend sh -c "migrate -path=/app/migrations -database='postgres://\$$DB_USER:\$$DB_PASSWORD@postgres:5432/hackathon?sslmode=disable' up"

# Revertir la última migración
migrate-down:
	docker compose exec backend sh -c "migrate -path=/app/migrations -database='postgres://\$$DB_USER:\$$DB_PASSWORD@postgres:5432/hackathon?sslmode=disable' down 1"

# Forzar versión específica (usar con cuidado)
# Uso: make migrate-force version=1
migrate-force:
	@if [ -z "$(version)" ]; then \
		echo "Error: Debes especificar una versión. Uso: make migrate-force version=1"; \
		exit 1; \
	fi
	docker compose exec backend sh -c "migrate -path=/app/migrations -database='postgres://\$$DB_USER:\$$DB_PASSWORD@postgres:5432/hackathon?sslmode=disable' force $(version)"

# Ver versión actual de la BD
migrate-version:
	docker compose exec backend sh -c "migrate -path=/app/migrations -database='postgres://\$$DB_USER:\$$DB_PASSWORD@postgres:5432/hackathon?sslmode=disable' version"

# Ayuda
help:
	@echo "Comandos de migración disponibles:"
	@echo "  make migrate-create name=nombre  - Crear nueva migración"
	@echo "  make migrate-up                  - Aplicar migraciones pendientes"
	@echo "  make migrate-down                - Revertir última migración"
	@echo "  make migrate-version             - Ver versión actual"
	@echo "  make migrate-force version=N     - Forzar versión (usar con cuidado)"
