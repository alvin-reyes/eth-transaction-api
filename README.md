
# Ethereum Transaction API

This is a Go application for managing Ethereum transactions, which provides API endpoints for managing accounts and transactions. The application supports database migration, seeding, and running an HTTP server with rate limiting capabilities. The following guide will help you set up the application for local development and explore its features.

## Table of Contents
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
  - [Database Migration](#database-migration)
  - [Database Seeding](#database-seeding)
- [Starting the Server](#starting-the-server)
- [CLI Features](#cli-features)
- [Rate Limiting](#rate-limiting)
- [Using SQLite](#using-sqlite)
- [API Endpoints](#api-endpoints)

## Installation

### Prerequisites
- [Go](https://golang.org/dl/) 1.20 or later
- [Git](https://git-scm.com/)
- SQLite (comes pre-installed with most systems, but can be installed via package managers if needed)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/ardata-tech/eth-transaction-api.git
   cd eth-transaction-api
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the application:
   ```bash
   go build -o eth-transaction-api
   ```

## Configuration

The application requires a `.env` file to specify settings such as the Ethereum RPC URL, rate limits, and server port. A sample `.env.example` file is provided, which you can rename to `.env` and modify according to your needs.

```bash
ETHERSCAN_API_KEY=
RATE_LIMIT=1
BURST_LIMIT=5
PORT=3000
```
### Environment Variables:
- **ETHERSCAN_API_KEY**: The API Key for Etherscan API.
- **RATE_LIMIT**: The number of requests per second allowed.
- **BURST_LIMIT**: The maximum burst of requests that can be handled.
- **PORT**: The port on which the server will run.

## Database Setup

The application uses SQLite as the database. This is a lightweight, file-based database that is easy to set up and use for local development.

### Database Migration (optional)

To set up the database schema, you need to run the migration command:

```bash
./eth-transaction-api --db:migrate
```

This will create the necessary tables for storing accounts and transactions.

### Database Seeding (optional)

To populate the database with initial data, use the seed command:

```bash
./eth-transaction-api --db:seed
```

This will insert sample data into the `accounts` table. The seeder automatically handles idempotency, so you can run it multiple times without duplicating data.

## Starting the Server

Once the database is set up, you can start the HTTP server to handle API requests:

```bash
./eth-transaction-api --server:start
```

By default, the server will run on port `3000`. You can change the port by modifying the configuration file or setting the `PORT` environment variable.

## CLI Features

The application supports the following command-line flags:

- `--db:migrate`: Migrates the database schema.
- `--db:seed`: Seeds the database with initial data.
- `--server:start`: Starts the HTTP server.

If no flag is provided, the application will display a usage message and exit.

## Rate Limiting

The application includes a rate limiter to control the number of requests per second. This is especially useful for preventing abuse or overloading the server.

The rate limiting is configured via the `RATE_LIMIT` and `BURST_LIMIT` variables in the environment file. The `RATE_LIMIT` defines how many requests per second are allowed, and the `BURST_LIMIT` defines how many requests can be handled in a burst.

## Using SQLite

The application uses SQLite as its database, which is stored in a file named `test.db`. SQLite is ideal for development and small-scale applications due to its simplicity and ease of setup.

### Managing the SQLite Database

- **Opening the Database**: The database is opened or created (if it doesn't exist) automatically when the application starts.
- **Migration**: Run `./eth-transaction-api --db:migrate` to create or update the database schema.
- **Seeding**: Run `./eth-transaction-api --db:seed` to insert initial data into the database.

The database file (`test.db`) will be created in the root directory of the project.

## API Endpoints

### Get Account Transactions

- **Endpoint**: `/accounts/{accountId}/transactions`
- **Method**: `GET`
- **Description**: Retrieves a list of transactions for the specified account.
- **Rate Limiting**: This endpoint is rate-limited as configured in the application.

Example request:

```bash
curl http://localhost:3000/accounts/9b3af3a7-51f1-49a7-aa3b-c700cf82a835/transactions
```

### Get Pooled Eth Transactions
**WIP**

### Get last 5 sETH transactions
**WIP**

