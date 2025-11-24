# Frontend Documentation

Welcome to the SportLink frontend documentation. This directory contains comprehensive documentation about the frontend architecture, technologies, and best practices.

## 📚 Documentation Files

### [Architecture](architecture.md)
Complete guide to the frontend architecture following Atomic Hexagonal Architecture pattern.

**Contents:**
- Overview of Atomic Hexagonal Architecture
- Main architectural layers (Domain, Infrastructure, UI)
- Core concepts: Ports, Adapters, Use Cases, Context
- UI organization with Atomic Design (Atoms, Molecules, Organisms, Pages)
- Directory structure
- Step-by-step guide for adding new features
- Testing strategy
- Architecture diagrams and code examples

**Read this if you want to understand:**
- How the frontend is architectured
- Why we use Hexagonal Architecture
- How Atomic Design organizes components
- Where to put new code

---

### [Technology Stack](technology.md)
Complete overview of all technologies, frameworks, and tools used.

**Contents:**
- TypeScript (5.2.2)
- React (18.2.0)
- Material-UI (5.15.0)
- Emotion (CSS-in-JS)
- React Router (6.30.2)
- Vite (5.0.8)
- ESLint and development tools
- Version management and dependencies
- Why each technology was chosen

**Read this if you want to know:**
- What technologies are used
- Which versions are installed
- Why these technologies were chosen
- How to configure the development environment

---

### [Best Practices](best-practice.md)
Coding standards, conventions, and guidelines for the project.

**Contents:**
- File and directory naming conventions
- Code organization standards
- Domain layer best practices
- Infrastructure layer guidelines
- UI layer best practices
- What to avoid (anti-patterns)
- Code style guide
- Performance best practices
- Testing best practices
- Checklist before committing

**Read this if you want to:**
- Follow project conventions
- Write consistent code
- Avoid common mistakes
- Know what's acceptable and what's not

---

## 🚀 Quick Start

### For New Developers

1. **Start with [Architecture](architecture.md)** to understand the big picture
2. **Review [Technology Stack](technology.md)** to know what tools are used
3. **Read [Best Practices](best-practice.md)** before writing code

### For Specific Tasks

**Adding a new feature?**
→ See [Architecture: Adding a New Feature](architecture.md#adding-a-new-feature)

**Need to understand a concept?**
→ See [Architecture: Core Concepts](architecture.md#core-architectural-concepts)

**Want to know about a technology?**
→ See [Technology Stack](technology.md)

**Not sure about naming?**
→ See [Best Practices: Code Organization](best-practice.md#code-organization)

**Writing tests?**
→ See [Architecture: Testing Strategy](architecture.md#testing-strategy)

---

## 📖 Architecture Overview

The frontend follows **Atomic Hexagonal Architecture**, combining:

1. **Hexagonal Architecture** - Ports and Adapters pattern
2. **Atomic Design** - UI component hierarchy
3. **Domain-Driven Design** - Use cases and business logic

### Key Principles

```
Domain Layer (Core)
    ↑ 
Dependencies point inward
    ↓
Infrastructure Layer (Adapters)
```

### Directory Structure

```
src/
├── features/              # Feature modules
│   └── team/
│       ├── domain/        # Business logic
│       ├── infrastructure/# Adapters
│       ├── context/       # DI
│       └── ui/            # Components
├── shared/                # Reusable code
│   ├── components/
│   └── types/
└── components/            # Layout components
```

---

## 🎯 Key Concepts

### Ports & Adapters
- **Port**: Interface defining contract (in domain)
- **Adapter**: Implementation of port (in infrastructure)

### Use Cases
- Business logic and rules
- Pure TypeScript, no framework dependencies
- Testable independently

### Context (Dependency Injection)
- Wires dependencies together
- Provides use cases to UI components

### Atomic Design
- **Atoms**: Basic components (buttons, inputs)
- **Molecules**: Combined atoms
- **Organisms**: Complex components
- **Pages**: Complete views

---

## 🔗 External References

### Architecture Inspiration

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [The Atomic Hexagonal Architecture — with React](https://newlight77.medium.com/the-atomic-hexagonal-architecture-on-the-frontend-with-react-6337a56e56e3)
- [Atomic Design](https://atomicdesign.bradfrost.com/)

### Technology Documentation

- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Material-UI](https://mui.com/)
- [Vite Guide](https://vitejs.dev/guide/)

---

## 📝 Documentation Maintenance

### When to Update

Update documentation when:
- Architecture changes
- New patterns are introduced
- Dependencies are upgraded
- Best practices evolve

### How to Update

1. Edit the relevant markdown file
2. Keep examples current with code
3. Update version numbers if changed
4. Maintain consistent formatting

---

## 🤔 Need Help?

### Finding Information

- **What is X?** → Check [Architecture](architecture.md) concepts section
- **How do I Y?** → Check [Architecture](architecture.md) adding features section
- **Should I Z?** → Check [Best Practices](best-practice.md)
- **What version?** → Check [Technology Stack](technology.md)

### Still Have Questions?

1. Search the documentation (Ctrl/Cmd + F)
2. Check code examples in documentation
3. Review similar existing code in the project
4. Ask the team

---

**Last Updated:** November 2025  
**Documentation Version:** 1.0.0

