# Architecture

How the Opal Labs API is put together: routes, data model, auth, caching and the GORM details that bite.
Read before changing endpoints, models or auth. Update when a route, model field or flow changes;
keep to facts that aren't obvious from a single file.

## Routes (`main.go`)
| Route | Purpose |
|---|---|
| `GET /health` | Wake-up/keep-alive ping; doesn't touch the DB |
| `POST /user-login` | Called server-side by the frontend's next-auth with Google tokens; upserts the user (new users get the welcome page) and returns `{token}` |
| `POST /user-logout` · `GET /user-get` | Revoke the Google token · current user's `{status, name, email, picture}` |
| `POST /page-create` · `GET /page-list` | Create a page (client generates `page_uuid`) · all the user's pages (favourites first, then most recently updated) |
| `GET /page-get/:page_uuid` | Owner: full page + elements + `sub_pages`. Anyone: a **public** page (and its view count goes up) |
| `POST /page-update` | Partial page update + optional full element list (see below) |
| `POST /page-delete` | Deletes the page and its sub-pages |

## Data model (`models/models.go`, auto-migrated on start)
- `users`: Google id/email/name/picture, `credentials` (Google tokens, JSONB), `status` ∈ freemium/premium/enterprise (default premium).
- `pages`: `page_uuid`, `page_name`, `is_root`, `parent_page_uuid`, `public_page`, `is_favourite`, `element_positions` (ordered element uuids), `view_count`, `date_view_count` (JSONB `{date: n}`), `etc` (JSONB page extras from the frontend: `{icon, cover, fullWidth, locked}`).
- `elements`: `element_uuid`, `type` (`varchar(20)`), `content` and `etc` (JSONB `{text}`), `size`.

## Page-update semantics (`controllers/page.go` `PageUpdate`)
- Only fields present are written: `page_name: ""` keeps the name, `etc` missing or `null` keeps the old value (`{}` clears it), `is_root`/`public_page`/`is_favourite` are pointers.
- When `elements` is present it is the complete ordered list: unknown uuids are created, missing ones deleted, `element_positions` rewritten. Without `elements`, blocks are untouched.
- Toggling `public_page` resets `view_count` and `date_view_count`.
- GORM gotcha: `Updates(PageUpdateObject)` maps fields by column name onto `models.Page`; zero values are skipped, so use `Omit(...)` for fields that must not change.

## Auth (`controllers/auth/auth.go`)
- `MakeTokenString` signs HS256 `{user_id, exp: 30 days}` with `SECRET`. `AuthenticateUser` reads `Authorization: Bearer` (cookie fallback for tests), verifies signature and expiry, and checks the user exists.
- Logout revokes the Google token only; issued JWTs stay valid until they expire (accepted gap).

## Caching (`api/caching`, `initializers/caching.go`)
- Optional Redis cache of `page-get` responses keyed by path; invalidated on page update/delete (and the parent page). Without `REDIS_ADDR` every helper is a no-op.

## CI (`.github/workflows/ci.yml`)
- Go build + vet, then a Postgres service container, the server started with `APP_ENV=development` and health-checked, and `go test -v ./tests/...` against `DOMAIN=http://localhost:8000`.
