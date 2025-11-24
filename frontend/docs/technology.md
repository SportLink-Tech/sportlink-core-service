# Frontend Technology Stack

This document describes all technologies, frameworks, libraries, and tools used in the SportLink frontend application.

---

## Core Technologies

### TypeScript

**Version:** 5.2.2

**What is TypeScript?**

TypeScript is a strongly typed programming language that builds on JavaScript, adding static type definitions.

**Why TypeScript over JavaScript?**

1. **Type Safety**: Catch errors at compile-time rather than runtime
2. **Better IDE Support**: Autocomplete, refactoring, and navigation
3. **Self-Documenting**: Types serve as inline documentation
4. **Scalability**: Easier to maintain large codebases
5. **Refactoring**: Confidence when changing code
6. **Team Collaboration**: Clear interfaces and contracts

**Benefits for this project:**
- Interfaces define clear contracts (Ports in Hexagonal Architecture)
- Use cases have well-defined input/output types
- Reduced runtime errors
- Better developer experience

**Configuration:**
- Strict mode enabled
- ES2020 target
- React JSX support
- Path aliases configured

---

## UI Framework

### React

**Version:** 18.2.0

**What is React?**

React is a JavaScript library for building user interfaces, maintained by Meta (Facebook) and a community of developers.

**Why React?**

1. **Component-Based**: Perfect for Atomic Design methodology
2. **Declarative**: UI is a function of state
3. **Large Ecosystem**: Abundant libraries and tools
4. **Performance**: Virtual DOM for efficient updates
5. **Community**: Huge community, extensive resources
6. **Maturity**: Battle-tested in production at scale

**React Features Used:**

#### Hooks
- **useState**: Component state management
- **useContext**: Accessing dependency injection context
- **useEffect**: Side effects (future use)
- **Custom Hooks**: `useTeamContext()` for feature access

#### Context API
- Used for Dependency Injection
- Provides use cases to UI components
- Replaces prop drilling

#### JSX/TSX
- TypeScript + JSX for type-safe components
- `.tsx` extension for components

**Why React for Hexagonal Architecture?**
- Context API perfect for Dependency Injection
- Hooks enable clean separation of concerns
- Component model aligns with Atomic Design

---

## Component Library

### Material-UI (MUI)

**Version:** 5.15.0

**What is Material-UI?**

Material-UI is a comprehensive React component library implementing Google's Material Design.

**Why Material-UI?**

1. **Mature**: Industry-standard, well-maintained
2. **Complete**: 100+ ready-to-use components
3. **Customizable**: Theming system for brand identity
4. **Accessible**: Built-in accessibility (ARIA)
5. **TypeScript**: First-class TypeScript support
6. **Documentation**: Excellent docs and examples
7. **Community**: Large community, many resources

**Components Used:**
- **Layout**: Container, Box, Stack, Grid
- **Inputs**: TextField, Button, MenuItem, Chip
- **Feedback**: Snackbar, Alert, CircularProgress
- **Navigation**: AppBar, Toolbar, Tabs, Menu
- **Surfaces**: Card, Paper
- **Typography**: Typography component
- **Icons**: @mui/icons-material

**Theme Customization:**

```typescript
const theme = createTheme({
  palette: {
    primary: { main: '#00C853' },    // Green
    secondary: { main: '#6A1B9A' },  // Purple
  },
  typography: {
    fontFamily: 'Inter, Roboto, sans-serif',
  },
  shape: { borderRadius: 12 },
})
```

**Styling Approach:**
- **sx prop**: Inline styles with theme access
- **CSS-in-JS**: Emotion (required by MUI)
- **Theme**: Centralized design tokens

---

## Styling

### Emotion

**Version:** 11.11.3 (React), 11.11.0 (Styled)

**What is Emotion?**

Emotion is a CSS-in-JS library designed for writing css styles with JavaScript.

**Why Emotion?**

- **Required by Material-UI**: MUI uses Emotion internally
- **Performance**: Fast and lightweight
- **TypeScript Support**: Full type safety for styles
- **Theme Access**: Direct access to MUI theme

**Packages:**
- `@emotion/react`: Core Emotion library
- `@emotion/styled`: styled-components API

**Note:** We primarily use MUI's `sx` prop rather than styled components for consistency.

---

## Routing

### React Router

**Version:** 6.30.2

**What is React Router?**

React Router is the standard routing library for React applications.

**Why React Router?**

1. **Industry Standard**: Most popular React routing solution
2. **Declarative**: Route configuration with JSX
3. **Dynamic**: Programmatic navigation, route params
4. **Nested Routes**: Support for complex layouts
5. **TypeScript Support**: Full type definitions

