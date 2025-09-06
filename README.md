# RediLink

A high-performance URL shortening service built with Go and Redis. RediLink provides a simple yet robust solution for creating and managing shortened URLs with built-in rate limiting and analytics tracking.

## Features

- **URL Shortening**: Convert long URLs into short, manageable links
- **Custom Short Codes**: Option to specify custom short codes for branded links
- **Automatic Expiration**: Configurable URL expiration (default: 24 hours)
- **Rate Limiting**: IP-based rate limiting (10 requests per 30 minutes)
- **Redis Storage**: High-performance Redis backend for fast URL resolution
- **Usage Analytics**: Built-in counter tracking for shortened URLs
- **RESTful API**: Clean REST API with comprehensive documentation
- **Docker Support**: Containerized deployment with Docker Compose

## Architecture

The application consists of two main components:

- **API Server**: Go-based HTTP server using Fiber framework
- **Database**: Redis instance for URL storage and rate limiting
- **Swagger**: API Reference at `/api/reference` 

## Prerequisites

- Docker and Docker Compose
- Go 1.22.2+ (for local development)
- Redis (if running without Docker)

## Quick Start

### Using Docker Compose (Recommended)

1. Clone the repository:
   ```bash
   git clone https://github.com/kunalsinghdadhwal/redilink
   cd redilink
   ```

2. Start the services:
   ```bash
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:3000`

### Local Development

1. Install dependencies:
   ```bash
   cd api
   go mod download
   ```

2. Set up environment variables:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Start Redis server locally

4. Run the application:
   ```bash
   go run main.go
   ```

## API Documentation

### Interactive Documentation

Access the interactive API documentation at:
- **Scalar UI**: `http://localhost:3000/api/reference`

### Endpoints

#### Shorten URL
```
POST /api/v1
```

Create a shortened URL from a long URL.

**Request Body:**
```json
{
  "url": "https://www.example.com",
  "short": "mylink",
  "expiry": 24
}
```

**Response:**
```json
{
  "url": "https://www.example.com",
  "short": "http://localhost:3000/abc123",
  "expiry": 24,
  "rate_limit": 9,
  "rate_limit_reset": 29
}
```

#### Resolve URL
```
GET /{short_code}
```

Redirect to the original URL using the shortened code.

**Example:**
```bash
curl -L http://localhost:3000/abc123
```

## Configuration

Environment variables can be configured in the `.env` file:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_ADDR` | Redis server address | `db:6379` |
| `DB_PASS` | Redis password | `""` |
| `APP_PORT` | Application port | `:3000` |
| `DOMAIN` | Base domain for shortened URLs | `localhost:3000` |
| `API_QUOTA` | Rate limit per IP | `10` |

## Rate Limiting

The API implements IP-based rate limiting:
- **Limit**: 10 requests per 30 minutes per IP address
- **Headers**: Rate limit information is included in responses
- **Reset**: Automatic reset after 30 minutes

## URL Expiration

URLs have configurable expiration times:
- **Default**: 24 hours
- **Custom**: Specify expiry in hours via API request
- **Cleanup**: Expired URLs are automatically removed by Redis TTL

## Development

### Project Structure

```
redilink/
├── api/
│   ├── main.go              # Application entry point
│   ├── routes/
│   │   ├── shorten.go       # URL shortening logic
│   │   └── resolve.go       # URL resolution logic
│   ├── database/
│   │   └── database.go      # Redis client configuration
│   ├── helpers/
│   │   └── helpers.go       # Utility functions
│   ├── docs/                # Generated API documentation
│   ├── Dockerfile
│   ├── go.mod
│   └── .env
├── db/
│   └── Dockerfile           # Redis configuration
└── docker-compose.yml
```

### Building from Source

```bash
cd api
go build -o redilink .
./redilink
```

### Running Tests

```bash
cd api
go test ./...
```

### Testing the API

#### Create a shortened URL:
```bash
curl -X POST -H "Content-type: application/json" \
  -d '{"url": "https://www.google.com"}' \
  http://localhost:3000/api/v1

# For formatted output:
curl -X POST -H "Content-type: application/json" \
  -d '{"url": "https://www.google.com"}' \
  http://localhost:3000/api/v1 | jq
```

#### Create a custom shortened URL:
```bash
curl -X POST -H "Content-type: application/json" \
  -d '{"url": "https://www.google.com", "short": "google", "expiry": 48}' \
  http://localhost:3000/api/v1 | jq
```

#### Access the shortened URL:
```bash
curl -L http://localhost:3000/{short_code}
```

## Docker Deployment

The application includes Docker configurations for production deployment:

### Building Images

```bash
# Build API image
docker build -t redilink-api ./api

# Build Redis image
docker build -t redilink-db ./db
```

### Production Deployment

For production deployment, consider:

1. **Environment Variables**: Use secure environment variable management
2. **SSL/TLS**: Configure HTTPS termination via reverse proxy
3. **Monitoring**: Implement application and infrastructure monitoring
4. **Backup**: Set up Redis data backup and recovery procedures
5. **Scaling**: Use Redis Cluster for horizontal scaling

## Security Considerations

- Rate limiting prevents abuse and ensures fair usage
- URL validation prevents malicious URL submission
- Domain restrictions can be configured to prevent self-referential loops
- Input sanitization protects against injection attacks

## Performance

- **Redis Storage**: Sub-millisecond URL resolution
- **Connection Pooling**: Efficient database connection management
- **Lightweight**: Minimal resource footprint
- **Concurrent**: Handles multiple requests simultaneously

## Monitoring

The application provides basic analytics:
- URL usage counter tracking
- Rate limit monitoring
- Request logging via Fiber middleware
