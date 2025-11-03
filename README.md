## Setup

1. Start PostgreSQL:

```bash
docker-compose up -d

2. Run migrations (automatically executed on startup)

3. Run the application:
  go mod tidy
  go run main.go

4. Environment Variables
  WORKERS: Number of worker goroutines (default: 4)

5. API Endpoints
* Orders
- POST /orders - Create new order
- GET /orders/:id - Get order details

* Jobs
- POST /jobs/settlement - Create settlement job
- GET /jobs/:id - Get job status
- POST /jobs/:id/cancel - Cancel job
- GET /downloads/:id.csv - Download settlement CSV


result test:
order_service_test.go:77: === RESULTS ===
    order_service_test.go:78: Successful orders: 100
    order_service_test.go:79: Failed orders: 400
    order_service_test.go:80: Final stock: 0
    order_service_test.go:81: Total orders in database: 100
--- PASS: TestConcurrentOrders (0.80s)
PASS
ok      github.com/anieswahdie1/order-product-api.git/internal/services      1.410s
```
