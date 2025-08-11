# App Store RSS Reviews Server

A Go backend service that polls iOS App Store Connect RSS feeds to fetch and store App Store RSS Reviews for a specific iOS app. The service provides a REST API to retrieve reviews.

## Features

- **Automated Polling**: Continuously polls App Store RSS feeds at configurable intervals
- **Data Persistence**: Stores reviews in JSON format with automatic state recovery
- **REST API**: Provides endpoints to fetch reviews with optional rating filtering

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   RSS Poller    │    │   HTTP Server    │    │  JSON Storage   │
│   (Cron Job)    │───▶│                  │───▶│  (File System)  │
└─────────────────┘    └──────────────────┘    └─────────────────┘
        │                        │
        ▼                        ▼
┌─────────────────┐    ┌──────────────────┐
│ App Store RSS   │    │  Mobile Client   │
│     Feed        │    │                  │
└─────────────────┘    └──────────────────┘
```

## Project Structure

```
server/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   └── config.go              # Configuration management
├── internal/
│   ├── api/
│   │   ├── handlers/          # HTTP request handlers
│   │   └── router.go          # Route definitions
│   ├── app/
│   │   └── app.go             # Business logic layer
│   ├── crons/
│   │   └── appstore_reviews_poller/  # RSS polling logic
│   ├── models/
│   │   └── appstore_review.go # Data models
│   └── repositories/
│       └── appstore_reviews.go # Data access layer
├── data/
│   └── reviews.json           # Persistent storage (auto-created)
├── go.mod                     # Go module definition
└── go.sum                     # Go module checksums
```

## Installation

1. **Clone the repository** (if not already done):

   ```bash
   git clone https://github.com/Gust4voSales/appstore-rss-reviews-app.git
   cd appstore-rss-reviews-app/server
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

## Configuration

The server can be configured using environment variables:

| Variable                   | Default             | Description                                  |
| -------------------------- | ------------------- | -------------------------------------------- |
| `PORT`                     | `8080`              | HTTP server port                             |
| `POLLING_INTERVAL_SECONDS` | `30`                | Interval between RSS feed polls (in seconds) |
| `APP_ID`                   | `389801252`         | App Id from App Store to subscribe to RSS    |
| `STORAGE_FILE_PATH`        | `data/reviews.json` | Path to JSON storage file                    |

**Note**: The current implementation is hardcoded to poll reviews for a specific app ID.

## Running the Server

### Development Mode

```bash
# Run with default configuration
go run cmd/main.go

# Run with custom environment variables
PORT=3000 POLLING_INTERVAL_SECONDS=30 go run cmd/main.go
```

### Using Docker (Optional)

Build and run:

```bash
docker build -t appstore-rss-reviews-server .
docker run -p 8080:8080 appstore-reviews-server
```

## API Endpoints

### Health Check

```
GET /health
```

**Response:**

```json
{
  "status": "ok",
  "uptime": "running"
}
```

### Get Reviews

```
GET /reviews
```

**Query Parameters:**

- `rating` (optional): Filter by rating (1-5)
- `hours` (optional): Filter reviews from last x hours. Defaults to 48h

**Example Requests:**

```bash
# Get all reviews
curl http://localhost:8080/reviews

# Get only 5-star reviews in last 24 hours
curl http://localhost:8080/reviews?rating=2&hours=24
```

**Response:**

```json
{
  "appId": "835599320",
  "count": 25,
  "reviews": [
    {
      "id": "review-id-123",
      "title": "Great app!",
      "content": "This app is amazing and works perfectly.",
      "author": "John Doe",
      "rating": 5,
      "updatedAt": "2024-01-15T10:30:00Z"
    }
  ],
  "lastHours": 24
}
```

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests in verbose mode
go test -v ./...
```
