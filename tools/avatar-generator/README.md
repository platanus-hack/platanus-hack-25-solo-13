# Avatar Generator con DALL-E

Generador de avatares para el sistema de gamificación de Lumera App usando OpenAI DALL-E-3.

## Características

- ✅ Genera 42 avatares temáticos con DALL-E-3
- ✅ 3 categorías: inicio (1⭐), logro (3-5⭐), compra (3-5⭐)
- ✅ Prompts optimizados por tier con lineamientos de estilo
- ✅ Descarga y almacenamiento automático de imágenes
- ✅ Inserción directa en base de datos con triggers de unlock
- ✅ Retry logic con 3 intentos por imagen
- ✅ Logging detallado y estadísticas finales

## Estructura del Proyecto

```
avatar-generator/
├── main.go              # Orchestration principal
├── generator/
│   ├── models.go       # Structs de Avatar y Stats
│   ├── db.go           # Conexión DB e inserts
│   └── openai_client.go # Cliente DALL-E-3
├── input/
│   └── avatar_config.json  # Configuración de 42 avatares
├── output/
│   ├── images/         # Imágenes generadas (.png)
│   └── failed_*.json   # Avatares fallidos (si hay)
├── .env                # Configuración
└── go.mod
```

## Distribución de Avatares

### Categoría "inicio" (1⭐) - 10 avatares
Mascotas cartoon genéricas disponibles desde el primer uso:
- Perrito Aventurero, Gatito Curioso, Conejito Estudioso
- Búho Sabio, Zorrito Lector, Osito Pensador
- Panda Alegre, Mapache Explorador, Pingüino Curioso, Ardilla Energética

### Categoría "logro" (3-5⭐) - 20 avatares

**Tier 3 (rare) - 8 avatares**
Relacionados a OAs específicos:
- Dragón Poeta (OA poesía), Fénix Narrador (OA narrativa)
- Unicornio Gramático, Lobo Argumentador, Lince Investigador
- Ciervo Orador, Tigre Dialogante, Delfín Comprensivo

**Tier 4 (epic) - 8 avatares**
Relacionados a niveles Bloom y rachas:
- León Analítico (Bloom nivel 4), Águila Evaluadora (Bloom nivel 5)
- Pantera Creativa (Bloom nivel 6), Elefante Recordador (Bloom nivel 1)
- Halcón Aplicador (Bloom nivel 3), Koala Comprensivo (Bloom nivel 2)
- Jaguar de Racha (7 días), Búho Constante (30 días)

**Tier 5 (legendary) - 4 avatares**
Logros excepcionales:
- Grifo Maestro (80% de Lenguaje), Dragón Ancestral (nivel 20)
- Fénix Eterno (2 materias completas), Quetzal Supremo (6 Bloom)

### Categoría "compra" (3-5⭐) - 12 avatares

**Tier 3 (rare) - 250-300 coins - 5 avatares**
- Caballero Medieval, Princesa Cuentista, Mago Místico
- Ninja Silencioso, Vikingo Épico

**Tier 4 (epic) - 500-650 coins - 4 avatares**
- Detective Sherlock, Capitán Pirata, Astronauta Soñador, Samurái Honorable

**Tier 5 (legendary) - 1200-1800 coins - 3 avatares**
- Emperador Galáctico, Guardián del Tiempo, Arquitecto de Mundos

## Configuración

Copiar `.env.example` a `.env` y configurar:

```bash
# OpenAI
OPENAI_API_KEY=your-api-key-here
OPENAI_MODEL=dall-e-3
OPENAI_TIMEOUT_SECONDS=60
OPENAI_MAX_RETRIES=3

# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=hackathon2025
DB_NAME=hackathon
DB_SSLMODE=disable
```

## Ejecución

### Opción 1: Docker (Recomendado)

```bash
cd /Users/johnny/git/lumera_app

docker run --rm \
  --network lumera_app_default \
  -v $(pwd)/avatar-generator:/app \
  -w /app \
  -e DB_HOST=postgres \
  --env-file avatar-generator/.env \
  golang:1.23-alpine \
  sh -c "go run main.go"
```

