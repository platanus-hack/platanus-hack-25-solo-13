# Deployment Guide - Lumera App

Esta guía explica cómo hacer deployment de Lumera App en diferentes ambientes.

## 🚀 Deployment Automatizado (Recomendado)

### Uso Básico

```bash
# Deployment completo (git + build + migraciones + verificación)
./scripts/deploy.sh

# Desarrollo local (sin git pull)
./scripts/deploy.sh --skip-git

# Deployment rápido (sin rebuild)
./scripts/deploy.sh --skip-build

# Saltarse migraciones (NO RECOMENDADO)
./scripts/deploy.sh --skip-migrations
```

### Ver Ayuda

```bash
./scripts/deploy.sh --help
```

## 📋 Qué Hace el Script Automáticamente

1. **Pre-deployment Checks**
   - Verifica que Docker esté corriendo
   - Verifica que docker compose esté instalado
   - Verifica que estés en el directorio correcto
   - Verifica que exista el archivo `.env`
   - Detecta cambios sin commitear en git (pregunta si continuar)

2. **Git Update** (opcional)
   - Obtiene el branch actual
   - Hace `git pull origin [branch]`
   - Puede saltearse con `--skip-git`

3. **Docker Build & Start**
   - Detiene contenedores existentes
   - Reconstruye imágenes (`--build`)
   - Levanta contenedores en modo detached
   - Espera a que estén saludables
   - Puede saltearse rebuild con `--skip-build`

4. **Database Migrations** (crítico)
   - Espera a que PostgreSQL esté listo
   - Carga credenciales desde `.env`
   - Muestra versión actual de migraciones
   - **Ejecuta migraciones pendientes automáticamente**
   - Muestra versión final
   - Incluye seed data (cursos, materias, OAs, preguntas)

5. **Verification**
   - Verifica endpoint de salud del backend
   - Muestra estado de contenedores
   - Lista servicios disponibles
   - Pregunta si quieres ver logs

## 🎯 Escenarios de Uso

### Desarrollo Local

```bash
# Primera vez
./scripts/deploy.sh --skip-git

# Cambios posteriores (ya tengo código actualizado)
./scripts/deploy.sh --skip-git

# Cambios solo de código (no Dockerfile)
./scripts/deploy.sh --skip-git --skip-build
```

### Servidor de Desarrollo/Staging

```bash
# SSH al servidor
ssh user@dev.lumera.com

cd /var/www/lumera_app

# Deployment completo
./scripts/deploy.sh

# El script automáticamente:
# - Hará git pull
# - Reconstruirá imágenes
# - Ejecutará migraciones
# - Verificará que todo funcione
```

### Producción

```bash
# SSH al servidor de producción
ssh user@lumera.com

cd /var/www/lumera_app

# Asegúrate de tener .env configurado
# DB_PASSWORD debe ser seguro
# JWT_SECRET debe estar configurado

# Deployment
./scripts/deploy.sh

# Verificar logs después
docker compose logs -f backend
```

## 🗄️ Migraciones Incluidas en Seed

El script ejecuta automáticamente estas migraciones que incluyen datos iniciales:

| Migración | Contenido | Cantidad |
|-----------|-----------|----------|
| 17 | Cursos y Materias | 1 curso, 2 materias, 1 relación |
| 18 | Objetivos Aprendizaje (Lengua) | 23 OAs |
| 19 | Bloom Objectives | 134 objectives |
| 20 | Preguntas | 672 preguntas |

**Total:** 672 preguntas listas para usar en el sistema de diagnóstico.

### Idempotencia

Todas las migraciones usan `ON CONFLICT DO UPDATE SET`, lo que significa:
- ✅ Puedes ejecutar el deployment múltiples veces
- ✅ No habrá errores de duplicados
- ✅ Datos existentes se actualizan si hay cambios
- ✅ Seguro para re-deploys

## 🔍 Troubleshooting

### Error: "Docker no está corriendo"

```bash
# macOS
open -a Docker

# Linux
sudo systemctl start docker
```

### Error: "Dirty database version"

```bash
# Ver qué migración está sucia
docker compose exec backend migrate -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' version

# Forzar a versión anterior (reemplaza N con el número correcto)
docker compose exec backend migrate -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' force N

# Re-ejecutar deployment
./scripts/deploy.sh --skip-git
```

### Error: "Backend no responde"

```bash
# Ver logs del backend
docker compose logs -f backend

# Verificar que el contenedor esté corriendo
docker compose ps

# Reiniciar solo el backend
docker compose restart backend

# Verificar salud
curl http://localhost:8080/api/health
```

### Error: "Migrations failed"

```bash
# Ver logs detallados
docker compose logs backend

# Conectar a la base de datos para investigar
docker compose exec postgres psql -U admin -d hackathon

# Dentro de psql:
# Ver tablas
\dt

# Ver versión de migraciones
SELECT * FROM schema_migrations;

# Salir
\q
```

### Resetear Todo Desde Cero

