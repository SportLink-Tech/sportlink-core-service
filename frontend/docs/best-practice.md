# Frontend Best Practices

This document outlines coding standards, conventions, and best practices for the SportLink frontend application.

---

## Code Organization

### File Naming Conventions

#### Components (TSX files)
```
✅ Correct:
- CreateTeamPage.tsx
- SportSelect.tsx
- TeamApiAdapter.tsx

❌ Incorrect:
- createTeamPage.tsx
- sport-select.tsx
- teamApiAdapter.tsx
```

**Rule:** Use **PascalCase** for all component and class files.

#### Non-Component Files (TS files)
```
✅ Correct:
- teamRepository.ts
- teamTypes.ts
- utils.ts

❌ Incorrect:
- TeamRepository.ts
- team-types.ts
- Utils.ts
```

**Rule:** Use **camelCase** for non-component TypeScript files.

#### Directories
```
✅ Correct:
- features/team/
- shared/components/
- domain/usecases/

❌ Incorrect:
- Features/Team/
- Shared/Components/
- Domain/UseCases/
```

**Rule:** Use **lowercase** for directories. Use hyphens for multi-word directories (rare).

---

### Export/Import Conventions

#### Named Exports (Preferred)

```typescript
✅ Correct:
// SportSelect.tsx
export function SportSelect(props) { ... }

// Importing
import { SportSelect } from './components/SportSelect'

❌ Avoid:
// SportSelect.tsx
export default function SportSelect(props) { ... }
```

**Why named exports?**
- Better refactoring support
- Explicit imports
- Tree-shaking optimization
- Easier to find usages

#### Default Exports (When Appropriate)

Use default exports only for:
- Page components (main export of a page)
- Feature modules
- App entry point

```typescript
✅ Acceptable:
// App.tsx
export default function App() { ... }

// TeamModule.tsx
export default TeamModule
```

#### Barrel Exports (Index Files)

Create index files to simplify imports:

```typescript
// features/team/domain/usecases/index.ts
export { CreateTeamUseCase } from './CreateTeamUseCase'
export { SearchTeamUseCase } from './SearchTeamUseCase'

// Usage:
import { CreateTeamUseCase, SearchTeamUseCase } from '../domain/usecases'
```

---

### Directory Structure Conventions

#### Feature Organization

```
features/[feature-name]/
├── domain/              # Pure business logic
│   ├── ports/
│   └── usecases/
├── infrastructure/      # External implementations
│   └── adapters/
├── context/             # Dependency injection
├── ui/                  # UI components
│   ├── pages/
│   ├── objects/
│   └── organisms/
└── [Feature]Module.tsx  # Feature wrapper
```

**Naming:**
- Feature directory: lowercase, singular (e.g., `team/` not `teams/`)
- Subdirectories: plural if they contain multiple items (`usecases/`, `adapters/`)

#### Shared Components Organization

```
shared/
├── components/
│   ├── atoms/           # Basic reusable components
│   ├── molecules/       # Combined atoms
│   └── organisms/       # Complex shared components
├── types/               # Shared TypeScript types
└── utils/               # Utility functions (future)
```

---

## Domain Layer Best Practices

### Keep Domain Logic Pure

```typescript
✅ Correct - Pure domain logic:
export class CreateTeamUseCase {
  constructor(private readonly repository: TeamRepository) {}
  
  async execute(request: CreateTeamRequest): Promise<Result> {
    // Pure business validation
    if (!request.name || request.name.trim().length === 0) {
      return { success: false, error: 'Name required' }
    }
    
    // Use repository through port
    const response = await this.repository.createTeam(request)
    return { success: true, data: response.data }
  }
}

❌ Incorrect - Framework dependencies:
export class CreateTeamUseCase {
  async execute(request: CreateTeamRequest) {
    // ❌ Don't import React
    const [state, setState] = useState()
    
    // ❌ Don't make HTTP calls directly
    const response = await fetch('/api/team', {...})
    
    // ❌ Don't access DOM
    document.getElementById('form').value = ''
  }
}
```

