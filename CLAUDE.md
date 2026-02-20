# MarketPulse

**Project:** MarketPulse - Real-time market data dashboard
**Tech Stack:** Angular 19 + Go Fiber + GCP Cloud Run
**Updated:** 2026-02-19

## Architecture

### Frontend
- Angular 19 with standalone components
- Signals for state management
- HttpClient for API communication
- SCSS for styling

### Backend
- Go 1.23 with Fiber framework
- RESTful API design
- Health check endpoint
- CORS enabled

### Deployment
- Containerized with Docker
- Google Cloud Run (serverless)
- GitHub Actions CI/CD

## Development Workflow

1. **Backend changes**: Edit in `backend/`, test with `go run main.go`
2. **Frontend changes**: Edit in `frontend/`, test with `npm start`
3. **Local full stack**: `docker-compose up --build`
4. **Deploy**: Push to main branch

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| GET | /api/hello | Hello message |

## File Locations

- Backend entry: `backend/main.go`
- Frontend entry: `frontend/src/main.ts`
- Component: `frontend/src/app/app.ts`
- Tests: `frontend/src/app/app.spec.ts`

## Environment

Backend port: 8080
Frontend port: 4200 (dev), 80 (production)

## Deployment Notes

Requires GitHub secrets:
- GCP_PROJECT_ID
- GCP_SA_KEY (service account JSON)
