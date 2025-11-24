# Frontend Architecture

## Overview

The SportLink frontend follows the **Atomic Hexagonal Architecture** pattern, which combines three powerful architectural approaches:

1. **Hexagonal Architecture** (Ports and Adapters)
2. **Atomic Design** (UI Component Hierarchy)
3. **Domain-Driven Design** (Use Cases and Business Logic)

This architecture was chosen to create a maintainable, testable, and scalable frontend application that isolates business logic from external concerns (UI frameworks, APIs, etc.).

### Why This Architecture?

**Problem it solves:**
- **Tight coupling**: Traditional frontend apps often mix business logic with UI components
- **Difficulty testing**: Business logic mixed with framework code is hard to test
- **Poor maintainability**: Changes in one area affect multiple unrelated parts
- **Scalability issues**: Adding features requires touching existing code

**How it improves code organization:**
- Clear separation of concerns across well-defined layers
- Business logic is framework-independent and easily testable
- UI components follow a clear hierarchy (Atomic Design)
- Features are self-contained modules
- Dependencies flow in one direction (toward the domain)

---

## Main Architectural Layers

The architecture consists of three primary layers:

### 1. **Domain Layer** (Core)
- **Contains**: Business logic, use cases, and interfaces (ports)
- **Characteristics**: Pure TypeScript, no framework dependencies
- **Location**: `features/[feature-name]/domain/`

### 2. **Infrastructure Layer**
- **Contains**: Concrete implementations (adapters) for external services
- **Characteristics**: Tightly coupled to APIs, libraries, frameworks
- **Location**: `features/[feature-name]/infrastructure/`

### 3. **UI Layer**
- **Contains**: React components following Atomic Design
- **Characteristics**: Organized by component complexity (atoms → molecules → organisms → pages)
- **Location**: `features/[feature-name]/ui/` and `shared/components/`

### Layer Interaction

```
┌─────────────────────────────────────────┐
│          UI Layer (Pages)               │  ← User Interaction
│      (React Components)                 │
└──────────────┬──────────────────────────┘
               │ uses Context Hook
               ↓
┌─────────────────────────────────────────┐
│        Domain Layer (Use Cases)         │  ← Core Business Logic
│      (Pure TypeScript)                  │
└──────────────┬──────────────────────────┘
               │ depends on (Port Interface)
               ↓
┌─────────────────────────────────────────┐
│         Port (Interface)                │  ← Contract
└──────────────┬──────────────────────────┘
               ↑ implements
┌─────────────────────────────────────────┐
│    Infrastructure Layer (Adapter)       │  ← External Services
│      (HTTP Calls, APIs)                 │
└─────────────────────────────────────────┘
```

### Dependency Direction

**All dependencies point INWARD toward the Domain Layer.**

- UI → Domain (uses use cases)
- Use Cases → Ports (depends on interfaces)
- Adapters → Ports (implements interfaces)
- **Never**: Domain → Infrastructure or Domain → UI

This implements the **Dependency Inversion Principle** (SOLID), ensuring that:
- The domain is independent of external concerns
- External implementations can be swapped without changing business logic
- Testing is simplified (mock adapters can replace real ones)

---

## Benefits of This Architecture

### 🔧 Maintainability
- Clear separation of concerns
- Changes in one layer don't affect others
- Easy to locate and modify specific functionality
- Self-documenting code structure

### 🧪 Testability
- Domain logic is pure and framework-independent
- Easy to unit test use cases without UI or API
- Adapters can be mocked for testing
- UI components can be tested with mocked use cases

### 📈 Scalability
- Add new features without modifying existing code
- Feature modules are self-contained
- Shared components are reusable
- Parallel development on different features

### 🧩 Modularity
- Features are completely isolated modules
- Each feature has its own domain, infrastructure, and UI
- Can be extracted to separate packages if needed
- Clear boundaries between features

### 🔄 Flexibility
- Easy to swap implementations (e.g., REST API → GraphQL)
- UI framework can be changed without touching business logic
- Multiple adapters for the same port (mock, real, cached)
- A/B testing different implementations

### 👥 Team Collaboration
- Different developers can work on different layers simultaneously
- Clear contracts (ports) define team interfaces
- Less merge conflicts due to isolated features
- Easier code reviews with clear responsibilities

---

## Core Architectural Concepts

### Domain Layer

The **Domain Layer** is the heart of the application. It contains all business logic and rules, completely independent of any framework or external library.

