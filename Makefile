# Development
dev:
	docker-compose up --build

dev-bg:
	docker-compose up -d --build

# Horizontal Scaling
INSTANCES ?= 3

dev-scale:
	docker-compose up --build --scale app=$(INSTANCES)

dev-scale-bg:
	docker-compose up -d --build --scale app=$(INSTANCES)

test-load:
	@for i in {1..10}; do curl -s http://localhost/health | grep instance; done

# Testing commands
test-service:
	docker-compose run --rm test go test -v ./internal/service/

test-repository:
	docker-compose -f docker-compose.test.yml down --remove-orphans
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --remove-orphans

test-handlers:
	docker-compose run --rm test go test -v ./internal/handlers/

test-repository-fresh:
	docker-compose -f docker-compose.test.yml down -v --remove-orphans
	docker-compose -f docker-compose.test.yml up --build --force-recreate --abort-on-container-exit

test-repository-quick:
	docker-compose -f docker-compose.test.yml run --rm test-runner

test-unit:
	make test-service
	make test-handlers

test-all:
	make test-service
	make test-repository
	make test-handlers

# Cleanup
clean-dev:
	docker-compose down -v

clean-test:  
	docker-compose -f docker-compose.test.yml down -v

clean:
	make clean-dev
	make clean-test

stop:
	docker-compose down
	docker-compose -f docker-compose.test.yml down

# All phony targets at the end (your preferred style)
.PHONY: build run dev dev-bg test-service test-repository test-handlers test-repository-fresh test-repository-quick test-all clean-dev clean-test clean stop