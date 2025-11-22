#!/bin/bash

# Test Question Bank and Diagnostic System Endpoints

echo "======================================"
echo "Testing Question Bank and Diagnostic System"
echo "======================================"

# First, login to get a token
echo -e "\n1. Login to get auth token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testquestions@lumera.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
  echo "❌ Failed to get auth token. Response: $LOGIN_RESPONSE"
  exit 1
fi

echo "✓ Got auth token: ${TOKEN:0:20}..."

# Test 1: Get question types (public)
echo -e "\n2. GET /api/question-types (should return 9 types)..."
TYPES_RESPONSE=$(curl -s http://localhost:8080/api/question-types)
TYPES_COUNT=$(echo $TYPES_RESPONSE | grep -o '"tipo"' | wc -l)
echo "   Found $TYPES_COUNT question types"
echo "   Response: $TYPES_RESPONSE" | head -c 200
echo "..."

# Test 2: Get all questions (should be empty initially)
echo -e "\n3. GET /api/questions (should be empty initially)..."
QUESTIONS_RESPONSE=$(curl -s http://localhost:8080/api/questions)
echo "   Response: $QUESTIONS_RESPONSE"

# Test 3: Create a multiple choice question (requires auth)
echo -e "\n4. POST /api/questions (create multiple choice question)..."

# First, get an OA Bloom objective ID
echo "   Getting an OA Bloom objective..."
OA_BLOOM=$(curl -s "http://localhost:8080/api/objetivos-aprendizaje" | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
echo "   Using OA Bloom ID: $OA_BLOOM"

CREATE_QUESTION_RESPONSE=$(curl -s -X POST http://localhost:8080/api/questions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"oa_bloom_objective_id\": $OA_BLOOM,
    \"tipo\": \"multiple_choice\",
    \"tipo_uso\": \"diagnostico\",
    \"question_data\": {
      \"pregunta\": \"Cual es la capital de Chile?\",
      \"opciones\": {
        \"A\": \"Valparaiso\",
        \"B\": \"Santiago\",
        \"C\": \"Concepcion\",
        \"D\": \"Antofagasta\"
      },
      \"explicacion\": \"Santiago es la capital de Chile desde 1541\"
    },
    \"validation_data\": {
      \"respuesta_correcta\": \"B\"
    },
    \"dificultad_relativa\": 2,
    \"tags\": [\"geografia\", \"chile\"]
  }")

