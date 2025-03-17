GYMONDO_BINARY=binary_file/gymondoApp

run_postgres:
	@echo "Stopping Docker containers (if running)..."
	docker compose down
	@echo "Building (if required) and starting Docker containers..."
	docker compose up -d postgres
	@echo "Postgres started"

## up_build: stops docker compose (if running), builds all projects and starts docker compose
up_build: build_gymondo
	@echo "stopping docker images (if running...)"
	docker-compose down
	@echo "building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "docker images built and started"

## down: stop docker compose
down:
	@echo "Stopping Docker compose..."
	docker compose down
	@echo "Done"

## tech_task: builds the tech task binary as a Linux executable
build_gymondo:
	@echo "building Gymondo binary..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${GYMONDO_BINARY} ./cmd
	@echo "build completed"

## running integration tests
test.integration:
	go test -tags=integration ./integration_tests -v 

## running UNIT tests
test.unit:
	go test ./...
