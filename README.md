# 🔥 PokeAPI REST Server — Clean Architecture, Caching, JWT Auth

A RESTful API server in Go that fetches data from [PokeAPI](https://pokeapi.co), caches it, and provides secure authentication using JWT.  
Designed with a focus on **Clean Architecture**, **scalability**, and **clarity**.

> 📌 This project was developed as a learning exercise and is not intended for production use.

---

## 🚀 Features

| Feature | Description |
|--------|-------------|
| ✅ **Clean Architecture** | Separation by handler/usecase/repository layers |
| ✅ **Dependency Injection (DI)** | Modular, testable, and maintainable components |
| ✅ **JWT Authentication** | Secure login/register with access control |
| ✅ **In-memory Cache** | Faster response times via TTL-based caching |
| ✅ **Retry Mechanism** | Exponential backoff for resilient HTTP requests |
| ✅ **Graceful Shutdown** | Proper resource cleanup on exit |
| ✅ **Dockerized** | Fully containerized setup using Docker Compose |
| ✅ **Postman Collection** | Postman tests included for easy API usage |

---

## 🏗️ Project Structure

```
/config             -> Application environment & configuration
/modules
  /pokemon          -> Domain logic for Pokémon (handler, usecase, repo, model)
  /user             -> Domain logic for user (handler, usecase, repo, model)
/pkg
  /authentication   -> Shared authentication logic
  /cache            -> In-memory cache abstraction using go-cache
  /httpclient       -> HTTP client with retry capability
  /jwt              -> JWT token utilities
  /middleware       -> Echo middlewares (e.g., JWT validator)
  /utils            -> Utility helpers
/server             -> HTTP server setup, routes, DI
main.go             -> Application entry point
stubby.yaml         -> Stubby configuration for mocking PokeAPI npm's module stubby
```

---

## 🔐 API Endpoints

### Public

| Method | Endpoint                       | Description          |
|--------|--------------------------------|----------------------|
| GET    | `/health`                      | Health check         |
| POST   | `/api/v1/auth/login`           | Login (returns JWT)  |
| POST   | `/api/v1/user/register`        | Register new user    |

### Protected (JWT required)

| Method | Endpoint                                  | Description                |
|--------|-------------------------------------------|----------------------------|
| GET    | `/api/v1/pokemon/:pokemonName`            | Retrieve Pokémon data      |
| GET    | `/api/v1/pokemon/:pokemonName/ability`    | Get Pokémon ability info   |
| GET    | `/api/v1/pokemon/random`                  | Get a random Pokémon       |

---

## 🧠 Architectural Decisions

- **Clean Architecture**: Modular, extensible structure
- **go-cache Abstraction**: Pluggable interface for future Redis support
- **In-memory User Store**: Simplified storage for prototype use
- **Retry Strategy**: Resilient handling of flaky external APIs
- **JWT Authentication**: Stateless token-based security with expiration

---

## 🧪 Test

Run all tests with:

```bash
go test ./...
```

---

## 🧰 How to Run (Work Instructions)

### 1. Copy the environment file

```bash
cp .example.env .env
```

### 2. Run in development mode (Go only)

```bash
go run main.go
```

### 3. Run using Docker Compose

```bash
docker compose up --build
```

---

## ⚙️ Environment Configuration (`.env`)
DEFAULT_CACHE_DURATION = 600 seconds as 10 minutes
CLEANUP_CACHE_INTERVAL = 60 seconds as 1 minute

```env
APP_NAME=PokemonAPI
APP_URL=0.0.0.0:3000
APP_STAGE=development
APP_BODY_LIMIT=10M

POKEMON_BASE_URL=https://pokeapi.co/api/v2

JWT_ACCESS_SECRET_KEY=tokensecret
JWT_ACCESS_DURATION=3600

DEFAULT_CACHE_DURATION=600
CLEANUP_CACHE_INTERVAL=60
```

---

## 📫 API Testing (via Postman)

Postman collection is included to simplify testing.  
Or test manually:

```bash
curl -X POST http://localhost:3000/api/v1/user/register   -H "Content-Type: application/json"   -d '{"username": "test", "password": "1234"}'
```

---

## 🧠 Technical Demonstration

- Modular software architecture
- Secure authentication practices (JWT)
- Cache and retry strategies for resilience
- Strong documentation and developer experience

---

## 🚀 Potential Improvements

- [ ] Add Redis-based distributed cache
- [ ] Integrate Swagger/OpenAPI for documentation
- [ ] Connect to persistent database (PostgreSQL or SQLite)
- [ ] Add CI/CD with GitHub Actions
- [ ] Apply rate limiting middleware
- [ ] Add centralized structured logging