QUESTION_ID=$(echo $CREATE_QUESTION_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')

if [ -z "$QUESTION_ID" ]; then
  echo "   ❌ Failed to create question. Response: $CREATE_QUESTION_RESPONSE"
  exit 1
fi

echo "   ✓ Created question with ID: $QUESTION_ID"

# Test 4: Get question by ID (should NOT include validation_data)
echo -e "\n5. GET /api/questions/$QUESTION_ID (should NOT show validation_data)..."
QUESTION_GET_RESPONSE=$(curl -s http://localhost:8080/api/questions/$QUESTION_ID)
echo "   Response: $QUESTION_GET_RESPONSE"

if echo "$QUESTION_GET_RESPONSE" | grep -q "validation_data"; then
  echo "   ❌ ERROR: validation_data is exposed (security issue!)"
else
  echo "   ✓ validation_data is properly hidden"
fi

# Test 5: Validate a correct answer
echo -e "\n6. POST /api/questions/$QUESTION_ID/validate (correct answer)..."
VALIDATE_CORRECT=$(curl -s -X POST http://localhost:8080/api/questions/$QUESTION_ID/validate \
  -H "Content-Type: application/json" \
  -d '{
    "user_answer": {
      "selected": "B"
    }
  }')
echo "   Response: $VALIDATE_CORRECT"

if echo "$VALIDATE_CORRECT" | grep -q '"is_correct":true'; then
  echo "   ✓ Correct answer validated successfully"
else
  echo "   ❌ Expected is_correct=true"
fi

# Test 6: Validate an incorrect answer
echo -e "\n7. POST /api/questions/$QUESTION_ID/validate (incorrect answer)..."
VALIDATE_WRONG=$(curl -s -X POST http://localhost:8080/api/questions/$QUESTION_ID/validate \
  -H "Content-Type: application/json" \
  -d '{
    "user_answer": {
      "selected": "A"
    }
  }')
echo "   Response: $VALIDATE_WRONG"

if echo "$VALIDATE_WRONG" | grep -q '"is_correct":false'; then
  echo "   ✓ Incorrect answer detected successfully"
  if echo "$VALIDATE_WRONG" | grep -q "correct_answer"; then
    echo "   ✓ Correct answer shown after failure (as expected)"
  fi
else
  echo "   ❌ Expected is_correct=false"
fi

# Test 7: Start diagnostic session
echo -e "\n8. POST /api/diagnostic-sessions (start diagnostic)..."

# Get a materia ID
MATERIA_ID=$(curl -s "http://localhost:8080/api/materias" | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
echo "   Using Materia ID: $MATERIA_ID"

DIAGNOSTIC_START=$(curl -s -X POST http://localhost:8080/api/diagnostic-sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"materia_id\": $MATERIA_ID
  }")

SESSION_ID=$(echo $DIAGNOSTIC_START | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')

if [ -z "$SESSION_ID" ]; then
  echo "   ❌ Failed to start diagnostic session. Response: $DIAGNOSTIC_START"
  exit 1
fi

echo "   ✓ Started diagnostic session with ID: $SESSION_ID"
echo "   Response: $DIAGNOSTIC_START"

# Test 8: Submit an answer to the diagnostic session
echo -e "\n9. POST /api/diagnostic-sessions/$SESSION_ID/answer (submit answer)..."
SUBMIT_ANSWER=$(curl -s -X POST http://localhost:8080/api/diagnostic-sessions/$SESSION_ID/answer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"question_id\": $QUESTION_ID,
    \"user_answer\": {
      \"selected\": \"B\"
    },
    \"tiempo_segundos\": 15
  }")

ANSWER_ID=$(echo $SUBMIT_ANSWER | grep -o '"answer_id":[0-9]*' | sed 's/"answer_id"://')

if [ -z "$ANSWER_ID" ]; then
  echo "   ❌ Failed to submit answer. Response: $SUBMIT_ANSWER"
else
  echo "   ✓ Submitted answer with ID: $ANSWER_ID"
  echo "   Response: $SUBMIT_ANSWER"
fi

# Test 9: Get session progress
echo -e "\n10. GET /api/diagnostic-sessions/$SESSION_ID (get progress)..."
SESSION_PROGRESS=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/diagnostic-sessions/$SESSION_ID)
echo "   Response: $SESSION_PROGRESS" | head -c 300
echo "..."

# Test 10: Complete diagnostic session
echo -e "\n11. POST /api/diagnostic-sessions/$SESSION_ID/complete..."
COMPLETE_DIAGNOSTIC=$(curl -s -X POST http://localhost:8080/api/diagnostic-sessions/$SESSION_ID/complete \
  -H "Authorization: Bearer $TOKEN")
echo "   Response: $COMPLETE_DIAGNOSTIC"

if echo "$COMPLETE_DIAGNOSTIC" | grep -q "completado"; then
  echo "   ✓ Diagnostic session completed successfully"
else
  echo "   ⚠ Diagnostic completed but response format unexpected"
fi

# Test 11: Get diagnostic results
echo -e "\n12. GET /api/diagnostic-sessions/$SESSION_ID/results..."
DIAGNOSTIC_RESULTS=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/diagnostic-sessions/$SESSION_ID/results)
echo "   Response: $DIAGNOSTIC_RESULTS"

echo -e "\n======================================"
echo "✅ All endpoint tests completed!"
echo "======================================"
echo ""
echo "Summary:"
echo "  - Question Types: $TYPES_COUNT types registered"
echo "  - Created Question ID: $QUESTION_ID"
echo "  - Diagnostic Session ID: $SESSION_ID"
echo "  - Answer ID: $ANSWER_ID"
echo ""
