#!/bin/bash

# Login
LOGIN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"testquestions@lumera.com","password":"password123"}')

echo "Login response:"
echo "$LOGIN"
echo ""

TOKEN=$(echo "$LOGIN" | grep -o '"token":"[^"]*' | sed 's/"token":"//')

echo "Extracted token:"
echo "$TOKEN"
echo ""

# Test diagnostic session
echo "Testing diagnostic session with token..."
curl -v -X POST http://localhost:8080/api/diagnostic-sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"materia_id": 2}' \
  2>&1 | grep -E "(> |< |HTTP|{)"