**Rules:**
- ❌ No React imports in domain layer
- ❌ No HTTP calls (use adapters)
- ❌ No DOM access
- ❌ No UI libraries
- ✅ Only TypeScript and domain types

---

### Use Dependency Injection

```typescript
✅ Correct - Constructor injection:
export class CreateTeamUseCase {
  constructor(private readonly repository: TeamRepository) {}
  
  async execute(data) {
    return this.repository.createTeam(data)
  }
}

❌ Incorrect - Direct instantiation:
export class CreateTeamUseCase {
  async execute(data) {
    // ❌ Don't create dependencies inside
    const repository = new TeamApiAdapter()
    return repository.createTeam(data)
  }
}
```

**Why?**
- Testability: Easy to inject mocks
- Flexibility: Swap implementations
- Decoupling: Use case doesn't know about concrete adapter

---

### Avoid Coupling

```typescript
✅ Correct - Depend on interface:
export class CreateTeamUseCase {
  constructor(private readonly repository: TeamRepository) {}
  //                                      ^^^^^^^^^^^^^^
  //                                      Interface (Port)
}

❌ Incorrect - Depend on concrete class:
export class CreateTeamUseCase {
  constructor(private readonly repository: TeamApiAdapter) {}
  //                                      ^^^^^^^^^^^^^^^
  //                                      Concrete class
}
```

**Rule:** Domain should depend only on **interfaces (ports)**, never concrete implementations.

---

## Infrastructure Layer Best Practices

### Adapters Can Use External Libraries

```typescript
✅ Correct - Adapter uses fetch:
export class TeamApiAdapter implements TeamRepository {
  async createTeam(request) {
    // ✅ OK to use fetch in adapter
    const response = await fetch('/api/team', {
      method: 'POST',
      body: JSON.stringify(request)
    })
    return response.json()
  }
}

// ✅ Could also use axios if needed:
import axios from 'axios'

export class TeamAxiosAdapter implements TeamRepository {
  async createTeam(request) {
    const response = await axios.post('/api/team', request)
    return response.data
  }
}
```

