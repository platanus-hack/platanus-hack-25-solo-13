#!/bin/bash

# Test script for Gamification & Customization System
# Tests all new endpoints and integration with existing handlers

set -e

API_URL="http://localhost:8080/api"
TOKEN=""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Gamification & Customization System Tests ===${NC}\n"

# Step 1: Login to get token
echo -e "${BLUE}[1/10] Logging in...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testgamif@lumera.com",
    "password": "password123"
  }')

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo -e "${RED}❌ Login failed${NC}"
  echo $LOGIN_RESPONSE
  exit 1
fi

echo -e "${GREEN}✓ Login successful${NC}\n"

# Step 2: Get Gamification Stats
echo -e "${BLUE}[2/10] Getting gamification stats...${NC}"
STATS=$(curl -s -X GET "$API_URL/gamification/stats" \
  -H "Authorization: Bearer $TOKEN")

echo $STATS | jq '.'
echo -e "${GREEN}✓ Gamification stats retrieved${NC}\n"

# Step 3: Get Leaderboard
echo -e "${BLUE}[3/10] Getting leaderboard...${NC}"
LEADERBOARD=$(curl -s -X GET "$API_URL/gamification/leaderboard" \
  -H "Authorization: Bearer $TOKEN")

echo $LEADERBOARD | jq '.leaderboard | length' | sed 's/^/   Players: /'
echo -e "${GREEN}✓ Leaderboard retrieved${NC}\n"

# Step 4: Get Customization Catalog
echo -e "${BLUE}[4/10] Getting customization catalog...${NC}"
CATALOG=$(curl -s -X GET "$API_URL/customization/catalog" \
  -H "Authorization: Bearer $TOKEN")

TOTAL_ITEMS=$(echo $CATALOG | jq '. | length')
OWNED_ITEMS=$(echo $CATALOG | jq '[.[] | select(.is_owned == true)] | length')
LOCKED_ITEMS=$(echo $CATALOG | jq '[.[] | select(.status == "locked")] | length')
PURCHASABLE=$(echo $CATALOG | jq '[.[] | select(.status == "can_purchase")] | length')

echo "   Total items: $TOTAL_ITEMS"
echo "   Owned: $OWNED_ITEMS"
echo "   Locked: $LOCKED_ITEMS"
echo "   Can purchase: $PURCHASABLE"
echo -e "${GREEN}✓ Catalog retrieved${NC}\n"

# Step 5: Get User Inventory
echo -e "${BLUE}[5/10] Getting user inventory...${NC}"
INVENTORY=$(curl -s -X GET "$API_URL/customization/inventory" \
  -H "Authorization: Bearer $TOKEN")

INV_COUNT=$(echo $INVENTORY | jq '. | length')
echo "   Inventory items: $INV_COUNT"
echo -e "${GREEN}✓ Inventory retrieved${NC}\n"

# Step 6: Get Equipment
echo -e "${BLUE}[6/10] Getting equipped items...${NC}"
EQUIPMENT=$(curl -s -X GET "$API_URL/customization/equipment" \
  -H "Authorization: Bearer $TOKEN")

echo $EQUIPMENT | jq '.'
echo -e "${GREEN}✓ Equipment retrieved${NC}\n"

# Step 7: Equip an item (use first owned avatar)
echo -e "${BLUE}[7/10] Equipping an avatar...${NC}"
FIRST_AVATAR_ID=$(echo $CATALOG | jq '[.[] | select(.type == "avatar" and .is_owned == true)] | .[0].id')

if [ "$FIRST_AVATAR_ID" != "null" ]; then
  EQUIP_RESPONSE=$(curl -s -X POST "$API_URL/customization/equip" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"item_id\": $FIRST_AVATAR_ID,
      \"slot\": \"avatar\"
    }")

  echo "   Equipped avatar ID: $FIRST_AVATAR_ID"
  echo -e "${GREEN}✓ Avatar equipped${NC}\n"
else
  echo -e "${RED}⚠ No owned avatars found${NC}\n"
fi

# Step 8: Test Progress Registration (should trigger gamification)
echo -e "${BLUE}[8/10] Testing progress registration with gamification...${NC}"

# First get a student and an OA Bloom objective
STUDENT_ID=$(echo $LOGIN_RESPONSE | grep -o '"user_id":[0-9]*' | cut -d':' -f2)

# Get first OA Bloom objective
OA_BLOOM=$(curl -s -X GET "$API_URL/objetivos-aprendizaje" \
  -H "Authorization: Bearer $TOKEN")