**What belongs in the Domain Layer:**
- Business rules and validation logic
- Use cases (application workflows)
- Port interfaces (contracts for external services)
- Domain entities and value objects

**What does NOT belong:**
- React components or hooks
- HTTP calls or API logic
- UI state management
- Framework-specific code

**Directory structure:**
```
features/team/domain/
├── ports/
│   └── TeamRepository.ts       # Interface defining contract
└── usecases/
    ├── CreateTeamUseCase.ts    # Business logic for creating teams
    └── SearchTeamUseCase.ts    # Business logic for searching teams
```

**Why framework-independent?**
- Business logic should outlive UI frameworks
- Easier to test without framework overhead
- Can be reused in different contexts (mobile, desktop, server)
- Changes in framework don't require rewriting business rules

---

### Ports (Interfaces)

**What is a Port?**

A **Port** is an interface that defines a contract between the domain layer and external services. It specifies **what** operations are needed without defining **how** they are implemented.

**Why use Ports?**
- **Dependency Inversion**: Domain defines what it needs; infrastructure provides it
- **Testability**: Easy to create mock implementations
- **Flexibility**: Multiple implementations possible (real, mock, cached)
- **Clear contracts**: Explicit agreement between layers

**Where are Ports defined?**

In the **Domain Layer** (`features/[feature]/domain/ports/`)

**Example Port:**

```typescript
// features/team/domain/ports/TeamRepository.ts
import { Team, CreateTeamRequest } from '../../../../shared/types/team'

export interface TeamRepository {
  createTeam(request: CreateTeamRequest): Promise<{ data: Team; status: number }>
  findTeam(sport: string, teamName: string): Promise<{ data: Team; status: number }>
}
```

**Key characteristics:**
- Pure TypeScript interface
- No implementation details
- Defines method signatures and return types
- Located in domain layer

---

### Use Cases

**What is a Use Case?**

A **Use Case** represents a single business operation or workflow. It encapsulates the business logic required to accomplish a specific user goal.

**What belongs in a Use Case:**
- Business rules and validation
- Workflow orchestration
- Domain logic
- Error handling for business rules

**What does NOT belong:**
- UI logic or rendering
- HTTP calls (done by adapters)
- Framework-specific code
- Data formatting for display

**How does a Use Case receive dependencies?**

Through **constructor injection**:

```typescript
// features/team/domain/usecases/CreateTeamUseCase.ts
export class CreateTeamUseCase {
  constructor(private readonly teamRepository: TeamRepository) {}
  
  async execute(request: CreateTeamRequest): Promise<Result> {
    // Business validation
    if (!request.name || request.name.trim().length === 0) {
      return { success: false, error: 'Team name is required' }
    }
    
    // Delegate to repository
    const response = await this.teamRepository.createTeam(request)
    
    // Business rule: Check status
    if (response.status === 201) {
      return { success: true, team: response.data }
    }
    
    return { success: false, error: 'Failed to create team' }
  }
}
```

**Key characteristics:**
- Single responsibility (one business operation)
- Depends only on ports (interfaces), not concrete implementations
- Returns domain results, not HTTP responses
- Pure business logic

---

### Infrastructure Layer

**What is the Infrastructure Layer?**

The **Infrastructure Layer** contains concrete implementations of ports. This is where the application interacts with the real world: APIs, databases, external services, etc.

**What belongs in Infrastructure:**
- HTTP client implementations
- API adapters
- Database access code
- External service integrations
- Framework-specific implementations

**Why is tight coupling acceptable here?**

Because this layer's **purpose** is to couple with external services. However, this coupling is **isolated** to this layer and doesn't affect the domain.

**Directory structure:**
```
features/team/infrastructure/
└── adapters/
    └── TeamApiAdapter.ts       # Implements TeamRepository
```

---

### Adapters

**What is an Adapter?**

An **Adapter** is a concrete implementation of a Port. It "adapts" an external service (like a REST API) to match the interface defined by the domain.

**Where are Adapters located?**

In the **Infrastructure Layer** (`features/[feature]/infrastructure/adapters/`)

**Example Adapter:**

