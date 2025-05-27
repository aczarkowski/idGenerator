# ID Generator

A high-performance, distributed unique ID generator service written in Go, inspired by Twitter's Snowflake algorithm. This service generates unique 64-bit IDs that are roughly time-ordered and can be used across distributed systems.

## Features

- **High Performance**: Thread-safe ID generation with minimal latency
- **Distributed**: Supports multiple worker nodes with unique worker IDs
- **Time-ordered**: Generated IDs contain timestamp information for rough ordering
- **Multiple Time Providers**: Support for both Epoch and Julian calendar time systems
- **REST API**: Simple HTTP endpoint for ID generation
- **Configurable**: Flexible configuration options for different deployment scenarios
- **Batch Generation**: Generate multiple IDs in a single request

## Architecture

The ID structure follows a 64-bit format:

```
| Unused | Timestamp (41 bits) | Worker ID (3 bits) | Thread ID (5 bits) | Counter (10 bits) |
```

- **Timestamp**: 41 bits for time (epoch or Julian)
- **Worker ID**: 3 bits supporting up to 8 worker nodes
- **Thread ID**: 5 bits for thread identification
- **Counter**: 10 bits for sequence within the same timestamp

## API Endpoints

### Generate IDs

```
GET /
```

**Query Parameters:**
- `numberOfIds` (optional): Number of IDs to generate (default: 1)

**Response:**
```json
{
  "ids": [1234567890123456789]
}
```

**Error Response:**
```json
{
  "error": "error message"
}
```

## Configuration

The service can be configured using command-line flags:

| Flag | Default | Description |
|------|---------|-------------|
| `--port` | 1323 | Port number for the HTTP server |
| `--workerId` | 1 | Unique worker ID (0-7) |
| `--timeProvider` | "epoch" | Time provider type ("epoch" or "julian") |
| `--offset` | 1420070400000 | Time offset for the provider |

## Time Providers

### Epoch Time Provider
Uses Unix epoch time in milliseconds with a configurable offset. Default offset corresponds to January 1, 2015.

### Julian Time Provider
Uses Julian calendar system with a custom time encoding that includes:
- Last 2 digits of the year
- Day of year
- Hour, minute, second, and millisecond

## Installation & Usage

### Prerequisites
- Go 1.24 or higher

### Build and Run

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd idGenerator
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the application:
   ```bash
   go build -o idGenerator
   ```

4. Run the service:
   ```bash
   ./idGenerator --port=8080 --workerId=1 --timeProvider=epoch
   ```

### Docker Usage

Build Docker image:
```bash
docker build -t id-generator .
```

Run container:
```bash
docker run -p 8080:1323 id-generator --workerId=1
```

## Examples

### Generate a single ID
```bash
curl http://localhost:1323/
```

### Generate multiple IDs
```bash
curl "http://localhost:1323/?numberOfIds=10"
```

## Performance

The service is designed for high-throughput scenarios and includes:
- Mutex-based thread safety
- Efficient counter management
- Automatic timestamp collision handling
- Batch generation support

## Testing

Run all tests:
```bash
go test ./...
```

Run benchmarks:
```bash
go test -bench=. ./generator
```

Run integration tests:
```bash
go test -v integration_test.go main.go
```

## Project Structure

```
├── main.go                     # Application entry point
├── generator/                  # Core ID generation logic
│   ├── worker.go              # Main worker implementation
│   ├── worker_test.go         # Unit tests
│   └── benchmark_test.go      # Performance benchmarks
├── handler/                   # HTTP handlers
│   ├── generator.go           # ID generation endpoint
│   └── generator_test.go      # Handler tests
├── middleware/                # Custom middleware
│   ├── generatorprovider.go   # Worker instance provider
│   └── generatorprovider_test.go
├── timeprovider/              # Time provider implementations
│   ├── timeprovider.go        # Interface definition
│   ├── epoch/                 # Epoch time provider
│   └── julian/                # Julian calendar provider
├── integration_test.go        # Integration tests
└── main_test.go               # Main function tests
```

## Dependencies

- [Echo](https://echo.labstack.com/) - Web framework for the REST API
- Go standard library for core functionality

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]