**Features Used:**

- **BrowserRouter**: HTML5 history API
- **Routes & Route**: Declarative route configuration
- **useNavigate**: Programmatic navigation
- **useLocation**: Current location access

**Route Configuration:**

```typescript
<Routes>
  <Route path="/" element={<SearchTeamPage />} />
  <Route path="/create" element={<CreateTeamPage />} />
</Routes>
```

**Type Definitions:**
- `@types/react-router-dom`: TypeScript types

---

## Build Tool

### Vite

**Version:** 5.0.8

**What is Vite?**

Vite is a modern build tool that provides a faster development experience for modern web projects.

**Why Vite over Webpack/CRA?**

1. **Fast**: Instant server start with native ESM
2. **Hot Module Replacement**: Lightning-fast HMR
3. **Optimized Builds**: Rollup for production
4. **Simple Config**: Minimal configuration needed
5. **Modern**: Built for ES modules
6. **Plugin Ecosystem**: Rich plugin system

**Benefits:**
- Development server starts in milliseconds
- Changes reflect instantly in browser
- Smaller production bundles
- Better developer experience

**Configuration:**

```typescript
// vite.config.ts
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})
```

**Features Used:**
- **Dev Server**: Fast development server
- **Proxy**: API proxying to backend
- **HMR**: Hot Module Replacement
- **Build**: Production optimization

**Plugin:**
- `@vitejs/plugin-react`: React support with Fast Refresh

---

## Development Dependencies

### ESLint

**Version:** 8.57.1

**What is ESLint?**

ESLint is a static code analysis tool for identifying problematic patterns in JavaScript/TypeScript code.

**Configuration:**
- React plugin: Enforces React best practices
- TypeScript plugin: TypeScript-specific rules
- React Hooks plugin: Validates hooks usage

**Rules:**
- `react-refresh/only-export-components`: Ensures HMR compatibility
- TypeScript strict rules enabled

### TypeScript ESLint

**Plugins:**
- `@typescript-eslint/eslint-plugin`: 6.14.0
- `@typescript-eslint/parser`: 6.14.0

**Purpose:**
- TypeScript-aware linting
- Catches type-related issues
- Enforces TypeScript best practices

---

## Package Management

### npm

**Version:** 10.8.2 (as detected)

**Minimum Node.js Version:** 18+

**Why npm?**

- Comes bundled with Node.js
- Widely used and well-documented
- Compatible with all dependencies
- Lock file ensures consistency

**Scripts:**

```json
{
  "dev": "vite",                    // Start development server
  "build": "tsc && vite build",     // Build for production
  "preview": "vite preview",        // Preview production build
  "lint": "eslint . --ext ts,tsx"   // Run linter
}
```

**Alternative Package Managers:**

While npm is used, the project is compatible with:
- **yarn**: `yarn install`, `yarn dev`
- **pnpm**: `pnpm install`, `pnpm dev`

---

## API Communication

### Fetch API

**What is used?**

Native browser **Fetch API** for HTTP requests.

**Why Fetch over Axios?**

1. **Native**: Built into modern browsers
2. **No Dependencies**: Reduces bundle size
3. **Modern**: Promise-based, async/await support
4. **Sufficient**: Meets all current needs

**Location:**

Fetch calls are isolated in **Adapters** (Infrastructure layer):

```typescript
// features/team/infrastructure/adapters/TeamApiAdapter.ts
async createTeam(request: CreateTeamRequest) {
  const response = await fetch('/api/team', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  })
  return response.json()
}
```

**Future Consideration:**

If more advanced features are needed (interceptors, request cancellation, progress tracking), consider:
- **Axios**: Feature-rich HTTP client
- **React Query**: Data fetching + caching library

---

## Type Definitions

### @types Packages

TypeScript type definitions for libraries without built-in types:

- `@types/react`: 18.2.43
- `@types/react-dom`: 18.2.17
- `@types/react-router-dom`: 5.3.3

**Purpose:**
- Provide TypeScript types for JavaScript libraries
- Enable autocomplete and type checking
- Catch errors at compile time

---

## Development Tools

### Browser DevTools

**Recommended:**
- **Chrome DevTools**: React Developer Tools extension
- **Redux DevTools**: (future if state management added)

### Code Editor

**Recommended:** Visual Studio Code

**Extensions:**
- ESLint: Real-time linting
- Prettier: Code formatting
- TypeScript: Built-in support
- ES7+ React/Redux snippets: Code snippets

---

## Testing (Future)