```typescript
// features/team/infrastructure/adapters/TeamApiAdapter.ts
import { TeamRepository } from '../../domain/ports/TeamRepository'

export class TeamApiAdapter implements TeamRepository {
  async createTeam(request: CreateTeamRequest): Promise<Response> {
    // Real HTTP call
    const response = await fetch('/api/team', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    })
    
    const data = await response.json()
    return { data, status: response.status }
  }
  
  async findTeam(sport: string, teamName: string): Promise<Response> {
    const response = await fetch(`/api/sport/${sport}/team/${teamName}`)
    const data = await response.json()
    return { data, status: response.status }
  }
}
```

**Can Adapters use external libraries?**

**Yes!** Adapters are in the infrastructure layer and can:
- Use fetch, axios, or any HTTP library
- Depend on external SDKs
- Use framework-specific features
- Be tightly coupled to APIs

**How do Adapters enable testing?**

You can create **mock adapters** for testing:

```typescript
export class MockTeamAdapter implements TeamRepository {
  async createTeam(request: CreateTeamRequest): Promise<Response> {
    // Return fake data
    return {
      data: { Name: request.name, Sport: request.sport, ...},
      status: 201
    }
  }
}
```

Then inject the mock instead of the real adapter in tests.

---

### Context Layer (Dependency Injection)

**What is the Context Layer?**

The **Context Layer** is responsible for:
1. Creating instances of adapters
2. Injecting adapters into use cases
3. Providing use cases to UI components via React Context

**Main responsibility:**

**Wiring dependencies** (Dependency Injection)

**Example Context Provider:**

```typescript
// features/team/context/TeamContext.tsx
import React, { createContext, useContext } from 'react'
import { CreateTeamUseCase } from '../domain/usecases/CreateTeamUseCase'
import { SearchTeamUseCase } from '../domain/usecases/SearchTeamUseCase'
import { TeamApiAdapter } from '../infrastructure/adapters/TeamApiAdapter'

interface TeamContextType {
  createTeamUseCase: CreateTeamUseCase
  searchTeamUseCase: SearchTeamUseCase
}

const TeamContext = createContext<TeamContextType | undefined>(undefined)

export function TeamProvider({ children }) {
  // 1. Create adapter (Infrastructure)
  const teamApiAdapter = new TeamApiAdapter()
  
  // 2. Inject adapter into use cases (Dependency Injection)
  const createTeamUseCase = new CreateTeamUseCase(teamApiAdapter)
  const searchTeamUseCase = new SearchTeamUseCase(teamApiAdapter)
  
  // 3. Provide to UI components
  return (
    <TeamContext.Provider value={{ createTeamUseCase, searchTeamUseCase }}>
      {children}
    </TeamContext.Provider>
  )
}

// Custom hook for UI components
export function useTeamContext(): TeamContextType {
  const context = useContext(TeamContext)
  if (!context) {
    throw new Error('useTeamContext must be used within TeamProvider')
  }
  return context
}
```

**How do UI components access the Context?**

Via the custom hook:

```typescript
function CreateTeamPage() {
  const { createTeamUseCase } = useTeamContext()
  
  const handleSubmit = async (data) => {
    const result = await createTeamUseCase.execute(data)
    // Handle result
  }
}
```

---

### Feature Module

**What is a Feature Module?**

A **Feature Module** is a wrapper component that encapsulates an entire feature's dependencies and context.

**What does it wrap?**

The feature's Context Provider:

```typescript
// features/team/TeamModule.tsx
import { TeamProvider } from './context/TeamContext'

export function TeamModule({ children }) {
  return <TeamProvider>{children}</TeamProvider>
}
```

**Why do we need Feature Modules?**

1. **Modularity**: Entire feature is self-contained
2. **Reusability**: Can be imported anywhere in the app
3. **Isolation**: Feature's dependencies are encapsulated
4. **Clarity**: Clear boundary of what belongs to the feature

**Where is it imported?**

At the application root or layout level:

```typescript
// App.tsx
import { TeamModule } from './features/team/TeamModule'

function App() {
  return (
    <TeamModule>
      {/* App routes and pages */}
    </TeamModule>
  )
}
```

---

## UI Layer and Atomic Design

### What is Atomic Design?

**Atomic Design** is a methodology for creating design systems by organizing UI components into a hierarchy based on their complexity, similar to atoms forming molecules in chemistry.

**Why use Atomic Design?**

- **Consistency**: Reusable components ensure UI consistency
- **Scalability**: Easy to build complex UIs from simple parts
- **Maintainability**: Small, focused components are easier to maintain
- **Reusability**: Lower-level components used across features
- **Team collaboration**: Clear naming and organization

