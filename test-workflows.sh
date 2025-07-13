#!/bin/bash

# Test all workflows with sample data and comprehensive examples

set -e

TASK_QUEUE="microservice-task-queue"

echo "🚀 Starting comprehensive workflow tests..."

# Check if worker is running
if ! curl -s http://localhost:9090/health > /dev/null 2>&1; then
    echo "❌ Worker not running. Start with: go run ."
    exit 1
fi

echo "✅ Worker is running"

# Register search attribute if not already registered
echo "🔧 Setting up search attributes..."
temporal operator search-attribute create --name userId --type Text --namespace claude 2>/dev/null || echo "  → userId search attribute already exists"

echo ""

# Test Order Workflows with different scenarios
echo "📦 Testing Order Workflows..."
echo "  → Processing regular order..."
temporal workflow start --task-queue "$TASK_QUEUE" --type "ProcessOrder.v1" --input '{"orderId":"order-001"}' --search-attribute 'userId="user-alice"'

echo "  → Processing priority order..."
temporal workflow start --task-queue "$TASK_QUEUE" --type "ProcessOrder.v1" --input '{"orderId":"order-002-priority"}' --search-attribute 'userId="user-bob"'

echo "  → Canceling order..."
temporal workflow start --task-queue "$TASK_QUEUE" --type "CancelOrder.v1" --input '{"orderId":"order-003-cancel"}' --search-attribute 'userId="user-charlie"'

echo ""

# Test Payment Workflows with different scenarios
echo "💳 Testing Payment Workflows..."
echo "  → Processing small payment..."
temporal workflow start --task-queue "$TASK_QUEUE" --type "ProcessPayment.v1" --input '{"paymentId":"payment-001","amount":29.99}' --search-attribute 'userId="user-alice"'

echo "  → Processing large payment..."
temporal workflow start --task-queue "$TASK_QUEUE" --type "ProcessPayment.v1" --input '{"paymentId":"payment-002","amount":299.99}' --search-attribute 'userId="user-bob"'

echo "  → Processing refund..."
temporal workflow start --task-queue "$TASK_QUEUE" --type "RefundPayment.v1" --input '{"paymentId":"payment-003"}' --search-attribute 'userId="user-charlie"'

echo "  → Processing international payment..."
temporal workflow start --task-queue "$TASK_QUEUE" --type "ProcessPayment.v1" --input '{"paymentId":"payment-004-intl","amount":149.50}' --search-attribute 'userId="user-diana"'

echo ""
echo "✅ All workflows started successfully!"
echo ""
echo "📊 Search and Monitor Examples:"
echo "  # List all workflows"
echo "  temporal workflow list"
echo ""
echo "  # Search by user"
echo "  temporal workflow list --query \"userId='user-alice'\""
echo "  temporal workflow list --query \"userId='user-bob'\""
echo "  temporal workflow list --query \"userId='user-charlie'\""
echo ""
echo "  # Search by workflow type"
echo "  temporal workflow list --query \"WorkflowType='ProcessOrder.v1'\""
echo "  temporal workflow list --query \"WorkflowType='ProcessPayment.v1'\""
echo ""
echo "  # Search by status"
echo "  temporal workflow list --query \"ExecutionStatus='Running'\""
echo "  temporal workflow list --query \"ExecutionStatus='Completed'\""
echo ""
echo "  # Combined searches"
echo "  temporal workflow list --query \"userId='user-alice' AND WorkflowType='ProcessPayment.v1'\""
echo ""
echo "🔍 Monitor workflow execution:"
echo "  temporal workflow show --workflow-id <WORKFLOW_ID>"
echo "  temporal workflow show --workflow-id <WORKFLOW_ID> --follow"