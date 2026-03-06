.PHONY: dev dev-api dev-web build build-web build-go serve clean convert-data docker docker-run

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

# Static build for GitHub Pages (no backend needed)
build-static:
	cd web && npm ci && npm run build:static
	cp web/dist/index.html web/dist/404.html

# Preview static build locally
preview-static: build-static
	cd web && npx vite preview --mode static --port 4173

# Docker build
docker:
	docker build -t heywood-tbs .

# Run in Docker (pass OPENAI_API_KEY or Azure env vars)
docker-run: docker
	docker run -p 8080:8080 --env-file .env heywood-tbs

clean:
	rm -rf bin/ web/dist/ web/node_modules/
