# Thunderdome Planning Poker - Development Guide

## Prerequisites

### Required Tools
- Go 1.24+ (Backend)
- Node.js 16+ (Frontend)
- Docker (Optional, for containerization)
- Swagger CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)
- make or task runner (Optional)

### Environment Setup
```bash
# Add Go tools to PATH (if not already done)
export PATH=$PATH:$(go env GOPATH)/bin

# Install swagger generator
go install github.com/swaggo/swag/cmd/swag@latest
```

## Building the Application

### Complete Build Process (Recommended)
```bash
# 1. Install frontend dependencies
cd ui/
npm install
cd ..

# 2. Build frontend assets
npm run build --prefix ui

# 3. Generate swagger documentation
swag init -g internal/http/http.go -o docs/swagger

# 4. Build backend binary
go build -o thunderdome-planning-poker
```

### Using Makefile (Alternative)
```bash
# Install all dependencies
make install

# Complete build (frontend + swagger + backend)
make build

# Development build with live reload
make dev
```

### Frontend Build (Standalone)

```bash
# Navigate to frontend directory
cd ui/

# Install dependencies
npm install

# Build production assets
npm run build

# Development with hot reload
npm run dev
```

### Backend Build (Standalone)

#### Prerequisites
```bash
# Generate swagger docs first (required)
swag init -g internal/http/http.go -o docs/swagger
```

#### Single Platform Build
```bash
# Build for current platform
go build -o thunderdome-planning-poker

# Windows (from any platform)
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o thunderdome-planning-poker.exe

# Linux (from any platform)  
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o thunderdome-planning-poker-linux

# macOS (from any platform)
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o thunderdome-planning-poker-macos

# ARM64 builds
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o thunderdome-planning-poker-linux-arm64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o thunderdome-planning-poker-macos-arm64
```

#### Cross-Platform Build Script
```bash
#!/bin/bash
# build-all.sh - Build for all platforms

# Generate docs first
swag init -g internal/http/http.go -o docs/swagger

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o release/thunderdome-planning-poker-windows-amd64.exe
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o release/thunderdome-planning-poker-windows-arm64.exe

# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o release/thunderdome-planning-poker-linux-amd64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o release/thunderdome-planning-poker-linux-arm64

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o release/thunderdome-planning-poker-macos-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o release/thunderdome-planning-poker-macos-arm64
```

### Docker Build

#### Prerequisites
- Ensure Docker is installed and running
- Frontend must be built first (handled automatically in Dockerfile)

#### Local Docker Build
```bash
# Build single-platform image
docker build -t thunderdome-planning-poker:latest -f build/Dockerfile .

# Build with version tag
docker build -t thunderdome-planning-poker:v1.0.0 -f build/Dockerfile --build-arg BUILD_VERSION=v1.0.0 .

# Test the Docker image
docker run --rm -p 8080:8080 thunderdome-planning-poker:latest
```

#### Multi-Platform Docker Build
```bash
# Setup buildx (one time)
docker buildx create --name multiarch --use
docker buildx inspect --bootstrap

# Build for multiple platforms
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t thunderdome-planning-poker:latest \
  -f build/Dockerfile \
  --push .

# Build and load locally (single platform)
docker buildx build \
  --platform linux/amd64 \
  -t thunderdome-planning-poker:latest \
  -f build/Dockerfile \
  --load .
```

## Development Workflow

### Local Development Setup
```bash
# Option 1: Using make
make dev  # Builds everything and starts with live reload

# Option 2: Manual setup
# Terminal 1 - Frontend dev server
cd ui/
npm run dev

# Terminal 2 - Backend with hot reload
swag init -g internal/http/http.go -o docs/swagger
go build -o thunderdome-planning-poker
HTTP_SECURE_PROTOCOL="false" SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN="localhost" COOKIE_SECURE="false" ./thunderdome-planning-poker live
```

### Running the Application
```bash
# Default configuration (requires database)
./thunderdome-planning-poker

# Development mode with local overrides
HTTP_SECURE_PROTOCOL="false" SMTP_ENABLED="false" DB_HOST="localhost" APP_DOMAIN="localhost" COOKIE_SECURE="false" ./thunderdome-planning-poker

# With custom config file
./thunderdome-planning-poker --config /path/to/config.json

# View available options
./thunderdome-planning-poker --help
```

## Docker Hub Deployment

### Repository Setup
1. Create repository on Docker Hub: `yourusername/thunderdome-planning-poker`
2. Enable automated builds (optional)
3. Configure build rules for tags and branches

### Image Tagging Strategy
```bash
# Development builds
docker tag thunderdome-planning-poker:latest yourusername/thunderdome-planning-poker:dev

# Release candidates  
docker tag thunderdome-planning-poker:latest yourusername/thunderdome-planning-poker:v1.2.3-rc1

# Stable releases
docker tag thunderdome-planning-poker:latest yourusername/thunderdome-planning-poker:v1.2.3
docker tag thunderdome-planning-poker:latest yourusername/thunderdome-planning-poker:latest
```

