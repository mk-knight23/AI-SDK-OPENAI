# AI-SDK-OPENAI

[![AI-SDK Ecosystem](https://img.shields.io/badge/AI--SDK-ECOSYSTEM-part%20of-blue)](https://github.com/mk-knight23/AI-SDK-ECOSYSTEM)
[![OpenAI](https://img.shields.io/badge/OpenAI-1.51.0-green)](https://openai.com/)
[![Angular](https://img.shields.io/badge/Angular-19-red)](https://angular.io/)
[![Go](https://img.shields.io/badge/Go-1.23-cyan)](https://go.dev/)

> **Framework**: OpenAI SDK (Assistants API & GPT-4o)
> **Stack**: Angular 19 + Go Fiber

---

## 🎯 Project Overview

**AI-SDK-OPENAI** showcases the OpenAI Assistants API with GPT-4o integration. It demonstrates function calling, code interpretation, file uploads, and real-time streaming for building production AI applications.

### Key Features

- 🤖 **Assistants API** - OpenAI's managed agent framework
- 🧩 **Function Calling** - Tool use and API integration
- 📁 **File Uploads** - Document analysis and RAG
- 💬 **Real-time Streaming** - Live response generation
- 🔍 **Code Interpreter** - Safe code execution environment

---

## 🛠 Tech Stack

| Technology | Purpose |
|-------------|---------|
| Angular 19 | Frontend framework |
| Go Fiber | Backend API |
| OpenAI SDK | LLM integration |
| Angular Material | UI components |
| WebSocket | Real-time updates |

---

## 🚀 Quick Start

```bash
# Frontend
cd frontend && npm install && ng serve

# Backend
cd backend && go run main.go
```

---

## 🔌 API Integrations

| Provider | Usage |
|----------|-------|
| OpenAI | Primary (Assistants, GPT-4o) |
| Azure OpenAI | Enterprise fallback |

---

## 📦 Deployment

**Google Cloud Run**

```bash
gcloud run deploy
```

---

## 📁 Project Structure

```
AI-SDK-SDK-OPENAI/
├── frontend/         # Angular application
├── backend/          # Go Fiber API
└── README.md
```

---

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

---


---

## 🏗️ Architecture

```mermaid
graph TB
    subgraph "Frontend"
        UI[User Interface]
    end
    
    subgraph "Backend"
        API[API Layer]
        Core[AI Framework]
        Providers[LLM Providers]
    end
    
    subgraph "Infrastructure"
        DB[(Database)]
        Cache[(Cache)]
    end
    
    UI -->|HTTP/WS| API
    API --> Core
    Core --> Providers
    API --> DB
    Core --> Cache
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /health | Health check |
| POST | /api/execute | Execute agent workflow |
| WS | /api/stream | WebSocket streaming |

---

## 🔧 Troubleshooting

### Common Issues

**Connection refused**
- Ensure backend is running
- Check port availability

**Authentication failures**
- Verify API keys in `.env`
- Check environment variables

**Rate limiting**
- Implement exponential backoff
- Reduce request frequency

---

## 📚 Additional Documentation

- [API Reference](docs/API.md) - Complete API documentation
- [Deployment Guide](docs/DEPLOYMENT.md) - Platform-specific deployment
- [Testing Guide](docs/TESTING.md) - Testing strategies and coverage
---


**Part of the [AI-SDK Ecosystem](https://github.com/mk-knight23/AI-SDK-ECOSYSTEM)**
