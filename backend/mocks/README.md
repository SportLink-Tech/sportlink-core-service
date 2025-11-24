# Mocks Directory

This directory contains auto-generated mock implementations for testing purposes.

## Structure

The mock files are organized to mirror the production code structure:

```
mocks/
├── api/
│   └── domain/
│       ├── player/
│       │   └── repository_mock.go
│       └── team/
│           └── repository_mock.go
└── README.md
```

## Generating Mocks

To regenerate the mocks, run:

```bash
make generate-mocks
```

This command uses [mockery](https://github.com/vektra/mockery) to automatically generate mocks for all interfaces in the domain layer.

## Usage in Tests

Import mocks using the new centralized path:

```go
import (
    pmocks "sportlink/mocks/api/domain/player"
    tmocks "sportlink/mocks/api/domain/team"
)
```

## Notes

- **Do not edit these files manually** - they are auto-generated
- Mock files are versioned in git to ensure CI/CD tests can run
- The mocks follow the same package structure as the production code for easy navigation
- If you modify an interface, regenerate the mocks with `make generate-mocks`