### Publishing to Docker Hub
```bash
# Login to Docker Hub
docker login

# Push specific version
docker push yourusername/thunderdome-planning-poker:v1.2.3

# Push latest
docker push yourusername/thunderdome-planning-poker:latest

# Push all tags
docker push --all-tags yourusername/thunderdome-planning-poker
```

### Multi-Platform Deployment
```bash
# Build and push multi-platform images
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  -t yourusername/thunderdome-planning-poker:v1.2.3 \
  -t yourusername/thunderdome-planning-poker:latest \
  -f build/Dockerfile \
  --push .
```

### GitHub Actions for Automated Deployment
Create `.github/workflows/docker.yml`:
```yaml
name: Build and Push Docker Image

on:
  push:
    tags:
      - 'v*'
  pull_request:
    branches: [ main ]

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v2
      
    - name: Login to Docker Hub
      uses: docker/login-action@v2
      with:
        username: ${{ secrets.DOCKERHUB_USERNAME }}
        password: ${{ secrets.DOCKERHUB_TOKEN }}
        
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: yourusername/thunderdome-planning-poker
        
    - name: Build and push
      uses: docker/build-push-action@v4
      with:
        context: .
        file: build/Dockerfile
        platforms: linux/amd64,linux/arm64
        push: ${{ github.event_name != 'pull_request' }}
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}
        build-args: BUILD_VERSION=${{ github.ref_name }}
```

## Testing Builds

### Verify Binary Builds
```bash
# Test Windows binary (on Windows)
.\thunderdome-planning-poker.exe --version

# Test Linux binary (on Linux/WSL)
./thunderdome-planning-poker-linux --version

# Test macOS binary (on macOS)
./thunderdome-planning-poker-macos --version
```

### Test Docker Image
```bash
# Run container locally
docker run --rm -p 8080:8080 \
  -e HTTP_SECURE_PROTOCOL=false \
  -e SMTP_ENABLED=false \
  -e DB_HOST=host.docker.internal \
  -e APP_DOMAIN=localhost \
  -e COOKIE_SECURE=false \
  thunderdome-planning-poker:latest

# Verify container health
docker run --rm thunderdome-planning-poker:latest --version
```

## Stop Game Feature Testing

### Manual Testing Checklist
- [ ] Game can be stopped via web interface
- [ ] Game can be stopped via API endpoint
- [ ] Database consistency after game stop
- [ ] User permissions respected for stop action
- [ ] Real-time updates to all connected users
- [ ] Proper cleanup of game state

### API Testing
```bash
# Stop a game via API (replace with actual game ID and auth)
curl -X POST http://localhost:8080/api/games/{gameId}/stop \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json"
```

## Troubleshooting

### Common Build Issues

#### Frontend Build Failures
```bash
# Clear npm cache and reinstall
cd ui/
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
npm run build
```

#### Go Build Failures
```bash
# Clean Go module cache
go clean -modcache
go mod download
go mod tidy

# Regenerate swagger docs
rm -rf docs/swagger
swag init -g internal/http/http.go -o docs/swagger
```

#### Docker Build Issues
```bash
# Clean Docker build cache
docker builder prune -a

# Build with no cache
docker build --no-cache -t thunderdome-planning-poker:latest -f build/Dockerfile .

# Check Docker daemon status
docker version
docker system info
```

### Platform-Specific Issues

#### Windows
- Ensure Go and Node.js are in PATH
- Use Git Bash or PowerShell for better compatibility
- May need to run Docker as administrator

#### macOS
- Install Xcode Command Line Tools: `xcode-select --install`
- Use Homebrew for dependency management
- Apple Silicon Macs: build ARM64 binaries for optimal performance

#### Linux
- Install build essentials: `sudo apt-get install build-essential`
- Ensure Docker group membership: `sudo usermod -aG docker $USER`

### Performance Optimization

#### Binary Size Reduction
```bash
# Build with size optimization
go build -ldflags="-w -s" -o thunderdome-planning-poker

# Further compression with UPX (optional)
upx --best thunderdome-planning-poker
```

#### Docker Image Optimization
- Multi-stage build already implemented in Dockerfile
- Uses scratch base image for minimal size
- Static binary compilation with CGO_ENABLED=0

## Contributing

### Before Submitting PRs
1. Run full build: `make build`
2. Test stop game feature manually
3. Verify Docker build works
4. Update documentation if needed
5. Follow existing code style and patterns

### Code Quality Checks
```bash
# Format Go code
goimports -w .
swag fmt

# Format frontend code
npm run format --prefix ui

# Run tests
go test ./...
npm test --prefix ui
```

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for detailed contribution guidelines.