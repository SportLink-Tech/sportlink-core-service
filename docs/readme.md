# SportLink Core Service Documentation

Welcome to the SportLink Core Service documentation. This directory contains comprehensive documentation about the API, its use cases, and architecture.

## 📚 Documentation Structure

```
docs/
├── readme.md                    # This file - documentation index
├── backend/
│   ├── readme.md               # Backend documentation index
│   ├── api-overview.md         # High-level business overview
│   ├── api-documentation.md    # Complete API reference
│   ├── architecture.md         # Backend Hexagonal Architecture
│   └── usecases/
│       ├── create-team.md      # Create Team use case
│       └── retrieve-team.md    # Retrieve Team use case
├── frontend/
│   ├── readme.md               # Frontend documentation index
│   ├── architecture.md         # Atomic Hexagonal Architecture
│   ├── technology.md           # Technology stack and versions
│   └── best-practice.md        # Best practices and conventions
└── templates/
    ├── summary.md              # Backend documentation template
    └── frontend-architecture.md # Frontend documentation template
```

## 🚀 Quick Start

### Backend (API)

1. **New to the API?** Start with the [API Overview](./backend/api-overview.md) to understand what SportLink Core Service does
2. **Looking for endpoints?** Check the [API Documentation](./backend/api-documentation.md) for complete endpoint reference
3. **Need implementation details?** See individual use cases in [backend/usecases/](./backend/usecases/)
4. **Want to understand backend architecture?** See [Backend Architecture](./backend/architecture.md)

### Frontend

1. **New to frontend?** Start with the [Frontend Index](./frontend/readme.md) for navigation
2. **Understanding architecture?** See [Frontend Architecture](./frontend/architecture.md)
3. **What technologies?** Check [Technology Stack](./frontend/technology.md)
4. **Coding standards?** Review [Best Practices](./frontend/best-practice.md)

## 📖 Documentation Structure

### [Backend Documentation](./backend/)
Complete backend/API documentation:
- **[Backend Index](./backend/readme.md)**: Navigation guide for backend docs
- **[API Overview](./backend/api-overview.md)**: Business perspective and capabilities
- **[API Documentation](./backend/api-documentation.md)**: Complete endpoint reference with curl examples
- **[Architecture](./backend/architecture.md)**: Hexagonal Architecture, directory structure, setup
- **[Use Cases](./backend/usecases/)**: Detailed documentation for each endpoint
  - [Create Team](./backend/usecases/create-team.md)
  - [Retrieve Team](./backend/usecases/retrieve-team.md)

### [Frontend Documentation](./frontend/)
Complete frontend documentation following Atomic Hexagonal Architecture:
- **[Frontend Index](./frontend/readme.md)**: Navigation guide for frontend docs
- **[Architecture](./frontend/architecture.md)**: Atomic Hexagonal Architecture explanation (Ports, Adapters, Use Cases, Atomic Design)
- **[Technology Stack](./frontend/technology.md)**: All technologies, frameworks, versions, and justifications (React, TypeScript, Material-UI, Vite)
- **[Best Practices](./frontend/best-practice.md)**: Coding standards, conventions, and guidelines

### [Templates](./templates/)
Documentation generation templates and guides:
- **[Backend Documentation Template](./templates/summary.md)**: Instructions for creating backend API documentation
- **[Frontend Architecture Template](./templates/frontend-architecture.md)**: Instructions for creating frontend architecture documentation

## 🔗 Quick Links

### Quick Backend Examples

**Create a team**:
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

**Retrieve a team**:
```bash
curl -X GET http://localhost:8080/sport/Paddle/team/Thunder%20Strikers
```

**Health check**:
```bash
curl http://localhost:8080/livez
```

More examples in [API Documentation](./backend/api-documentation.md) and [Use Cases](./backend/usecases/)

## 🎯 API Use Cases at a Glance

| Use Case | Endpoint | Method | Documentation |
|----------|----------|--------|---------------|
| [Create Team](./backend/usecases/create-team.md) | `/team` | POST | Create a new sports team |
| [Retrieve Team](./backend/usecases/retrieve-team.md) | `/sport/:sport/team/:team` | GET | Get team details |

## 📊 Supported Features

- ✅ Multi-sport support (Football, Paddle, Tennis)
- ✅ Skill-based categorization (Unranked, L1-L7)
- ✅ Team management (create, retrieve)
- ✅ Player validation
- ✅ Team statistics tracking
- ✅ Health monitoring endpoints

## 🛠️ For Developers

### Adding New Documentation

When adding new use cases or features:

1. Follow the [documentation template](./templates/summary.md)
2. Create a new file in `docs/usecases/` for each use case
3. Include Mermaid diagrams for component interactions
4. Add curl examples for all scenarios (success and errors)
5. Update the [api-documentation.md](./api-documentation.md) use case inventory

### Documentation Standards

- Use clear, concise language
- Provide working curl examples
- Include component interaction diagrams
- Document all error cases
- Keep examples realistic and testable
- **Use lowercase file names with hyphens (kebab-case)**

## 📝 Contributing to Documentation

To contribute to this documentation:

1. Ensure all code examples are tested
2. Follow the existing documentation structure
3. Use Mermaid for diagrams
4. Keep use case documentation in separate files
5. Update the main inventory when adding new use cases
6. **Always use lowercase file names with hyphens**

## 🔍 Finding Information

### Backend
- **Business questions**: See [API Overview](./backend/api-overview.md)
- **How to use endpoints**: See [API Documentation](./backend/api-documentation.md)
- **Implementation details**: See individual use case files in [backend/usecases/](./backend/usecases/)
- **Backend architecture**: See [Backend Architecture](./backend/architecture.md)
- **Error handling**: See individual use case documentation or [API Documentation](./backend/api-documentation.md#error-handling)

### Frontend
- **Architecture concepts**: See [Frontend Architecture](../frontend/docs/architecture.md)
- **What technology to use**: See [Technology Stack](../frontend/docs/technology.md)
- **Coding standards**: See [Best Practices](../frontend/docs/best-practice.md)
- **Adding features**: See [Architecture: Adding Features](../frontend/docs/architecture.md#adding-a-new-feature)
- **Navigation**: Start with [Frontend Index](../frontend/docs/readme.md)

## 📞 Need Help?

If you can't find what you're looking for:

### Backend
1. Check the [Backend Index](./backend/readme.md) for navigation
2. Review [API Documentation](./backend/api-documentation.md) for endpoint reference
3. Check individual [use case documentation](./backend/usecases/)
4. Refer to [Backend Architecture](./backend/architecture.md)

### Frontend
1. Start with the [Frontend Index](../frontend/docs/readme.md)
2. Check [Frontend Architecture](../frontend/docs/architecture.md) for concepts
3. Review [Best Practices](../frontend/docs/best-practice.md) for coding standards

### Still need help?
Contact the development team

---

**Last Updated**: November 2025  
**API Version**: 1.0.0

