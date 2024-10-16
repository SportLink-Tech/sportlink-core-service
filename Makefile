.PHONY: set-up
set-up:
	go mod tidy
	go mod vendor

.PHONY: env-up
env-up:
	-docker rmi sportlink-core-service
	docker-compose up --build

.PHONY: test
test:
	go clean --testcache
	go test -cover -parallel 4 ./... | grep -v '?'

.PHONY: generate-mocks
generate-mocks:
	@# Buscar el archivo que contiene la interfaz 'Repository'
	$(eval REPO_FILE=$(shell grep -rl "type Repository interface" api/domain/team | sed 's|.*/||' | sed 's|\.go||'))
	@# Ejecutar mockery usando el nombre de archivo encontrado como prefijo
	mockery --name=Repository --dir=api/domain/team --output=api/domain/team/mocks --outpkg=mocks --case=underscore --filename=$(REPO_FILE)_mock.go

	@# Buscar el archivo que contiene la interfaz 'Repository'
	$(eval REPO_FILE=$(shell grep -rl "type Repository interface" api/domain/player | sed 's|.*/||' | sed 's|\.go||'))
	@# Ejecutar mockery usando el nombre de archivo encontrado como prefijo
	mockery --name=Repository --dir=api/domain/player --output=api/domain/player/mocks --outpkg=mocks --case=underscore --filename=$(REPO_FILE)_mock.go