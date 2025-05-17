
# Challenge Auction - Fechamento Automático de Leilões 🕒

Este projeto implementa uma funcionalidade de **fechamento automático de leilões** após um tempo definido via variável de ambiente.

## 🧱 Tecnologias

- Golang
- MongoDB
- Docker / Docker Compose
- Clean Architecture

## ▶️ Como Rodar o Projeto

### 1. Pré-requisitos

- Docker
- Docker Compose
- Go 1.18+

### 2. Configurar Variáveis de Ambiente

Crie um arquivo `.env` com a seguinte variável:

```env
AUCTION_DURATION_SECONDS=5
```

### 3. Subir o Projeto

```bash
make run
```

### 4. Executar Testes

```bash
make test
```

## 📂 Funcionalidade Implementada

Ao criar um leilão (`auction`), o sistema inicia uma goroutine que:
- Lê a variável `AUCTION_DURATION_SECONDS`.
- Aguarda o tempo especificado.
- Atualiza o status do leilão para `CLOSED`.

## 🧪 Teste Automatizado

Arquivo: `internal/infra/database/auction/create_auction_test.go`

Valida que o leilão é fechado automaticamente após o tempo.

## 🧼 Boas Práticas

- Clean Architecture
- Go Routines
- Testes Automatizados
- Makefile

---

Desenvolvido para o desafio Go Expert 🚀
