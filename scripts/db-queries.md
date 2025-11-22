# Queries SQL Útiles

Colección de queries SQL para desarrollo rápido.

## Ejecución Rápida

```bash
# Ejecutar query directamente
make db-query SQL="SELECT * FROM users LIMIT 5"

# Conectar a shell interactivo
make db-shell
```

## Queries Comunes

### Ver todas las tablas
```bash
make db-query SQL="\dt"
```

### Contar registros
```bash
make db-query SQL="SELECT COUNT(*) FROM users"
```

### Ver últimos registros
```bash
make db-query SQL="SELECT * FROM users ORDER BY created_at DESC LIMIT 10"
```

### Buscar por campo
```bash
make db-query SQL="SELECT * FROM users WHERE email LIKE '%@example.com'"
```

### Limpiar tabla (desarrollo)
```bash
make db-query SQL="TRUNCATE users CASCADE"
```

### Ver estructura de tabla
```bash
make db-query SQL="\d users"
```

## Queries de Análisis

### Registros por fecha
```sql
SELECT
    DATE(created_at) as date,
    COUNT(*) as count
FROM users
GROUP BY DATE(created_at)
ORDER BY date DESC;
```

Ejecutar:
```bash
make db-query SQL="SELECT DATE(created_at) as date, COUNT(*) as count FROM users GROUP BY DATE(created_at) ORDER BY date DESC"
```

### Búsqueda de texto
```sql
SELECT * FROM users
WHERE name ILIKE '%john%'
OR email ILIKE '%john%';
```

### Actualizar registros
```sql
UPDATE users
SET name = 'Updated Name'
WHERE email = 'alice@example.com';
```

### Eliminar registros antiguos
```sql
DELETE FROM users
WHERE created_at < NOW() - INTERVAL '30 days';
```

## Tips

### Formato bonito en psql
```bash
docker compose exec postgres psql -U admin -d hackathon -c "\x" -c "SELECT * FROM users LIMIT 1"
```

### Exportar a CSV
```bash
docker compose exec postgres psql -U admin -d hackathon -c "\COPY (SELECT * FROM users) TO STDOUT WITH CSV HEADER" > users.csv
```

### Ver tamaño de tablas
```sql
SELECT
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

### Ver queries lentas (si tienes muchos datos)
```sql
SELECT
    query,
    calls,
    total_time,
    mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;
```
