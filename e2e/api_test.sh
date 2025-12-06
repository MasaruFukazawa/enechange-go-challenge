#!/bin/bash
# E2E API Test Script

set -e

BASE_URL="${API_URL:-http://localhost:8070}"
PASSED=0
FAILED=0

test_case() {
    local name="$1"
    local expected="$2"
    local actual="$3"

    if [ "$expected" = "$actual" ]; then
        echo "✅ PASS: $name"
        ((PASSED++))
    else
        echo "❌ FAIL: $name (expected: $expected, got: $actual)"
        ((FAILED++))
    fi
}

echo "=== E2E API Tests ==="
echo "Base URL: $BASE_URL"
echo ""

# Test 1: 正常リクエスト - 旭川付近50km
echo "--- Test 1: Normal request (Asahikawa, 50km) ---"
RESULT=$(curl -s "${BASE_URL}/api/locations?latitude=43.796611&longitude=142.375917&radius=50")
COUNT=$(echo "$RESULT" | jq '. | length')
test_case "Returns locations" "true" "$([ "$COUNT" -gt 0 ] && echo true || echo false)"
test_case "Has 11 locations" "11" "$COUNT"

# Test 2: EVSEが含まれている
echo ""
echo "--- Test 2: Response contains EVSEs ---"
HAS_EVSE=$(echo "$RESULT" | jq '[.[].evses | length > 0] | any')
test_case "Locations have EVSEs" "true" "$HAS_EVSE"

# Test 3: JSONフォーマット確認
echo ""
echo "--- Test 3: JSON format validation ---"
FIRST_LOC=$(echo "$RESULT" | jq '.[0]')
HAS_ID=$(echo "$FIRST_LOC" | jq 'has("id")')
HAS_ADDRESS=$(echo "$FIRST_LOC" | jq 'has("address")')
HAS_COORDS=$(echo "$FIRST_LOC" | jq 'has("coordinates")')
test_case "Has id field" "true" "$HAS_ID"
test_case "Has address field" "true" "$HAS_ADDRESS"
test_case "Has coordinates field" "true" "$HAS_COORDS"

# Test 4: 遠隔地は結果なし
echo ""
echo "--- Test 4: Remote area returns empty ---"
REMOTE_RESULT=$(curl -s "${BASE_URL}/api/locations?latitude=0.000000&longitude=0.000000&radius=1")
REMOTE_COUNT=$(echo "$REMOTE_RESULT" | jq '. | length')
test_case "Remote area is empty" "0" "$REMOTE_COUNT"

# Test 5: 必須パラメータ欠落
echo ""
echo "--- Test 5: Missing required parameters ---"
LAT_ERROR=$(curl -s "${BASE_URL}/api/locations?longitude=142.375917" | jq -r '.error')
test_case "Missing latitude error" "latitude is required" "$LAT_ERROR"

LON_ERROR=$(curl -s "${BASE_URL}/api/locations?latitude=43.796611" | jq -r '.error')
test_case "Missing longitude error" "longitude is required" "$LON_ERROR"

# Test 6: 不正フォーマット
echo ""
echo "--- Test 6: Invalid format ---"
INVALID_LAT=$(curl -s "${BASE_URL}/api/locations?latitude=invalid&longitude=142.375917" | jq -r '.error')
test_case "Invalid latitude format" "Invalid latitude format" "$INVALID_LAT"

# Summary
echo ""
echo "=== Summary ==="
echo "Passed: $PASSED"
echo "Failed: $FAILED"

if [ "$FAILED" -gt 0 ]; then
    exit 1
fi
