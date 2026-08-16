# XKCD Helper

XKCD Helper is a Go web service for downloading, indexing, and searching XKCD comics. It keeps comic metadata in PostgreSQL, builds an inverted index for search, provides HTTP API endpoints, and includes a small browser UI for logging in and searching comics.

## Installation

1. Clone the repository:

```sh
git clone https://github.com/arinamklvch/xkcd-helper.git
cd xkcd-helper
```

2. Create a local `.env` file. The values below are suitable for local development:

```sh
cat > .env <<'EOF'
DATABASE_URL=postgres://admin:admin@localhost:5433/xkcd?sslmode=disable
JWT_SECRET_KEY=xkcdSecretKey
EOF
```

Other application settings can be changed in config.yml.

3. Run the application:

```sh
make up
```

The service starts on `http://localhost:8081` by default. On startup, it runs database migrations and downloads any missing XKCD comics.

## Usage

### Use the web UI in your browser

1. After the application starts, create a local user:

```sh
docker compose exec postgres psql -U admin -d xkcd \
  -c "INSERT INTO users (login, password, role) VALUES ('admin', 'admin', 'admin');"
```

2. Open `http://localhost:8081` in your browser.

3. Log in through `http://localhost:8081/login`, then search comics from the browser.

### Use Swagger UI
Open the Swagger UI at:

```text
http://localhost:8081/swagger/
```

### Call the API directly

1. Get a JWT token:

```sh
TOKEN=$(curl -i -s -X POST http://localhost:8081/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"admin","password":"admin"}' \
  | sed -n 's/^Set-Cookie: JWT_token=\([^;]*\).*/\1/p')
```

2. Load a range of comics:

```sh
curl "http://localhost:8081/load-comics?from=1&to=10" \
  -H "Authorization: Bearer $TOKEN"
```

3. Get the latest stored comic:

```sh
curl "http://localhost:8081/last-comic" \
  -H "Authorization: Bearer $TOKEN"
```

4. Update the local comic database. This endpoint requires an admin user:

```sh
curl -X PUT "http://localhost:8081/update" \
  -H "Authorization: Bearer $TOKEN"
```
