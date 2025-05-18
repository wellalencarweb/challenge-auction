
# 🏷️ Full Cycle Auction - GoExpert

Sistema de leilão desenvolvido em Go como parte do curso Go Expert da Full Cycle.

## 🚀 Funcionalidade Adicionada
Leilões agora são **encerrados automaticamente** após um tempo configurado via variável de ambiente (`AUCTION_DURATION_SECONDS`).

## 🧰 Tecnologias Utilizadas

- Golang
- MongoDB
- Docker & Docker Compose
- Go Modules
- Logrus (logger)
- Testes com `testing` padrão do Go

## 📁 Estrutura do Projeto

```
.
├── cmd/auction             # Entry point da aplicação
├── configuration           # Configuração de DB, logger, etc
├── internal
│   ├── entity              # Entidades de domínio
│   ├── usecase             # Casos de uso
│   └── infra
│       └── database
│           └── auction     # Repositórios de Auction (CRUD)
├── docker-compose.yml      # Orquestração com MongoDB
├── Makefile                # Comandos úteis de build, run e test
└── README.md
```

## ⚙️ Variáveis de Ambiente

As variáveis devem ser configuradas no arquivo `.env` dentro de `cmd/auction/.env`:

```env
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=auction_test_db
AUCTION_DURATION_SECONDS=60
```

## 🐳 Rodando com Docker

```bash
make up          # Sobe os containers
make down        # Para os containers
```

## ▶️ Executando o Projeto

```bash
make run
```

## 🧪 Executando os Testes

```bash
make test
```

## ✅ Comportamento Esperado

Ao criar um novo leilão, o sistema automaticamente inicia uma goroutine que aguarda o tempo configurado e encerra o leilão ao fim do prazo. Isso é feito sem afetar os lances já existentes, pois a verificação de status já está implementada na rotina de bid.

---

Desenvolvido com 💻 por Full Cycle GoExpert.
