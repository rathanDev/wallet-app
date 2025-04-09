# Wallet App

A simple centralized wallet application built in Go, supporting basic financial operations such as deposit, withdrawal, and transfer of funds between wallets. It also allows users to check balances and view transaction history.

---

## Technologies Used

- Golang
- Gin Gonic
- GORM
- Viper
- Logrus
- PostgreSQL
- Docker (for database only)

---

## Available APIs

| API                     | Method | Endpoint                                |
|--------------------------|--------|------------------------------------------|
| Create Wallet            | POST   | `/wallets`                               |
| Get Wallets by User ID   | GET    | `/wallets/user/:userId`                  |
| Deposit Money            | POST   | `/wallets/:walletId/deposit`             |
| Withdraw Money           | POST   | `/wallets/:walletId/withdraw`            |
| Transfer Money           | POST   | `/wallets/:walletId/transfer`            |
| Get Balance              | GET    | `/wallets/:walletId/balance`             |
| Get Transactions         | GET    | `/wallets/:walletId/transactions`        |

---

## Assumptions or Decisions

- Only a single currency is supported (SGD assumed).
- All amounts are stored in cents to avoid floating point issues.  (100 cents = 1.00 SGD)
- User registration, authentication, and authorization are skipped for simplicity.
- Zero-amount transactions are allowed for now.

---

## Database Schema

### table - wallets 
id | user_id  | balance | created_at | updated_at 

### table - transactions 
id | wallet_id |  amount  | counterparty_wallet_id | trx_type | group_id | created_at

---

## 🧑‍💻 How to Run

### 🔄 Clone the repository
git clone https://github.com/rathanDev/wallet-app.git
cd wallet-app

### 🔄 Start Postgres database
If using docker, use the below command
docker run --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=app_db \
  -p 5432:5432 \
  -d postgres

### Update config.yaml if necessary

### Run the app 
go run main.go

### Run unit tests 
* Unit tests - API and Server logic 
    go test .\test\service\wallet_service_test.go -v
* Race condition test 
    go test .\test\service\wallet_service_race_condition_test.go -v

---

## Completed Features
* Wallet creation and listing APIs
* Deposit, Withdraw, Transfer APIs
* Get wallet balance API
* Get wallet transactions API
* Persistent database logic via PostgreSQL and GORM
* Race condition-safe operations:
    1. All money operations (deposit, withdraw, transfer) are wrapped in database transactions
    2. Relevant rows (wallet records) are explicitly locked using SELECT ... FOR UPDATE to ensure safe concurrent access
    3. This prevents issues like double spending or inconsistent balances during high concurrency
* Tests for edge cases, error handling, and race conditions 


## Areas for Improvement
Add Redis for caching
Implement API pagination and filters for transactions
Add Swagger/OpenAPI documentation
Improve error types and validation messages
Add retry/rollback logic for failed transactions
Use Docker Compose for full app + DB orchestration

---

## Time spent
10 - 12 hours 


## How should reviewer view my code
Code is organized by layers: controller -> service -> repo 
      main ---------> route -> controller -> service -> repo
      init db
      init config

Service contains core logic and race condition safety

Tests cover both successful and error scenarios


## Does it follow engineering best practices?
yes, it does

1. Layered architecture: Clear separation of concerns between controller (API layer), service (business logic), and repository (data access).

2. Transaction-safe operations: All wallet operations are wrapped in database transactions with record-level locking (SELECT FOR UPDATE) to ensure data integrity and avoid race conditions under concurrency.

3. Input validation: API inputs are validated to ensure correctness and prevent invalid operations (e.g., negative amounts).

4. Error handling: Consistent and informative error messages are returned from all layers, making debugging and debugging easier.

5. Code readability: Code is organized, formatted for maintainability.

6. Testing: Includes unit tests covering core logic and edge cases, along with dedicated tests for detecting race conditions.

7. Configuration management: Externalized configuration using Viper allows easy switching between environments without code changes.

8. Logging: Logrus is used for structured logging to aid in observability and debugging.