### Component Hierarchy

```
Atoms → Molecules → Organisms → Objects → Pages
(smallest)                           (largest)
```

---

### Atoms

**What are Atoms?**

Atoms are the **basic building blocks** of the UI. They are the smallest, most fundamental components that cannot be broken down further without losing their meaning.

**Characteristics:**
- Single-purpose components
- No business logic
- Highly reusable
- Usually wrap or style basic HTML elements

**Examples:**
- `SportSelect` - Dropdown for selecting a sport
- `CategorySelect` - Dropdown for selecting a category
- `Button` - Styled button component
- `TextField` - Input field
- `Label` - Text label

**Location:** `shared/components/atoms/`

**Can Atoms depend on other Atoms?**

Generally **no**. Atoms should be independent. If an atom needs another atom, consider if it's actually a **molecule**.

**Should Atoms contain business logic?**

**No.** Atoms should be pure presentation components. Business logic belongs in use cases.

**Example Atom:**

```typescript
// shared/components/atoms/SportSelect.tsx
import { TextField, MenuItem } from '@mui/material'
import { Sport } from '../../types/team'

interface SportSelectProps {
  value: Sport
  onChange: (value: Sport) => void
  required?: boolean
  sports: Sport[]
}

export function SportSelect({ value, onChange, required, sports }: SportSelectProps) {
  return (
    <TextField
      select
      label="Deporte"
      value={value}
      onChange={(e) => onChange(e.target.value as Sport)}
      required={required}
      fullWidth
    >
      {sports.map((sport) => (
        <MenuItem key={sport} value={sport}>
          {sport}
        </MenuItem>
      ))}
    </TextField>
  )
}
```

---

### Molecules

**What are Molecules?**

Molecules are **combinations of atoms** that work together as a unit. They are still relatively simple but more functional than atoms.

**How do Molecules differ from Atoms?**

- **Atoms**: Single-purpose, minimal functionality
- **Molecules**: Combination of atoms with specific purpose

**Examples:**
- Search bar (input atom + button atom)
- Form field with label and error message
- Card header (title atom + icon atom + button atom)

**Location:** `shared/components/molecules/`

**When to create a Molecule vs Atom?**

Create a **Molecule** when:
- You're combining multiple atoms
- The combination has a specific, reusable purpose
- The combination appears in multiple places

Keep as **Atoms** when:
- It's a single UI element
- It doesn't need other components to be useful

---

### Organisms

**What are Organisms?**

Organisms are **complex UI components** composed of molecules and/or atoms. They form distinct sections of an interface.

**Complexity level:**

Organisms are more complex than molecules and can:
- Have their own state
- Contain multiple molecules
- Represent a complete section of UI

**Examples:**
- Navigation bar (logo + menu items + user avatar)
- Data table (headers + rows + pagination)
- Modal dialog (header + body + footer + actions)

**Location:** `shared/components/organisms/` or `features/[feature]/ui/organisms/`

**Can Organisms use Context?**

**Yes**, organisms can access context for shared state, but they should **not** contain business logic or call use cases directly. That's the responsibility of pages.

---

### Objects (Featured Objects)

**What are Objects?**

Objects are **feature-specific components** that may contain business logic related to that feature. They are similar to organisms but scoped to a specific feature.

**How do Objects differ from Organisms?**

- **Organisms**: Shared, reusable across features, no business logic
- **Objects**: Feature-specific, may contain feature logic

**Are Objects feature-specific or shared?**

**Feature-specific.** They live within a feature's directory structure.

**Location:** `features/[feature]/ui/objects/`

**Example:**
- `TeamSearchForm` - Form specific to team searching
- `TeamStatisticsCard` - Card showing team statistics
- `PlayerList` - List of players for a team

---

### Pages

**What are Pages?**

Pages are **complete views** that represent a route in the application. They orchestrate organisms, molecules, and atoms to create a full user experience.

**Role of a Page:**
- Top-level component for a route
- Orchestrates smaller components
- **Uses use cases from Context**
- Handles user interactions
- Manages page-level state

**Do Pages use Use Cases from Context?**

**Yes!** This is one of the key responsibilities of pages:

```typescript
export function CreateTeamPage() {
  const { createTeamUseCase } = useTeamContext()
  
  const handleSubmit = async (data) => {
    const result = await createTeamUseCase.execute(data)
    // Handle result, show success/error
  }
  
  return (
    // Render UI
  )
}
```

