# ArcRank

ArcRank is a lightweight leaderboard + player search service built with **Go** and the **Gin** web framework.

It is designed to demonstrate:

* REST API development in Go
* MongoDB operations
* Search capabilities powered by Elasticsearch
* Simple service-layer architecture
* Containerized load testing with Grafana k6

ArcRank sustained ~400 update requests/sec with 50 concurrent users. Median latency remained under 10ms, with p95 under 100ms.

---

## Features

* Create and manage players
* Retrieve player profiles
* Update player data
* Query leaderboard rankings
* Username search
* Health check endpoint

---

## Project Structure

```text
.
├── cmd/                # Application entrypoint
├── config/             # Configuration loading
├── internal/
│   ├── db/             # MongoDB + Elasticsearch clients
│   ├── handler/        # HTTP handlers (Gin)
│   ├── model/          # Data models
│   ├── route/          # Router setup
│   └── service/        # Core business logic
├── scripts/            # Utilities (seeding players)
├── k6/                 # Load testing scripts
├── docker-compose.yaml # Local infra + k6 runner
└── Dockerfile          # Server container build
```

Server logic lives in `internal/`, following standard Go project conventions.

---

## Getting Started

### Prerequisites

* Go 1.24+
* Podman (or Docker) + Compose
* Python 3.9+ (optional, for seeding players)

---

## Running with Compose

ArcRank is fully containerized, including its dependencies and load testing tools.

The recommended workflow is:

1. Start infrastructure (MongoDB + Elasticsearch)
2. Start the API server
3. Seed sample players
4. Run load tests with k6

---

### 1. Start Infrastructure

Run MongoDB and Elasticsearch first:

```bash
docker-compose up -d mongo es
```

---

### 2. Start the Server

Once infra is ready, start the API server:

```bash
docker-compose up -d server
```

The API will be available at:

```
http://localhost:8080
```

---

### 3. Seed Sample Players (Optional)

Before running load tests, you may want to populate the database.

The project includes a Python seeding script:

```bash
python scripts/seed_players.py
```

This creates ~10,000 random players and writes their IDs into:

```
k6/scripts/ids.json
```

Required Python dependency:

```bash
pip install requests
```

---

## API Endpoints

### Health

```http
GET /health
```

---

### Players

```http
POST   /players/
GET    /players/:id
PATCH  /players/:id
```

---

### Leaderboard

```http
GET /leaderboard/top
```

---

### Search

```http
GET /search?q=<username>&limit=50
```

---

## Load Testing (k6)

ArcRank includes containerized load testing using **k6**.

Because k6 runs inside Docker, you can:

* Reproduce performance tests consistently
* Control CPU/memory limits via Compose
* Benchmark the server under realistic constraints

Load test scripts live in:

```
k6/scripts/
```

### Available Tests

* `leaderboard.js` — leaderboard read throughput
* `search.js` — username search performance
* `update.js` — concurrent player updates

---

### Running Load Tests

Run a test using the bundled k6 container:

```bash
docker-compose run --rm k6 run /scripts/leaderboard.js
```

Example:

```bash
docker-compose run --rm k6 run /scripts/search.js
```

---

## Notes

This project intentionally keeps the architecture simple:

* No repository abstraction layer
* Minimal transaction handling
* Direct service-to-database interaction

The focus is on clarity and demonstrating core backend fundamentals.