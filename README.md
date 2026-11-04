<div align="center">

# 🧠 Opal Labs: Go REST API for a Notion-Style Notes App (Backend)

The **Go + Gin + PostgreSQL** REST API behind [Opal Labs](https://opal-labs.vercel.app/), a Notion-like note-taking app: pages, nested subpages, rich elements, publishing, and JWT auth with Google sign-in.

<img src="Cover Image.jpg" alt="Opal Labs backend: Go REST API cover image" width="65%" />

<br />

[🌐 Live Demo](https://opal-labs.vercel.app/) •
[🖥 Frontend Repository](https://github.com/AbdulDevHub/Opal-Labs-Frontend)

<br />

[![Go CI](https://github.com/AbdulDevHub/Opal-Labs-Backend/actions/workflows/ci.yml/badge.svg)](https://github.com/AbdulDevHub/Opal-Labs-Backend/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?logo=redis&logoColor=white)
![OAuth](https://img.shields.io/badge/Auth-Google%20OAuth-red)
![MIT License](https://img.shields.io/badge/License-MIT-green)

</div>

---

## ✨ Introduction

This is a **public copy of a private repository** where I collaborated with a **team of 7 developers** to build a Notion-like productivity platform.

This repository contains the **backend** services for the application.  
For the frontend implementation, visit 👉 **[Opal Labs Frontend](https://github.com/AbdulDevHub/Opal-Labs-Frontend)**

### 🧱 Backend Tech Stack

- **Language:** Go (Gin web framework, GORM)
- **Database:** PostgreSQL (Neon)
- **Caching:** Redis (optional)
- **Authentication:** Google OAuth 2.0 + JWT bearer tokens
- **Testing:** Go tests in GitHub Actions against a throwaway Postgres
- **Hosting:** Render
- **Methodology:** Agile

---

## 🛠 Setup

### 🔐 Environment Variables

Create a `.env` file in the project root (a missing `.env` is fine in production, where the host injects the variables):

```bash
# "development" = localhost URLs + test token allowed. Leave unset in production.
APP_ENV="development"
PORT=8000

# Must match SECRET in the frontend (signs the session token)
SECRET="change-me"

# Production only (ignored when APP_ENV=development)
DOMAIN="https://your-backend.onrender.com"
FRONTEND_URL="https://your-frontend.vercel.app"  # comma-separate several origins; trailing slashes are ignored

# Neon connection string. When set, the DB_* values are ignored.
DATABASE_URL="postgresql://user:password@ep-xxxx.region.aws.neon.tech/neondb?sslmode=require"

# Local Postgres alternative
DB_USER="postgres"
DB_PASSWORD="password"
DB_NAME="Website"
DB_HOST="localhost"
DB_PORT=5432
DB_SSL_MODE="disable"

# Redis (optional; caching is disabled if unreachable)
REDIS_ADDR="localhost:6379"
REDIS_PASSWORD=""
REDIS_DB=0
```

### 🔑 Authentication

`POST /user-login` returns a signed session token (`{ "token": "..." }`) that holds the user's ID and expires in 30 days. The frontend sends it on every request as `Authorization: Bearer <token>`. A cookie is not used because the frontend and backend are hosted on different sites, and browsers don't send cookies across them.

---

### 🔁 CompileDaemon (Hot Reload)

To enable hot-reload during development, install **CompileDaemon**:

📦 [https://pkg.go.dev/github.com/githubnemo/compiledaemon](https://pkg.go.dev/github.com/githubnemo/compiledaemon)

---

## 🐘 PostgreSQL Setup

### Local Installation

1. Download PostgreSQL:
   👉 [https://www.enterprisedb.com/downloads/postgres-postgresql-downloads](https://www.enterprisedb.com/downloads/postgres-postgresql-downloads)

2. During setup:

   - Set password to: `password`
   - Leave all other options as default

3. Open **pgAdmin4** and create a database named:

   ```text
   Website
   ```

---

### ☁️ Neon (Hosted Option)

ElephantSQL has shut down, so the hosted database is now [Neon](https://neon.tech) (free tier, wakes itself up after idling).

1. Create a free account and a new project at 👉 [https://console.neon.tech](https://console.neon.tech)
2. Click **Connect**, then copy the connection string (it ends in `?sslmode=require`)
3. Set it as `DATABASE_URL` in your `.env` (and in the Render dashboard for the deployed backend)

Tables are created automatically on first start (GORM auto-migration).

---

## ⚡ Redis Setup

### Local Installation

1. Download Redis:
   👉 [https://redis.io/download](https://redis.io/download)

2. Ensure Redis is running on:

   ```text
   localhost:6379
   ```

3. Defaults (already set in `.env.example`):

   - Password: `""` (empty)
   - Database: `0`

If your setup differs, update:

```env
REDIS_ADDR
REDIS_PASSWORD
REDIS_DB
```

---

### ☁️ Hosted Redis (RedisLabs)

1. Create an account:
   👉 [https://redislabs.com/](https://redislabs.com/)

2. Create a **free** Redis database

Update your `.env` file:

```bash
REDIS_ADDR="<Public Endpoint>"
REDIS_PASSWORD="<Default User Password>"
REDIS_DB=0
```

ℹ️ These values can be found under the **Configuration** tab in RedisLabs.

---

### 📦 Redis Usage Notes

- The backend will still run **without Redis**, but caching will be disabled.
- Start local Redis:

  ```bash
  redis-server
  ```

- Stop Redis:

  ```bash
  redis-cli shutdown
  ```

#### Common Redis Commands

- `SET key value`
- `GET key`
- `DEL key`
- `FLUSHALL`

---

## 🔐 Google OAuth 2.0

Follow the setup instructions provided in the team drive:

👉 [https://drive.google.com/drive/folders/1PWzpsJGXIDA_RnRRoEcJe_U5yvGC6s_U?usp=sharing](https://drive.google.com/drive/folders/1PWzpsJGXIDA_RnRRoEcJe_U5yvGC6s_U?usp=sharing)

---

## 🚀 Running the Server

```bash
go build
go run .
```

### Hot Reload Mode

```bash
CompileDaemon -command="./backend"
```

---

## 🧪 Running Tests

With the server running, open another terminal:

```bash
cd tests
go test -v
```

Run a specific test:

```bash
go test -run "test_function_name"
```

---

## ☁️ Testing Hosted Backend

Replace:

```text
localhost:8000
```

With:

```text
https://opal-labs-backend-muej.onrender.com
```

### Example

```text
https://opal-labs-backend-muej.onrender.com/page-get/123
```

---

## 🤝 Contributing

Contributions are welcome!

- Open an issue for bugs or feature requests
- Submit a PR for improvements

---

## 📄 License

This project is licensed under the **MIT License**.
See the [LICENSE](LICENSE) file for details.
