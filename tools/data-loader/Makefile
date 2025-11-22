.PHONY: help run build deps clean test example

help: ## Mostrar ayuda
	@echo "Comandos disponibles:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

deps: ## Descargar dependencias
	go mod download
	go mod tidy

build: ## Compilar el binario
	go build -o data-loader main.go

run: ## Ejecutar con archivo por defecto (input/oas.csv)
	go run main.go

example: ## Ejecutar con archivo de ejemplo (input/oas_example.csv)
	go run main.go --input=input/oas_example.csv

custom: ## Ejecutar con archivo custom (usar: make custom FILE=input/mi_archivo.csv)
	go run main.go --input=$(FILE)

clean: ## Limpiar binarios y archivos generados
	rm -f data-loader
	rm -f output/failed_oas_*.json

test: ## Verificar que el código compila
	go build -v ./...

setup: ## Setup inicial: copiar .env.example a .env
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "✅ Archivo .env creado. Por favor edítalo con tus credenciales."; \
	else \
		echo "⚠️  .env ya existe, no se sobrescribe."; \
	fi

check-env: ## Verificar que .env esté configurado
	@if [ ! -f .env ]; then \
		echo "❌ Error: .env no existe. Ejecuta 'make setup' primero."; \
		exit 1; \
	fi
	@if ! grep -q "OPENAI_API_KEY=sk-" .env; then \
		echo "⚠️  Advertencia: OPENAI_API_KEY parece no estar configurada en .env"; \
	fi
	@if ! grep -q "DB_PASSWORD=.." .env; then \
		echo "⚠️  Advertencia: DB_PASSWORD parece no estar configurada en .env"; \
	fi
	@echo "✅ .env existe y parece configurado"