```bash
# ADVERTENCIA: Esto borra TODOS los datos

# Detener y eliminar volúmenes
docker compose down -v

# Re-deployment completo
./scripts/deploy.sh --skip-git

# Esto creará base de datos limpia y aplicará todas las migraciones
```

## 📊 Verificación Post-Deployment

### 1. Verificar Servicios

```bash
docker compose ps

# Todos deben estar "Up" y "healthy"
```

### 2. Verificar Backend

```bash
curl http://localhost:8080/api/health

# Debe retornar:
# {"status":"ok","database":"connected","timestamp":"..."}
```

### 3. Verificar Migraciones

```bash
docker compose exec backend migrate -path=/app/migrations \
  -database='postgres://admin:hackathon2025@postgres:5432/hackathon?sslmode=disable' version

# Debe mostrar: 20
```

### 4. Verificar Datos

```bash
# Conectar a base de datos
docker compose exec postgres psql -U admin -d hackathon

# Verificar datos cargados
SELECT COUNT(*) FROM cursos;           -- Debe ser 1
SELECT COUNT(*) FROM materias;         -- Debe ser 2
SELECT COUNT(*) FROM objetivos_aprendizaje; -- Debe ser 23
SELECT COUNT(*) FROM oa_bloom_objectives;   -- Debe ser 134
SELECT COUNT(*) FROM questions;        -- Debe ser 672
```

### 5. Verificar Frontend

```bash
# Abrir en navegador
open http://localhost:5173

# Debe cargar la aplicación Svelte
```

## 🔐 Variables de Entorno

El script lee credenciales desde `.env`:

```env
# Database
DB_USER=admin
DB_PASSWORD=tu_password_aqui
DB_NAME=hackathon
DB_HOST=postgres
DB_PORT=5432

# Backend
PORT=8080
JWT_SECRET=tu_secret_seguro_aqui

# CORS (separado por comas)
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:4321
```

**IMPORTANTE para producción:**
- Cambia `DB_PASSWORD` a algo seguro
- Genera un `JWT_SECRET` aleatorio y largo
- Actualiza `ALLOWED_ORIGINS` con tu dominio real

## 📝 Logs y Monitoreo

```bash
# Ver todos los logs
docker compose logs -f

# Solo backend
docker compose logs -f backend

# Solo frontend
docker compose logs -f frontend

# Solo postgres
docker compose logs -f postgres

# Ver últimas 100 líneas
docker compose logs --tail=100 backend

# Ver logs desde hace 10 minutos
docker compose logs --since=10m backend
```

## 🎛️ Comandos Manuales (Avanzado)

Si necesitas ejecutar pasos individuales sin el script:

```bash
# 1. Git pull
git pull origin main

# 2. Rebuild containers
docker compose down
docker compose up -d --build

# 3. Run migrations
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:PASSWORD@postgres:5432/hackathon?sslmode=disable' up

# 4. Verify
curl http://localhost:8080/api/health
docker compose ps
```

## 🌐 Deployment en Cloud

### AWS EC2 / DigitalOcean / Linode

1. Crear instancia con Docker instalado
2. Configurar firewall (puertos 80, 443, 8080)
3. SSH a la instancia
4. Clonar repositorio
5. Configurar `.env` con credenciales de producción
6. Ejecutar `./scripts/deploy.sh`
7. Opcional: Configurar nginx como reverse proxy

### Docker Swarm / Kubernetes

El script está diseñado para `docker compose`. Para orquestadores:
- Adaptar `docker-compose.yml` a formato de swarm/k8s
- Ejecutar migraciones como Job/Init Container
- Usar ConfigMaps/Secrets para `.env`

## 💡 Tips y Mejores Prácticas

1. **Siempre revisa los logs después del deployment:**
   ```bash
   ./scripts/deploy.sh && docker compose logs -f backend
   ```

2. **Usa `--skip-git` en desarrollo local:**
   ```bash
   # Evita conflictos con cambios locales
   ./scripts/deploy.sh --skip-git
   ```

3. **Backup antes de deployment en producción:**
   ```bash
   # Backup de base de datos
   docker compose exec postgres pg_dump -U admin hackathon > backup_$(date +%Y%m%d_%H%M%S).sql

   # Luego deploy
   ./scripts/deploy.sh
   ```

4. **Usa tags de git para releases:**
   ```bash
   git tag -a v1.0.0 -m "Release 1.0.0"
   git push origin v1.0.0

   # En servidor
   git checkout v1.0.0
   ./scripts/deploy.sh --skip-git
   ```

5. **Monitorea después del deployment:**
   ```bash
   # Terminal 1: Logs
   docker compose logs -f backend

   # Terminal 2: Health checks
   watch -n 5 'curl -s http://localhost:8080/api/health | jq'
   ```

## 🆘 Soporte

Si tienes problemas:

1. Revisa esta guía de troubleshooting
2. Consulta [CLAUDE.md](../CLAUDE.md) para documentación completa
3. Revisa logs: `docker compose logs -f backend`
4. Verifica estado: `docker compose ps`
5. En último caso: resetear todo con `docker compose down -v && ./scripts/deploy.sh`

---

**Última actualización:** 2025-11-22
