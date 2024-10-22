# Transaction Processing Application

This application processes incoming POST requests to update a user’s balance based on transactions. It is built using **Golang** and **PostgreSQL**, and runs inside Docker containers.

## Features

- Process POST requests to adjust user balance based on "win" or "lost" transactions.
- Ensure that each transaction is processed only once.
- Every few minutes, the app cancels the latest 10 odd-numbered transactions.
- Dockerized for easy setup.

## Prerequisites

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## How to Run

1. **Clone the repository:**

   ```bash
   git clone https://github.com/IrakliGiorgadze/transaction-processor.git
   cd transaction-processor

2. **Build and run the app using Docker Compose:**

   ```bash
   docker-compose up --build

3. **Access the app at:**

 * http://localhost:8080