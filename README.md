# User Management API

REST and gRPC user management API built with Go, MongoDB, JWT authentication, and a hexagonal architecture.

## How to Run

Prerequisites: `docker` and `make`.

```bash
make start
```

This builds the API, starts MongoDB, creates the unique email index, and exposes:

- HTTP: `:8080`
- gRPC: `:8081`
- MongoDB: `:27017`

Configurations are in [`config/config.yml`](config/config.yml).

To stop the services and delete the MongoDB volume:

```bash
make stop
```

## Unit Test

Unit tests use Go's `testing` package and mocked service/repository interfaces using [mockgen](https://github.com/uber-go/mock).

```bash
make test
```

## HTTP API

`POST /api/auth/register` and `POST /api/auth/login` are public. All `/api/users` routes require a bearer token.

### Register

```bash
curl -s http://localhost:8080/api/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Alice","email":"alice@example.com","password":"password123"}'
```

Response (`201`):

```json
{ "code": "0000", "desc": "success" }
```

### Login

```bash
curl -s http://localhost:8080/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password123"}'
```

Response (`200`):

```json
{
  "code": "0000",
  "desc": "success",
  "data": { "access_token": "<token>" }
}
```

Set token from the login response:

```bash
TOKEN='<token>'
```

### Create a user

```bash
curl -s http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Bob","email":"bob@example.com","password":"password123"}'
```

Response (`201`):

```json
{ "code": "0000", "desc": "success" }
```

### List users

```bash
curl -s http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN"
```

Response (`200`):

```json
{
  "code": "0000",
  "desc": "success",
  "data": [
    {
      "id": "6a81e2c979de1337600aec9a",
      "name": "Alice",
      "email": "alice@example.com",
      "created_at": "2026-08-16T16:18:17.307Z"
    }
  ]
}
```

### Get a user

```bash
curl -s "http://localhost:8080/api/users/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
```

Response (`200`):

```json
{
  "code": "0000",
  "desc": "success",
  "data": {
    "id": "6a81e2c979de1337600aec9a",
    "name": "Alice",
    "email": "alice@example.com",
    "created_at": "2026-08-16T16:18:17.307Z"
  }
}
```

### Update a user

At least one of `name` or `email` is required.

```bash
curl -s -X PATCH "http://localhost:8080/api/users/$USER_ID" \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"name":"Alice Smith","email":"alice.smith@example.com"}'
```

Response (`200`):

```json
{ "code": "0000", "desc": "success" }
```

### Delete a user

```bash
curl -s -X DELETE "http://localhost:8080/api/users/$USER_ID" \
  -H "Authorization: Bearer $TOKEN"
```

Response (`200`):

```json
{ "code": "0000", "desc": "success" }
```

### Error responses

| HTTP status | Code | Description | When |
| --- | --- | --- | --- |
| `400` | `0001` | `bad request` | Invalid body, email, or ID |
| `401` | `0001` | `bad request` | Missing or invalid JWT |
| `401` | `0003` | `invalid email or password` | Invalid login credentials |
| `404` | `0001` | `bad request` | User not found |
| `409` | `0002` | `user already exists` | Duplicate email |
| `500` | `9999` | `internal server error` | Unexpected server error |

## gRPC API

The bonus gRPC service provides `CreateUser` and `GetUser`. Install `grpcurl` and use the JWT returned by the HTTP login endpoint.

### CreateUser

```bash
grpcurl -plaintext \
  -import-path . \
  -proto proto/user.proto \
  -H "authorization: Bearer $TOKEN" \
  -d '{"name":"Carol","email":"carol@example.com","password":"password123"}' \
  localhost:8081 userproto.UserService/CreateUser
```

Response:

```json
{}
```

### GetUser

```bash
grpcurl -plaintext \
  -import-path . \
  -proto proto/user.proto \
  -H "authorization: Bearer $TOKEN" \
  -d "{\"id\":\"$USER_ID\"}" \
  localhost:8081 userproto.UserService/GetUser
```

Response:

```json
{
  "user": {
    "id": "6a81e2c979de1337600aec9a",
    "name": "Alice",
    "email": "alice@example.com",
    "createdAt": "2026-08-16T16:18:17.307Z"
  }
}
```

## Design decisions and assumptions

- `domain`, `service`, and `port` contain business rules and interfaces; HTTP, gRPC, and MongoDB are adapters.
- No `internal` and `cmd` for easier file navigation (personal preference). Use `handler` as inbound and `adapter` as outbound.
- Use `yml` config to support more data types. We can later override the `yml` config with `env` if desired.
- Passwords are stored as bcrypt hashes. Emails are trimmed, lowercased, and protected by a MongoDB unique index.
- User IDs are MongoDB ObjectID hex strings. Password hashes are never returned.
- Registration and Login are public; CRUD and all gRPC methods require authentication. Any authenticated user can manage users (no access control).
- User listing is intentionally unpaginated for the scope of this exercise.
- HTTP and gRPC requests are logged with method/path or RPC method and execution time.
- SIGINT and SIGTERM trigger graceful shutdown of both servers, job, and the MongoDB connection.
- Omit unit test for `mongodb` because it gives less value compared to integration test.
