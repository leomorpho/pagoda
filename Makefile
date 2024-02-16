# Determine if you have docker-compose or docker compose installed locally
# If this does not work on your system, just set the name of the executable you have installed
DCO_BIN := $(shell { command -v docker-compose || command -v docker compose; } 2>/dev/null)

define Comment
	- Run `make help` to see all the available options.
endef

.PHONY: help
help: ## Show this help message.
	@echo "Available options:"
	@echo
	@awk 'BEGIN {FS = ":.*?## "}; /^[a-zA-Z_-]+:.*?## .*$$/ { printf "\033[36m%-30s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo
	@echo "To see the details of each command, run: make <command>"

.PHONY: db
db: ## Connect to the primary database
	docker exec -it pagoda_db psql postgresql://admin:admin@localhost:5432/app

.PHONY: db-test
db-test: ## Connect to the test database (you must run tests first before running this)
	docker exec -it pagoda_db psql postgresql://admin:admin@localhost:5432/app_test

.PHONY: cache
cache: ## Connect to the primary cache
	docker exec -it pagoda_cache redis-cli

.PHONY: cache-clear
cache-clear: ## Clear the primary cache
	docker exec -it pagoda_cache redis-cli flushall


.PHONY: cache-test
cache-test: ## Connect to the test cache
	docker exec -it pagoda_cache redis-cli -n 1

.PHONY: ent-install
ent-install: ## Install Ent code-generation module
	go get -d entgo.io/ent/cmd/ent

.PHONY: ent-gen
ent-gen: ## Generate Ent code
	go generate ./ent

.PHONY: ent-new
ent-new: ## Create a new Ent entity
	go run entgo.io/ent/cmd/ent new $(name)


.PHONY: up
up: ## Start the Docker containers
	$(DCO_BIN) up -d
	sleep 3


.PHONY: stop
stop: ## Stop the Docker containers
	$(DCO_BIN) stop


.PHONY: down
down: ## Drop the Docker containers to wipe all data
	$(DCO_BIN) down

.PHONY: reset
reset: ## Rebuild Docker containers to wipe all data
	$(DCO_BIN) down
	make up


.PHONY: build-js
build-js: ## Build JS/Svelte assets
	npm run build


.PHONY: build-js
watch-js: ## Build JS/Svelte assets (auto reload changes)
	npm run watch 

watch-css: ## Build CSS assets (auto reload changes)
	npx tailwindcss -i ./styles/styles.css -o ./static/styles_bundle.css --watch


.PHONY: run
watch-go: ## Run the application with air (auto reload changes)
	clear
	air

watch: 
	overmind start

.PHONY: test
test: ## Run all tests
	go test -p 1 ./...


.PHONY: test
testall: ## Run all tests
	go test -count=1 -p 1 ./...

.PHONY: cover
cover: ## Run the Go coverage tool on the codebase
	@echo "Running tests with coverage..."
	@go test -coverprofile=/tmp/coverage.out -count=1 -p 1  ./...
	@echo "Generating HTML coverage report..."
	@go tool cover -html=/tmp/coverage.out


.PHONY: worker
worker: ## Run the worker
	clear
	go run cmd/worker/main.go

.PHONY: workerui
workerui: ## Run the worker asynq dash
	asynq dash

.PHONY: check-updates
check-updates: ## Check for direct dependency updates
	go list -u -m -f '{{if not .Indirect}}{{.}}{{end}}' all | grep "\["

.PHONY: test-e2e
e2e: ## Run Playwright tests
	@echo "Running end-to-end tests..."
	@cd e2e_tests && npm install && npx playwright test

.PHONY: test-e2e
e2eui: ## Run Playwright tests
	@echo "Running end-to-end tests..."
	@cd e2e_tests && npm install && npx playwright test --ui

.PHONY: codegen
codegen: ## Generate Playwright tests interactively
	@echo "Running Playwright codegen for URL http://localhost:8000..."
	@cd e2e_tests && npx playwright codegen http://localhost:8000