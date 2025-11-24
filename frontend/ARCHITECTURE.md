# Frontend Architecture

## Atomic Hexagonal Architecture

This frontend follows the **Atomic Hexagonal Architecture** pattern, which combines:
- **Hexagonal Architecture** (Ports & Adapters)
- **Atomic Design** (UI Component Hierarchy)
- **Domain-Driven Design** (Use Cases and Business Logic)

Reference: [The Atomic Hexagonal Architecture on the Frontend with React](https://newlight77.medium.com/the-atomic-hexagonal-architecture-on-the-frontend-with-react-6337a56e56e3)

## Architecture Layers

### 1. **Domain Layer** (`features/team/domain/`)

Contains pure business logic, independent of frameworks.

#### Ports (`ports/`)
Interfaces that define contracts for external services.
```typescript
// features/team/domain/ports/TeamRepository.ts
export interface TeamRepository {
  createTeam(request: CreateTeamRequest): Promise<{ data: Team; status: number }>
  findTeam(sport: string, teamName: string): Promise<{ data: Team; status: number }>
}
```

#### Use Cases (`usecases/`)
Business logic and rules.
```typescript
// features/team/domain/usecases/CreateTeamUseCase.ts
export class CreateTeamUseCase {
  constructor(private readonly teamRepository: TeamRepository) {}
  
  async execute(request: CreateTeamRequest): Promise<Result> {
    // Business logic here
  }
}
```

### 2. **Infrastructure Layer** (`features/team/infrastructure/`)

Implements ports with real-world adapters.

#### Adapters (`adapters/`)
Concrete implementations tightly coupled to external APIs/libraries.
```typescript
// features/team/infrastructure/adapters/TeamApiAdapter.ts
export class TeamApiAdapter implements TeamRepository {
  async createTeam(request: CreateTeamRequest): Promise<Response> {
    // Real HTTP call to backend API
    return fetch('/api/team', { ... })
  }
}
```

### 3. **Context Layer** (`features/team/context/`)

Manages Dependency Injection (DI).

```typescript
// features/team/context/TeamContext.tsx
export function TeamProvider({ children }) {
  // Create adapter
  const teamApiAdapter = new TeamApiAdapter()
  
  // Inject adapter into use cases (DI)
  const createTeamUseCase = new CreateTeamUseCase(teamApiAdapter)
  
  // Provide to UI
  return <TeamContext.Provider value={{ createTeamUseCase }}>
    {children}
  </TeamContext.Provider>
}
```

### 4. **UI Layer** (`features/team/ui/`)

Following **Atomic Design** principles:

#### Atoms (`shared/components/atoms/`)
Basic, reusable components
- `SportSelect.tsx` - Sport selection dropdown
- `CategorySelect.tsx` - Category selection dropdown

#### Molecules (`shared/components/molecules/`)
Combinations of atoms
- (To be added as needed)

#### Organisms (`features/team/ui/organisms/`)
Complex components combining molecules
- (To be added as needed)

#### Pages (`features/team/ui/pages/`)
Complete views that use Context and Use Cases
- `CreateTeamPage.tsx` - Uses `CreateTeamUseCase`
- `SearchTeamPage.tsx` - Uses `SearchTeamUseCase`

### 5. **Feature Module** (`features/team/TeamModule.tsx`)

Wraps the feature with its Provider for modularity.

```typescript
// features/team/TeamModule.tsx
export function TeamModule({ children }) {
  return <TeamProvider>{children}</TeamProvider>
}
```

## Directory Structure

```
frontend/src/
├── features/
│   └── team/                           # Team Feature Module
│       ├── domain/                     # Domain Layer (pure business logic)
│       │   ├── ports/                  # Interfaces (contracts)
│       │   │   └── TeamRepository.ts
│       │   └── usecases/               # Use Cases (business rules)
│       │       ├── CreateTeamUseCase.ts
│       │       └── SearchTeamUseCase.ts
│       ├── infrastructure/             # Infrastructure Layer
│       │   └── adapters/               # Concrete implementations
│       │       └── TeamApiAdapter.ts
│       ├── context/                    # Dependency Injection
│       │   └── TeamContext.tsx
│       ├── ui/                         # UI Layer (Atomic Design)
│       │   ├── pages/                  # Feature Pages
│       │   │   ├── CreateTeamPage.tsx
│       │   │   └── SearchTeamPage.tsx
│       │   ├── objects/                # Featured Objects (future)
│       │   └── organisms/              # Complex UI components (future)
│       └── TeamModule.tsx              # Feature Module wrapper
├── shared/                             # Shared across features
│   ├── components/
│   │   ├── atoms/                      # Basic reusable components
│   │   │   ├── SportSelect.tsx
│   │   │   └── CategorySelect.tsx
│   │   └── molecules/                  # Combined atoms (future)
│   └── types/
│       └── team.ts                     # Shared types
└── components/                         # Layout components
    └── Layout.tsx
```

## Dependency Flow

```
UI Components (Pages)
    ↓ uses
Use Cases (Domain Logic)
    ↓ depends on (interface)
Ports (Contracts)
    ↑ implements
Adapters (Real Implementation)
```

**Key Principle**: Dependencies point **inward** toward the domain.

## Benefits

1. **Modularity**: Each feature is self-contained in its module
2. **Testability**: Domain logic is pure and easy to test
3. **Flexibility**: Easy to swap implementations (mock adapters for testing)
4. **Maintainability**: Clear separation of concerns
5. **Scalability**: Add new features without touching existing code
6. **Reusability**: Atomic components can be reused across features

## Adding a New Feature

1. Create feature directory: `features/my-feature/`
2. Define **Port** (interface) in `domain/ports/`
3. Create **Use Case** in `domain/usecases/`
4. Implement **Adapter** in `infrastructure/adapters/`
5. Setup **Context** with DI in `context/`
6. Build **UI Components** following Atomic Design
7. Create **Feature Module** wrapper
8. Import in App.tsx

## Testing Strategy

- **Domain Layer**: Unit tests (pure logic, no dependencies)
- **Adapters**: Integration tests (mock HTTP)
- **UI Components**: Component tests (with mocked use cases)
- **E2E**: Full flow tests

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Atomic Design](https://atomicdesign.bradfrost.com/)
- [The Atomic Hexagonal Architecture — on the frontend — with React](https://newlight77.medium.com/the-atomic-hexagonal-architecture-on-the-frontend-with-react-6337a56e56e3)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)