While not currently implemented, the recommended testing stack would be:

### Jest
- **Purpose**: Test runner and assertion library
- **Why**: Industry standard for React testing

### React Testing Library
- **Purpose**: Component testing
- **Why**: Tests behavior, not implementation details

### Mock Service Worker (MSW)
- **Purpose**: API mocking
- **Why**: Mock API calls in tests without changing code

---

## Architecture-Specific Choices

### Why These Technologies Fit Hexagonal Architecture

#### TypeScript
- Interfaces define clear **Ports**
- Types ensure correct **Adapter** implementations
- Compile-time validation of dependencies

#### React + Context
- Context API perfect for **Dependency Injection**
- Hooks enable separation of concerns
- Component model supports **Atomic Design**

#### Material-UI
- Pre-built **Atoms** and **Molecules**
- Consistent design system
- Accessibility built-in

#### Vite
- Fast feedback loop for development
- Modern tooling for modern architecture
- Simple configuration

---

## Version Management

### Node.js Version

**Minimum:** 18.0.0  
**Recommended:** 20.x LTS

**Why Node 18+?**
- Fetch API available natively
- ES2022 features support
- Better performance
- Long-term support

### Dependency Updates

**Strategy:**
- **Minor/Patch**: Update regularly for bug fixes
- **Major**: Evaluate breaking changes before updating
- **Security**: Apply security patches immediately

**Tools:**
- `npm outdated`: Check for updates
- `npm audit`: Security vulnerabilities
- Dependabot: Automated dependency updates (if enabled)

---

## Performance Considerations

### Bundle Size Optimization

- **Tree Shaking**: Vite automatically removes unused code
- **Code Splitting**: React.lazy() for route-based splitting (future)
- **Material-UI**: Import only used components

### Development Performance

- **Vite HMR**: Instant feedback on changes
- **TypeScript**: Incremental compilation
- **ESLint**: Fast linting with caching

### Runtime Performance

- **React 18**: Concurrent features, automatic batching
- **Virtual DOM**: Efficient UI updates
- **Memoization**: useMemo/useCallback (when needed)

---

## Environment Configuration

### Development

```bash
npm run dev
```

- Port: 3000
- API Proxy: → localhost:8080
- Hot Reload: Enabled
- Source Maps: Enabled

### Production

```bash
npm run build
npm run preview
```

- Minified bundle
- Optimized assets
- Production React build
- No source maps

---

## Integration with Backend

### API Proxy Configuration

Vite proxies `/api/*` requests to the backend:

```
Frontend (localhost:3000)
    ↓ /api/team
Vite Proxy
    ↓ rewrites to /team
Backend (localhost:8080)
```

**Benefits:**
- No CORS issues in development
- Same-origin requests
- Transparent to frontend code

**Production:**

In production, both frontend and backend would typically:
- Be served from the same domain
- Or use CORS headers for cross-origin requests

---

## Summary

### Technology Stack Overview

| Category | Technology | Version | Purpose |
|----------|-----------|---------|---------|
| Language | TypeScript | 5.2.2 | Type-safe JavaScript |
| UI Framework | React | 18.2.0 | Component-based UI |
| Components | Material-UI | 5.15.0 | Design system |
| Styling | Emotion | 11.11.x | CSS-in-JS |
| Routing | React Router | 6.30.2 | Navigation |
| Build Tool | Vite | 5.0.8 | Dev server & bundler |
| Linting | ESLint | 8.57.1 | Code quality |
| Package Manager | npm | 10.8.2 | Dependencies |

### Why This Stack?

1. **Modern**: Latest stable versions of proven technologies
2. **Type-Safe**: TypeScript throughout
3. **Fast**: Vite for instant feedback
4. **Professional**: Material-UI for polished UI
5. **Scalable**: Architecture supports growth
6. **Maintainable**: Clear separation of concerns
7. **Testable**: All layers can be tested independently

---

## References

### Official Documentation

- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [React Documentation](https://react.dev/)
- [Material-UI Documentation](https://mui.com/)
- [Emotion Documentation](https://emotion.sh/docs/introduction)
- [React Router Documentation](https://reactrouter.com/)
- [Vite Documentation](https://vitejs.dev/)
- [ESLint Documentation](https://eslint.org/)

### Learning Resources

- [React TypeScript Cheatsheet](https://react-typescript-cheatsheet.netlify.app/)
- [Material-UI Templates](https://mui.com/material-ui/getting-started/templates/)
- [Vite Guide](https://vitejs.dev/guide/)

---

**Last Updated:** November 2025  
**Version:** 1.0.0

