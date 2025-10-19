#!/bin/bash

# ArasAuth Setup Script
# This script helps set up the ArasAuth service for development and production

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    # Check Go
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        print_success "Go $GO_VERSION found"
    else
        print_error "Go is not installed. Please install Go 1.22 or later."
        exit 1
    fi
    
    # Check Docker
    if command_exists docker; then
        print_success "Docker found"
    else
        print_warning "Docker not found. Docker is required for containerized deployment."
    fi
    
    # Check Docker Compose
    if command_exists docker-compose; then
        print_success "Docker Compose found"
    else
        print_warning "Docker Compose not found. Required for containerized deployment."
    fi
    
    # Check PostgreSQL (optional for local development)
    if command_exists psql; then
        print_success "PostgreSQL client found"
    else
        print_warning "PostgreSQL client not found. Required for local development."
    fi
}

# Function to generate build information
generate_build_info() {
    local version=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    local build_time=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    local git_commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    echo "BUILD_VERSION=$version"
    echo "BUILD_TIME=$build_time"
    echo "GIT_COMMIT=$git_commit"
}

# Function to setup environment
setup_environment() {
    print_status "Setting up environment..."
    
    if [ ! -f .env ]; then
        if [ -f .env.example ]; then
            cp .env.example .env
            print_success "Created .env file from .env.example"
            
            # Add build information to .env
            print_status "Adding build information to .env..."
            echo "" >> .env
            echo "# Build Information (for Docker builds)" >> .env
            generate_build_info >> .env
            
            print_success "Build information added to .env"
            print_warning "Please update .env with your configuration"
        else
            print_error ".env.example not found"
            exit 1
        fi
    else
        print_warning ".env file already exists"
        
        # Check if build info exists in .env
        if ! grep -q "BUILD_VERSION" .env; then
            print_status "Adding build information to existing .env..."
            echo "" >> .env
            echo "# Build Information (for Docker builds)" >> .env
            generate_build_info >> .env
            print_success "Build information added to .env"
        fi
    fi
}

# Function to install dependencies
install_dependencies() {
    print_status "Installing Go dependencies..."
    go mod download
    go mod tidy
    print_success "Dependencies installed"
}

# Function to build the application
build_application() {
    print_status "Building application..."
    go build -o bin/aras-auth ./cmd/server
    go build -o bin/migrate ./cmd/migrate
    print_success "Application built successfully"
}

# Function to build Docker image with version info
build_docker_versioned() {
    print_status "Building Docker image with version information..."
    
    # Generate build info
    local version=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
    local build_time=$(date -u +%Y-%m-%dT%H:%M:%SZ)
    local git_commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    print_status "Build info: Version=$version, Time=$build_time, Commit=$git_commit"
    
    docker build \
        --build-arg BUILD_VERSION="$version" \
        --build-arg BUILD_TIME="$build_time" \
        --build-arg GIT_COMMIT="$git_commit" \
        -t aras-auth:latest .
    
    print_success "Docker image built successfully with version info"
}

# Function to run database migrations
run_migrations() {
    print_status "Running database migrations..."
    
    if [ -f bin/migrate ]; then
        ./bin/migrate up
        print_success "Database migrations completed"
    else
        print_error "Migration binary not found. Please build the application first."
        exit 1
    fi
}

# Function to start with Docker
start_docker() {
    print_status "Starting services with Docker Compose..."
    docker-compose up -d
    print_success "Services started"
    print_status "Waiting for services to be ready..."
    sleep 10
    
    # Run migrations and print version
    print_status "Running database migrations..."
    docker-compose exec aras_auth ./aras-auth --version
    print_success "Database migrations completed"
    
    print_success "ArasAuth is running at http://localhost:7600"
    print_status "Health check: http://localhost:7600/health"   
}

# Function to start locally
start_local() {
    print_status "Starting application locally..."
    print_warning "Make sure PostgreSQL is running and configured"
    go run ./cmd/server
}

# Function to run tests
run_tests() {
    print_status "Running tests..."
    go test -v ./...
    print_success "Tests completed"
}

# Function to clean up
cleanup() {
    print_status "Cleaning up..."
    rm -rf bin/
    go clean
    print_success "Cleanup completed"
}

# Function to show help
show_help() {
    echo "ArasAuth Setup Script"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  check       Check prerequisites"
    echo "  setup       Setup environment and install dependencies"
    echo "  build       Build the application"
    echo "  docker-build Build Docker image with version info"
    echo "  migrate     Run database migrations"
    echo "  start       Start the application locally"
    echo "  docker      Start with Docker Compose"
    echo "  test        Run tests"
    echo "  clean       Clean build artifacts"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 setup    # Setup environment and dependencies"
    echo "  $0 docker   # Start with Docker Compose"
    echo "  $0 start    # Start locally (requires PostgreSQL)"
}

# Main script logic
case "${1:-help}" in
    check)
        check_prerequisites
        ;;
    setup)
        check_prerequisites
        setup_environment
        install_dependencies
        print_success "Setup completed successfully"
        ;;
    build)
        install_dependencies
        build_application
        ;;
    docker-build)
        build_docker_versioned
        ;;
    migrate)
        run_migrations
        ;;
    start)
        start_local
        ;;
    docker)
        start_docker
        ;;
    test)
        run_tests
        ;;
    clean)
        cleanup
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac


