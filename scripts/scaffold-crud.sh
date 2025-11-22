#!/bin/bash
# Script para generar CRUD completo de una entidad
# Uso: ./scripts/scaffold-crud.sh EntityName field1:type field2:type

set -e

if [ "$#" -lt 2 ]; then
    echo "❌ Uso: $0 EntityName field1:type field2:type ..."
    echo ""
    echo "Ejemplo: $0 Product name:string price:decimal stock:int"
    echo ""
    echo "Tipos soportados: string, text, int, decimal, bool, timestamp"
    exit 1
fi

ENTITY=$1
shift
FIELDS=("$@")

# Convertir a minúsculas para tabla
TABLE=$(echo "$ENTITY" | tr '[:upper:]' '[:lower:]')s

echo "🚀 Generando CRUD para: $ENTITY"
echo "📋 Tabla: $TABLE"
echo "📊 Campos: ${FIELDS[@]}"
echo ""

# Crear migración
echo "1️⃣  Creando migración..."
MIGRATION_NAME="create_${TABLE}_table"

# Generar SQL para campos
SQL_FIELDS="id SERIAL PRIMARY KEY,"
GO_FIELDS=""

for field_def in "${FIELDS[@]}"; do
    IFS=':' read -r field type <<< "$field_def"
    field_lower=$(echo "$field" | tr '[:upper:]' '[:lower:]')
    field_camel=$(echo "$field" | sed 's/_\([a-z]\)/\U\1/g' | sed 's/^\([a-z]\)/\U\1/')

    case $type in
        string)
            SQL_FIELDS="$SQL_FIELDS\n    $field_lower VARCHAR(255) NOT NULL,"
            GO_FIELDS="$GO_FIELDS\n    $field_camel string \`json:\"$field_lower\" gorm:\"not null\"\`"
            ;;
        text)
            SQL_FIELDS="$SQL_FIELDS\n    $field_lower TEXT,"
            GO_FIELDS="$GO_FIELDS\n    $field_camel string \`json:\"$field_lower\"\`"
            ;;
        int)
            SQL_FIELDS="$SQL_FIELDS\n    $field_lower INTEGER NOT NULL DEFAULT 0,"
            GO_FIELDS="$GO_FIELDS\n    $field_camel int \`json:\"$field_lower\" gorm:\"default:0\"\`"
            ;;
        decimal)
            SQL_FIELDS="$SQL_FIELDS\n    $field_lower DECIMAL(10,2) NOT NULL,"
            GO_FIELDS="$GO_FIELDS\n    $field_camel float64 \`json:\"$field_lower\" gorm:\"type:decimal(10,2)\"\`"
            ;;
        bool)
            SQL_FIELDS="$SQL_FIELDS\n    $field_lower BOOLEAN NOT NULL DEFAULT false,"
            GO_FIELDS="$GO_FIELDS\n    $field_camel bool \`json:\"$field_lower\" gorm:\"default:false\"\`"
            ;;
        timestamp)
            SQL_FIELDS="$SQL_FIELDS\n    $field_lower TIMESTAMP,"
            GO_FIELDS="$GO_FIELDS\n    $field_camel time.Time \`json:\"$field_lower\"\`"
            ;;
        *)
            echo "⚠️  Tipo desconocido: $type, usando VARCHAR"
            SQL_FIELDS="$SQL_FIELDS\n    $field_lower VARCHAR(255),"
            GO_FIELDS="$GO_FIELDS\n    $field_camel string \`json:\"$field_lower\"\`"
            ;;
    esac
done

SQL_FIELDS="$SQL_FIELDS\n    created_at TIMESTAMP NOT NULL DEFAULT NOW(),\n    updated_at TIMESTAMP NOT NULL DEFAULT NOW()"
GO_FIELDS="$GO_FIELDS\n    CreatedAt time.Time \`json:\"created_at\"\`\n    UpdatedAt time.Time \`json:\"updated_at\"\`"

# Crear archivos de migración
docker compose exec backend sh -c "migrate create -ext sql -dir /app/migrations -seq $MIGRATION_NAME"

