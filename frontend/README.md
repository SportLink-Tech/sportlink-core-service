# SportLink Frontend

Frontend application for SportLink Core Service built with React and TypeScript.

## Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Material UI (MUI)** - Component library and design system
- **Emotion** - CSS-in-JS styling (required by MUI)
- **Vite** - Build tool and dev server
- **ESLint** - Code linting

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

```bash
npm install
```

### Development

Run the development server:

```bash
npm run dev
```

The application will be available at `http://localhost:3000`

The dev server is configured with a proxy that forwards `/api/*` requests to the backend at `http://localhost:8080`

### Build

Create a production build:

```bash
npm run build
```

### Preview Production Build

Preview the production build locally:

```bash
npm run preview
```

### Linting

Run ESLint:

```bash
npm run lint
```

## Architecture

This project follows **Atomic Hexagonal Architecture**, combining:
- Hexagonal Architecture (Ports & Adapters)
- Atomic Design (UI Component Hierarchy)
- Domain-Driven Design (Use Cases)

See [ARCHITECTURE.md](./ARCHITECTURE.md) for detailed documentation.

## Project Structure

```
frontend/
├── src/
│   ├── features/                    # Feature modules
│   │   └── team/                    # Team feature
│   │       ├── domain/              # Business logic (pure)
│   │       │   ├── ports/           # Interfaces
│   │       │   └── usecases/        # Use cases
│   │       ├── infrastructure/      # External adapters
│   │       │   └── adapters/        # API implementations
│   │       ├── context/             # Dependency Injection
│   │       ├── ui/                  # UI Components (Atomic Design)
│   │       │   └── pages/           # Feature pages
│   │       └── TeamModule.tsx       # Feature wrapper
│   ├── shared/                      # Shared resources
│   │   ├── components/              # Atomic components
│   │   │   ├── atoms/               # Basic components
│   │   │   └── molecules/           # Combined atoms
│   │   └── types/                   # Shared types
│   ├── components/                  # Layout components
│   ├── App.tsx                      # Main app
│   └── main.tsx                     # Entry point
├── public/                          # Static assets
├── ARCHITECTURE.md                  # Architecture documentation
└── package.json                     # Dependencies
```

## API Integration

The Vite dev server is configured to proxy API requests:

- Frontend: `http://localhost:3000`
- API requests: `http://localhost:3000/api/*` → `http://localhost:8080/*`

Example:
```typescript
// This will call http://localhost:8080/team
fetch('/api/team')
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