**Location:** `features/[feature]/ui/pages/`

**What should NOT be in a Page?**
- Business logic (belongs in use cases)
- HTTP calls (belongs in adapters)
- Reusable components (extract to organisms/molecules/atoms)
- Complex validation logic (belongs in use cases)

**Example Page:**

```typescript
// features/team/ui/pages/SearchTeamPage.tsx
export function SearchTeamPage() {
  const { searchTeamUseCase } = useTeamContext()
  const [team, setTeam] = useState<Team | null>(null)
  
  const handleSearch = async (sport: string, name: string) => {
    const result = await searchTeamUseCase.execute(sport, name)
    if (result.success) {
      setTeam(result.team)
    }
  }
  
  return (
    <Box>
      <SearchForm onSubmit={handleSearch} />
      {team && <TeamDetails team={team} />}
    </Box>
  )
}
```

---

## Directory Structure

### Complete Feature Structure

```
features/
└── team/                               # Team Feature
    ├── domain/                         # ✅ Pure Business Logic
    │   ├── ports/
    │   │   └── TeamRepository.ts       # Interface (contract)
    │   └── usecases/
    │       ├── CreateTeamUseCase.ts    # Create team logic
    │       └── SearchTeamUseCase.ts    # Search team logic
    ├── infrastructure/                 # ✅ External Implementations
    │   └── adapters/
    │       └── TeamApiAdapter.ts       # API implementation
    ├── context/                        # ✅ Dependency Injection
    │   └── TeamContext.tsx             # DI container
    ├── ui/                             # ✅ UI Components
    │   ├── pages/
    │   │   ├── CreateTeamPage.tsx      # Create page
    │   │   └── SearchTeamPage.tsx      # Search page
    │   ├── objects/                    # (future) Feature objects
    │   └── organisms/                  # (future) Complex components
    └── TeamModule.tsx                  # ✅ Feature wrapper
```

### Shared Structure

```
shared/
├── components/                         # Reusable components
│   ├── atoms/                          # Basic components
│   │   ├── SportSelect.tsx
│   │   ├── CategorySelect.tsx
│   │   └── ...
│   ├── molecules/                      # Combined atoms
│   │   └── (future)
│   └── organisms/                      # Complex shared components
│       └── (future)
└── types/                              # Shared TypeScript types
    └── team.ts
```

### Naming Conventions

**Files:**
- Use PascalCase for component files: `CreateTeamPage.tsx`
- Use camelCase for non-component files: `teamRepository.ts`
- Add `.tsx` extension for files with JSX
- Add `.ts` extension for pure TypeScript

**Directories:**
- Use lowercase with hyphens: `use-cases/` (if multi-word)
- Use singular for most directories: `context/`, `domain/`
- Use plural for collections: `usecases/`, `adapters/`, `pages/`

**Components:**
- PascalCase: `CreateTeamPage`, `SportSelect`
- Descriptive names indicating purpose
- Suffix with type if helpful: `CreateTeamUseCase`, `TeamApiAdapter`

---

## Adding a New Feature

### Step-by-Step Guide

**Example:** Adding a "Player" feature

#### Step 1: Create Directory Structure

```bash
mkdir -p src/features/player/{domain/{ports,usecases},infrastructure/adapters,context,ui/pages}
```

#### Step 2: Define Port (Interface)

```typescript
// features/player/domain/ports/PlayerRepository.ts
export interface PlayerRepository {
  createPlayer(data: CreatePlayerRequest): Promise<Player>
  findPlayer(id: string): Promise<Player>
}
```

#### Step 3: Create Use Case(s)

```typescript
// features/player/domain/usecases/CreatePlayerUseCase.ts
export class CreatePlayerUseCase {
  constructor(private readonly playerRepository: PlayerRepository) {}
  
  async execute(data: CreatePlayerRequest): Promise<Result> {
    // Business validation
    if (!data.name) {
      return { success: false, error: 'Name required' }
    }
    
    // Delegate to repository
    const player = await this.playerRepository.createPlayer(data)
    return { success: true, player }
  }
}
```

#### Step 4: Implement Adapter

```typescript
// features/player/infrastructure/adapters/PlayerApiAdapter.ts
export class PlayerApiAdapter implements PlayerRepository {
  async createPlayer(data: CreatePlayerRequest): Promise<Player> {
    const response = await fetch('/api/player', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    return response.json()
  }
}
```

