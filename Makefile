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
	@tailwindcss -c ./config/.tailwind.config.js -i ./web/assets/styles/main.css -o ./web/public/styles/main.css --minify --watch

css:
	@tailwindcss -c ./config/.tailwind.config.js -i ./web/assets/styles/main.css -o ./web/public/styles/main.css --minify


seed:
	@echo "Seeding..."
	@go run ./cmd/seeder

test:
	@go test -v ./...
	

test-cov:
	@go test -v -converprofile=coverage.txt ./...

gen-mocks:
	@mockery
