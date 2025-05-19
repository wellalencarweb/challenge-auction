# Auction API

Projeto de leilões com fechamento automático utilizando goroutines.

## Rodando o projeto com Docker

```bash
make build
make run
```

A API ficará disponível em: `http://localhost:8080`

## Testando os requisitos

### 1. Criação de leilão com tempo definido

Use o endpoint `POST /auction` com o campo `duration_seconds` (ex: 60) para criar um leilão com duração automática.

### 2. Verificando encerramento automático

Após o tempo, o leilão será automaticamente encerrado. Você pode verificar isso com o endpoint:

```http
GET /auction/:auctionId
```

O campo `"status"` deverá retornar `"closed"` após o tempo definido.

### 3. Buscando o vencedor do leilão

```http
GET /auction/winner/:auctionId
```

Retorna o lance vencedor (caso haja lances).

## Documentação da API

A pasta [`docs/`](./docs) contém uma collection do Postman com exemplos de uso da API.
Importe em seu Postman e teste os endpoints com facilidade.
