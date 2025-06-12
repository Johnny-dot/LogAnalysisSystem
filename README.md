# LogAnalysisSystem

This project is a minimal skeleton for a Web-based log analysis backend. It follows the architecture proposed in the project description and demonstrates how to start building a system similar to Bugly or Sentry.

## Features

- **Upload endpoint**: `/api/upload` accepts `.tar.gz` files, extracts them and parses each line into memory.
- **Analyze endpoint**: `/api/analyze` returns the parsed log entries. Query with `?level=ERROR` to filter by level.

The backend is implemented in Go using the Gin framework. Parsed log entries are stored in memory for demonstration purposes.

## Getting Started

1. Install Go 1.20 or later.
2. Initialize dependencies:
   ```bash
   cd backend
   go mod tidy
   ```
   (Requires internet access to download Gin.)
3. Run the server:
   ```bash
   go run ./cmd
   ```
4. Upload logs:
   ```bash
   curl -F "file=@sample.tar.gz" http://localhost:8080/api/upload
   ```
5. Query log entries:
   ```bash
   curl http://localhost:8080/api/analyze?level=ERROR
   ```

## Next Steps

- Persist parsed logs in a database such as ClickHouse.
- Add aggregation and dashboarding via a React or Vue frontend under `frontend/`.
- Containerize with Docker and provide Kubernetes manifests.

This repository only includes the backend skeleton. It can be extended to implement the full workflow described in the detailed plan.
