.PHONY: local_deploy

local_deploy:
	-docker rmi sportlink-core-service
	docker-compose up --build