.PHONY: help
help:
	@echo "Comandos disponibles:"
	@echo ""
	@echo "Full Stack:"
	@echo "  make up                  - Levantar backend + frontend (detached)"
	@echo "  make down                - Detener backend + frontend"
	@echo "  make logs                - Ver logs de los contenedores"
	@echo "  make rebuild             - Reconstruir y levantar infraestructura"
	@echo ""
	@echo "Backend:"
	@echo "  make test                - Ejecutar tests del backend"
	@echo "  make coverage-report     - Generar reporte de cobertura"
	@echo "  make lint                - Ejecutar linter"
	@echo "  make set-up              - Ejecutar go mod tidy y vendor"
	@echo "  make install-dependencies - Instalar herramientas de desarrollo"
	@echo "  make generate-mocks      - Generar mocks para testing"
	@echo ""
	@echo "Frontend:"
	@echo "  make frontend-install    - Instalar dependencias del frontend"
	@echo "  make frontend-dev        - Levantar frontend en modo desarrollo"
	@echo "  make frontend-build      - Compilar frontend para producción"
	@echo ""
	@echo "Legacy:"
	@echo "  make dev                 - Levantar backend + frontend (foreground)"
	@echo "  make stop                - Detener todo"

.PHONY: set-up

.PHONY: install-dependencies
install-dependencies:
	@echo "Installing dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
	go install github.com/vektra/mockery/v2@v2.46.1
	go install golang.org/x/tools/cmd/stringer@latest

.PHONY: set-up
set-up:
	cd backend && go mod tidy
	cd backend && go mod vendor

.PHONY: up
up:
	@echo "Starting infrastructure locally..."
	docker-compose up -d
	@echo "Waiting for backend to be ready..."
	@sleep 5
	@echo "Starting frontend in background..."
	@cd frontend && nohup npm run dev > ../frontend.log 2>&1 & echo $$! > ../frontend.pid
	@echo "✓ Backend running on http://localhost:8080"
	@echo "✓ Frontend running on http://localhost:3000"
	@echo "✓ Frontend logs: tail -f frontend.log"

.PHONY: down
down:
	@echo "Stopping local infrastructure..."
	docker-compose down
	@echo "Stopping frontend..."
	@if [ -f frontend.pid ]; then \
		kill $$(cat frontend.pid) 2>/dev/null || true; \
		rm -f frontend.pid; \
	fi
	@-pkill -f "vite" || true
	@rm -f frontend.log
	@echo "✓ All services stopped"

.PHONY: logs
logs:
	@echo "Showing logs..."
	docker-compose logs -f

.PHONY: rebuild
rebuild:
	@echo "Rebuilding infrastructure..."
	@$(MAKE) down
	-docker rmi sportlink-core-service
	docker-compose up -d --build
	@echo "Waiting for backend to be ready..."
	@sleep 5
	@echo "Starting frontend in background..."
	@cd frontend && nohup npm run dev > ../frontend.log 2>&1 & echo $$! > ../frontend.pid
	@echo "✓ Backend running on http://localhost:8080"
	@echo "✓ Frontend running on http://localhost:3000"
	@echo "✓ Frontend logs: tail -f frontend.log"

.PHONY: env-up
env-up:
	@echo "Creating infrastructure locally..."
	-docker rmi sportlink-core-service
	docker-compose up

.PHONY: env-down
env-down:
	@echo "Destroying local infrastructure..."
	docker-compose down

.PHONY: coverage-report
coverage-report:
	@echo "Running tests and creating coverage report..."
	@cd backend && go test ./... -coverprofile=coverage_full.out > /dev/null 2>&1
	@cd backend && grep -Ev "(/mocks/|/dev/|/cmd/)" coverage_full.out > coverage.out
	@cd backend && go tool cover -html=coverage.out -o coverage.html > /dev/null 2>&1
	@cd backend && echo "Total code coverage: $$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}')" && rm -f coverage_full.out coverage.out

.PHONY: lint
lint:
	@echo "Running GolangCI-Lint..."
	@cd backend && -golangci-lint run -v --fix --out-format json cmd/... api/... > ../golangci_lint.json

.PHONY: test
test:
	@echo "Running tests..."
	@cd backend && go clean --testcache
	@cd backend && go test -cover -parallel 4 ./... | grep -v '?'

.PHONY: generate-mocks
generate-mocks:
	@echo "Generating mocks..."
	@# Buscar el archivo que contiene la interfaz 'Repository' para team
	$(eval REPO_FILE=$(shell grep -rl "type Repository interface" backend/api/domain/team | sed 's|.*/||' | sed 's|\.go||'))
	@# Ejecutar mockery usando el nombre de archivo encontrado como prefijo
	cd backend && mockery --name=Repository --dir=api/domain/team --output=mocks/api/domain/team --outpkg=mocks --case=underscore --filename=$(REPO_FILE)_mock.go

	@# Buscar el archivo que contiene la interfaz 'Repository' para player
	$(eval REPO_FILE=$(shell grep -rl "type Repository interface" backend/api/domain/player | sed 's|.*/||' | sed 's|\.go||'))
	@# Ejecutar mockery usando el nombre de archivo encontrado como prefijo
	cd backend && mockery --name=Repository --dir=api/domain/player --output=mocks/api/domain/player --outpkg=mocks --case=underscore --filename=$(REPO_FILE)_mock.go
	@echo "Mocks generated successfully in backend/mocks/"

# Frontend commands
.PHONY: frontend-install
frontend-install:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

.PHONY: frontend-dev
frontend-dev:
	@echo "Starting frontend in development mode..."
	cd frontend && npm run dev

.PHONY: frontend-build
frontend-build:
	@echo "Building frontend for production..."
	cd frontend && npm run build

# Full stack commands
.PHONY: dev
dev:
	@echo "Starting full stack (backend + frontend)..."
	@$(MAKE) up
	@echo "Waiting for backend to be ready..."
	@sleep 5
	@echo "Backend is ready on http://localhost:8080"
	@echo "Starting frontend on http://localhost:3000"
	@cd frontend && npm run dev

.PHONY: stop
stop:
	@echo "Stopping all services..."
	@$(MAKE) down
	@echo "Killing frontend process..."
	@-pkill -f "vite" || true
	@echo "All services stopped"