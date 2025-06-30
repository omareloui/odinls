watch:
	@air -c config/.air.toml

build: templ css
	@echo "Building the binary..."
	@go build -o ./tmp/main ./cmd/odinls

templ:
	@echo "Generating the templates..."
	@templ generate



start: build
	@echo "Starting the service..."
	@./tmp/main

css-dev:
	@pnpm css:dev

css:
	@pnpm css


seed:
	@echo "Seeding..."
	@go run ./cmd/seeder

test:
	@go test -v ./...
	

test-cov:
	@go test -v -converprofile=coverage.txt ./...

gen-mocks:
	@mockery
