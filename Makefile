# Development
dev:
	docker-compose up --build

dev-bg:
	docker-compose up -d --build

# Testing

# Service Unit Tests
test-service:
	docker-compose run --rm test go test -v ./internal/service/

# Repository Integration Tests
# Always recreate containers to ensure fresh runs (keeps volume)
test-repository:
	docker-compose -f docker-compose.test.yml down --remove-orphans
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --remove-orphans

# Alternative test command that forces recreation (removes volumes)
test-repository-fresh:
	docker-compose -f docker-compose.test.yml down -v --remove-orphans
	docker-compose -f docker-compose.test.yml up --build --force-recreate --abort-on-container-exit

# Run tests against existing test database 
test-repository-quick:
	docker-compose -f docker-compose.test.yml run --rm test-runner

# All Tests
test-all:
	make test-service
	make test-repository

# Cleanup
clean-dev:
	docker-compose down -v

clean-test:  
	docker-compose -f docker-compose.test.yml down -v

clean:
	make clean-dev
	make clean-test

# Soft cleanup (keep data)
stop:
	docker-compose down
	docker-compose -f docker-compose.test.yml down

.PHONY: dev dev-bg test-service test-repository test-repository-fresh test-repository-quick test-all clean clean-dev clean-test stop