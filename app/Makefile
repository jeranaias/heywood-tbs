.PHONY: dev dev-api dev-web build build-web build-go serve clean convert-data

# Run API + Vite dev server concurrently
dev:
	@echo "Starting Heywood dev environment..."
	@echo "API: http://localhost:8080"
	@echo "Web: http://localhost:5173"
	$(MAKE) dev-api &
	$(MAKE) dev-web

# Go API server in dev mode
dev-api:
	cd cmd/server && go run . -dev -port 8080 -data ../../data

# Vite dev server with HMR
dev-web:
	cd web && npm run dev

# Production build
build: build-web build-go

build-web:
	cd web && npm ci && npm run build

build-go:
	CGO_ENABLED=0 go build -o bin/heywood ./cmd/server

# Run the production binary
serve: build
	cd bin && ./heywood -port 8080 -data ../data

# Convert CSV source data to JSON
convert-data:
	python scripts/csv-to-json.py

clean:
	rm -rf bin/ web/dist/ web/node_modules/
