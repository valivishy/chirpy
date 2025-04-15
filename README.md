# Chirpy

Chirpy is a mini Twitter-like web application written in Go, designed as a learning project.

## Features

- User registration and login with hashed passwords
- JWT-based authentication
- Posting and retrieving "chirps"
- Optional sort parameter for listing chirps (asc/desc)
- RESTful JSON API
- Basic webhook handling
- Admin endpoints

## Setup

1. **Install dependencies:**

   ```bash
   
   go mod tidy
   ```

2. **Set up Postgres (with trust auth):**

   ```bash
   
   docker run --name chirpy-db \
     -e POSTGRES_HOST_AUTH_METHOD=trust \
     -e POSTGRES_DB=chirpy_test \
     -p 5432:5432 \
     -d postgres:15
   ```

3. **Run migrations:**

   ```bash
   
   goose -dir migrations postgres "postgres://postgres@localhost:5432/chirpy_test?sslmode=disable" up
   ```

4. **Run the app:**

   ```bash
   
   go run main.go
   ```

## Testing

 ```bash
 
 go test ./...
 ```

To run tests with a local Postgres instance, make sure the `chirpy_test` database is running and reachable at `localhost:5432`.

## API Endpoints

- `POST /api/users` — Register
- `POST /api/login` — Authenticate
- `POST /api/chirps` — Post a chirp (JWT required)
- `GET /api/chirps?sort=asc|desc` — List chirps
- `POST /admin/reset` — Reset DB (test only)
- `GET /admin/metrics` — App metrics (test only)

## Environment

You can create a `.env` file in the `/tests` folder for test configs:

 ```env

 DB_URL=postgres://postgres@localhost:5432/chirpy_test?sslmode=disable
 PLATFORM=TEST
 SECRET=...
 POLKA_KEY=...
 ```