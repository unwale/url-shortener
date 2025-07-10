# URL Shortener

A modern, production-ready URL shortener service built with Go, PostgreSQL, and Redis.

<div align="center">

![CI/CD](https://github.com/unwale/url-shortener/actions/workflows/ci.yaml/badge.svg)
[![codecov](https://codecov.io/github/unwale/url-shortener/graph/badge.svg?token=1K56JZJFLG)](https://codecov.io/github/unwale/url-shortener)
![License](https://img.shields.io/github/license/unwale/url-shortener)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?e&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?&logo=docker&logoColor=white)

</div>

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Environment Variables](#environment-variables)
  - [Running Locally](#running-locally)
  - [Running with Docker](#running-with-docker)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

This project allows users to shorten long URLs, redirect to original URLs, and track usage statistics. The backend is written in Go, with PostgreSQL for persistent storage and Redis for caching.

---

## Features

- Shorten long URLs with optional custom aliases
- Redirect to original URLs
- Track click statistics
- Caching with Redis
- Persistent storage with PostgreSQL
- Dockerized for easy deployment

---

## Architecture

The service follows a clean architecture with clear separation of concerns:

- **API Layer:** Handles HTTP requests and responses
- **Service Layer:** Business logic for URL shortening and redirection
- **Persistence Layer:** PostgreSQL for data storage, Redis for caching
- **Migrations:** Managed with [migrate](https://github.com/golang-migrate/migrate)
- **Code Generation:** SQL queries generated with [sqlc](https://github.com/kyleconroy/sqlc)


---

## Tech Stack

- **Go** (Golang)
- **PostgreSQL**
- **Redis**
- **Docker & Docker Compose**
- **sqlc** (type-safe SQL in Go)
- **migrate** (database migrations)
- **slog** (structured logging)

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Go 1.24+](https://go.dev/) (for local development)

### Environment Variables

Create a `.env` file in the project root with the following variables:

```
POSTGRES_DB=your_db
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
```

### Running

1. Clone the repository:
    ```sh
    git clone https://github.com/unwale/url-shortener.git
    cd url-shortener
    ```

2. Start services with Docker Compose:
    ```sh
    docker compose up --build
    ```

3. The API will be available at [http://localhost:8080](http://localhost:8080).

All services (API, PostgreSQL, Redis, migrations) are orchestrated via `docker-compose.yaml`. No manual setup required.

---

## API Endpoints

| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| POST   | `/api/shorten`   | Shorten a new URL          |
| GET    | `/:short_code`   | Redirect to original URL   |
| GET    | `/api/stats/:id` | Get statistics for a URL   |


---

## Testing

To run **unit** tests locally:

```sh
make test-unit
```

To run **integration** tests locally:
```sh
make test-integration
```

---

## Project Structure

```
db/           # Database migrations and queries
db/sqlc/      # Generated Go code from SQL queries
cmd/          # Application entrypoint
internal/     # Application logic
Dockerfile
docker-compose.yaml
sqlc.yaml
```

---

## Contributing

Contributions are welcome! Please open issues or submit pull requests.

---

## License

This project is licensed under the [MIT](./LICENSE.md)