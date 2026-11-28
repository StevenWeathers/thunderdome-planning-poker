# AGENTS.md

Follow existing conventions first: read the closest comparable feature and match it. If this file and the code disagree,
follow the newest comparable code and report the discrepancy.

## Project Overview

Thunderdome is an open source agile collaboration app: planning poker, retrospectives, story mapping and team check-ins,
with teams and organizations. A Go 1.27 backend (`net/http`, PostgreSQL 17, websockets) serves a Svelte 5 + TypeScript
single page app (Vite, Tailwind CSS **3**), built with Node 26 and shipped inside one Go binary.

- Priorities: SaaS first, self-hosted second, always open source. Require only Postgres; no cloud-specific features.
- **Do not add AI-driven features.** The project does not accept them (see [README](README.md)).

## Repository Map

- `main.go`, `cmd/`: Cobra CLI (`serve`, `migrate`, `config`); `cmd/serve.go` wires services together.
- `thunderdome/`: domain types and sentinel errors shared by all layers.
- `internal/db/<domain>/`: data services (SQL). `internal/db/migrations/`: Goose migrations.
- `internal/http/`: REST handlers, routes (`http.go`), middleware, service interfaces (`types.go`).
- `internal/http/{poker,retro,storyboard,checkin}/`: websocket services, built on `internal/wshub/`.
- `internal/config/`: Viper configuration. Other `internal/` packages: email, auth helpers, sanitizing, webhooks.
- `ui/src/`: `pages/` (routes), `components/` (by domain; shared in `global/`, `forms/`, `table/`), `types/`, `i18n/`.
  Read [ui/AGENTS.md](ui/AGENTS.md) before UI work.
- `e2e/`: Playwright tests (`tests/`) and page objects and fixtures (`fixtures/`).
- `docs/`: guides and generated `swagger/`. `build/`: Dockerfiles, Helm charts, dev Caddyfile.

Go tests sit beside the source. UI tests are `*.test.ts` beside the source or in `__tests__/`.

## Commands

