# 🚀 Desafio: Implementação de Fechamento Automático de Leilões

## 🎯 Objetivo
Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

## 📋 Contexto
Use o repositório `labs-auction-goexpert` como base. Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

## ⚙️ Funcionalidades a Serem Desenvolvidas

### 1. Função de Cálculo de Tempo
**Desenvolver uma função que:** 
- Calcule o tempo do leilão baseado em parâmetros previamente definidos em variáveis de ambiente
- Utilize variáveis de ambiente para configuração flexível

### 2. Goroutine de Fechamento Automático
**Implementar uma nova goroutine que:**
- Valide a existência de leilões vencidos (que o tempo já se esgotou)
- Realize o update fechando o leilão (auction)
- Execute periodicamente baseado em intervalo configurável

### 3. Testes Automatizados
**Criar testes para validar:**
- Se o fechamento está acontecendo de forma automatizada
- Cenários de concorrência e condições de corrida
- Comportamento em diferentes intervalos de tempo

## 📍 Foco de Implementação
**Arquivo principal:** `internal/infra/database/auction/create_auction.go`

**Atenção especial para:**
- Trabalho com concorrência (goroutines)
- Mecanismos de sincronização thread-safe
- Análise do cálculo de intervalo existente na rotina de criação de bid

## 🛠️ Tecnologias e Conceitos
- Go Routines para concorrência
- Variáveis de ambiente para configuração
- Testes automatizados em Go
- Synchronization patterns

## 📦 Entrega Esperada

### ✅ Código-Fonte
- Implementação completa das funcionalidades
- Código limpo e bem documentado
- Tratamento adequado de erros

### 📚 Documentação
- Explicação de como rodar o projeto em ambiente dev
- Instruções de configuração das variáveis de ambiente
- Guia de execução dos testes

### 🐋 Containerização
- Dockerfile para construção da aplicação
- docker-compose.yml com PostgreSQL e aplicação
- Scripts de inicialização e migração
- Health checks para dependências

## 💡 Dicas Importantes

### Concorrência
- Implemente solução robusta para trabalho concorrente
- Analise como o cálculo de intervalo é feito na criação de bids
- Considere condições de corrida e race conditions

### Performance
- Evite locking desnecessário no banco de dados
- Use batch operations para múltiplos leilões
- Implemente backoff para erros temporários

### Boas Práticas
- Siga os padrões do projeto base
- Mantenha a consistência com o código existente
- Documente as novas funcionalidades

## 🔧 Variáveis de Ambiente Sugeridas
```bash
AUCTION_DURATION=5m              # Duração padrão do leilão
AUCTION_CHECK_INTERVAL=30s       # Intervalo de verificação
AUCTION_BATCH_SIZE=10            # Tamanho do lote para processamento