### Opción 2: Local (requiere Go 1.23+)

```bash
cd avatar-generator
go mod download
go run main.go
```

## Personalización

Para agregar/modificar avatares, editar `input/avatar_config.json`:

```json
{
  "avatares": [
    {
      "nombre": "Nuevo Avatar",
      "descripcion": "Descripción del avatar",
      "categoria": "inicio|logro|compra",
      "tier": 1-5,
      "rarity": "common|rare|epic|legendary",
      "precio_puntos": 0,
      "is_default": true|false,
      "unlock_trigger": {
        "trigger_type": "default|oa_complete|bloom_mastery|streak|level|coins|multiple",
        "trigger_key": "oa_35",
        "display_text": "Texto mostrado al usuario",
        "extra_data": { /* datos adicionales */ }
      },
      "dalle_prompt": "Prompt base para DALL-E"
    }
  ]
}
```

## Lineamientos de Prompts por Tier

El generador aplica automáticamente estilos según el tier:

- **Tier 1**: Simple design, basic shapes, primary colors, minimalist
- **Tier 2**: More detailed, colorful accessories, expressive features
- **Tier 3**: Detailed design, thematic accessories, rich palette, dynamic pose
- **Tier 4**: Highly detailed, magical elements, subtle glow, heroic pose
- **Tier 5**: Extremely detailed, prominent magical effects, epic pose, legendary appearance

**Estilo base (todos)**: Friendly cartoon mascot, educational theme, clean design, white background, vibrant colors, appealing to middle school students (14-18 years)

## Estadísticas de Generación

Al finalizar, el generador muestra:
- Duración total
- Avatares generados exitosamente vs fallidos
- Distribución por categoría (inicio, logro, compra)
- Distribución por tier (1⭐ a 5⭐)

## Integración con Backend

Los avatares se insertan automáticamente en:
- `customization_items` (catálogo global)
- `item_unlock_triggers` (eventos de desbloqueo)

Las imágenes se guardan en: `avatar-generator/output/images/`

**Para producción**: Mover imágenes a `backend/static/avatares/` o subirlas a CDN.

## API de Customización (Ya implementada)

El backend ya tiene endpoints para trabajar con los avatares:

```
GET  /api/customization/catalog      # Catálogo completo con ownership
GET  /api/customization/inventory    # Items desbloqueados del usuario
GET  /api/customization/equipment    # Avatar/frame equipado
POST /api/customization/equip        # Equipar avatar
POST /api/customization/purchase     # Comprar con coins
GET  /api/customization/notifications # Notificaciones de unlocks
```

## Costos Estimados

**DALL-E-3 standard 1024x1024**: ~$0.040 USD por imagen

- 42 avatares × $0.040 = **~$1.68 USD** (generación completa)
- Con fallos y retries: **~$2-3 USD** máximo

## Troubleshooting

**Error: "no images generated"**
- Verificar que OPENAI_API_KEY sea válida
- Verificar cuota de OpenAI API

**Error: "failed to insert avatar"**
- Verificar que la migración 000016 esté aplicada
- Verificar conexión a la base de datos

**Imágenes no se descargan**
- Verificar permisos de escritura en `output/images/`
- Verificar conectividad de red

## Próximos Pasos

1. Ejecutar generador con los 42 avatares configurados
2. Mover imágenes a `backend/static/avatares/` (opcional)
3. Configurar backend para servir imágenes estáticas
4. Probar API de customización desde el frontend
5. Implementar UI de selección de avatares en el frontend

## Notas de Desarrollo

- **Patrón similar a question-generator**: Mismo approach con retry logic
- **Sistema de gamificación ya existe**: Solo agregamos contenido
- **Tier numérico (1-5)**: Mapea a rarities pero más flexible
- **Unlock triggers**: Sistema event-driven ya implementado
- **Modularidad**: Fácil agregar más avatares editando el JSON