#### Step 5: Setup Context with DI

```typescript
// features/player/context/PlayerContext.tsx
export function PlayerProvider({ children }) {
  const playerApiAdapter = new PlayerApiAdapter()
  const createPlayerUseCase = new CreatePlayerUseCase(playerApiAdapter)
  
  return (
    <PlayerContext.Provider value={{ createPlayerUseCase }}>
      {children}
    </PlayerContext.Provider>
  )
}
```

#### Step 6: Build UI Components

```typescript
// features/player/ui/pages/CreatePlayerPage.tsx
export function CreatePlayerPage() {
  const { createPlayerUseCase } = usePlayerContext()
  
  const handleSubmit = async (data) => {
    const result = await createPlayerUseCase.execute(data)
    // Handle result
  }
  
  return <form onSubmit={handleSubmit}>...</form>
}
```

#### Step 7: Create Feature Module

```typescript
// features/player/PlayerModule.tsx
export function PlayerModule({ children }) {
  return <PlayerProvider>{children}</PlayerProvider>
}
```

#### Step 8: Import in App

```typescript
// App.tsx
import { PlayerModule } from './features/player/PlayerModule'

function App() {
  return (
    <PlayerModule>
      <TeamModule>
        {/* Routes */}
      </TeamModule>
    </PlayerModule>
  )
}
```

---

## Testing Strategy

### Domain Layer Testing

**Ports:**
- No testing needed (interfaces don't have implementation)

**Use Cases:**
- **Unit tests** with mocked ports
- Test business logic thoroughly
- Test all validation rules
- Test error handling

**Example:**

```typescript
describe('CreateTeamUseCase', () => {
  it('should return error if name is empty', async () => {
    const mockRepo = new MockTeamRepository()
    const useCase = new CreateTeamUseCase(mockRepo)
    
    const result = await useCase.execute({ name: '', sport: 'Football' })
    
    expect(result.success).toBe(false)
    expect(result.error).toBe('Team name is required')
  })
})
```

### Infrastructure Layer Testing

**Adapters:**
- **Integration tests** with mocked HTTP
- Use libraries like `msw` or `nock` to mock API calls
- Test error handling
- Test response parsing

**Example:**

```typescript
describe('TeamApiAdapter', () => {
  it('should call correct API endpoint', async () => {
    const adapter = new TeamApiAdapter()
    
    // Mock fetch
    global.fetch = jest.fn(() =>
      Promise.resolve({
        json: () => Promise.resolve({ Name: 'Test Team' }),
        status: 201
      })
    )
    
    await adapter.createTeam({ name: 'Test', sport: 'Football' })
    
    expect(fetch).toHaveBeenCalledWith('/api/team', expect.any(Object))
  })
})
```

### UI Layer Testing

**Atoms:**
- **Component tests** (React Testing Library)
- Test rendering with different props
- Test user interactions
- No business logic to test

**Pages:**
- **Component tests** with mocked use cases
- Test user workflows
- Test error states
- Test loading states

**Example:**

```typescript
describe('CreateTeamPage', () => {
  it('should call use case on form submit', async () => {
    const mockUseCase = {
      execute: jest.fn(() => Promise.resolve({ success: true }))
    }
    
    render(
      <TeamContext.Provider value={{ createTeamUseCase: mockUseCase }}>
        <CreateTeamPage />
      </TeamContext.Provider>
    )
    
    // Fill form and submit
    fireEvent.click(screen.getByText('Create Team'))
    
    expect(mockUseCase.execute).toHaveBeenCalled()
  })
})
```

---

## References

### Architecture Patterns

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) by Alistair Cockburn
- [The Atomic Hexagonal Architecture — on the Frontend — with React](https://newlight77.medium.com/the-atomic-hexagonal-architecture-on-the-frontend-with-react-6337a56e56e3) by Kong To
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) by Robert C. Martin

### Design Methodologies

- [Atomic Design](https://atomicdesign.bradfrost.com/) by Brad Frost
- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html) by Martin Fowler

### SOLID Principles

- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID) on Wikipedia
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)

### React & TypeScript

- [React Documentation](https://react.dev/)
- [React Hooks](https://react.dev/reference/react/hooks)
- [TypeScript Documentation](https://www.typescriptlang.org/docs/)

---

**Last Updated:** November 2025  
**Version:** 1.0.0

