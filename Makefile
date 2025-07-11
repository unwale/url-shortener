.PHONY: test-integration test-unit check-env

test-unit:
	go test ./... -coverprofile=coverage.txt

check-env:
	@if [ ! -f .test.env ]; then \
		echo "--> Error: .test.env file not found. Create it to run integration tests."; \
		exit 1; \
	fi

test-integration: check-env
	@set -a; \
	source .test.env; \
	set +a; \
	docker compose -f docker-compose.test.yaml up -d --build; \
	go test --tags=integration -p=1 ./... -coverprofile=coverage.txt; \
	EXIT_CODE=$$?; \
	docker compose -f docker-compose.test.yaml down --remove-orphans; \
	exit $$EXIT_CODE