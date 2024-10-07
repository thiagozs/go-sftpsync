# Makefile for go-sftpsync

# Variables
BINARY_NAME=go-sftpsync
SOURCE_DIR=./main.go
BUILD_DIR=./bin

# Default target: Build binary
all: build

# Compile and generate the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SOURCE_DIR)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean the build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleanup completed."

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Rebuild the project
rebuild: clean build

# Help
help:
	@echo "Makefile for $(BINARY_NAME)"
	@echo "Usage:"
	@echo "  make         Build the project"
	@echo "  make build   Compile and generate the binary"
	@echo "  make clean   Clean the build artifacts"
	@echo "  make run     Build and run the project"
	@echo "  make rebuild Clean and rebuild the project"
	@echo "  make help    Show this help message"

