# todo — CLI TODO Tracker with Notification Daemon

# Default: show available recipes
default:
    @just --list

# ─── Build ───────────────────────────────────────────────

# Build the binary
build:
    go build -o todo .

# Build with stripped symbols (smaller binary)
build-release:
    CGO_ENABLED=0 go build -ldflags="-s -w" -o todo .

# Install to ~/.local/bin
install: build-release
    cp todo ~/.local/bin/todo
    @echo "✓ Installed to ~/.local/bin/todo"

# Clean build artifacts
clean:
    rm -f todo
    go clean -cache -testcache

# ─── Test ────────────────────────────────────────────────

# Run all tests
test:
    go test ./... -count=1

# Run tests with verbose output
test-v:
    go test ./... -v -count=1

# Run tests with coverage
test-cover:
    go test ./... -coverprofile=coverage.out -count=1
    go tool cover -func=coverage.out
    @echo ""
    @echo "To view HTML report: just test-cover-html"

# Open coverage report in browser
test-cover-html: test-cover
    go tool cover -html=coverage.out -o coverage.html
    @echo "✓ Coverage report: coverage.html"

# ─── Lint & Format ──────────────────────────────────────

# Format all Go files
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Run all checks (fmt, vet, test)
check: fmt vet test

# ─── Docker ─────────────────────────────────────────────

# Build Docker image
docker-build:
    docker compose build

# Start daemon container
docker-up:
    docker compose up -d

# Stop daemon container
docker-down:
    docker compose down

# View daemon container logs
docker-logs:
    docker compose logs -f

# Rebuild and restart daemon container
docker-restart: docker-build
    docker compose down
    docker compose up -d

# ─── Run ─────────────────────────────────────────────────

# Run the daemon locally (default 60s interval)
daemon interval="60":
    go run . daemon -i {{interval}}

# Run any todo command (e.g. just run add "My task" -p high)
run *ARGS:
    go run . {{ARGS}}

# ─── Dev Helpers ─────────────────────────────────────────

# Show current config
config:
    @go run . config get postback_url
    @go run . config get ntfy_url

# List all TODOs
list:
    @go run . list -a

# Tidy go modules
tidy:
    go mod tidy

# Show project stats
stats:
    @echo "Go files:"
    @find . -name '*.go' | wc -l
    @echo "Lines of Go:"
    @find . -name '*.go' | xargs wc -l | tail -1
    @echo "Test files:"
    @find . -name '*_test.go' | wc -l
    @echo "Binary size:"
    @ls -lh todo 2>/dev/null | awk '{print $5}' || echo "(not built)"
