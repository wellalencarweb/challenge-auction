conf ?= cmd/auction/.env
include $(conf)
export $(shell sed 's/=.*//' $(conf))



## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@rm -f coverage.*



## ----- COMPOSE
.PHONY: up
up: ## Put the compose containers up
	@docker-compose up -d

.PHONY: down
down: ## Put the compose containers down
	@docker-compose down


.PHONY: build
build: ## Build the application
	@go build -o bin/app cmd/auction/main.go

.PHONY: run
run: ## Run the application
	@go run cmd/auction/main.go

.PHONY: test
test: ## Run tests
	@go test ./... -v

.PHONY: lint
lint: ## Run go vet and lint
	@go vet ./...
