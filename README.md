Certainly! Here's the README documentation converted to English:

markdown

# Go-Rich

This project is an example of how to use RabbitMQ and Go to process exchange rates.

## Overview

The project consists of three main parts:

1. **Server (API)**: A Go API that receives requests to fetch exchange rates.

2. **Sender (Sender)**: Sends exchange rate requests to a RabbitMQ server.

3. **Worker (Worker)**: A Go worker that consumes messages from RabbitMQ, fetches exchange rates, and sends the results back to the API.

## Requirements

- [Go](https://golang.org/dl/)
- [RabbitMQ](https://www.rabbitmq.com/download.html)

## Configuration

Before running the project, you need to create a free account at [CurrencyFreaks](https://currencyfreaks.com/) to obtain an API access key.

Create a .env file in the api folder with the following content:

API_KEY=YOUR_API_KEY

bash


## Running the Project

1. **Start RabbitMQ**:
   ```bash
   # Latest RabbitMQ 3.12
   docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

    Start the Worker:

    bash

go run ./workers/worker.go

Start the Server:

bash

    go run ./api/server.go

Make a GET request to the API at the following URL:

bash

http://localhost:8080/latest?currency=USD

Replace USD with the desired currency.
Contributions

Feel free to contribute to this project. Just open an issue or send a pull request.
License

This project is licensed under the MIT License.