Run from the repository root. [Task](https://taskfile.dev) is the preferred runner; the `Makefile` is out of date.

```bash
# Install
task install                       # Go modules, npm ci in ui and e2e, Playwright browsers

# Run (Postgres required)
docker compose up -d db            # Postgres only
task dev                           # format, build and serve on http://localhost:8080
docker compose up --build          # whole stack in containers

# Build
task build                         # format, then build UI and binary
npm run build --prefix ui          # UI only

# Format and generate
task format                        # all formatting; also regenerates Swagger and i18n types
task gen-swag                      # docs/swagger
task gen-i8n                       # i18n types (the task name really is "i8n")
npm run format --prefix ui         # oxfmt write
npm run format --prefix e2e

# Check (as in CI)
golangci-lint run --timeout 3m
npm run prettier --prefix ui       # oxfmt check, despite the name
npm run prettier --prefix e2e

# Go tests
go test ./thunderdome ./internal/...
go test ./internal/http -run TestHandleTeamUpdateUser -v
TEST_DB_HOST=localhost go test ./internal/db   # database tests, skipped without it; they create and drop their own database

# UI tests (Chromium via Playwright)
npm test --prefix ui               # full suite with coverage
cd ui && npx vitest run src/components/team/UpdateUser.test.ts
cd ui && npx vitest run src/components/team/UpdateUser.test.ts -t "modal dialog"

# E2E tests (app must already be running on localhost:8080)
cd e2e && npm test
cd e2e && npx playwright test tests/team/checkin.spec.ts -g "redirects to login"
cd e2e && npm run test:api         # tests tagged @api; test:ui runs the rest

# Migrations
go tool goose -dir internal/db/migrations create <short_description> sql
```

## Code Style & Conventions

### Go

- Follow Effective Go and Go Code Review Comments. Format with `goimports`; lint rules are in `.golangci.yml`.
- Prefer the standard library. Justify any new dependency. Don't add a web framework or ORM.
- Keep `thunderdome/` free of database and HTTP concerns. Define limits there as constants and errors as
  `ErrX = errors.New("SCREAMING_SNAKE_CODE")`. Use camelCase `json` tags.
- Take `context.Context` first in new I/O code and use the `...Context` database methods.
- Build a feature as a slice: domain type, migration, data service, interface in `internal/http/types.go`, handlers,
  routes in `http.go`, wiring in `cmd/serve.go`. Reference: `usertitle.go` in `thunderdome/`, `internal/db/team/` and
  `internal/http/`.

**Data access**

- Use `database/sql` with parameterized SQL against `thunderdome.<table>`. Put new logic in queries and Go, not new
  stored procedures.
- Scope writes by owner and record (`WHERE organization_id = $1 AND id = $2`).
- Map `sql.ErrNoRows` and Postgres `23505` to domain errors; wrap other errors with the action name.
- Return empty lists as `make([]*T, 0)`, and check `rows.Err()`.
- Sanitize user supplied HTML with the injected `HTMLSanitizerPolicy` before storing it.

**HTTP**

- Write handlers as `func (s *Service) handleXxx() http.HandlerFunc`.
- Validate path IDs with `validate.Var(id, "required,uuid")`; trim body strings, then `validate.Struct`.
- Respond with `s.Success` / `s.Failure`. Return client facing errors with `Errorf` and the codes in `error.go`; other
  errors reach clients as "Internal error.". Log with `s.Logger.Ctx(ctx)` and `zap` fields; never log secrets.
- Add Swagger annotations to every endpoint.
- Protect every non-public route. Order matters: `a.userOnly(a.teamUserOnly(a.teamAdminOnly(h)))`.
- Register team resources on each applicable path: `/api/teams/{teamId}`, `/api/organizations/{orgId}/teams/{teamId}`,
  `/api/organizations/{orgId}/departments/{departmentId}/teams/{teamId}`.
- Match neighboring routes for feature flags and `subscribed...Only` middleware.
- Register websocket events in the game's `New`; add facilitator-only events to the facilitator set too.
- Add settings in `internal/config` with a default, pass them through `cmd/serve.go`, and document them in
  `docs/CONFIGURATION.md`.

### Svelte and TypeScript

- Format with `oxfmt` (`ui/.oxfmtrc.json`: 2 spaces, single quotes, semicolons, trailing commas, width 120). There is no
  ESLint.
- Use `<script lang="ts">` with runes (`$props` and a `Props` interface, `$state`, `$derived`), snippets, `onclick` and
  callback props. Don't copy legacy `export let` or `on:` code.
- Name components `PascalCase.svelte`. Put shared types in `ui/src/types/`; avoid new `any`.
- Order imports as nearby files do: Svelte, components and modules, then `import type` (observed, not enforced).
- Call APIs through the `xfetch` prop, not `fetch`. It rejects with `[error, response]`, so read
  `error?.[1]?.status`. Report results through the `notifications` prop.
- Reuse `global/`, `forms/` and `table/` components. Pass `apiPrefix` to components shared across org, department and
  team scopes.
- Define routes in `ui/src/config.ts` and register them in `App.svelte`. Honor `PathPrefix` and `AppConfig` flags.
- Style with Tailwind and `dark:` variants; allow for right-to-left (Persian). Keep labels and `aria-label`s.
- Never hardcode user facing text. Use `$LL.key()`, add the key to **all** of
  `ui/src/i18n/{de,en,es,fa,fr,it,pt,ru}/index.ts`, then run `task gen-i8n`.

## Testing & Verification

Add or update tests for every feature and bug fix. No coverage threshold is configured.

- **Go:** table driven tests with `t.Run` and `testify`. Test handlers with `httptest` and stubs that embed the service
  interface and override only what the test calls.
- **UI:** `vitest-browser-svelte` `render` with `vitest/browser`; `ui/vitest.setup.ts` loads English.
- **E2E:** use page objects and fixtures; keep tests atomic and parallel safe. Setup seeds the database directly
  (`DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASS`; defaults match `docker-compose.yml`).

Before finishing:

1. Run `task format` and keep the regenerated Swagger and i18n files that belong to your change.
2. Backend changes: run Go tests and `golangci-lint`.
3. Frontend changes: run `npm test --prefix ui` and both `oxfmt` checks.
4. Update docs for changed behavior, APIs or configuration.
5. State which checks failed or were not run.

CI also builds the app and runs the whole Playwright suite against Postgres 17, plus CodeQL.

## Agent Guardrails & Gotchas

**Don't edit by hand**

- `docs/swagger/*`: generated and committed; use `task gen-swag`.
- `ui/src/i18n/i18n-*.ts`: generated; use `task gen-i8n`.
- `ui/dist/`, `ui/coverage/`, `e2e/storage/`, `e2e/playwright-report/`: build and test output.
- `go.sum`, `package-lock.json`: change only through dependency commands.
- Merged migrations: never modify; add a new one with `Up` and `Down` sections.
- `ui/src/lib/`: vendored router, route parser and socket client; change only if the task requires it.

**Quirks**

- Go compilation embeds `ui/dist`. Without a UI build, create a placeholder `ui/dist/index.html` before Go tests or lint.
- `go test ./...` also picks up a Go package inside `ui/node_modules`; use the package list under Commands.
- `golangci-lint` must be built with Go 1.27 or newer, or it refuses to load. CI uses v2.15.
- `npm run check --prefix ui` (type check) currently reports hundreds of errors and is not run by CI. Don't treat it as a
  gate or fix unrelated errors; avoid adding new ones in files you touch.
- `TestHandleSocketClose` in `internal/wshub` has failed on Windows with no Go changes. Whether that is expected is
  unconfirmed, and CI runs on Linux. Report the failure; don't change the test unless the task is about it.
- `task` build, dev and format commands rewrite files. Review `git status` afterwards.
- "battle" means poker game, "plan" means poker story, "warrior" means user. Don't rename them in passing.
- JSON request bodies under `/api/` are size limited by middleware.
- Scripts use a CSP nonce. New inline scripts or external origins need a policy change in `internal/http/http.go`.
- Avoid breaking the public REST API; that requires a major version.

**Environment**

- Config comes from `config.yaml` or environment variables ([CONFIGURATION](docs/CONFIGURATION.md)).
- `db.host` defaults to `db`; set `DB_HOST=localhost` outside Compose. `task dev` also sets `APP_DOMAIN=localhost`,
  `COOKIE_SECURE=false`, `HTTP_SECURE_PROTOCOL=false` and `SMTP_ENABLED=false`.
- Real deployments must override `COOKIE_HASHKEY`, `CONFIG_AES_HASHKEY` and database credentials. Never commit secrets.

**Git and releases**

- Commit or push only when asked. Target `main`.
- Follow the issue style in history: `feat-<issue> Short description`, body `This resolves #<issue> by ...`.
- Keep PRs small, keep every section of the PR template, and add screenshots for visual changes.
- Keep comments brief and about why. Review all generated code ([CONTRIBUTING](docs/CONTRIBUTING.md)).
- The maintainer cuts releases with git tags; GitHub Actions then runs GoReleaser. The maintainer also manages
  deployments. Don't create or push tags, publish releases or deploy.
- Self-hosting is covered in [INSTALLATION](docs/INSTALLATION.md).
