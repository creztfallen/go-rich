# Go-Rich

Este projeto é um exemplo de como usar RabbitMQ e Go para processar taxas de câmbio.

## Visão Geral

O projeto consiste em três partes principais:

1. **Server (API)**: Uma API Go que recebe solicitações para buscar taxas de câmbio.

2. **Sender (Sender)**: Envia solicitações de taxas de câmbio para um servidor RabbitMQ.

3. **Worker (Worker)**: Um worker Go que consome mensagens do RabbitMQ, busca as taxas de câmbio e envia os resultados de volta para a API.

## Requisitos

- [Go](https://golang.org/dl/)
- [RabbitMQ](https://www.rabbitmq.com/download.html)

## Configuração

Antes de executar o projeto, é necessário criar uma conta gratuita no [CurrencyFreaks](https://currencyfreaks.com/) para ter acesso a uma chave de acesso a API deles.

API_KEY=YOUR_API_KEY

## Executando o Projeto

1. **Inicie o RabbitMQ**:
   ```bash
   # latest RabbitMQ 3.12
    docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

3. **Inicie o Worker**:

   ```bash
   
   go run ./workers/worker.go
4. **Inicie o servidor**:

  ´´´bash
  
    go run ./api/server.go

Faça uma Solicitação de Taxa de Câmbio:

Faça uma solicitação GET para a API na seguinte URL:

bash

    http://localhost:8080/latest?currency=USD

    Substitua USD pela moeda desejada.

Contribuições

Sinta-se à vontade para contribuir para este projeto. Basta abrir uma issue ou enviar um pull request.
Licença

Este projeto está licenciado sob a MIT License.
