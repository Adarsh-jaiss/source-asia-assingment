# Source Asia - Backend Assignment

This repository contains the backend assignment for Source Asia, implementing a rate-limited API and a product catalog with media management. The service is written in Go and uses the Gin HTTP framework.

## Project Structure
- **api/models**: Contains data transfer objects and internal models.
- **api/repository**: Contains the in-memory data store for rate limits and product catalog.
- **api/controllers**: Contains HTTP handlers and business logic.
- **routes**: Defines the Gin routes wiring everything together.
- **main.go**: The entrypoint that initializes dependencies and starts the server.

## Features & Implementation Notes

### Part 1 - Rate-limited API
- **Approach**: A fixed-window rate limiter is used (1 minute). The window resets completely when the current time exceeds the `WindowStart + 1 minute`.
- **Data Store**: An in-memory map `map[string]*UserRateLimit` protected by `sync.RWMutex` ensures thread-safe updates under concurrent load.
- **Endpoints**:
  - `POST /api/v1/request`: Accepts a request. 
    - Success: `201 Created`
    - Rate Limit Exceeded: `429 Too Many Requests`
    - Invalid Input: `400 Bad Request`
  - `GET /api/v1/stats`: Returns rate-limiting statistics for all users. The response contains both the `accepted_count_current_window` (which resets) and `rejected_count_total` (cumulative).

### Part 2 - Product Catalog with Media
- **Approach**: In-memory CRUD with separate structures to optimize performance.
- **Data Model**: The repository maintains the full product information internally. 
- **List vs Detail optimization**: When `GET /api/v1/products` is called, it iterates and returns a mapped list of `ProductSummaryResponse` structs (excluding the large `ImageURLs` and `VideoURLs` slices). This satisfies the requirement of returning thousands of products without serializing 10,000+ strings. The endpoint uses pagination (`limit` and `offset` queries).
- **Endpoints**:
  - `POST /api/v1/products`: Creates a new product. 
    - Validates duplicate SKU -> `409 Conflict`
    - Validates URLs and limit (max 20 per array) -> `400 Bad Request`
  - `GET /api/v1/products?limit=20&offset=0`: Retrieves paginated list of products without loading full URLs.
  - `GET /api/v1/products/:id`: Returns full details of the product including all URLs.
  - `POST /api/v1/products/:id/media`: Appends new media to the product.


### Validation & Business Rules
- **Name and SKU**: Must be non-empty strings. Lead/trail spaces are trimmed. Empty names or SKUs return `400 Bad Request`.
- **SKU Uniqueness**: Creating a product with an existing SKU returns `409 Conflict`.
- **URL Validation**: 
  - Must use `http://` or `https://` schemes.
  - Maximum character length per URL is **2048 characters**.
  - Any URL violating these rules returns `400 Bad Request`.
- **Batch Limits**: A maximum of **20 URLs** are allowed in each media array (`image_urls`/`video_urls`) per request. Exceeding this returns `400 Bad Request`.

### JSON Schemas & Response Formats

#### Rate Limiting Request Confirmations (POST `/api/v1/request`)
- **Success (201 Created)**:
  ```json
  {
    "success": true,
    "data": {
      "message": "Request accepted"
    }
  }
  ```

#### Rate Limiting Stats (GET `/api/v1/stats`)
- **Success (200 OK)**:
  ```json
  {
    "success": true,
    "data": {
      "stats": [
        {
          "user_id": "user123",
          "accepted_count_current_window": 3,
          "rejected_count_total": 0
        }
      ]
    }
  }
  ```
  *Note: `accepted_count_current_window` resets every fixed 1-minute window. `rejected_count_total` is cumulative across windows.*

#### Standard Error Response (e.g. 400, 429, 409, 500)
- **Format**:
  ```json
  {
    "success": false,
    "error": {
      "code": "ERROR_CODE_STRING",
      "message": "Human-readable description"
    }
  }
  ```

### Production Limitations
- Since this is an in-memory implementation, all state is lost upon process restart.
- For a multi-instance deployment (like Kubernetes), this local rate-limiter and product store would not be synced across instances. 
- **What to change for Production**:
  - Use **Redis** for distributed rate limiting (e.g., using a sliding window with sorted sets).
  - Use **PostgreSQL** to store products and a separate table for media items. List queries would use `LIMIT / OFFSET` directly in SQL, avoiding loading large collections in memory.
  - A real CDN would be integrated on the frontend or processing layer.


## How to Run

### Option 1: Running Locally (Go)

1. **Prerequisites**: Ensure Go (>= 1.20) is installed.
2. **Download modules**:
   ```bash
   go mod tidy
   ```
3. **Run the server**:
   ```bash
   go run main.go
   ```
   OR 
   
   ```bash
   make run
   ```
   The server will start on port `8080`.

### Option 2: Running with Docker Compose

If you prefer to run the application containerized using Docker Compose:
```bash
docker-compose up --build -d
```
The server will start on port `8080`.

## Running Tests

To execute the automated unit and integration tests located in the `tests/` directory, run:
```bash
go test -v ./tests/...
```

## Swagger Documentation

An interactive API documentation (Swagger UI) is available. Once the server is running, you can access it in your browser at:
**[http://localhost:8080/api/docs/index.html](http://localhost:8080/api/docs/index.html)**

This allows reviewers to easily explore and test the endpoints directly from the browser.

## Example Usage

### 1. Test Rate Limiter (POST /request)
```bash
curl -X POST http://localhost:8080/api/v1/request \
     -H "Content-Type: application/json" \
     -d '{"user_id": "user123", "payload": {"foo": "bar"}}'
```
*Run this more than 5 times in a minute to see the `429 Too Many Requests` response.*

### 2. Get Stats (GET /stats)
```bash
curl http://localhost:8080/api/v1/stats
```

### 3. Create Product (POST /products)
```bash
curl -X POST http://localhost:8080/api/v1/products \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Widget A",
       "sku": "SKU-001",
       "image_urls": ["https://cdn.example.com/img1.jpg", "https://cdn.example.com/img2.jpg"],
       "video_urls": []
     }'
```

### 4. Append Media (POST /products/:id/media)
*(Replace `<ID>` with the ID returned from the creation step)*
```bash
curl -X POST http://localhost:8080/api/v1/products/<ID>/media \
     -H "Content-Type: application/json" \
     -d '{
       "video_urls": ["https://cdn.example.com/demo.mp4"]
     }'
```

### 5. List Products (GET /products)
```bash
curl "http://localhost:8080/api/v1/products?limit=10&offset=0"
```

### 6. Get Product Details (GET /products/:id)
```bash
curl "http://localhost:8080/api/v1/products/<ID>"
```

## AI Tools
* AI (Gemini via Agent) was used to assist in writing standard boilerplate, outlining the implementation plan, formatting code, and producing this README.
