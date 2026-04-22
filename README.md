# Real Time Chat

A real-time chat application built with Go, PostgreSQL, WebSockets, and a vanilla JavaScript frontend. The project provides a REST API for authentication and a WebSocket endpoint for live messaging across multiple chat rooms.

---

## Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [Running with Docker](#running-with-docker)
  - [Running Locally](#running-locally)
- [Database Migrations](#database-migrations)
- [API Reference](#api-reference)
  - [Authentication](#authentication)
  - [WebSocket](#websocket)
- [WebSocket Message Format](#websocket-message-format)
- [Frontend](#frontend)
- [Middleware](#middleware)
- [Swagger Documentation](#swagger-documentation)

---

## Features

- User registration and login with JWT authentication
- Password reset via email token
- Real-time messaging using WebSockets
- Multiple chat rooms with isolated broadcasts
- Rate limiting per IP address
- CORS support
- Swagger UI for API documentation
- Static frontend served directly from the Go server

---

## Tech Stack

| Layer       | Technology                        |
|-------------|-----------------------------------|
| Language    | Go 1.25                           |
| Database    | PostgreSQL 16 (via pgx/v5)        |
| WebSockets  | gorilla/websocket v1.5.3          |
| Auth        | golang-jwt/jwt v5                 |
| Hashing     | golang.org/x/crypto (bcrypt)      |
| Config      | joho/godotenv                     |
| API Docs    | swaggo/swag + swaggo/http-swagger |
| Frontend    | Vanilla JavaScript, HTML, CSS     |
| Container   | Docker Compose                    |

---

## Project Structure

```
real-time-chat/
├── cmd/
│   └── real-time-chat/
│       └── main.go               # Application entry point
├── config/
│   └── config.go                 # Environment variable loading
├── docs/
│   ├── docs.go                   # Generated Swagger docs
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── auth/
│   │   ├── jwt.go                # JWT generation and validation
│   │   ├── service.go            # Business logic: register, reset password
│   │   └── user.go               # User authentication
│   ├── chat/
│   │   ├── client.go             # WebSocket client: ReadPump, WritePump
│   │   ├── hub.go                # Hub: manages rooms and broadcasts
│   │   ├── message.go            # Message struct
│   │   └── room.go               # Room struct
│   ├── db/
│   │   ├── db.go                 # PostgreSQL connection pool and queries
│   │   └── migrations/           # SQL migration files
│   └── models/
│       └── user.go               # User model
├── server/
│   ├── auth_handlers.go          # HTTP handlers: register, login, password reset
│   ├── http.go                   # Route registration
│   ├── middleware.go             # CORS, rate limiting, auth middleware
│   └── websocket.go              # WebSocket upgrade and handler
├── web/
│   ├── index.html                # Single page frontend
│   ├── app.js                    # Frontend logic: auth + WebSocket chat
│   └── styles.css                # Styles
├── tests/                        # Integration and unit tests
├── docker-compose.yml            # PostgreSQL container
├── go.mod
├── go.sum
└── .env                          # Environment variables (not committed)
```

---

## Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) and Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI

Install golang-migrate:

```sh
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Make sure `$HOME/go/bin` is in your `$PATH`:

```sh
export PATH="$PATH:$HOME/go/bin"
```

---

### Environment Variables

Create a `.env` file in the root of the project:

```env
DB_USER=chatuser
DB_PASSWORD=chatpass
DB_HOST=localhost
DB_PORT=5432
DB_NAME=chatdb

JWT_SECRET=your_secret_key_here
```

---

### Running with Docker

Start the PostgreSQL database using Docker Compose:

```sh
docker compose up -d
```

This starts a PostgreSQL 16 container on port `5432` with the credentials defined in `docker-compose.yml`.

---

### Running Locally

1. Install dependencies:

```sh
go mod tidy
```

2. Apply database migrations:

```sh
migrate -database "postgres://chatuser:chatpass@localhost:5432/chatdb?sslmode=disable" -path internal/db/migrations up
```

3. Start the server:

```sh
go run ./cmd/real-time-chat/
```

The server will be available at `http://localhost:8000`.

---

## Database Migrations

Migrations are located in `internal/db/migrations/` and managed with `golang-migrate`.

**Apply all migrations:**

```sh
migrate -database "postgres://chatuser:chatpass@localhost:5432/chatdb?sslmode=disable" -path internal/db/migrations up
```

**Rollback the last migration:**

```sh
migrate -database "postgres://chatuser:chatpass@localhost:5432/chatdb?sslmode=disable" -path internal/db/migrations down 1
```

**Create a new migration:**

```sh
migrate create -ext sql -dir internal/db/migrations your_migration_name
```

This generates two files:
- `xxxx_your_migration_name.up.sql` — changes to apply
- `xxxx_your_migration_name.down.sql` — changes to revert

---

## API Reference

### Authentication

All REST endpoints are protected with CORS middleware and rate limiting (10 requests per minute per IP).

---

#### Register

```
POST /register
```

**Request body:**

```json
{
  "username": "john",
  "email": "john@example.com",
  "password": "secret123"
}
```

**Responses:**

| Status | Description           |
|--------|-----------------------|
| 201    | User registered       |
| 400    | Invalid request body  |
| 500    | Internal server error |

---

#### Login

```
POST /login
```

**Request body:**

```json
{
  "username": "john",
  "password": "secret123"
}
```

**Response:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

| Status | Description              |
|--------|--------------------------|
| 200    | Returns JWT token        |
| 401    | Invalid credentials      |
| 500    | Could not generate token |

---

#### Forgot Password

```
POST /forgot-password
```

**Request body:**

```json
{
  "email": "john@example.com"
}
```

**Response:**

```
reset link sent
```

| Status | Description      |
|--------|------------------|
| 200    | Email processed  |
| 400    | Invalid request  |

---

#### Reset Password

```
POST /reset-password
```

**Request body:**

```json
{
  "token": "your_reset_token",
  "new_password": "newSecret456"
}
```

| Status | Description               |
|--------|---------------------------|
| 200    | Password reset successful |
| 400    | Invalid request           |
| 401    | Invalid or expired token  |

---

### WebSocket

```
GET /ws?room_id={id}&user_id={id}
```

**Query parameters:**

| Parameter | Type   | Required | Description               |
|-----------|--------|----------|---------------------------|
| room_id   | int    | yes      | ID of the room to join    |
| user_id   | string | yes      | ID or username of the user|

**Example connection (JavaScript):**

```js
const ws = new WebSocket("ws://localhost:8000/ws?room_id=1&user_id=john");

ws.onopen = () => console.log("Connected");
ws.onmessage = event => console.log(JSON.parse(event.data));
ws.send("Hello room!");
```

---

## WebSocket Message Format

All messages broadcast to clients are JSON encoded:

```json
{
  "sender_id": "john",
  "room_id": 1,
  "content": "Hello room!",
  "time": 1745280000
}
```

| Field     | Type   | Description                          |
|-----------|--------|--------------------------------------|
| sender_id | string | ID of the user who sent the message  |
| room_id   | int    | ID of the room                       |
| content   | string | Text content of the message          |
| time      | int64  | Unix timestamp                       |

---

## Frontend

The frontend is a single-page application served from the `web/` folder at `http://localhost:8000`.

**Views:**

- **Login** — Authenticate with username and password. JWT is stored in `localStorage`.
- **Register** — Create a new account.
- **Forgot Password** — Request a password reset link by email.
- **Chat** — Join any room by entering a numeric room ID. Messages are displayed in real time.

**How rooms work:**

- Enter a room ID in the sidebar and click "Join".
- The client connects to the WebSocket endpoint with the room ID and your user ID.
- Messages sent by other users in the same room appear instantly.
- Rooms are created automatically on the server when the first user joins.

---

## Middleware

| Middleware         | Description                                               |
|--------------------|-----------------------------------------------------------|
| CORSMiddleware     | Sets permissive CORS headers for all origins              |
| RateLimitMiddleware| Limits each IP to 10 requests per minute                 |
| AuthMiddleware     | Validates JWT token from the Authorization header         |

---

## Swagger Documentation

The interactive API documentation is available at:

```
http://localhost:8000/swagger/
```

To regenerate the Swagger docs after modifying handler annotations:

```sh
swag init -g cmd/real-time-chat/main.go
```
