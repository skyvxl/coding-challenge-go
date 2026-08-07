# Subscription Service

A small REST service for managing recurring user subscriptions. It was built as
part of my personal Go challenges to practice backend development and explore
tools commonly used in Go services.

## Features

- CRUD operations for subscriptions
- Total subscription cost calculation with filters
- PostgreSQL storage with Goose migrations and SQLC
- Swagger API documentation
- Structured request logging
- Docker Compose setup

## Run with Docker

```bash
docker compose up --build
```

The API will be available at `http://localhost:8090`, and Swagger UI at
`http://localhost:8090/swagger/index.html`.

To stop the project:

```bash
docker compose down
```

## Tests

```bash
go test ./...
```
