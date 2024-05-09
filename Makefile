dev:
	@air

build: 
	@echo "Building the binary..."
	@go build -o ./tmp/main ./cmd/odinls

start: build
	@echo "Starting the service..."
	@./tmp/main
