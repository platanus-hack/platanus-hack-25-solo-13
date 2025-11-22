# Lumera - Plataforma Educativa Personalizada

<img src="./project-logo.png" alt="Project Logo" width="200" />

**Platanus Hack 25** - Track: ✨ consumer AI

Sistema de diagnóstico adaptativo y personalización educativa con IA para estudiantes de enseñanza media.

## Integrantes

- Jonathan Olivares ([@jcoruiz](https://github.com/jcoruiz))

---

## 🚀 Quick Start

### Requisitos Previos
- Docker Desktop instalado y corriendo
- Git

### Deployment Automatizado

El proyecto incluye un script de deployment que automatiza **todo** el proceso:

```bash
# Deployment completo (git pull + build + migraciones + verificación)
./scripts/deploy.sh

# Desarrollo local (sin git pull)
./scripts/deploy.sh --skip-git

# Ver todas las opciones
./scripts/deploy.sh --help
```

**¿Qué hace automáticamente?**
1. ✅ Actualiza código desde git
2. ✅ Construye y levanta contenedores Docker
3. ✅ Ejecuta migraciones de base de datos (seed data incluido)
4. ✅ Verifica que todo esté funcionando

**Servicios disponibles después del deployment:**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432
- Adminer (DB Admin): http://localhost:8088

### Deployment Manual (sin script)

Si prefieres hacerlo paso a paso:

```bash
# 1. Levantar servicios
docker compose up -d --build

# 2. Ejecutar migraciones
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:hackathon2025@postgres:5432/hackathon?sslmode=disable' up

# 3. Verificar estado
docker compose ps
curl http://localhost:8080/api/health
```

---

## 📊 Datos Precargados

El proyecto incluye migraciones de seed con datos educativos reales:

- **1 Curso**: Primero Medio
- **2 Materias**: Matemáticas y Lengua y Literatura
- **23 Objetivos de Aprendizaje** (Lengua y Literatura)
- **134 Bloom Objectives** vinculados a los OAs
- **672 Preguntas** del banco de preguntas

Estos datos se cargan automáticamente al ejecutar las migraciones.

---

## 🛠️ Stack Tecnológico

- **Backend**: Go 1.23 + Chi Router + GORM
- **Frontend**: Svelte 5 + Vite + Tailwind CSS
- **Base de Datos**: PostgreSQL 16
- **Infraestructura**: Docker + Docker Compose
- **Migraciones**: golang-migrate

---

## 📁 Estructura del Proyecto

```
lumera_app/
├── backend/
│   ├── cmd/main.go              # Entry point del backend
│   ├── internal/
│   │   ├── db/                  # Conexión PostgreSQL
│   │   ├── handlers/            # HTTP handlers
│   │   ├── models/              # GORM models
│   │   └── middleware/          # Auth middleware
│   └── migrations/              # Migraciones SQL (incluyendo seed data)
├── frontend/
│   ├── src/
│   │   ├── App.svelte           # Componente principal
│   │   └── components/          # Componentes Svelte
│   └── vite.config.js           # Configuración Vite
├── scripts/
│   ├── deploy.sh                # Script de deployment automatizado
│   └── scaffold-crud.sh         # Generador de CRUD
├── docker-compose.yml           # Orquestación de servicios
└── CLAUDE.md                    # Documentación completa para desarrollo
```

---

## 🔧 Comandos Útiles

### Docker Compose

```bash
# Ver logs de todos los servicios
docker compose logs -f

# Ver logs de un servicio específico
docker compose logs -f backend
docker compose logs -f frontend

# Reiniciar un servicio
docker compose restart backend

# Detener todo
docker compose down

# Ver estado de servicios
docker compose ps
```

### Base de Datos

```bash
# Conectar a PostgreSQL
docker compose exec postgres psql -U admin -d hackathon

# Ejecutar query directo
docker compose exec postgres psql -U admin -d hackathon -c "SELECT COUNT(*) FROM users"

# Ver versión de migraciones
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:hackathon2025@postgres:5432/hackathon?sslmode=disable' version
```

### Migraciones

```bash
# Aplicar migraciones pendientes
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:hackathon2025@postgres:5432/hackathon?sslmode=disable' up

# Rollback última migración
docker compose exec backend migrate \
  -path=/app/migrations \
  -database='postgres://admin:hackathon2025@postgres:5432/hackathon?sslmode=disable' down 1

# Crear nueva migración
make -C backend migrate-create name=nombre_de_migracion
```

---

## 🎯 Características Principales

### Sistema de Diagnóstico Adaptativo
- Evaluación inicial para determinar nivel del estudiante
- Algoritmo adaptativo que ajusta dificultad en tiempo real
- Generación de perfil de aprendizaje basado en taxonomía de Bloom

### Banco de Preguntas
- 672 preguntas categorizadas por OA y nivel de Bloom
- Múltiples tipos: selección múltiple, verdadero/falso, relacionar conceptos
- Validación automática de respuestas
- Tracking de uso y dificultad

### Sistema de Gamificación
- Monedas, experiencia y nivel del estudiante
- Sistema de logros desbloqueables
- Leaderboard social
- Progreso visual por materia y OA

### Personalización de Avatar
- Catálogo de items (rostros, accesorios, fondos)
- Sistema de compra con monedas ganadas
- Inventario y equipamiento
- Desbloqueo por logros especiales

---

## 📚 Documentación Completa

Para documentación detallada de desarrollo, arquitectura y mejores prácticas:

**👉 Ver [CLAUDE.md](./CLAUDE.md)**

Incluye:
- Guías de desarrollo backend y frontend
- Sistema de migraciones
- Arquitectura de base de datos
- Generador CRUD automático
- Mejores prácticas para el hackathon
- Troubleshooting común

---

## 🔐 Credenciales por Defecto

**PostgreSQL:**
- Usuario: `admin`
- Password: `hackathon2025`
- Base de datos: `hackathon`

**Adminer (opcional):**
- URL: http://localhost:8088
- System: PostgreSQL
- Server: postgres
- Username: admin
- Password: hackathon2025
- Database: hackathon

---

## 🧪 Testing de Endpoints

```bash
# Health check
curl http://localhost:8080/api/health

# Registro de usuario
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@lumera.com","password":"password123","nombre":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@lumera.com","password":"password123"}'

# Ver OAs disponibles
curl http://localhost:8080/api/objetivos-aprendizaje

# Ver preguntas (con filtros opcionales)
curl "http://localhost:8080/api/questions?materia_id=2&bloom_level=3"
```

---

## 🚢 Deployment en Producción

### Servidor VPS/Cloud

1. Clonar repositorio:
   ```bash
   git clone [repo-url]
   cd lumera_app
   ```

2. Configurar variables de entorno:
   ```bash
   cp .env.example .env
   # Editar .env con credenciales de producción
   ```

3. Ejecutar deployment:
   ```bash
   ./scripts/deploy.sh
   ```

4. Configurar reverse proxy (nginx/caddy) si es necesario

### Variables de Entorno Importantes

```env
# Database
DB_USER=admin
DB_PASSWORD=[cambiar-en-producción]
DB_NAME=hackathon
DB_HOST=postgres
DB_PORT=5432

# Backend
PORT=8080
JWT_SECRET=[generar-secret-seguro]

# CORS (ajustar para producción)
ALLOWED_ORIGINS=http://localhost:5173,https://tu-dominio.com
```

---

## 📝 Checklist Pre-Hackathon

- ✅ Proyecto configurado y corriendo
- ✅ Migraciones de seed ejecutadas
- ✅ Sistema de autenticación funcionando
- ✅ Banco de preguntas cargado
- ✅ Frontend conectado al backend
- ✅ Script de deployment automatizado
- ✅ Documentación completa en CLAUDE.md

---

## 🎉 Hackathon Submission Info

**Submission Deadline:** 23rd Nov, 9:00 AM, Chile time

**Track:** ✨ consumer AI

**Estado del proyecto:**
- ✅ Nombre y descripción en platanus-hack-project.json
- ✅ Logo 1000x1000 PNG (max 500kb)
- ✅ README conciso y directo al punto

---

## 📞 Soporte

Para problemas durante el desarrollo:

1. Revisar [CLAUDE.md](./CLAUDE.md) - Sección de troubleshooting
2. Ver logs: `docker compose logs -f backend`
3. Verificar salud: `curl http://localhost:8080/api/health`
4. Resetear todo: `docker compose down -v && ./scripts/deploy.sh`

---

**¡Buena suerte en el hackathon! 🚀**
