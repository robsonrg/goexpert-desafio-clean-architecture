# Desafio Pós Go Expert - Clean Architecture

> Este projeto contém a solução para o desafio de Clean Architecture da pós-graduação Go Expert da FullCycle.

(...) precisará criar o usecase de listagem das orders.

Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL

Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.

Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.

---

# Como executar a aplicação

1. Iniciar mysql

```sh
docker compose up -d
```

2. Executar o aplicação

```sh
cd cmd
go run main.go
```

# Testando a aplicação

### Web server

> Teste usando o plugin REST Client, do VS Code, com os arquivos em `/api` ou via cURL: 

1. Criando um pedido

```sh
curl --location 'http://localhost:8000/order' \
--header 'Content-Type: application/json' \
--data '{
    "id": "1",
    "price": 35.5,
    "tax": 0.1
}'

```

2. Listando os pedidos

```sh
curl --location 'http://localhost:8000/order'
```

### gRPC server

> Instale o gRPCurl para testar: https://github.com/fullstorydev/grpcurl

1. Criando um pedido

```sh
grpcurl -plaintext -d '{"id":"2","price": 89.9, "tax": 0.5}' localhost:50051 pb.OrderService/CreateOrder
```

2. Listando os pedidos

```sh
grpcurl -plaintext -d '{}' localhost:50051 pb.OrderService/ListOrders
```

### GraphQL server

> Acesse o playgrand para testar: http://localhost:8080

1. Criando um pedido

```graphql
mutation createOrder {
  createOrder(input: { id: "3", price: 70, tax: 1 }) {
    id
    price
    tax
    finalPrice
  }
}
```

2. Listando os pedidos

```graphql
query  {
  listOrders {
    id
    price
    tax
    finalPrice
  }
}
```

# Dev

### gRPC

Pré-requisitos

- Protocol buffer compiler, protoc, version 3 \
https://protobuf.dev/installation

- Plugin `gRPC in Go` \
https://grpc.io/docs/languages/go/quickstart

Configurando no projeto

```sh
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto
```

### GraphQL

Configurando no projeto

```sh
go run github.com/99designs/gqlgen init
go run github.com/99designs/gqlgen generate
```