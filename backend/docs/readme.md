# Backend Documentation

Welcome to the SportLink backend documentation. This directory contains comprehensive documentation about the API, architecture, use cases, and implementation details.

## 📚 Documentation Files

### [API Overview](api-overview.md)
High-level business perspective of the SportLink Core Service.

**Contents:**
- What is SportLink Core Service
- Business domain (sports management and social networking)
- Intended users (mobile apps, facility management, organizers)
- Main capabilities (team management, player management, multi-sport support)
- Problems it solves
- Technical approach (Hexagonal Architecture)
- Supported sports and skill categories

**Read this if you want to understand:**
- What the service does from a business perspective
- Who uses the service and why
- What problems it solves
- High-level capabilities

---

### [API Documentation](api-documentation.md)
Complete API reference with all endpoints and use cases.

**Contents:**
- Base URL and authentication
- Common headers
- Complete use case inventory with links
- Quick curl examples for each endpoint
- Health check endpoints (`/livez`, `/readyz`, `/metrics`)
- Error handling conventions
- Data models and schemas (Team, Player, Statistics)
- Getting started guide
- Version history

**Read this if you want to:**
- Use the API endpoints
- See curl examples
- Understand request/response formats
- Know available endpoints
- Check health and metrics

---

### [Architecture](architecture.md)
Technical architecture documentation following Hexagonal Architecture.

**Contents:**
- Hexagonal Architecture (Ports and Adapters) principles
- Directory structure and responsibilities
- Layer separation (domain, application, infrastructure)
- Development environment setup
- LocalStack configuration
- CloudFormation templates

**Read this if you want to understand:**
- How the backend is architected
- Where to put new code
- Hexagonal Architecture principles
- Directory organization
- Infrastructure setup

---

### [Use Cases](./usecases/)
Detailed documentation for each API use case with diagrams and examples.

#### [Create Team](usecases/create-team.md)
Complete documentation for the team creation endpoint.

**Contents:**
- Business goal and purpose
- HTTP endpoint: `POST /team`
- Input/output specifications
- Validations and business rules
- Mermaid diagram showing component interactions
- Curl examples (success and error cases)
- Error handling details

#### [Retrieve Team](usecases/retrieve-team.md)
Complete documentation for the team retrieval endpoint.

**Contents:**
- Business goal and purpose
- HTTP endpoint: `GET /sport/:sport/team/:team`
- Input/output specifications
- Query parameters
- Mermaid diagram showing component interactions
- Curl examples with URL encoding
- Error scenarios and known issues

**Read use case docs to:**
- Understand how a specific endpoint works
- See component interaction diagrams
- Get working curl examples
- Know all error cases

---

## 🚀 Quick Start

### For New Developers

1. **Start with [API Overview](api-overview.md)** to understand what the service does
2. **Review [Architecture](architecture.md)** to understand the code structure
3. **Check [API Documentation](api-documentation.md)** for endpoint reference
4. **Read individual [Use Cases](./usecases/)** for detailed implementation info

### For API Users

1. **Read [API Overview](api-overview.md)** for business context
2. **Check [API Documentation](api-documentation.md)** for endpoint list
3. **Copy curl examples** from use case documentation

### For Specific Tasks

**Need to call an endpoint?**
→ See [API Documentation](api-documentation.md) for quick reference

**Want to understand how an endpoint works?**
→ See specific use case documentation in [usecases/](./usecases/)

**Adding a new feature?**
→ See [Architecture](architecture.md) for structure guidelines

**Setting up development environment?**
→ See [Architecture: Development Environment](architecture.md)

---

## 📖 Architecture Overview

The backend follows **Hexagonal Architecture** (Ports and Adapters pattern).

### Key Principles

```
┌─────────────────────────────────────┐
│     Infrastructure (REST API)       │  ← Adapters
│         (Controllers)               │
└──────────────┬──────────────────────┘
               │
               ↓
┌─────────────────────────────────────┐
│      Application Layer              │  ← Use Cases
│         (Use Cases)                 │
└──────────────┬──────────────────────┘
               │
               ↓
┌─────────────────────────────────────┐
│       Domain Layer                  │  ← Core Business Logic
│   (Entities, Repositories)          │
└─────────────────────────────────────┘
```

### Directory Structure

```
api/
├── application/           # Use cases (business workflows)
│   ├── player/
│   │   └── usecases/
│   └── team/
│       └── usecases/
├── domain/               # Core business logic
│   ├── common/
│   ├── player/
│   │   ├── entity.go
│   │   └── repository.go
│   └── team/
│       ├── entity.go
│       └── repository.go
└── infrastructure/       # External adapters
    ├── config/
    ├── persistence/
    │   ├── dynamodb/
    │   ├── player/
    │   └── team/
    └── rest/
        ├── player/
        └── team/
```

