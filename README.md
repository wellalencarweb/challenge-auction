
# Desafio Go Routines - Abertura e Fechamento do Leilão

Este projeto é parte do desafio da Pós Go Expert da Full Cycle e implementa uma rotina de fechamento automático de leilões utilizando goroutines.

## ✅ Funcionalidade Implementada

- Leilões agora são automaticamente encerrados após um tempo configurado via variável de ambiente (`AUCTION_DURATION_SECONDS`).
- A verificação roda periodicamente em uma goroutine.
- Teste automatizado incluso para validar o fechamento automático.

---

## 🚀 Como rodar o projeto

### Pré-requisitos

- Docker
- Docker Compose

### Passos

```bash
# Subir os serviços
make up
```

O serviço Go será iniciado junto ao MongoDB.

---

## ⚙️ Variáveis de Ambiente

As variáveis estão no arquivo `.env`. A principal variável adicionada é:

```env
AUCTION_DURATION_SECONDS=30
```

Ela define o tempo (em segundos) que um leilão pode permanecer aberto.

---

## 🧪 Rodando os testes

Após o ambiente estar no ar, execute:

```bash
make test
```

Isso executará os testes unitários, incluindo o teste de fechamento automático de leilões.

---

## 🗂 Estrutura Relevante

- `internal/infra/database/auction/create_auction.go`: Contém a lógica de monitoramento automático dos leilões.
- `internal/infra/database/auction/auction_monitor_test.go`: Teste automatizado criado para este desafio.
- `docker-compose.yml`: Define os serviços da aplicação e MongoDB.

---

## 📦 Encerrando

```bash
make down
```
