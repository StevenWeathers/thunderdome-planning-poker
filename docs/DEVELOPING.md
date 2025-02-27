# Development

## Code Standards

### Frontend

- Follow [Typescript](https://www.typescriptlang.org/) and [Svelte](https://svelte.dev/) best practices

### Go

- [Code Review Comments guide](https://go.dev/wiki/CodeReviewComments) from the Go project should be
  followed

## Developing locally

There are a few ways to run the application locally, the easiest way is to
use [Docker Compose](https://docs.docker.com/compose/), but you can also run
the application directly on your machine provided you have the dependencies installed.

- [Building and running with Docker Compose](#building-and-running-with-docker-compose)
- [Build and Run without Docker Compose](#build-and-run-without-docker-compose)
- [Creating SQL Migrations](#creating-sql-migrations)
- [Adding new Localizations](#adding-new-localizations)
- [Format code before commit](#format-code-before-commit)

## Building and running with Docker Compose

```
docker-compose up --build
```

## Build and Run without Docker Compose

### Prerequisites

- [Postgres](https://www.postgresql.org/download/)
- [Go](https://golang.org/dl/)
- [Node.js](https://nodejs.org/en/download/)
- [Task](https://taskfile.dev/#/) (recommended) is used to manage development tasks in the project.
- [Goose](https://github.com/pressly/goose) (recommended) is used to manage SQL migrations
- [Caddy Server](https://caddyserver.com/) (optional) is used to run locally with HTTPS

### Install dependencies

With [Task](https://taskfile.dev/#/) running the `task install` command will install dependencies

OR manually run the following commands

```
go mod download
cd ui && npm install
```

#### Build and run

Postgres must be running with the same configuration as in the `docker-compose.yml` file.

With [Task](https://taskfile.dev/#/) run the `task dev` command to build and run the application then visit
[http://localhost:8080](http://localhost:8080) in your browser.

Or `task dev-secure` to run with HTTPS (using Caddy) then
visit [https://thunderdome.localhost](https://thunderdome.localhost) in your browser.

### Format code before commit

- Using the above dev Task(s) will automatically handle this for you.

If you've setup Task simply `task format` to format all go and frontend code before commit, otherwise `make format` is
still available for those with make installed.

- not doing so can result in a failed build on GitHub

### Restful API Changes

- Using the above dev Task(s) will automatically handle this for you.

The restful API is documented using swagger comments in the Go code, any changes to that documentation require
regenerating the docs with the following commands and committing the updated docs with the changes.

With [Task](https://taskfile.dev/#/) `task gen-swag` will regenerate the swagger docs

OR manually run the following commands

```bash
go tool swag fmt
go tool swag init -g internal/http/http.go -o docs/swagger
```

### Creating SQL Migrations

Generate new migration file using the following command, replacing `SHORT_DESCRIPTIVE_FILNAME` with a short descriptive
name for the migration example `create_poker_table`.

```
go tool goose -dir internal/db/migrations create SHORT_DESCRIPTIVE_FILNAME sql
```

### Adding new Localizations

**Thunderdome** supports Locale selection on the UI (Default en-US)

Adding new locale's involves just a few steps.

1. First by copying the `ui/src/i18n/en/index.ts` into the new locale folder
   at `ui/src/i18n/{locale}/index.ts` using the two letter locale code for the directory name and translating all
   the values.
2. Second, the locale will need to be added to the locales list used by switcher component
   in ```ui/config.js``` ```locales``` object
3. Run `npm run locales` in `ui` directory to generate the new locale types used by the build process
4. commit changes and open PR