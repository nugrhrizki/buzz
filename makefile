.PHONY: all info server-build client-build server-run client-run install-template-dependencies clean

# Print information about available commands
info:
	$(info ------------------------------------------)
	$(info -                KinStack                -)
	$(info ------------------------------------------)
	$(info This Makefile helps you manage the project.)
	$(info )
	$(info Available commands:)
	$(info - server-build:  Build the Golang project.)
	$(info - client-build:  Build the SolidJS project.)
	$(info - server-run:    Run the Golang project. (development mode))
	$(info - client-run:    Run the SolidJS project. (development mode))
	$(info - all:           Run all commands (SolidJSBuild, GoBuild).)
	$(info - clean:         Clean build artifacts.)
	$(info )
	$(info Usage: make <command>)

# Default target
all: client-build server-build

# Run the Golang project
server-run: client-build
	@echo "=== Running Server ==="
	@air -c .air.toml

# Build the Golang project
server-build: client-build
	@echo "=== Building Server ==="
	@go build -o app -v ./cmd/web/main.go

# Run the SolidJS project
client-run: install-template-dependencies
	@echo "=== Running Client ==="
	@if command -v pnpm >/dev/null; then \
		pnpm run -C ./web dev; \
	else \
		npm run --prefix ./web dev; \
	fi

# Build the SolidJS project
client-build: install-template-dependencies
	@echo "=== Building Client ==="
	@if command -v pnpm >/dev/null; then \
		pnpm run -C ./web build; \
	else \
		npm run --prefix ./web build; \
	fi

# Install template dependencies
install-template-dependencies:
	@if command -v pnpm >/dev/null; then \
		pnpm install -C ./web; \
	else \
		npm install --prefix ./web; \
	fi

# Clean build artifacts
clean:
	@echo "=== Cleaning build artifacts ==="
	@rm -f app
	@rm -rf web/dist
