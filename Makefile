.PHONY: local_deploy

env-up:
	-docker rmi sportlink-core-service
	docker-compose up --build

test:
	go test -cover ./...

create-mock:
	mockery --name=Repository --dir=api/domain/team --output=api/domain/team/mocks --outpkg=mocks --case=underscore
