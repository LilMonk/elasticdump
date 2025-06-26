.ONESHELL:
.DEFAULT_GOAL := help

MAIN_FILE := main.go
BINARY := bin/elasticdump

##@ Commands
.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <command>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the project
	@echo "Building..."
	@start=$$(date +%s); \
	go build -o ${BINARY} ${MAIN_FILE}; \
	end=$$(date +%s); \
	echo "Build completed in $$(($${end}-$${start})) seconds"

.PHONY: test
test: ## Run the tests
	go clean -testcache
	go test -v ./...

.PHONY: coverage
coverage: ## Run the tests with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: coverage-detailed
coverage-detailed: ## Run tests with detailed coverage report
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out

.PHONY: test-clean
test-clean: ## Clean test cache and coverage files
	go clean -testcache
	rm -f coverage.out coverage.html

.PHONY: format
format: ## Format the code
	go fmt ./...

.PHONY: lint
lint: ## Lint the code
	golangci-lint run

.PHONY: clean
clean: test-clean ## Clean the build artifacts
	rm -f bin/*

.PHONY: install
install: ## Install the binary
	go install .

# Cross-compilation targets
.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@start=$$(date +%s); \
	GOOS=linux GOARCH=amd64 go build -o bin/elasticdump-linux-amd64 . ; \
	end=$$(date +%s); \
	echo "Linux build completed in $$(($${end}-$${start})) seconds"

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "Building for Windows..."
	@start=$$(date +%s); \
	GOOS=windows GOARCH=amd64 go build -o bin/elasticdump-windows-amd64.exe . ; \
	end=$$(date +%s); \
	echo "Windows build completed in $$(($${end}-$${start})) seconds"

.PHONY: build-darwin
build-darwin: ## Build for macOS
	@echo "Building for macOS..."
	@start=$$(date +%s); \
	GOOS=darwin GOARCH=amd64 go build -o bin/elasticdump-darwin-amd64 . ; \
	end=$$(date +%s); \
	echo "macOS build completed in $$(($${end}-$${start})) seconds"

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## Build for all platforms
	@echo "All builds completed."

# Development
.PHONY: dev-deps
dev-deps: ## Install development dependencies
	go mod download
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: docker-build
docker-build: ## Build the docker image
	docker build --no-cache -t $(DOCKER_IMAGE_NAME) . -f $(DOCKER_FILE)

.PHONY: docker-up
docker-up: ## Start the services defined in docker-compose.yaml
	docker compose -f $(DOCKER_COMPOSE_FILE) up -d

.PHONY: docker-down
docker-down: ## Stop the services defined in docker-compose.yaml
	docker compose -f $(DOCKER_COMPOSE_FILE) down

.PHONY: git-release
git-release: ## Create a new git release
	@echo "Creating a new git release..."
	@read -p "Enter the version (e.g., v1.0.0): " version; \
	if [ -z "$$version" ]; then \
		echo "Version cannot be empty"; \
		exit 1; \
	fi; \
	git tag -a "$$version" -m "Release $$version"; \
	git push origin "$$version"; \
	echo "Release $$version created and pushed to origin."