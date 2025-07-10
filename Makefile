.PHONY: test-integration test-unit check-env


check-env:
	@if [ ! -f .test.env ]; then \
		echo "--> Error: .test.env file not found."; \
		exit 1; \
	fi


include .test.env

export $(shell sed 's/=.*//' .test.env)


test-integration: check-env
	docker compose -f docker-compose.test.yaml up -d --build
	go test --tags=integration -p=1 ./... -coverprofile=coverage.txt 
	docker compose -f docker-compose.test.yaml down --remove-orphans

test-unit:
	go test ./... -coverprofile=coverage.txt