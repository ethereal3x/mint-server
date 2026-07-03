.PHONY: help deps gen gen-local gen-remote proto openapi lint test build compose-up compose-down migrate

help: ## show make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## update buf dependency lock file
	buf dep update

gen: deps ## generate code locally when plugins exist, fallback to buf remote plugins
	@if command -v protoc-gen-go >/dev/null 2>&1 && \
		command -v protoc-gen-go-grpc >/dev/null 2>&1 && \
		command -v protoc-gen-grpc-gateway >/dev/null 2>&1 && \
		command -v protoc-gen-openapiv2 >/dev/null 2>&1; then \
		echo "using local protoc plugins"; \
		$(MAKE) gen-local; \
	else \
		echo "local protoc plugins are incomplete, fallback to buf remote plugins"; \
		$(MAKE) gen-remote; \
	fi

gen-local: ## generate protobuf, grpc, gateway and openapi code with local plugins
	buf generate --template buf.local.gen.yaml

gen-remote: ## generate protobuf, grpc, gateway and openapi code with buf remote plugins
	buf generate --template buf.remote.gen.yaml

proto: gen ## generate protobuf, grpc, gateway and openapi code

openapi: gen ## generate openapi documentation

lint: ## run golangci-lint and buf lint
	golangci-lint run ./...
	buf lint

test: ## run unit tests
	go test ./...

build: ## build all services
	mkdir -p bin
	@for dir in cmd/*; do \
		if [ -d "$$dir" ] && ls "$$dir"/*.go >/dev/null 2>&1; then \
			service=$$(basename "$$dir"); \
			go build -o "bin/$$service" "./$$dir"; \
		fi; \
	done

compose-up: ## start local dependencies
	docker compose -f deploy/compose/docker-compose.yml up -d

compose-down: ## stop local dependencies
	docker compose -f deploy/compose/docker-compose.yml down

migrate: ## run database migrations manually (apply SQL files in order)
	@echo "Apply migrations/mysql/*.sql to your MySQL instance in numeric order"
