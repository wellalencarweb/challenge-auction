# ---------- CONFIGURATION
ENV_FILE := cmd/auction/.env
ifneq ("$(wildcard $(ENV_FILE))","")
    include $(ENV_FILE)
    export $(shell sed 's/=.*//' $(ENV_FILE))
else
    $(error Environment file $(ENV_FILE) not found. Please ensure it exists.)
endif

# ---------- UTILS
.PHONY: help
help: ## Display the list of available make commands with descriptions
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ---------- COMPOSE
.PHONY: up
up: ## Start the Docker Compose containers in detached mode
	docker compose up -d

.PHONY: down
down: ## Stop and remove the Docker Compose containers
	docker compose down

.PHONY: restart
restart: down up ## Restart the Docker Compose containers by stopping and starting them again

.PHONY: logs
logs: ## Tail the logs of Docker Compose containers in real-time
	docker compose logs -f

# ---------- TESTS
.PHONY: create
create: ## Send a POST request to create an auction
	curl -X POST http://localhost:8080/auction \
		-H "Content-Type: application/json" \
		-d '{"product_name": "celular", "category": "eletronicos", "description": "Iphone 1 8GB", "condition": 0}'

.PHONY: list
list: ## Send a GET request to list auctions with status=0
	curl -X GET http://localhost:8080/auction?status=0 \
		-H "Content-Type: application/json"

.PHONY: check
check: ## Send a GET request to check auctions with status=1
	curl -X GET http://localhost:8080/auction?status=1 \
		-H "Content-Type: application/json"

.PHONY: clear
clear: ## Remove project-specific Docker containers, images, and volumes
	@echo "Stopping and removing challenge-auction containers..."
	@docker compose down -v || true
	
	@echo "Removing project volumes..."
	@docker volume rm challenge-auction_mongo-data 2>/dev/null || true
	
	@echo "Removing project images..."
	@docker images -q challenge-auction-app mongodb | xargs -r docker rmi -f
	
	@echo "Cleaning up unused project networks..."
	@docker network ls --filter name=challenge-auction --format "{{.Name}}" | xargs -r docker network rm 2>/dev/null || true
	
	@echo "Project cleanup completed."

.PHONY: test
test: ## Run a specific test file or package.
	go test -v internal/infra/database/auction/create_auction_test.go

# ---------- COLORS
BLUE := \033[1;34m
GREEN := \033[1;32m
RED := \033[1;31m
YELLOW := \033[1;33m
NC := \033[0m # No Color

define print_step
	@echo "$(BLUE)\n📌 $(1)...$(NC)"
endef

define print_success
	@echo "$(GREEN)✅ $(1)$(NC)"
endef

define print_wait
	@echo "$(YELLOW)⏳ Aguardando $(1) segundos...$(NC)"
endef

.PHONY: setup
setup: ## Configuração completa do ambiente: inicialização, criação e teste do leilão
	@echo "$(BLUE)🚀 Iniciando setup do ambiente de leilões$(NC)"
	
	$(call print_step,"Iniciando containers")
	@$(MAKE) up
	$(call print_success,"Containers iniciados")
	$(call print_wait,5)
	@sleep 5
	
	$(call print_step,"Criando novo leilão")
	@$(MAKE) create
	$(call print_success,"Leilão criado")
	$(call print_wait,5)
	@sleep 5
	
	$(call print_step,"Listando leilões ativos")
	@$(MAKE) list
	$(call print_success,"Lista de leilões ativos obtida")
	$(call print_wait,5)
	@sleep 5
	
	$(call print_step,"Verificando leilões fechados")
	@$(MAKE) check
	$(call print_success,"Verificação de leilões fechados concluída")
	$(call print_wait,5)
	@sleep 5
	
	$(call print_step,"Executando testes")
	@$(MAKE) test
	$(call print_success,"Testes finalizados")
	
	@echo "$(GREEN)\n✨ Setup concluído com sucesso!$(NC)"
