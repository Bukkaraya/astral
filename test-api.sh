#!/bin/bash

# Test API endpoint for triggering order workflow

set -e

API_ENDPOINT="http://localhost:8080/api/workflows/order/process"

echo "🚀 Testing order processing API endpoint..."

# Check if API server is running
if ! curl -s ${API_ENDPOINT} > /dev/null 2>&1; then
    echo "❌ API server not running. Start with: go run ."
    exit 1
fi

echo "✅ API server is running"
echo ""

# Test Order Processing
echo "📦 Testing Order Processing API..."

echo "  → Processing order with user ID..."
curl -X POST ${API_ENDPOINT} \
  -H "Content-Type: application/json" \
  -d '{"orderId":"api-order-001","userId":"user-alice"}' \
  | jq '.'

echo ""
echo "  → Processing order without user ID..."
curl -X POST ${API_ENDPOINT} \
  -H "Content-Type: application/json" \
  -d '{"orderId":"api-order-002"}' \
  | jq '.'

echo ""
echo "✅ API endpoint tested!"
echo ""
echo "📊 Check workflows with:"
echo "  temporal workflow list --query \"userId='user-alice'\""
echo "  temporal workflow list"