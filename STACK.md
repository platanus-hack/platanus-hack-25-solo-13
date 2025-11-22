# Stack Técnico - Platanus Hack 25

## Resumen Ejecutivo

Stack optimizado para hackathon de 33.5 horas. Prioridad: **velocidad de desarrollo** sobre todo.

## Arquitectura

```
Backend:  Go 1.23
Frontend: Svelte 5 + Vite
Database: PostgreSQL 16
Proxy:    Caddy 2.10
Deploy:   Docker Compose + GitHub Actions → VPS
```

## Backend: Go

**Versión instalada:** Go 1.23.5 (`/usr/local/go/bin/go`)

**Librerías core:**
- `github.com/go-chi/chi/v5` - Router HTTP
- `gorm.io/gorm` + `gorm.io/driver/postgres` - ORM
- `github.com/go-chi/cors` - CORS middleware
- `github.com/joho/godotenv` - Variables de entorno

**Hot reload:**
- `github.com/air-verse/air` - Auto-recompilación en desarrollo

**Estructura sugerida:**
```
backend/
├── cmd/main.go              # Entry point
├── internal/
│   ├── handlers/            # HTTP handlers
│   ├── models/              # GORM models
│   └── db/                  # Database connection
├── go.mod
├── .air.toml                # Config hot reload
└── .env
```

**Por qué Go:**
- Compila en 1-2 segundos (vs 5-10 min Rust)
- stdlib HTTP excelente (no necesitas frameworks pesados)
- Hot reload con air
- Deployment trivial (binario único)

## Frontend: Svelte + Vite

**Versión Node instalada:** v24.11.1

**Dependencies:**
- `svelte@^5.17.0` - Framework
- `vite@^6.0.7` - Build tool (HMR instantáneo)
- `@sveltejs/vite-plugin-svelte@^5.0.4`
- `tailwindcss@^3.4.17` - Estilos
- `gsap@^3.12.7` - Animaciones

**Estructura sugerida:**
```
frontend/
├── src/
│   ├── App.svelte
│   ├── main.js
│   ├── app.css              # Tailwind imports
│   ├── lib/                 # Utilidades
│   └── components/
├── index.html
├── vite.config.js
├── svelte.config.js
├── tailwind.config.js
└── package.json
```

**Por qué Svelte + Vite:**
- Compila a vanilla JS (sin virtual DOM overhead)
- HMR < 1 segundo
- Sintaxis más simple que React
- Vite build time: 2-5 segundos (vs Next.js 30-60s)

## Animaciones: GSAP

**Por qué GSAP (no Framer Motion):**
- Framework-agnostic (funciona con Svelte, React, vanilla)
- Performance superior
- Timeline controls precisos
- Demos más impresionantes

**Uso básico en Svelte:**
```javascript
import { onMount } from 'svelte';
import gsap from 'gsap';

onMount(() => {
  gsap.from('.element', { duration: 0.8, y: 50, opacity: 0 });
});
```

## Database: PostgreSQL

**Versión:** PostgreSQL 16 (Alpine en Docker)

**Conexión local:**
```
Host:     localhost (postgres en Docker network)
Port:     5432
User:     admin
Password: [definir en .env]
Database: hackathon
```

**Por qué Postgres (no Supabase):**
- Zero vendor lock-in
- Control total
- Mismo setup local y producción (Docker Compose)
- Migración trivial a managed DB después del hackathon

## Reverse Proxy: Caddy

**Versión instalada:** Caddy v2.10.2 (`/usr/bin/caddy`)

**Por qué Caddy (no Nginx):**
- Configuración más simple
- HTTPS automático
- Recarga sin downtime
- Menos líneas de config que Nginx

**Uso básico (Caddyfile):**
```
:80 {
    reverse_proxy /api/* localhost:8080
    reverse_proxy /* localhost:5173
}
```

## Containerización: Docker

**Versión instalada:**
- Docker 29.0.2
- Docker Compose v2.40.3

**Servicios en docker-compose.yml:**
```yaml
services:
  postgres:   # PostgreSQL 16-alpine
  backend:    # Go app con air (hot reload)
  frontend:   # Svelte con Vite dev server
```

