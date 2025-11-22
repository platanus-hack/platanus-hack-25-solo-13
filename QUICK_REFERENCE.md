# Quick Reference - Lumera App

Guía rápida de comandos para desarrollo veloz en hackathon.

## 🚀 Inicio Rápido

```bash
make help          # Ver todos los comandos
make up            # Levantar todo
make logs          # Ver logs
make down          # Detener todo
```

## 🗄️ Base de Datos

### Queries Rápidos

```bash
# Ejecutar query SQL
make db-query SQL="SELECT * FROM users"

# Conectar a shell interactivo
make db-shell

# Cargar datos de prueba
make db-seed

# Resetear BD (⚠️ elimina datos)
make db-reset
```

### Ejemplos de Queries Comunes

```bash
# Contar registros
make db-query SQL="SELECT COUNT(*) FROM users"

# Ver últimos registros
make db-query SQL="SELECT * FROM users ORDER BY created_at DESC LIMIT 10"

# Buscar por email
make db-query SQL="SELECT * FROM users WHERE email LIKE '%@example.com'"

# Listar tablas
make db-query SQL="\dt"

# Describir tabla
make db-query SQL="\d users"
```

## 🔄 Migraciones

```bash
# Crear nueva migración
make migrate-create name=add_products_table

# Aplicar migraciones pendientes
make migrate-up

# Revertir última migración (solo desarrollo)
make migrate-down

# Ver versión actual
make migrate-version
```

## ⚡ Generar CRUD Completo

**Ahorra 10-15 minutos por entidad**

```bash
# Sintaxis
./scripts/scaffold-crud.sh EntityName field1:type field2:type

# Ejemplo: Productos
./scripts/scaffold-crud.sh Product name:string price:decimal stock:int

# Ejemplo: Órdenes
./scripts/scaffold-crud.sh Order status:string total:decimal user_id:int

# Luego:
make migrate-up                    # 1. Aplicar migración
# 2. Copiar rutas del output al main.go
docker compose restart backend     # 3. Reiniciar
```

**Tipos soportados:**
- `string` - VARCHAR(255)
- `text` - TEXT largo
- `int` - INTEGER
- `decimal` - DECIMAL(10,2)
- `bool` - BOOLEAN
- `timestamp` - TIMESTAMP

## 📝 Logs

```bash
make logs              # Todos
make backend-logs      # Backend
make frontend-logs     # Frontend
make postgres-logs     # PostgreSQL
```

## 🧪 Testing

```bash
# Probar health check
make test-api

# Curl manual
curl http://localhost:8080/api/health

# Probar endpoint creado
curl http://localhost:8080/api/products

# POST con datos
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Product 1","price":99.99,"stock":10}'
```

## 🔧 Servicios

```bash
make up                # Levantar
make down              # Detener
make restart           # Reiniciar
make ps                # Ver estado
make rebuild           # Rebuild completo
make clean             # Limpiar todo (⚠️ elimina datos)
```

## 🐚 Shells

```bash
make db-shell          # PostgreSQL psql
make backend-shell     # Shell en contenedor backend
```

## 📊 URLs Importantes

- **Frontend:** http://localhost:5173
- **Backend API:** http://localhost:8080
- **Health Check:** http://localhost:8080/api/health
- **PostgreSQL:** localhost:5432 (user: admin, pass: ver .env, db: hackathon)

## 🔥 Workflow Típico

### Agregar nueva feature (ej: Products)

```bash
# 1. Generar CRUD completo
./scripts/scaffold-crud.sh Product name:string price:decimal stock:int

# 2. Aplicar migración
make migrate-up

# 3. Agregar rutas (copiar del output del script)
# Editar backend/cmd/main.go:
# r.Route("/api/products", func(r chi.Router) {
#     r.Get("/", handlers.GetProducts)
#     r.Post("/", handlers.CreateProduct)
#     r.Get("/{id}", handlers.GetProduct)
#     r.Put("/{id}", handlers.UpdateProduct)
#     r.Delete("/{id}", handlers.DeleteProduct)
# })

# 4. Reiniciar backend
docker compose restart backend

# 5. Probar
curl http://localhost:8080/api/products

# 6. Cargar datos de prueba si es necesario
make db-seed
```

### Debugging

```bash
# Ver logs en tiempo real
make backend-logs

# Ejecutar query para ver datos
make db-query SQL="SELECT * FROM products"

# Conectar a BD para investigar
make db-shell
```

### Resetear ambiente

```bash
# Resetear solo BD (⚠️ elimina datos)
make db-reset

# Limpiar todo y empezar de cero (⚠️⚠️)
make clean
make up
make migrate-up
make db-seed
```

## 💡 Tips

1. **Siempre ejecuta `make help` si olvidas un comando**
2. **Usa `make db-query` para queries rápidos sin conectarte a psql**
3. **El scaffold CRUD ahorra mucho tiempo - úsalo**
4. **`make db-seed` después de `make db-reset` para tener datos**
5. **Ver `scripts/db-queries.md` para más ejemplos de queries**
6. **Commit frecuentemente - las migraciones son versionadas**

## 🚨 Troubleshooting Rápido

```bash
# Backend no compila
docker compose exec backend go mod tidy
docker compose restart backend

# Frontend no conecta con backend
# Verificar proxy en frontend/vite.config.js
# Verificar CORS en backend/cmd/main.go

# PostgreSQL no inicia
make ps                          # Ver estado
docker compose logs postgres     # Ver logs

# Puerto en uso
make down                        # Detener todo
# Cambiar puertos en docker-compose.yml si es necesario

# BD corrupta
make db-reset                    # Resetear todo
make migrate-up                  # Aplicar migraciones
make db-seed                     # Cargar datos
```

---

**Para documentación completa ver:**
- `README.md` - Guía general
- `CLAUDE.md` - Guía para Claude Code
- `backend/migrations/README.md` - Migraciones detalladas
- `scripts/db-queries.md` - Queries SQL útiles
