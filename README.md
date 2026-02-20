# OpenAI SDK

A modern full-stack application built with Angular 19 and Go Fiber, deployed on Google Cloud Run.

## Tech Stack

- **Frontend**: Angular 19 with standalone components
- **Backend**: Go with Fiber framework
- **Deployment**: Google Cloud Run (serverless containers)
- **CI/CD**: GitHub Actions

## Project Structure

```
.
├── frontend/          # Angular 19 application
│   ├── src/          # Source code
│   ├── Dockerfile    # Frontend container
│   └── nginx.conf    # Nginx configuration
├── backend/          # Go Fiber API
│   ├── main.go       # Entry point
│   ├── Dockerfile    # Backend container
│   └── go.mod        # Go dependencies
├── .github/
│   └── workflows/
│       └── ci.yml    # CI/CD pipeline
├── Dockerfile        # Combined build
└── docker-compose.yml # Local development
```

## Quick Start

### Prerequisites

- Node.js 20+
- Go 1.23+
- Docker (optional, for containerized development)

### Backend

```bash
cd backend
go mod download
go run main.go
```

The API will be available at `http://localhost:8080`

Endpoints:
- `GET /health` - Health check
- `GET /api/hello` - Hello world message

### Frontend

```bash
cd frontend
npm install
npm start
```

The app will be available at `http://localhost:4200`

### Docker Compose (Full Stack)

```bash
docker-compose up --build
```

- Frontend: http://localhost
- Backend: http://localhost:8080

## Deployment

### GCP Cloud Run

1. Set up GCP project and enable Cloud Run API
2. Create a service account with Cloud Run Admin and Storage Admin roles
3. Add GCP_SA_KEY and GCP_PROJECT_ID to GitHub secrets
4. Push to main branch to trigger deployment

### Manual Deployment

```bash
# Build and push images
gcloud builds submit --tag gcr.io/PROJECT_ID/marketpulse-backend ./backend
gcloud builds submit --tag gcr.io/PROJECT_ID/marketpulse-frontend ./frontend

# Deploy to Cloud Run
gcloud run deploy marketpulse-backend --image gcr.io/PROJECT_ID/marketpulse-backend --platform managed
gcloud run deploy marketpulse-frontend --image gcr.io/PROJECT_ID/marketpulse-frontend --platform managed
```

## Environment Variables

### Backend
- `PORT` - Server port (default: 8080)

### Frontend
- `API_URL` - Backend API URL

## Development

### Running Tests

Backend:
```bash
cd backend
go test ./...
```

Frontend:
```bash
cd frontend
npm test
```

### Code Style

- Go: Standard `gofmt`
- TypeScript: ESLint with Angular rules

## License

MIT
