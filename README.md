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

```