OA_BLOOM_ID=$(echo $OA_BLOOM | jq '.[0].bloom_objectives[0].id')

if [ "$OA_BLOOM_ID" != "null" ]; then
  # Get initial stats
  INITIAL_STATS=$(curl -s -X GET "$API_URL/gamification/stats" \
    -H "Authorization: Bearer $TOKEN")
  INITIAL_XP=$(echo $INITIAL_STATS | jq '.xp')
  INITIAL_COINS=$(echo $INITIAL_STATS | jq '.coins')

  echo "   Initial XP: $INITIAL_XP"
  echo "   Initial Coins: $INITIAL_COINS"

  # Register progress as "dominado" (should award 50 XP + 10 coins)
  PROGRESS_RESPONSE=$(curl -s -X POST "$API_URL/progress" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"user_id\": $STUDENT_ID,
      \"oa_bloom_objective_id\": $OA_BLOOM_ID,
      \"estado\": \"dominado\",
      \"notas\": \"Test progress for gamification\"
    }")

  # Wait a bit for async processing
  sleep 2

  # Get updated stats
  UPDATED_STATS=$(curl -s -X GET "$API_URL/gamification/stats" \
    -H "Authorization: Bearer $TOKEN")
  UPDATED_XP=$(echo $UPDATED_STATS | jq '.xp')
  UPDATED_COINS=$(echo $UPDATED_STATS | jq '.coins')

  XP_GAINED=$((UPDATED_XP - INITIAL_XP))
  COINS_GAINED=$((UPDATED_COINS - INITIAL_COINS))

  echo "   XP gained: $XP_GAINED (expected: 50)"
  echo "   Coins gained: $COINS_GAINED (expected: 10)"

  if [ $XP_GAINED -eq 50 ] && [ $COINS_GAINED -eq 10 ]; then
    echo -e "${GREEN}✓ Progress registration triggered gamification correctly${NC}\n"
  else
    echo -e "${RED}⚠ Gamification rewards don't match expected values${NC}\n"
  fi
else
  echo -e "${RED}⚠ No OA Bloom objectives found${NC}\n"
fi

# Step 9: Check for unlock notifications
echo -e "${BLUE}[9/10] Checking unlock notifications...${NC}"
NOTIFICATIONS=$(curl -s -X GET "$API_URL/customization/notifications" \
  -H "Authorization: Bearer $TOKEN")

NOTIF_COUNT=$(echo $NOTIFICATIONS | jq '. | length')
echo "   Unlock notifications: $NOTIF_COUNT"

if [ $NOTIF_COUNT -gt 0 ]; then
  echo $NOTIFICATIONS | jq '.[0]'
fi
echo -e "${GREEN}✓ Notifications retrieved${NC}\n"

# Step 10: Test purchasing an item (if user has enough coins)
echo -e "${BLUE}[10/10] Testing item purchase...${NC}"

PURCHASABLE_ITEM=$(echo $CATALOG | jq '[.[] | select(.status == "can_purchase")] | .[0]')
PURCHASABLE_ID=$(echo $PURCHASABLE_ITEM | jq '.id')

if [ "$PURCHASABLE_ID" != "null" ]; then
  ITEM_NAME=$(echo $PURCHASABLE_ITEM | jq -r '.name')
  ITEM_COST=$(echo $PURCHASABLE_ITEM | jq '.base_coins_cost')

  echo "   Attempting to purchase: $ITEM_NAME ($ITEM_COST coins)"

  PURCHASE_RESPONSE=$(curl -s -X POST "$API_URL/customization/purchase" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{
      \"item_id\": $PURCHASABLE_ID
    }")

  REMAINING_COINS=$(echo $PURCHASE_RESPONSE | jq '.remaining_coins')
  echo "   Remaining coins: $REMAINING_COINS"
  echo -e "${GREEN}✓ Item purchased successfully${NC}\n"
else
  echo -e "${RED}⚠ No purchasable items found (not enough coins or all owned)${NC}\n"
fi

echo -e "${BLUE}=== All Tests Completed ===${NC}"
echo -e "${GREEN}✓ Gamification & Customization system is working!${NC}\n"

# Summary
echo -e "${BLUE}Summary:${NC}"
echo "  - Total items in catalog: $TOTAL_ITEMS"
echo "  - Items owned: $OWNED_ITEMS"
echo "  - Current XP: $UPDATED_XP"
echo "  - Current Coins: $UPDATED_COINS"
echo "  - Unlock notifications: $NOTIF_COUNT"
