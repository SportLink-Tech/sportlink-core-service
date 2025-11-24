# SportLink Core Service - API Documentation

## Table of Contents

1. [Overview](#overview)
2. [Base URL](#base-url)
3. [Authentication](#authentication)
4. [Common Headers](#common-headers)
5. [Use Cases](#use-cases)
6. [Health Check Endpoints](#health-check-endpoints)
7. [Error Handling](#error-handling)
8. [Data Models](#data-models)

---

## Overview

For a comprehensive understanding of what SportLink Core Service does and the business problems it solves, please refer to the [API Overview](api-overview.md).

**Quick Summary**: SportLink Core Service is a RESTful API designed to facilitate the organization and management of sports activities, enabling users to create and manage sports teams across multiple sports (Paddle, Football, Tennis) with skill-based categorization.

---

## Base URL

```
http://localhost:8080
```

**Production URL**: (To be provided upon deployment)

---

## Authentication

Currently, the API does not require authentication. This may change in future versions.

---

## Common Headers

### Request Headers

```
Content-Type: application/json
Accept: application/json
```

### Response Headers

```
Content-Type: application/json
```

---

## Use Cases

This section provides an inventory of all available use cases in the SportLink Core Service. Each use case has detailed documentation including component interaction diagrams and curl examples.

### 1. [Create Team](usecases/create-team.md)

**Endpoint**: `POST /team`  
**Description**: Creates a new sports team with specified members, sport type, and skill category. Validates that all team members exist in the system before creating the team.

**Key Features**:
- Multi-sport support (Football, Paddle)
- Skill category assignment (Unranked, L1-L7)
- Team member validation
- Automatic statistics initialization

**Quick Example**:
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

[📖 View full documentation →](usecases/create-team.md)

---

### 2. [Retrieve Team](usecases/retrieve-team.md)

**Endpoint**: `GET /sport/:sport/team/:team`  
**Description**: Retrieves detailed information about a specific team by sport type and team name, including all team members and statistics.

**Key Features**:
- Sport-specific team lookup
- Complete team member details
- Team statistics (wins, losses, draws)
- Case-sensitive name matching

**Quick Example**:
```bash
curl -X GET http://localhost:8080/sport/Paddle/team/Thunder%20Strikers \
  -H "Accept: application/json"
```

[📖 View full documentation →](usecases/retrieve-team.md)

---

## Health Check Endpoints

The API provides health check endpoints for monitoring and orchestration systems:

### Liveness Check

**Endpoint**: `GET /livez`  
**Description**: Indicates whether the application is running  
**Response**: `200 OK` if the service is alive

```bash
curl -X GET http://localhost:8080/livez
```

### Readiness Check

**Endpoint**: `GET /readyz`  
**Description**: Indicates whether the application is ready to serve traffic  
**Response**: `200 OK` if the service is ready

```bash
curl -X GET http://localhost:8080/readyz
```

### Metrics

**Endpoint**: `GET /metrics`  
**Description**: Prometheus metrics endpoint for monitoring  
**Response**: Metrics in Prometheus format

```bash
curl -X GET http://localhost:8080/metrics
```

---

## Error Handling

All API errors follow a consistent JSON structure:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error description"
  }
}
```

### Common Error Codes

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `INVALID_REQUEST_FORMAT` | 400 | The request body is malformed or missing required parameters |
| `VALIDATION_ERROR` | 400 | One or more fields failed validation rules |
| `USE_CASE_EXECUTION_FAILED` | 500 | Business logic validation failed or database operation error |

### HTTP Status Codes

The API uses standard HTTP status codes:

| Status Code | Meaning | When Used |
|-------------|---------|-----------|
| `200 OK` | Success | Successful GET requests |
| `201 Created` | Created | Successful POST requests that create resources |
| `400 Bad Request` | Client Error | Invalid request format or validation errors |
| `404 Not Found` | Not Found | Resource or endpoint not found |
| `500 Internal Server Error` | Server Error | Use case execution failures or database errors |

---

## Data Models

### Sport Type

Enumeration of supported sports:

```
"Football" | "Paddle" | "Tennis"
```

### Category

Skill level classification (integer):

```
0  = Unranked
1  = L1 (Beginner)
2  = L2
3  = L3
4  = L4
5  = L5
6  = L6
7  = L7 (Advanced)
```

### Team Entity

```json
{
  "Name": "string",
  "Category": "integer (0-7)",
  "Stats": {
    "Wins": "integer",
    "Losses": "integer",
    "Draws": "integer"
  },
  "Sport": "string (Football|Paddle|Tennis)",
  "Members": [
    {
      "ID": "string",
      "Category": "integer (0-7)",
      "Sport": "string (Football|Paddle|Tennis)"
    }
  ]
}
```

### Player Entity

```json
{
  "ID": "string",
  "Category": "integer (0-7)",
  "Sport": "string (Football|Paddle|Tennis)"
}
```

### Statistics

```json
{
  "Wins": "integer",
  "Losses": "integer",
  "Draws": "integer"
}
```

---

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for local development)
- AWS LocalStack (for local DynamoDB)

### Running Locally

1. **Start LocalStack with DynamoDB**:
```bash
docker-compose up -d
```

2. **Run the application**:
```bash
make run
```

3. **Test the API**:
```bash
# Health check
curl http://localhost:8080/livez

# Create a team
curl -X POST http://localhost:8080/team \
  -H "Content-Type: application/json" \
  -d '{
    "sport": "Paddle",
    "name": "Test Team",
    "category": 2
  }'
```

---

## Architecture

This service follows **Hexagonal Architecture** (Ports and Adapters). For more details about the architecture, please see [Architecture Documentation](architecture.md).

**Key Principles**:
- Domain logic is isolated from infrastructure concerns
- Dependencies point inward toward the domain
- Easy to test and maintain
- Flexible and adaptable to different data sources

---

## Additional Resources

- [API Overview](api-overview.md) - High-level business perspective
- [Architecture Documentation](architecture.md) - Technical architecture details
- [Use Case: Create Team](usecases/create-team.md) - Detailed create team documentation
- [Use Case: Retrieve Team](usecases/retrieve-team.md) - Detailed retrieve team documentation

---

## Support and Contribution

For issues, questions, or contributions, please refer to the project's repository.

---

## Version History

- **v1.0.0** - Initial release with team management capabilities
  - Create Team endpoint
  - Retrieve Team endpoint
  - Multi-sport support (Football, Paddle, Tennis)
  - Skill-based categorization (L1-L7)

---

## Future Enhancements

Planned features for future versions:

- Player CRUD operations via REST endpoints
- Team update and deletion
- Advanced search and filtering
- Match scheduling and results
- Player statistics aggregation
- Authentication and authorization
- Rate limiting
- Pagination for list endpoints

