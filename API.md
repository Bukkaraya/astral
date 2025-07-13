# Workflow API Endpoint

The service provides a simple HTTP endpoint to trigger the order processing workflow through the order orchestrator.

## Endpoint

### Process Order
```http
POST http://localhost:8080/api/workflows/order/process
Content-Type: application/json

{
  "orderId": "order-123",
  "userId": "user-alice"  // optional - for search attributes
}
```

## Response Format

The endpoint returns:
```json
{
  "workflowId": "process-order-order-123-1641234567",
  "runId": "abc123-def456-ghi789", 
  "message": "Order processing workflow started for order order-123"
}
```

## Testing

1. Start the service:
   ```bash
   go run .
   ```

2. Test the endpoint:
   ```bash
   ./test-api.sh
   ```

3. Or test manually:
   ```bash
   curl -X POST http://localhost:8080/api/workflows/order/process \
     -H "Content-Type: application/json" \
     -d '{"orderId":"test-123","userId":"user-alice"}'
   ```

## Search & Monitor

- Use `userId` field to enable searching by user
- View workflows: `temporal workflow list --query "userId='user-alice'"`
- Monitor execution: `temporal workflow show --workflow-id <workflow-id> --follow`