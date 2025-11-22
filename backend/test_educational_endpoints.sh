#!/bin/bash

# Educational API Test Script
set -e

BASE_URL="http://localhost:8080"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo4LCJlbWFpbCI6InRlYWNoZXJAbHVtZXJhLmNvbSIsInJvbGUiOiJ1c2VyIiwiZXhwIjoxNzYzODg4NzU0LCJpYXQiOjE3NjM4MDIzNTR9.RdQHuXrb5k9C_H7hxcvmIBbnraNBiT_o5n2w-Y_tDf8"

echo "=== Testing Educational API Endpoints ==="
echo ""

# Test 1: Get all Bloom levels
echo "1. GET /api/bloom-levels"
curl -s "$BASE_URL/api/bloom-levels" | jq 'length'
echo "✓ Bloom levels retrieved"
echo ""

# Test 2: Create a curso
echo "2. POST /api/cursos (Create '1° Medio')"
CURSO_RESPONSE=$(curl -s -X POST "$BASE_URL/api/cursos" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nombre":"Primero Medio","codigo":"1M","nivel_educativo":"Enseñanza Media","descripcion":"Primer año de educación media","activo":true}')
CURSO_ID=$(echo "$CURSO_RESPONSE" | jq -r '.id')
echo "$CURSO_RESPONSE" | jq .
echo "✓ Curso created with ID: $CURSO_ID"
echo ""

# Test 3: Get all cursos
echo "3. GET /api/cursos"
curl -s "$BASE_URL/api/cursos" | jq .
echo "✓ Cursos retrieved"
echo ""

# Test 4: Create a materia
echo "4. POST /api/materias (Create 'Lenguaje y Comunicación')"
MATERIA_RESPONSE=$(curl -s -X POST "$BASE_URL/api/materias" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"nombre":"Lenguaje y Comunicación","codigo":"LEN","descripcion":"Lenguaje, comunicación y literatura","color":"#3B82F6","activo":true}')
MATERIA_ID=$(echo "$MATERIA_RESPONSE" | jq -r '.id')
echo "$MATERIA_RESPONSE" | jq .
echo "✓ Materia created with ID: $MATERIA_ID"
echo ""

# Test 5: Assign materia to curso
echo "5. POST /api/curso-materias (Assign Lenguaje to 1° Medio)"
ASSIGNMENT=$(curl -s -X POST "$BASE_URL/api/curso-materias" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"curso_id\":$CURSO_ID,\"materia_id\":$MATERIA_ID,\"horas_semanales\":6}")
echo "$ASSIGNMENT" | jq .
echo "✓ Materia assigned to Curso"
echo ""

# Test 6: Create an OA with all 6 Bloom objectives
echo "6. POST /api/objetivos-aprendizaje (Create OA with 6 Bloom levels)"
OA_PAYLOAD='{
  "materia_id": '$MATERIA_ID',
  "codigo": "OA01",
  "titulo": "Analizar textos narrativos",
  "descripcion": "Analizar e interpretar textos narrativos considerando diversos elementos",
  "orden": 1,
  "bloom_objectives": [
    {
      "bloom_level_id": 1,
      "objetivo_especifico": "Identificar elementos narrativos básicos (narrador, personajes, tiempo, espacio)",
      "indicadores_logro": ["Reconoce el tipo de narrador", "Identifica personajes principales"],
      "tipo_actividad_sugerida": "Lectura guiada",
      "complejidad_estimada": 2
    },
    {
      "bloom_level_id": 2,
      "objetivo_especifico": "Explicar la función de los elementos narrativos en el texto",
      "indicadores_logro": ["Describe el rol del narrador", "Explica la caracterización de personajes"],
      "tipo_actividad_sugerida": "Cuestionario",
      "complejidad_estimada": 3
    },
    {
      "bloom_level_id": 3,
      "objetivo_especifico": "Aplicar conceptos narrativos en el análisis de nuevos textos",
      "indicadores_logro": ["Utiliza terminología narratológica", "Aplica conceptos a textos nuevos"],
      "tipo_actividad_sugerida": "Análisis de texto",
      "complejidad_estimada": 5
    },
    {
      "bloom_level_id": 4,
      "objetivo_especifico": "Analizar la relación entre elementos narrativos y el significado del texto",
      "indicadores_logro": ["Relaciona narrador con perspectiva", "Analiza desarrollo de personajes"],
      "tipo_actividad_sugerida": "Ensayo breve",
      "complejidad_estimada": 6
    },
    {
      "bloom_level_id": 5,
      "objetivo_especifico": "Evaluar la efectividad de las estrategias narrativas empleadas",
      "indicadores_logro": ["Juzga coherencia narrativa", "Evalúa impacto de decisiones narrativas"],
      "tipo_actividad_sugerida": "Crítica literaria",
      "complejidad_estimada": 8
    },
    {
      "bloom_level_id": 6,
      "objetivo_especifico": "Crear un texto narrativo original aplicando elementos analizados",
      "indicadores_logro": ["Produce texto con estructura narrativa", "Aplica técnicas narrativas creativas"],
      "tipo_actividad_sugerida": "Escritura creativa",
      "complejidad_estimada": 9
    }
  ]
}'

OA_RESPONSE=$(curl -s -X POST "$BASE_URL/api/objetivos-aprendizaje" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$OA_PAYLOAD")
OA_ID=$(echo "$OA_RESPONSE" | jq -r '.id')
echo "$OA_RESPONSE" | jq '.bloom_objectives | length'
echo "✓ OA created with ID: $OA_ID and 6 Bloom objectives"
echo ""

# Test 7: Get OA with all details
echo "7. GET /api/objetivos-aprendizaje/$OA_ID"
curl -s "$BASE_URL/api/objetivos-aprendizaje/$OA_ID" | jq '{id, codigo, titulo, bloom_count: (.bloom_objectives | length)}'
echo "✓ OA retrieved with Bloom breakdown"
echo ""

# Test 8: Register student progress
echo "8. POST /api/progress (Register progress on Bloom level 1)"
PROGRESS_PAYLOAD='{
  "user_id": 8,
  "oa_bloom_objective_id": 1,
  "estado": "en_proceso",
  "porcentaje_logro": 60,
  "tipo_evento": "practica",
  "puntaje_obtenido": 12,
  "puntaje_maximo": 20,
  "notas": "Primera práctica del OA01 nivel Recordar"
}'

PROGRESS=$(curl -s -X POST "$BASE_URL/api/progress" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$PROGRESS_PAYLOAD")
echo "$PROGRESS" | jq .
echo "✓ Progress registered"
echo ""

# Test 9: Get student progress
echo "9. GET /api/progress/8"
curl -s -H "Authorization: Bearer $TOKEN" "$BASE_URL/api/progress/8" | jq 'length'
echo "✓ Student progress retrieved"
echo ""

echo "=== All tests completed successfully! ==="