**Comandos esenciales:**
```bash
docker compose up           # Levantar todo
docker compose down         # Bajar todo
docker compose build        # Rebuild
docker compose logs -f      # Ver logs
```

## Deploy: GitHub Actions + VPS

**VPS actual:**
- IP: [definir]
- OS: Ubuntu 25.04
- RAM: 12GB
- Disk: 93GB disponibles

**Flujo de deploy:**
1. `git push origin main`
2. GitHub Actions ejecuta workflow
3. SSH al VPS
4. `git pull`
5. `docker compose down && docker compose build && docker compose up -d`

**Secrets necesarios en GitHub:**
- `VPS_HOST` - IP del servidor
- `VPS_USER` - Usuario SSH (ubuntu)
- `VPS_SSH_KEY` - Llave privada SSH

## Variables de Entorno

**Backend (.env):**
```bash
PORT=8080
DB_HOST=postgres
DB_USER=admin
DB_PASSWORD=cambiar_esto
DB_NAME=hackathon
DB_PORT=5432
```

**Frontend:**
- Vite usa variables prefijadas con `VITE_`
- Proxy a backend configurado en `vite.config.js`

## Comandos de Desarrollo

**Iniciar proyecto completo:**
```bash
cd /ruta/proyecto
docker compose up
```

**Backend solo (manual):**
```bash
cd backend
air  # hot reload automático
```

**Frontend solo (manual):**
```bash
cd frontend
npm run dev
```

## Comparación de Velocidad

| Métrica          | Next.js + Supabase | Go + Svelte + VPS |
|------------------|-------------------|-------------------|
| Setup inicial    | 30 min            | 15 min            |
| Hot reload       | 2-5s              | <1s               |
| Build time       | 30-60s            | 2s                |
| Deploy time      | 3-5 min           | 1-2 min           |
| Control total    | ❌                | ✅                |
| Vendor lock-in   | ⚠️ Supabase       | ✅ Zero           |

## Filosofía del Stack

1. **Go backend:** Compilación rápida, deployment simple
2. **Svelte frontend:** HMR instantáneo, menos boilerplate que React
3. **GSAP:** Animaciones impresionantes para demos
4. **Postgres en Docker:** Desarrollo = Producción
5. **Caddy:** Configuración mínima, HTTPS gratis
6. **VPS propio:** Sin sorpresas de vendor, control absoluto

## Antipatrones a Evitar

- ❌ No usar frameworks pesados (Next.js, NestJS)
- ❌ No implementar auth compleja (usar JWT simple o magic links)
- ❌ No over-engineer (microservicios, arquitecturas complejas)
- ❌ No confiar en servicios externos que puedan fallar durante demo

## Prioridades en Hackathon

1. **Velocidad de desarrollo** > Escalabilidad
2. **Demo visual** > Features complejas
3. **Stack conocido** > Tecnología nueva y cool
4. **Simplicidad** > Arquitectura "correcta"

## Servicios Instalados en VPS

- ✅ Docker 29.0.2
- ✅ Docker Compose v2.40.3
- ✅ Go 1.23.5
- ✅ Node.js v24.11.1
- ✅ Caddy v2.10.2

**PATH configurado:**
```bash
/usr/local/go/bin  # Go binaries
~/go/bin           # Go tools (air, etc)
```

## Próximos Pasos (al iniciar proyecto)

1. Crear directorio del proyecto: `/home/ubuntu/projects/[nombre-proyecto]`
2. Inicializar Go module: `go mod init github.com/platanus-hack-25/[nombre]`
3. Instalar dependencias Go
4. Inicializar frontend: `npm create vite@latest`
5. Configurar docker-compose.yml
6. Crear GitHub repo y configurar Actions
7. Primera iteración: Backend API + Frontend básico + DB connection
8. Verificar hot reload funcionando en ambos

## Contacto/Debug

- Backend logs: `docker compose logs -f backend`
- Frontend logs: `docker compose logs -f frontend`
- DB logs: `docker compose logs -f postgres`
- Caddy logs: `sudo journalctl -u caddy -f`

---

**Última actualización:** 2025-11-22
**Servidor:** VPS Ubuntu 25.04 (vps-3df5f1f7)