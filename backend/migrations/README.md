# Database Migrations

Este directorio contiene las migraciones de base de datos usando golang-migrate.

## Comandos Rápidos

```bash
# Crear nueva migración
make migrate-create name=create_products_table

# Aplicar todas las migraciones pendientes
make migrate-up

# Revertir la última migración
make migrate-down

# Ver versión actual de la BD
make migrate-version
```

## Estructura de Archivos

Cada migración tiene dos archivos:

- `XXXXXX_nombre.up.sql` - Aplicar cambios (crear tablas, columnas, etc.)
- `XXXXXX_nombre.down.sql` - Revertir cambios (rollback)

Ejemplo:
```
000001_init_schema.up.sql
000001_init_schema.down.sql
000002_add_products.up.sql
000002_add_products.down.sql
```

## Crear Nueva Migración

```bash
# Esto crea dos archivos: XXXXXX_nombre.up.sql y XXXXXX_nombre.down.sql
make migrate-create name=add_products_table
```

Luego edita ambos archivos:

**XXXXXX_add_products_table.up.sql:**
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**XXXXXX_add_products_table.down.sql:**
```sql
DROP TABLE IF EXISTS products;
```

## Aplicar Migraciones en Producción

### Opción 1: Desde Docker Compose (recomendado para VPS)

```bash
# En el VPS, después de git pull
docker compose exec backend migrate -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' up
```

### Opción 2: Archivo script de deploy

Crear `scripts/deploy.sh`:
```bash
#!/bin/bash
git pull
docker compose build backend
docker compose up -d
docker compose exec backend migrate -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' up
```

## Verificar Estado

```bash
# Ver qué versión está aplicada
make migrate-version

# Ver en la BD qué migraciones se han ejecutado
docker compose exec postgres psql -U admin -d hackathon -c "SELECT * FROM schema_migrations;"
```

## Rollback (Emergencia)

```bash
# Revertir última migración
make migrate-down

# Forzar una versión específica (usar solo si hay problemas)
make migrate-force version=1
```

## Troubleshooting

### Error "Dirty database version"

Si una migración falló a medias:

```bash
# Ver versión actual
make migrate-version

# Forzar la versión limpia
make migrate-force version=N  # N = última versión conocida buena
```

### Recrear BD desde cero (DESARROLLO SOLO)

```bash
docker compose down
docker volume rm lumera_app_postgres_data
docker compose up -d
make migrate-up
```

## Best Practices

1. **Nunca edites migraciones ya aplicadas en producción**
   - Crea una nueva migración para hacer cambios

2. **Siempre escribe el .down.sql**
   - Necesitas poder hacer rollback en emergencias

3. **Prueba rollbacks en desarrollo**
   ```bash
   make migrate-up
   make migrate-down
   make migrate-up
   ```

4. **Una migración = un cambio conceptual**
   - ✅ `add_users_table` - crear tabla users
   - ✅ `add_email_to_users` - agregar columna
   - ❌ `update_everything` - demasiado amplio

5. **Usa transacciones cuando sea posible**
   ```sql
   BEGIN;
   -- tus cambios aquí
   COMMIT;
   ```
