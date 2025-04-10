# Car Dealership Web App (Go)

## Description
This is a car dealership web application built with Go, following Clean Architecture principles. It supports user registration, authentication, car browsing, reservations, test drives, and purchasing via deposit balance.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/ginsengq/GO-Project.git
cd GO-Project
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Set Up Environment Variables

Create a .env file based on .env.example:

```bash
cp .env.example .env
```

Edit the .env file as needed, specifying your database credentials and other config values.

### 4. Run the Server

```bash
go run cmd/app/main.go
```

### 5. Example API Endpoints (Use Postman or curl)

- POST /auth/sign-up
- POST /auth/sign-in
- GET /cars
- POST /orders
- etc.

## Project Structure

- `cmd/` — entry point (main.go)
- `internal/` — business logic (usecases, delivery, repositories)
- `pkg/` — reusable libraries/utilities
- `entity/` — domain models
- `.env` — environment configuration

## Technologies

- Go 1.21+
- PostgreSQL
- Clean Architecture (Dr. Marten)
- REST API
- JWT Auth

## Running Tests

```bash
go test ./...
```

## Author

Galymzhankyzy Diana, 2025

# GO-Project