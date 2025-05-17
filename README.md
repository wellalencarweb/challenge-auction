
# 🏷️ Challenge Auction - Fechamento Automático de Leilões

Este projeto implementa uma nova funcionalidade para fechamento automático de leilões, com base em tempo configurável via variável de ambiente.

---

## 🚀 Funcionalidade Adicionada

- Ao criar um leilão, uma **goroutine** é iniciada automaticamente.
- Essa goroutine aguarda o tempo configurado (`AUCTION_DURATION_SECONDS`) e **fecha o leilão automaticamente**, atualizando seu status para `CLOSED`.

---

## ⚙️ Como Rodar o Projeto (Dev Environment)

### 1. Clone o repositório

```bash
git clone https://github.com/wellalencarweb/challenge-auction.git
cd challenge-auction
```

### 2. Configure a variável de ambiente

Edite ou crie um `.env` (ou configure diretamente no `docker-compose.yml`):

```env
AUCTION_DURATION_SECONDS=10
```

> Tempo em segundos que o leilão ficará aberto após sua criação.

### 3. Suba a aplicação com Docker

```bash
docker-compose up --build
```

---

## 🧪 Executando os Testes

### Pré-requisitos

- Ter o MongoDB rodando localmente em `mongodb://localhost:27017`
  (isso é feito pelo próprio `docker-compose`).

### Executar os testes:

```bash
docker exec -it challenge-auction-app go test ./...
```

Ou diretamente localmente, com Go instalado:

```bash
go test ./internal/infra/database/auction/...
```

O teste principal está em:

```
internal/infra/database/auction/create_auction_test.go
```

---

## 📁 Estrutura de Código Alterada

- `internal/infra/database/auction/create_auction.go`: Lógica de goroutine adicionada
- `internal/infra/database/auction/create_auction_test.go`: Teste de fechamento automático

---

## 🧠 Dica Técnica

- A goroutine usa `time.Sleep` com base no tempo da variável de ambiente para aguardar e então fechar o leilão.
- O fechamento é feito com `UpdateOne` no MongoDB, garantindo que apenas leilões com status `OPENED` sejam alterados.

---

## 📞 Suporte

Em caso de dúvidas, envie uma issue ou mensagem para o responsável do repositório.

---

Desenvolvido com 💡 por Go Expert 🚀
