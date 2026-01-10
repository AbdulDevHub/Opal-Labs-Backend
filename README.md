# 🧠 Opal Labs — Backend

Backend services powering the Opal Labs productivity platform.

<img src="Cover Image.jpg" width="65%" />

<br />

[🖥 Frontend Repository](https://github.com/AbdulDevHub/Opal-Labs-Frontend)

<br />

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

- **Language:** Go
- **Database:** PostgreSQL
- **Caching:** Redis
- **Authentication:** Google OAuth 2.0
- **Testing:** Go testing framework
- **Methodology:** Agile

---

## 🛠 Setup

### 🔐 Environment Variables

The `.env.example` file contains a template for configuration.

```bash
cp .env.example .env
````

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

### ☁️ ElephantSQL (Hosted Option)

If you prefer not to run PostgreSQL locally:

1. Create an account:
   👉 [https://customer.elephantsql.com/login](https://customer.elephantsql.com/login)

2. Click **Create New Instance**

3. Select plan: **Tiny Turtle (Free)**

Update the following values in your `.env` file:

```bash
DB_USER="<User & Default Database>"
DB_PASSWORD="<Password>"
DB_NAME="<User & Default Database>"
DB_HOST="<Server>"
DB_PORT=5432
DB_SSL_MODE="disable"
```

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
