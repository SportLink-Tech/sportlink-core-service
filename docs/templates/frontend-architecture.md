# Frontend Architecture Documentation Template

This document provides the structure and instructions for generating comprehensive frontend architecture documentation.

---

## Important Notes

- **Output Location**: All generated documentation files must be placed in the `docs/frontend/` directory
- **File Naming Convention**: Use singular form for file names (e.g., `architecture.md`, `technology.md`, not `architectures.md`)
- **Language**: All documentation must be written in English
- **Organization**: Separate architectural concepts from technology/framework documentation

---

## 1. Architecture Overview Document

**Instructions:** Create a high-level architecture document that explains the architectural pattern used in the frontend.

**File**: `docs/frontend/architecture.md`

### Questions to Answer:

#### 1.1 What architectural pattern does the frontend follow?
- What is the name of the architecture?
- Why was this architecture chosen?
- What problem does it solve?
- How does it improve code organization?

#### 1.2 What are the main architectural layers?
- What are the core layers in this architecture?
- How do they interact with each other?
- What is the direction of dependencies between layers?
- Which layer is the most important (core)?

#### 1.3 What are the benefits of this architecture?
List concisely (bullet points):
- Maintainability benefits
- Testability benefits
- Scalability benefits
- Modularity benefits
- Flexibility benefits
- Team collaboration benefits

#### 1.4 Dependency Flow
- Draw or explain the dependency flow diagram
- Which direction do dependencies point?
- What principle is being applied? (e.g., Dependency Inversion)
- How does this prevent coupling?

---

## 2. Core Architectural Concepts

**Instructions:** For each core concept in the architecture, create detailed explanations.

### 2.1 Domain Layer

**File Section**: `docs/frontend/architecture.md` (section within the main document)

#### Questions to Answer:
- What is the Domain Layer?
- What belongs in the Domain Layer?
- Why should the Domain Layer be framework-independent?
- What are the subdirectories within the Domain Layer?

### 2.2 Ports (Interfaces)

**Questions to Answer:**
- What is a Port?
- Why do we use Ports?
- Where are Ports defined? (Which layer?)
- What do Ports contain? (Methods, interfaces, contracts?)
- Give an example of a Port interface
- How do Ports relate to the Dependency Inversion Principle?

### 2.3 Use Cases

**Questions to Answer:**
- What is a Use Case?
- What kind of logic belongs in a Use Case?
- What does a Use Case depend on?
- How does a Use Case receive its dependencies?
- Give an example of a Use Case implementation
- What should NOT be in a Use Case?

### 2.4 Infrastructure Layer

**Questions to Answer:**
- What is the Infrastructure Layer?
- What belongs in the Infrastructure Layer?
- Why is tight coupling acceptable here?
- What are the subdirectories within the Infrastructure Layer?

### 2.5 Adapters

**Questions to Answer:**
- What is an Adapter?
- What does an Adapter implement?
- Where are Adapters located? (Which layer?)
- Give an example of an Adapter implementation
- Can Adapters use external libraries/frameworks?
- How do Adapters enable testing? (Mock adapters)

### 2.6 Context Layer (Dependency Injection)

**Questions to Answer:**
- What is the Context Layer?
- What is its main responsibility?
- How does it implement Dependency Injection?
- What does the Context provide to UI components?
- Show an example of a Context Provider
- How do UI components access the Context?

### 2.7 Feature Module

**Questions to Answer:**
- What is a Feature Module?
- What does it wrap?
- Why do we need Feature Modules?
- How does it enforce modularity?
- Where is it imported in the application?

---

## 3. UI Layer and Atomic Design

**Instructions:** Explain the UI organization following Atomic Design principles.

**File Section**: `docs/frontend/architecture.md` (section within the main document)

### 3.1 What is Atomic Design?

**Questions to Answer:**
- What is Atomic Design methodology?
- Why do we use Atomic Design for UI components?
- What are the levels of component granularity?

### 3.2 Atoms

**Questions to Answer:**
- What are Atoms?
- Give examples of Atomic components
- Where are Atoms located in the directory structure?
- Can Atoms depend on other Atoms?
- Should Atoms contain business logic?
- Example: List 3-5 Atomic components

### 3.3 Molecules

**Questions to Answer:**
- What are Molecules?
- How do Molecules differ from Atoms?
- Give examples of Molecular components
- Where are Molecules located in the directory structure?
- When should you create a Molecule vs an Atom?

### 3.4 Organisms

**Questions to Answer:**
- What are Organisms?
- What level of complexity do they represent?
- Give examples of Organism components
- Where are Organisms located in the directory structure?
- Can Organisms use Context?

### 3.5 Objects (Featured Objects)

**Questions to Answer:**
- What are Objects in the context of features?
- How do Objects differ from Organisms?
- Are Objects feature-specific or shared?
- Where are Objects located in the directory structure?

### 3.6 Pages

**Questions to Answer:**
- What are Pages?
- What is the role of a Page component?
- Do Pages use Use Cases from Context?
- Where are Pages located in the directory structure?
- Give an example of a Page component
- What should NOT be in a Page component?

---

## 4. Directory Structure Documentation

**Instructions:** Document the complete directory structure with explanations.

**File Section**: `docs/frontend/architecture.md` (section within the main document)

### Questions to Answer:

#### 4.1 Feature Directory Structure
Explain the structure of a feature module:
```
features/
└── [feature-name]/
    ├── domain/
    ├── infrastructure/
    ├── context/
    ├── ui/
    └── [Feature]Module.tsx
```

For each subdirectory, explain:
- What goes in this directory?
- What are the naming conventions?
- What are example files?

#### 4.2 Shared Directory Structure
Explain what goes in the shared directory:
```
shared/
├── components/
│   ├── atoms/
│   ├── molecules/
│   └── organisms/
└── types/
```

#### 4.3 Complete Tree
Provide a complete directory tree with annotations for each level.

---

## 5. Technology Stack Document

**Instructions:** Create a separate document for technologies and frameworks.

**File**: `docs/frontend/technology.md`

### 5.1 Core Technologies

**Questions to Answer:**
- What programming language is used?
- What version of the language?
- Why was this language chosen?

### 5.2 Framework and Libraries

**Questions to Answer for each major technology:**

#### React
- What version?
- Why React?
- What React features are used? (Hooks, Context, etc.)

#### TypeScript
- What version?
- Why TypeScript over JavaScript?
- What are the benefits?

#### UI Framework (Material-UI)
- What UI framework is used?
- What version?
- Why was this framework chosen?
- What components does it provide?

#### Routing
- What routing library is used?
- What version?
- How is routing configured?

#### Build Tool (Vite)
- What build tool is used?
- Why Vite over other tools?
- What are the benefits?

### 5.3 Development Dependencies

**Questions to Answer:**
- What linter is used?
- What testing framework is used? (if applicable)
- What other dev tools are configured?

### 5.4 Package Management

**Questions to Answer:**
- What package manager is used? (npm, yarn, pnpm)
- What is the minimum Node.js version required?

---

## 6. Adding a New Feature

**Instructions:** Provide a step-by-step guide for adding a new feature following the architecture.

**File Section**: `docs/frontend/architecture.md` (section within the main document)

### Questions to Answer:

1. What are the steps to add a new feature?
2. In what order should you create files?
3. What naming conventions should be followed?
4. How do you register the feature in the app?
5. Provide a concrete example (step-by-step)

**Step-by-step template:**
1. Create feature directory structure
2. Define Port (interface)
3. Create Use Case(s)
4. Implement Adapter
5. Setup Context with Dependency Injection
6. Build UI Components (Atoms → Pages)
7. Create Feature Module
8. Import and use in App

---

## 7. Testing Strategy

**Instructions:** Document how different layers should be tested.

**File Section**: `docs/frontend/architecture.md` or separate `testing.md`

### Questions to Answer:

#### 7.1 Domain Layer Testing
- How should Ports be tested?
- How should Use Cases be tested?
- What should be mocked?
- Example test case

#### 7.2 Infrastructure Layer Testing
- How should Adapters be tested?
- Should you mock HTTP calls?
- Example test case

#### 7.3 UI Layer Testing
- How should Atoms be tested?
- How should Pages be tested?
- Should you mock Use Cases?
- Example test case

---

## 8. Best Practices

**Instructions:** Document best practices and conventions.

**File**: `docs/frontend/best-practice.md` (singular)

### Questions to Answer:

#### 8.1 Code Organization
- File naming conventions
- Directory naming conventions
- Export/Import conventions

#### 8.2 Domain Layer Best Practices
- Keep domain logic pure
- No framework dependencies
- Use dependency injection
- Avoid coupling

#### 8.3 UI Layer Best Practices
- Keep components small and focused
- Reuse Atoms across features
- Use TypeScript types properly
- Handle errors gracefully

#### 8.4 What to Avoid
- Don't put business logic in UI components
- Don't import infrastructure in domain
- Don't skip the Context layer
- Don't mix concerns

---

## 9. References and Resources

**Instructions:** Provide links to external resources and references.

**File Section**: All documentation files should include a references section

### Questions to Answer:

1. What articles/books inspired this architecture?
2. What are the official documentation links?
3. What are useful learning resources?
4. Internal documentation links

**Key References to Include:**
- Hexagonal Architecture article
- Atomic Design book/website
- Domain-Driven Design resources
- React official documentation
- Material-UI documentation
- TypeScript handbook

---

## 10. Documentation File Organization

The final documentation structure should be:

```
docs/
└── frontend/
    ├── architecture.md          # Main architecture document
    ├── technology.md            # Technologies and frameworks
    ├── best-practice.md         # Best practices and conventions
    └── testing.md               # Testing strategy (optional)
```

---

## Checklist for Complete Documentation

Before considering the frontend documentation complete, verify:

- [ ] `architecture.md` explains the Atomic Hexagonal Architecture
- [ ] All core concepts are explained (Ports, Adapters, Use Cases, Context)
- [ ] UI layer with Atomic Design is documented (Atoms, Molecules, etc.)
- [ ] Directory structure is clearly documented with examples
- [ ] Dependency flow diagram is included
- [ ] `technology.md` lists all technologies with versions
- [ ] Each technology has a justification (why it was chosen)
- [ ] Step-by-step guide for adding a new feature is provided
- [ ] Best practices document exists
- [ ] All files use singular naming convention
- [ ] All documentation is in English
- [ ] Internal links between documents work
- [ ] External references are provided with URLs

---

## Notes

- Keep explanations concise but complete
- Use diagrams where helpful (Mermaid syntax supported)
- Provide code examples for each concept
- Link related concepts together
- Update documentation when architecture evolves
- Use consistent terminology throughout all documents

