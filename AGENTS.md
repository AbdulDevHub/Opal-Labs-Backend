# AGENTS.md

Context for AI coding agents working on the Opal Labs **backend**: the Go REST API behind the notes app
(the Next.js frontend is the separate `Opal-Labs-Frontend` repo). Read this first, then `context/` for depth.
After any meaningful change to this project, update this file, `context/PROGRESS_TRACKER.md`, and the
"Project Overview" section of `README.md` to reflect the current state — do this proactively, without being asked.

## Overview
REST API for pages, nested sub-pages, page elements (blocks), publishing with view counts and Google sign-in.
Deployed on Render (https://opal-labs-backend-muej.onrender.com, free tier, sleeps when idle) with Postgres on Neon.

## Tech stack
Go 1.21, Gin, GORM + Postgres (`pgtype.JSONB` columns), golang-jwt v4, Google OAuth token check, optional Redis cache, testify.

## Commands
- Run: `go run .` (reads `.env` if present; listens on `PORT`, 8000 locally) · Hot reload: CompileDaemon (see README)
- Build/vet: `go build ./... && go vet ./...`
- Tests: `go test ./tests/...` against a **running** server with `APP_ENV=development` and `DOMAIN` set (CI starts Postgres + the server, see `.github/workflows/ci.yml`).
- Verify: go build ./... && go vet ./...
- Go may not be on PATH on the dev machine; install it or use a portable toolchain.

## Conventions & constraints
- Auth: `Authorization: Bearer <JWT {user_id, exp 30d}>` (cookie fallback only for tests). `auth.AuthenticateUser` also checks the user still exists.
- The test token `Only_for_testing1200332` works only when `APP_ENV=development`; never enable that in production.
- Startup fails fast without `SECRET` or `FRONTEND_URL` (CORS; trailing slashes trimmed, comma-separated list allowed). Redis is optional: no `REDIS_ADDR` → caching helpers no-op.
- Element `type` is `varchar(20)`; the frontend owns type names and styling strings (`content.text`, `etc.text`). Don't validate or rename them here.
- Page-update semantics: missing fields are left unchanged (`page_name: ""`, `null`/missing `etc`); `elements` present = full replacement list in order. Toggling `public_page` resets view counts.
- New-user welcome page comes from `controllers/templates/welcome/*.json`; the frontend's guest copy (`lib/data/welcomeTemplate.ts`) must match.
- Commits use gitmoji + conventional type; pushing straight to `main` is fine (Render deploys on push).

## Boundaries
- `.env*` files are permission-blocked: don't read/write them; document env vars in `README.md`.
- The repo is public: never commit secrets.

## More
- `README.md` — plain-language overview, env vars, setup, testing.
- `context/PROGRESS_TRACKER.md` — current state, next steps, session log.
- `context/ARCHITECTURE.md` — routes, data model, auth flow, caching, GORM gotchas.
