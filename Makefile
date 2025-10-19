# Image configuration
IMAGE_NAME = ghcr.io/aras-group-co/aras-auth

.PHONY: help build run test clean docker-build-versioned docker-push-versioned docker-run docker-stop migrate-up migrate-down

# Default target
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application locally"
	@echo "  test           - Run tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-build-versioned - Build Docker image with version info"
	@echo "  docker-push-versioned - Push versioned Docker image to registry (clean tags only)"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose services"
	@echo "  migrate-up     - Run database migrations up"
	@echo "  migrate-down   - Run database migrations down"
	@echo "  dev            - Run in development mode"

# Build the application
build:
	go build -o bin/aras-auth ./cmd/server

# Run the application locally
run:
	go run ./cmd/server

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image with version info
docker-build-versioned:
	@VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	echo "Building image: $(IMAGE_NAME):$$VERSION"; \
	docker build \
		--build-arg BUILD_VERSION=$$VERSION \
		--build-arg BUILD_TIME=$$(date -u +%Y-%m-%dT%H:%M:%SZ) \
		--build-arg GIT_COMMIT=$$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
		-t $(IMAGE_NAME):$$VERSION \
		.

# Push versioned Docker image to registry (only clean git tags)
docker-push-versioned:
	@VERSION=$$(git describe --tags --always --dirty 2>/dev/null || echo "dev"); \
	echo "Version detected: $$VERSION"; \
	if echo "$$VERSION" | grep -q 'dirty' || echo "$$VERSION" | grep -q '\-g' || [ "$$VERSION" = "dev" ]; then \
		echo "❌ ERROR: Cannot push non-release version: $$VERSION"; \
		echo ""; \
		echo "Only clean git tags can be pushed to registry."; \
		echo ""; \
		echo "Current issues:"; \
		if echo "$$VERSION" | grep -q 'dirty'; then \
			echo "  - You have uncommitted changes"; \
		fi; \
		if echo "$$VERSION" | grep -q '\-g'; then \
			echo "  - You have commits after the last tag"; \
		fi; \
		if [ "$$VERSION" = "dev" ]; then \
			echo "  - No git tags found"; \
		fi; \
		echo ""; \
		echo "To fix:"; \
		echo "  1. Commit all changes: git add . && git commit"; \
		echo "  2. Create a tag: git tag v1.x.x"; \
		echo "  3. Push tag: git push origin v1.x.x"; \
		exit 1; \
	fi; \
	echo "✅ Pushing clean version: $$VERSION"; \
	docker push $(IMAGE_NAME):$$VERSION

# Run with Docker Compose
docker-run:
	docker-compose up -d

# Stop Docker Compose services
docker-stop:
	docker-compose down

# Run database migrations up
migrate-up:
	go run ./cmd/migrate up

# Run database migrations down
migrate-down:
	go run ./cmd/migrate down

# Run in development mode
dev:
	docker-compose -f docker-compose.dev.yml up

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate mocks (if using mockgen)
mocks:
	mockgen -source=internal/domain/user.go -destination=internal/mocks/user_mock.go
	mockgen -source=internal/domain/group.go -destination=internal/mocks/group_mock.go
	mockgen -source=internal/domain/role.go -destination=internal/mocks/role_mock.go
	mockgen -source=internal/domain/permission.go -destination=internal/mocks/permission_mock.go