**Rules:**
- ✅ Adapters can import external libraries
- ✅ Can use framework-specific features
- ✅ Tight coupling is OK here (it's isolated)

---

### Handle Errors in Adapters

```typescript
✅ Correct - Adapter handles HTTP errors:
export class TeamApiAdapter implements TeamRepository {
  async createTeam(request) {
    const response = await fetch('/api/team', {...})
    
    if (!response.ok) {
      const error = await response.json()
      throw new Error(error.message || 'API error')
    }
    
    return response.json()
  }
}

❌ Incorrect - Ignore errors:
export class TeamApiAdapter implements TeamRepository {
  async createTeam(request) {
    const response = await fetch('/api/team', {...})
    // ❌ What if response is 400 or 500?
    return response.json()
  }
}
```

**Rules:**
- ✅ Check response status
- ✅ Parse and throw meaningful errors
- ✅ Let errors propagate to use case

---

## UI Layer Best Practices

### Keep Components Small and Focused

```typescript
✅ Correct - Single responsibility:
export function SportSelect({ value, onChange, sports }) {
  return (
    <TextField select value={value} onChange={onChange}>
      {sports.map(sport => (
        <MenuItem key={sport} value={sport}>{sport}</MenuItem>
      ))}
    </TextField>
  )
}

❌ Incorrect - Too many responsibilities:
export function TeamForm() {
  // ❌ Handles form state
  const [sport, setSport] = useState()
  const [name, setName] = useState()
  const [players, setPlayers] = useState([])
  
  // ❌ Makes API calls
  const handleSubmit = async () => {
    await fetch('/api/team', {...})
  }
  
  // ❌ Validation logic
  const validate = () => {...}
  
  // ❌ 300 lines of JSX
  return <form>...</form>
}
```

**Rules:**
- ✅ One component, one purpose
- ✅ Extract reusable parts to atoms/molecules
- ✅ Keep components under 200 lines
- ✅ Use cases handle business logic, not components

---

### Don't Put Business Logic in Components

```typescript
✅ Correct - Business logic in use case:
export function CreateTeamPage() {
  const { createTeamUseCase } = useTeamContext()
  
  const handleSubmit = async (data) => {
    // ✅ Use case handles validation and logic
    const result = await createTeamUseCase.execute(data)
    if (result.success) {
      // Just handle UI feedback
      showSuccessMessage()
    }
  }
}

❌ Incorrect - Business logic in component:
export function CreateTeamPage() {
  const handleSubmit = async (data) => {
    // ❌ Validation in component
    if (!data.name || data.name.length < 3) {
      return setError('Name too short')
    }
    
    // ❌ Business rules in component
    if (data.players.length > 10) {
      return setError('Too many players')
    }
    
    // ❌ HTTP calls in component
    await fetch('/api/team', {...})
  }
}
```

**Rules:**
- ❌ No validation rules in components
- ❌ No HTTP calls in components
- ❌ No business rules in components
- ✅ Only UI state and presentation logic

---

### Use TypeScript Types Properly

```typescript
✅ Correct - Typed props:
interface SportSelectProps {
  value: Sport
  onChange: (value: Sport) => void
  sports: Sport[]
  required?: boolean
}

export function SportSelect({ 
  value, 
  onChange, 
  sports, 
  required = false 
}: SportSelectProps) {
  // Implementation
}

❌ Incorrect - Any types:
export function SportSelect(props: any) {
  // ❌ No type safety
}
```

**Rules:**
- ✅ Define interface for component props
- ✅ Type function parameters and returns
- ✅ Use specific types, avoid `any`
- ✅ Optional props should have default values

---

### Handle Errors Gracefully

```typescript
✅ Correct - Show user-friendly errors:
export function CreateTeamPage() {
  const [error, setError] = useState<string | null>(null)
  
  const handleSubmit = async (data) => {
    setError(null)
    const result = await createTeamUseCase.execute(data)
    
    if (!result.success) {
      setError(result.error || 'An error occurred')
      return
    }
    
    showSuccess()
  }
  
  return (
    <>
      {error && <Alert severity="error">{error}</Alert>}
      <form onSubmit={handleSubmit}>...</form>
    </>
  )
}

❌ Incorrect - Ignore errors or crash:
export function CreateTeamPage() {
  const handleSubmit = async (data) => {
    // ❌ Uncaught promise rejection
    const result = await createTeamUseCase.execute(data)
    showSuccess() // ❌ What if it failed?
  }
}
```

**Rules:**
- ✅ Always handle async errors
- ✅ Show user-friendly messages
- ✅ Use try/catch or .catch()
- ✅ Never let errors silently fail

---

### Reuse Atoms Across Features

```typescript
✅ Correct - Shared atoms:
// shared/components/atoms/SportSelect.tsx
export function SportSelect({ value, onChange }) {
  // Reusable across all features
}

// Used in team feature
import { SportSelect } from '../../../shared/components/atoms/SportSelect'

// Used in player feature
import { SportSelect } from '../../../shared/components/atoms/SportSelect'

❌ Incorrect - Duplicate atoms:
// features/team/ui/atoms/SportSelect.tsx
export function SportSelect() { ... }

// features/player/ui/atoms/SportSelect.tsx
export function SportSelect() { ... } // ❌ Duplication!
```

**Rules:**
- ✅ Atoms go in `shared/components/atoms/`
- ✅ Feature-specific components go in `features/[feature]/ui/`
- ✅ If used in 2+ features → move to shared

---

## What to Avoid

### ❌ Don't Import Infrastructure in Domain

```typescript
❌ WRONG:
// features/team/domain/usecases/CreateTeamUseCase.ts
import { TeamApiAdapter } from '../../infrastructure/adapters/TeamApiAdapter'
//                            ^^^^^^^^^^^^^^^^^^ WRONG LAYER!

export class CreateTeamUseCase {
  private adapter = new TeamApiAdapter() // ❌ Creates concrete dependency
}

✅ CORRECT:
// features/team/domain/usecases/CreateTeamUseCase.ts
import { TeamRepository } from '../ports/TeamRepository'
//                            ^^^^^^^^^^ Port (interface)

export class CreateTeamUseCase {
  constructor(private readonly repository: TeamRepository) {}
  //                                      ^^^^^^^^^^^^^^^^^^ Interface!
}
```

---

### ❌ Don't Skip the Context Layer

```typescript
❌ WRONG:
// App.tsx
import { TeamApiAdapter } from './features/team/infrastructure/...'
import { CreateTeamUseCase } from './features/team/domain/...'

// ❌ Wiring dependencies manually everywhere
const adapter = new TeamApiAdapter()
const useCase = new CreateTeamUseCase(adapter)

✅ CORRECT:
// App.tsx
import { TeamModule } from './features/team/TeamModule'

function App() {
  return (
    <TeamModule>  {/* ✅ Context handles DI internally */}
      {/* Routes */}
    </TeamModule>
  )
}
```

**Why?**
- Context centralizes dependency wiring
- Easy to swap implementations (mock for tests)
- Clear injection point

---

### ❌ Don't Mix Concerns

```typescript
❌ WRONG - Everything in one place:
// CreateTeamPage.tsx
export function CreateTeamPage() {
  // ❌ Validation logic
  const validateTeam = (data) => {
    if (!data.name) return 'Name required'
    if (data.players.length > 10) return 'Too many'
  }
  
  // ❌ HTTP calls
  const createTeam = async (data) => {
    await fetch('/api/team', { method: 'POST', ... })
  }
  
  // ❌ Complex business rules
  const calculateCategory = (players) => {
    // 50 lines of logic
  }
}

✅ CORRECT - Separated concerns:
// Use Case handles business logic
// Adapter handles HTTP calls
// Component only handles UI

export function CreateTeamPage() {
  const { createTeamUseCase } = useTeamContext()
  
  const handleSubmit = async (data) => {
    const result = await createTeamUseCase.execute(data)
    // Only UI concerns
  }
}
```

---

### ❌ Don't Use Magic Numbers/Strings

```typescript
❌ WRONG:
if (response.status === 201) { ... }
if (sport === 'football') { ... }
if (category >= 0 && category <= 7) { ... }

✅ CORRECT:
const HTTP_STATUS = {
  CREATED: 201,
  OK: 200,
  BAD_REQUEST: 400
} as const

const SPORTS = {
  FOOTBALL: 'Football',
  PADDLE: 'Paddle',
  TENNIS: 'Tennis'
} as const

const CATEGORY_RANGE = {
  MIN: 0,
  MAX: 7
} as const

if (response.status === HTTP_STATUS.CREATED) { ... }
if (sport === SPORTS.FOOTBALL) { ... }
```

---

## Code Style

### Formatting

Use consistent formatting:

```typescript
✅ Correct:
export function Component({ value, onChange }: Props) {
  const [state, setState] = useState(initial)
  
  const handleClick = () => {
    setState(newValue)
  }
  
  return (
    <div>
      <Button onClick={handleClick}>
        Click Me
      </Button>
    </div>
  )
}
```

**Rules:**
- 2 spaces indentation
- Semicolons at end of statements
- Single quotes for strings
- Trailing commas in objects/arrays

### Naming

```typescript
✅ Correct naming:
// Components: PascalCase
export function CreateTeamPage() {}
export function SportSelect() {}

// Functions/variables: camelCase
const handleSubmit = () => {}
const teamName = 'Thunder'

// Constants: UPPER_SNAKE_CASE
const MAX_PLAYERS = 10
const API_BASE_URL = '/api'

// Types/Interfaces: PascalCase
interface TeamRepository {}
type Sport = 'Football' | 'Paddle'

// Private class fields: _prefixed (optional)
class UseCase {
  private _repository: Repository
}
```

---

### Comments

```typescript
✅ Good comments - Explain WHY:
// Normalize Members to empty array because backend returns null
const teamData = {
  ...response.data,
  Members: response.data.Members || []
}

// Apply business rule: Teams can't exceed 10 players
if (players.length > MAX_TEAM_SIZE) {
  return error
}

❌ Bad comments - Explain WHAT (code already says this):
// Set the team name
setTeamName(name)

// Loop through players
players.forEach(player => ...)
```

**Rules:**
- Comment **why**, not **what**
- Document business rules
- Explain workarounds
- Add TODO with ticket number

---

## Performance Best Practices

### Avoid Unnecessary Re-renders

```typescript
✅ Correct - Memoized callbacks:
const handleChange = useCallback((value: Sport) => {
  setSport(value)
}, [])

<SportSelect value={sport} onChange={handleChange} />

✅ Correct - Memoized values:
const expensiveValue = useMemo(() => {
  return calculateComplexValue(data)
}, [data])
```

### Lazy Load Routes (Future)

```typescript
✅ Correct - Code splitting:
const CreateTeamPage = lazy(() => import('./pages/CreateTeamPage'))

<Suspense fallback={<Loading />}>
  <Routes>
    <Route path="/create" element={<CreateTeamPage />} />
  </Routes>
</Suspense>
```

---

## Testing Best Practices

### Test Domain Logic Thoroughly

```typescript
✅ Correct - Test use cases:
describe('CreateTeamUseCase', () => {
  it('should validate required fields', async () => {
    const mockRepo = new MockTeamRepository()
    const useCase = new CreateTeamUseCase(mockRepo)
    
    const result = await useCase.execute({ name: '', sport: 'Football' })
    
    expect(result.success).toBe(false)
    expect(result.error).toBe('Team name is required')
  })
  
  it('should create team when valid', async () => {
    // Test successful path
  })
  
  it('should handle repository errors', async () => {
    // Test error handling
  })
})
```

**Rules:**
- ✅ Test all business rules
- ✅ Test happy path
- ✅ Test error cases
- ✅ Use mock repositories

---

### Test Components with Mocked Use Cases

```typescript
✅ Correct - Mock use cases in component tests:
describe('CreateTeamPage', () => {
  it('should call use case on submit', async () => {
    const mockUseCase = {
      execute: jest.fn(() => Promise.resolve({ success: true }))
    }
    
    render(
      <TeamContext.Provider value={{ createTeamUseCase: mockUseCase }}>
        <CreateTeamPage />
      </TeamContext.Provider>
    )
    
    fireEvent.submit(screen.getByRole('form'))
    
    expect(mockUseCase.execute).toHaveBeenCalled()
  })
})
```

---

## Summary Checklist

Before committing code, verify:

### Domain Layer
- [ ] No framework imports (React, etc.)
- [ ] Uses dependency injection
- [ ] Depends only on ports (interfaces)
- [ ] Business rules are clear

### Infrastructure Layer
- [ ] Implements ports correctly
- [ ] Handles errors properly
- [ ] Can use external libraries

### UI Layer
- [ ] Components are small and focused
- [ ] No business logic in components
- [ ] Types are properly defined
- [ ] Errors are handled gracefully

### General
- [ ] File naming follows conventions
- [ ] Imports are clean
- [ ] No code duplication
- [ ] Tests are written (for use cases)

---

## References

- [React Best Practices](https://react.dev/learn/thinking-in-react)
- [TypeScript Best Practices](https://www.typescriptlang.org/docs/handbook/declaration-files/do-s-and-don-ts.html)
- [Clean Code](https://www.oreilly.com/library/view/clean-code-a/9780136083238/)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)

---

**Last Updated:** November 2025  
**Version:** 1.0.0