# Esperar a que se creen los archivos
sleep 1

# Encontrar los archivos creados
UP_FILE=$(ls backend/migrations/*_${MIGRATION_NAME}.up.sql | tail -1)
DOWN_FILE=$(ls backend/migrations/*_${MIGRATION_NAME}.down.sql | tail -1)

# Escribir migración UP
cat > "$UP_FILE" << EOF
CREATE TABLE IF NOT EXISTS $TABLE (
$(echo -e "$SQL_FIELDS")
);

CREATE INDEX idx_${TABLE}_created_at ON $TABLE(created_at);
EOF

# Escribir migración DOWN
cat > "$DOWN_FILE" << EOF
DROP INDEX IF EXISTS idx_${TABLE}_created_at;
DROP TABLE IF EXISTS $TABLE;
EOF

echo "✅ Migración creada: $UP_FILE"

# Crear modelo GORM
echo "2️⃣  Creando modelo GORM..."
mkdir -p backend/internal/models
mkdir -p backend/internal/handlers
MODEL_FILE="backend/internal/models/${TABLE%s}.go"
cat > "$MODEL_FILE" << EOF
package models

import "time"

type $ENTITY struct {
    ID uint \`json:"id" gorm:"primaryKey"\`
$(echo -e "$GO_FIELDS")
}
EOF

echo "✅ Modelo creado: $MODEL_FILE"

# Crear handler
echo "3️⃣  Creando handler..."
HANDLER_FILE="backend/internal/handlers/${TABLE}.go"
cat > "$HANDLER_FILE" << EOF
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/platanus-hack-25/lumera_app/internal/db"
    "github.com/platanus-hack-25/lumera_app/internal/models"
)

// Get all ${TABLE}
func Get${ENTITY}s(w http.ResponseWriter, r *http.Request) {
    var items []models.$ENTITY
    db.DB.Find(&items)
    json.NewEncoder(w).Encode(items)
}

// Get single ${TABLE%s}
func Get${ENTITY}(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var item models.$ENTITY

    if err := db.DB.First(&item, id).Error; err != nil {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(item)
}

// Create ${TABLE%s}
func Create${ENTITY}(w http.ResponseWriter, r *http.Request) {
    var item models.$ENTITY

    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    db.DB.Create(&item)
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(item)
}

// Update ${TABLE%s}
func Update${ENTITY}(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var item models.$ENTITY

    if err := db.DB.First(&item, id).Error; err != nil {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }

    if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    db.DB.Save(&item)
    json.NewEncoder(w).Encode(item)
}

// Delete ${TABLE%s}
func Delete${ENTITY}(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    if err := db.DB.Delete(&models.$ENTITY{}, id).Error; err != nil {
        http.Error(w, "Not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
EOF

echo "✅ Handler creado: $HANDLER_FILE"

# Instrucciones finales
echo ""
echo "🎉 CRUD generado exitosamente!"
echo ""
echo "📋 Próximos pasos:"
echo ""
echo "1. Aplicar migración:"
echo "   make migrate-up"
echo ""
echo "2. Agregar rutas en backend/cmd/main.go:"
echo "   r.Route(\"/api/${TABLE}\", func(r chi.Router) {"
echo "       r.Get(\"/\", handlers.Get${ENTITY}s)"
echo "       r.Post(\"/\", handlers.Create${ENTITY})"
echo "       r.Get(\"/{id}\", handlers.Get${ENTITY})"
echo "       r.Put(\"/{id}\", handlers.Update${ENTITY})"
echo "       r.Delete(\"/{id}\", handlers.Delete${ENTITY})"
echo "   })"
echo ""
echo "3. Reiniciar backend:"
echo "   docker compose restart backend"
echo ""
echo "4. Probar endpoints:"
echo "   GET    http://localhost:8080/api/${TABLE}"
echo "   POST   http://localhost:8080/api/${TABLE}"
echo "   GET    http://localhost:8080/api/${TABLE}/1"
echo "   PUT    http://localhost:8080/api/${TABLE}/1"
echo "   DELETE http://localhost:8080/api/${TABLE}/1"
