.PHONY: test-integration test-unit

test-unit:
	go test ./... -coverprofile=coverage.txt

test-integration:
	@touch .test.env; \
	set -a; \
	source .test.env; \
	set +a; \
	docker compose -f docker-compose.test.yaml up -d --build; \
	go test --tags=integration -p=1 ./... -coverprofile=coverage.txt; \
	EXIT_CODE=$$?; \
	docker compose -f docker-compose.test.yaml down --remove-orphans; \
	exit $$EXIT_CODE