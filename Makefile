.PHONY: test build release-patch release-minor release-major

# Test and build
test:
	go test -v ./...

build:
	go build -v ./...

# Release commands
release-patch:
	@./release.sh patch

release-minor:
	@./release.sh minor

release-major:
	@./release.sh major

# Quick aliases
patch: release-patch
minor: release-minor
major: release-major

# Help
help:
	@echo "Available commands:"
	@echo "  make test         - Run tests"
	@echo "  make build        - Build project"
	@echo "  make patch        - Create patch release (x.x.X)"
	@echo "  make minor        - Create minor release (x.X.0)"
	@echo "  make major        - Create major release (X.0.0)"