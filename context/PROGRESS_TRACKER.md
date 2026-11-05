# Progress Tracker

Where the Opal Labs API stands and what's next, so a new session can pick up without rediscovery.
Update "Current state" and "Next steps" and add a dated session-log entry after each work session.
Collapse old log entries into "Earlier history" once they stop mattering. Frontend state lives in the frontend repo's tracker.

## Current state
- Deployed on Render from `main`, Postgres on Neon; CI green (build, vet, API tests against a throwaway Postgres).
- Render env: `SECRET`, `DATABASE_URL` (Neon, direct not pooled), `FRONTEND_URL` (https://opal-labs.vercel.app). No Redis in production.
- `page-update` accepts page `etc` (icon, cover, fullWidth, locked) since Nov 2026; covered by `TestUpdatePageEtc`.

## In progress
- Nothing.

## Next steps
1. Upload `.github/og-image.png` as the repo's social preview (GitHub Settings → General → Social preview; no API for it).
2. Verify on production with a signed-in user that page icons/covers persist after reload.
3. Optional: invalidate issued JWTs on logout (e.g. a token version per user); today they last 30 days.
4. Optional: a `POST /page-duplicate` so the frontend doesn't need create + update round trips.

## Session log
### 2026-11-04 (late night)
- Agent context moved into this repo (`AGENTS.md`, `CLAUDE.md`, `context/`) from the old shared workspace folder. Added `.github/og-image.png` (1280×640) and used it as the README cover.
- Welcome page's last line now says pages save automatically (the frontend autosaves).

### 2026-11-04 (night)
- `PageUpdateObject` gained `Etc`; it is written only when present (`Status == Present`). New `TestUpdatePageEtc` (two test fixes: create needs `element_positions`, and the cleanup `defer` must wrap the request in a func).

### Earlier history
- Revival on 2026-11-04: Bearer-token auth replaced cross-site cookies (JWT `{user_id, exp}` instead of embedding Google tokens; user existence checked on auth); ElephantSQL (dead) replaced by Neon; Redis made optional after a startup crash on empty `REDIS_DB`; CORS accepts trailing slashes and comma-separated origins and fails fast when unset; CI rebuilt around a Postgres service container.