---

## 🎯 Available Endpoints

### Team Management

| Endpoint | Method | Purpose | Documentation |
|----------|--------|---------|---------------|
| `/team` | POST | Create a new team | [Create Team](usecases/create-team.md) |
| `/sport/:sport/team/:team` | GET | Retrieve team by sport and name | [Retrieve Team](usecases/retrieve-team.md) |

### Health & Monitoring

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/livez` | GET | Liveness check |
| `/readyz` | GET | Readiness check |
| `/metrics` | GET | Prometheus metrics |

---

## 💾 Technology Stack

### Core

- **Language**: Go 1.21+
- **Framework**: Gin (HTTP framework)
- **Database**: DynamoDB (via AWS SDK v2)
- **Architecture**: Hexagonal (Ports and Adapters)

### Development

- **LocalStack**: Local AWS services (DynamoDB, SQS)
- **Docker**: Containerization
- **Docker Compose**: Local environment orchestration

### Tools

- **GolangCI-Lint**: Code linting
- **Mockery**: Mock generation for testing
- **Go Modules**: Dependency management

---

## 🔗 Data Models

### Sport Types

```
Football | Paddle | Tennis
```

### Category (Skill Levels)

```
0 = Unranked
1 = L1 (Beginner)
2 = L2
3 = L3
4 = L4
5 = L5
6 = L6
7 = L7 (Advanced)
```

### Team Entity

```json
{
  "Name": "string",
  "Category": 0-7,
  "Stats": {
    "Wins": 0,
    "Losses": 0,
    "Draws": 0
  },
  "Sport": "Football|Paddle|Tennis",
  "Members": [
    {
      "ID": "string",
      "Category": 0-7,
      "Sport": "Football|Paddle|Tennis"
    }
  ]
}
```

---

## 📝 API Examples

### Create Team

```bash
curl -X POST http://localhost:8080/team \
  -H "Content-Type: application/json" \
  -d '{
    "sport": "Paddle",
    "name": "Thunder Strikers",
    "category": 3,
    "players": ["player-001", "player-002"]
  }'
```

### Retrieve Team

```bash
curl -X GET http://localhost:8080/sport/Paddle/team/Thunder%20Strikers
```

### Health Check

```bash
curl http://localhost:8080/livez
```

---

## 🚦 Error Handling

All API errors follow a consistent format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable description"
  }
}
```

### Common Error Codes

- `INVALID_REQUEST_FORMAT` (400): Malformed request
- `VALIDATION_ERROR` (400): Validation failed
- `USE_CASE_EXECUTION_FAILED` (500): Business logic or database error

---

## 🏃 Running the Backend

### With Docker (Recommended)

```bash
# Start all services (backend + LocalStack)
make env-up

# Stop services
make env-down
```

### Locally (Development)

```bash
# Install dependencies
make set-up

# Run tests
make test

# Run linter
make lint

# Generate mocks
make generate-mocks
```

---

## 🧪 Testing

```bash
# Run all tests
make test

# Generate coverage report
make coverage-report
```

---

## 🔍 Finding Information

- **What does the API do?** → [API Overview](api-overview.md)
- **How do I call an endpoint?** → [API Documentation](api-documentation.md)
- **How does a specific endpoint work?** → [Use Cases](./usecases/)
- **Where should I put my code?** → [Architecture](architecture.md)
- **What's the database schema?** → [Architecture](architecture.md) and Use Case docs
- **How do I set up LocalStack?** → [Architecture](architecture.md)

---

## 📚 External References

### Architecture

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) by Alistair Cockburn
- [Ports and Adapters Pattern](https://herbertograca.com/2017/09/14/ports-adapters-architecture/)

### Go

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)

### AWS

- [AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/docs/)
- [DynamoDB Documentation](https://docs.aws.amazon.com/dynamodb/)
- [LocalStack Documentation](https://docs.localstack.cloud/)

---

## 📝 Documentation Maintenance

### When to Update

Update documentation when:
- New endpoints are added
- API contracts change
- Architecture evolves
- Error codes are added/modified

### How to Update

1. Update the relevant documentation file
2. Add/update use case documentation for new endpoints
3. Update curl examples if needed
4. Maintain Mermaid diagrams
5. Keep version history current

---

**Last Updated:** November 2025  
**API Version:** 1.0.0  
**Go Version:** 1.21+

