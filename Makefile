# Development
dev:
	docker-compose up --build

dev-bg:
	docker-compose up -d --build

# Testing
# Always recreate containers to ensure fresh runs (keeps volume)
test:
	docker-compose -f docker-compose.test.yml down --remove-orphans
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --remove-orphans

# Alternative test command that forces recreation (removes volumes)
test-fresh:
	docker-compose -f docker-compose.test.yml down -v --remove-orphans
	docker-compose -f docker-compose.test.yml up --build --force-recreate --abort-on-container-exit

# Start only the test database for development
test-db:
	docker-compose -f docker-compose.test.yml up -d postgres-test

# Run tests against existing test database 
test-quick:
	docker-compose -f docker-compose.test.yml run --rm test-runner

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

.PHONY: dev dev-bg test test-db clean-dev clean-test clean stop